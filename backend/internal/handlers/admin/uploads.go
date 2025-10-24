package admin

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tanydotai/tanyai/backend/internal/auth"
	"github.com/tanydotai/tanyai/backend/internal/config"
	"github.com/tanydotai/tanyai/backend/internal/httpapi"
	"github.com/tanydotai/tanyai/backend/internal/middleware"
	"github.com/tanydotai/tanyai/backend/internal/storage"
)

// UploadsHandler processes secure media uploads for admin users.
type UploadsHandler struct {
	storage storage.ObjectStorage
	policy  config.UploadConfig
	logger  *log.Logger
}

// NewUploadsHandler constructs UploadsHandler.
func NewUploadsHandler(store storage.ObjectStorage, policy config.UploadConfig, logger *log.Logger) *UploadsHandler {
	if logger == nil {
		logger = log.New(os.Stdout, "", 0)
	}
	return &UploadsHandler{storage: store, policy: policy, logger: logger}
}

var mimeExtensions = map[string]string{
	"image/jpeg":    ".jpg",
	"image/png":     ".png",
	"image/webp":    ".webp",
	"image/svg+xml": ".svg",
}

// Create handles secure image uploads and returns a public URL.
func (h *UploadsHandler) Create(c *gin.Context) {
	started := time.Now()
	claims, _ := middleware.GetClaims(c)

	if h.storage == nil {
		httpapi.RespondError(c, http.StatusInternalServerError, httpapi.ErrorCodeInternal, "storage not configured", nil)
		return
	}

	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, h.policy.MaxBytes)
	if err := c.Request.ParseMultipartForm(h.policy.MaxBytes); err != nil {
		if isMaxBytesError(err) {
			httpapi.RespondError(c, http.StatusRequestEntityTooLarge, httpapi.ErrorCodeValidation, fmt.Sprintf("file exceeds %d bytes", h.policy.MaxBytes), nil)
			return
		}
		httpapi.RespondError(c, http.StatusBadRequest, httpapi.ErrorCodeValidation, "invalid multipart payload", nil)
		return
	}
	defer func() {
		if c.Request.MultipartForm != nil {
			_ = c.Request.MultipartForm.RemoveAll()
		}
	}()

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		httpapi.RespondError(c, http.StatusBadRequest, httpapi.ErrorCodeValidation, "file field is required", nil)
		return
	}
	defer file.Close()

	buf := &bytes.Buffer{}
	if _, err := io.Copy(buf, file); err != nil {
		httpapi.RespondError(c, http.StatusBadRequest, httpapi.ErrorCodeValidation, "failed to read upload", nil)
		return
	}

	data := buf.Bytes()
	if int64(len(data)) > h.policy.MaxBytes {
		httpapi.RespondError(c, http.StatusRequestEntityTooLarge, httpapi.ErrorCodeValidation, fmt.Sprintf("file exceeds %d bytes", h.policy.MaxBytes), nil)
		return
	}
	if len(data) == 0 {
		httpapi.RespondError(c, http.StatusBadRequest, httpapi.ErrorCodeValidation, "file is empty", nil)
		return
	}

	headerType := normalizeMIME(header.Header.Get("Content-Type"))
	detected := normalizeMIME(http.DetectContentType(data))
	if detected == "application/octet-stream" && headerType != "" {
		detected = headerType
	}

	looksLikeSVG := isLikelySVG(data) || headerType == "image/svg+xml" || detected == "image/svg+xml"
	if looksLikeSVG {
		detected = "image/svg+xml"
	}

	if !h.isAllowed(detected) {
		httpapi.RespondError(c, http.StatusUnsupportedMediaType, httpapi.ErrorCodeValidation, "unsupported media type", nil)
		return
	}

	if headerType != "" && headerType != detected && !(looksLikeSVG && headerType == "image/svg+xml") {
		httpapi.RespondError(c, http.StatusUnsupportedMediaType, httpapi.ErrorCodeValidation, "content type mismatch", nil)
		return
	}

	if detected == "image/svg+xml" {
		if !h.policy.AllowSVG {
			httpapi.RespondError(c, http.StatusUnsupportedMediaType, httpapi.ErrorCodeValidation, "svg uploads are disabled", nil)
			return
		}
		sanitized, err := sanitizeSVG(data)
		if err != nil {
			httpapi.RespondError(c, http.StatusBadRequest, httpapi.ErrorCodeValidation, "invalid svg payload", nil)
			return
		}
		data = sanitized
	}

	ext := mimeExtensions[detected]
	if ext == "" {
		ext = filepath.Ext(header.Filename)
	}
	if ext == "" {
		ext = ".bin"
	}

	now := time.Now().UTC()
	key := fmt.Sprintf("uploads/%04d/%02d/%02d/%s%s", now.Year(), now.Month(), now.Day(), uuid.NewString(), ext)

	publicURL, err := h.storage.Put(c.Request.Context(), key, data, detected)
	if err != nil {
		httpapi.RespondError(c, http.StatusInternalServerError, httpapi.ErrorCodeInternal, "failed to store file", nil)
		h.logUpload(map[string]interface{}{
			"route":      c.FullPath(),
			"method":     c.Request.Method,
			"error":      err.Error(),
			"mime":       detected,
			"size":       len(data),
			"key":        key,
			"user_id":    userIDFromClaims(claims),
			"latency_ms": time.Since(started).Milliseconds(),
		})
		return
	}

	h.logUpload(map[string]interface{}{
		"route":      c.FullPath(),
		"method":     c.Request.Method,
		"mime":       detected,
		"size":       len(data),
		"key":        key,
		"url":        publicURL,
		"user_id":    userIDFromClaims(claims),
		"latency_ms": time.Since(started).Milliseconds(),
	})

	httpapi.RespondData(c, http.StatusCreated, gin.H{
		"url":         publicURL,
		"key":         key,
		"contentType": detected,
		"size":        len(data),
	})
}

func (h *UploadsHandler) isAllowed(mime string) bool {
	normalized := normalizeMIME(mime)
	for _, allowed := range h.policy.AllowedMIME {
		if normalizeMIME(allowed) == normalized {
			if normalized == "image/svg+xml" && !h.policy.AllowSVG {
				return false
			}
			return true
		}
	}
	return false
}

func isMaxBytesError(err error) bool {
	var maxErr *http.MaxBytesError
	if errors.As(err, &maxErr) {
		return true
	}
	return strings.Contains(strings.ToLower(err.Error()), "request body too large")
}

func normalizeMIME(value string) string {
	v := strings.ToLower(strings.TrimSpace(value))
	if v == "" {
		return ""
	}
	if idx := strings.Index(v, ";"); idx >= 0 {
		v = v[:idx]
	}
	return v
}

func isLikelySVG(data []byte) bool {
	trimmed := strings.TrimSpace(string(data))
	if len(trimmed) == 0 {
		return false
	}
	trimmedLower := strings.ToLower(trimmed)
	return strings.HasPrefix(trimmedLower, "<svg")
}

func userIDFromClaims(claims *auth.Claims) string {
	if claims == nil {
		return ""
	}
	if claims.UserID == uuid.Nil {
		return ""
	}
	return claims.UserID.String()
}

func (h *UploadsHandler) logUpload(fields map[string]interface{}) {
	payload, err := json.Marshal(fields)
	if err != nil {
		h.logger.Printf("{\"error\":\"log marshal failed\",\"message\":%q}", err.Error())
		return
	}
	h.logger.Println(string(payload))
}

func sanitizeSVG(data []byte) ([]byte, error) {
	decoder := xml.NewDecoder(bytes.NewReader(data))
	decoder.Strict = false

	var buf bytes.Buffer
	encoder := xml.NewEncoder(&buf)
	var skipDepth int
	seenRoot := false

	for {
		token, err := decoder.Token()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}

		switch t := token.(type) {
		case xml.StartElement:
			name := strings.ToLower(t.Name.Local)
			if skipDepth > 0 {
				skipDepth++
				continue
			}
			if name == "script" {
				skipDepth = 1
				continue
			}
			if !seenRoot {
				if name != "svg" {
					return nil, errors.New("svg: root element must be <svg>")
				}
				seenRoot = true
			}
			attrs := sanitizeSVGAttributes(t.Attr)
			if err := encoder.EncodeToken(xml.StartElement{Name: t.Name, Attr: attrs}); err != nil {
				return nil, err
			}
		case xml.EndElement:
			if skipDepth > 0 {
				skipDepth--
				continue
			}
			if err := encoder.EncodeToken(t); err != nil {
				return nil, err
			}
		case xml.CharData:
			if skipDepth > 0 {
				continue
			}
			if err := encoder.EncodeToken(t); err != nil {
				return nil, err
			}
		case xml.Comment, xml.Directive, xml.ProcInst:
			// skip dangerous nodes
		}
	}

	if !seenRoot {
		return nil, errors.New("svg: missing <svg> root element")
	}

	if err := encoder.Flush(); err != nil {
		return nil, err
	}

	sanitized := bytes.TrimSpace(buf.Bytes())
	if len(sanitized) == 0 {
		return nil, errors.New("svg: empty content after sanitization")
	}

	return sanitized, nil
}

func sanitizeSVGAttributes(attrs []xml.Attr) []xml.Attr {
	cleaned := make([]xml.Attr, 0, len(attrs))
	for _, attr := range attrs {
		name := strings.ToLower(attr.Name.Local)
		space := strings.ToLower(attr.Name.Space)
		value := strings.TrimSpace(attr.Value)

		if strings.HasPrefix(name, "on") {
			continue
		}
		lower := strings.ToLower(value)
		if strings.Contains(lower, "javascript:") {
			continue
		}
		if (name == "href" || (space == "xlink" && name == "href")) && isExternalReference(value) {
			continue
		}
		if name == "style" && containsExternalStyle(value) {
			continue
		}

		cleaned = append(cleaned, attr)
	}
	return cleaned
}

func isExternalReference(value string) bool {
	trimmed := strings.Trim(strings.TrimSpace(strings.ToLower(value)), "\"'")
	if trimmed == "" {
		return false
	}
	if strings.HasPrefix(trimmed, "#") {
		return false
	}
	if strings.HasPrefix(trimmed, "data:") {
		return true
	}
	if strings.HasPrefix(trimmed, "http://") || strings.HasPrefix(trimmed, "https://") {
		return true
	}
	if strings.HasPrefix(trimmed, "//") {
		return true
	}
	if strings.HasPrefix(trimmed, "javascript:") {
		return true
	}
	return false
}

func containsExternalStyle(value string) bool {
	lower := strings.ToLower(value)
	if strings.Contains(lower, "@import") || strings.Contains(lower, "javascript:") {
		return true
	}

	rest := lower
	for {
		idx := strings.Index(rest, "url(")
		if idx < 0 {
			break
		}
		rest = rest[idx+4:]
		end := strings.Index(rest, ")")
		if end < 0 {
			break
		}
		target := strings.Trim(rest[:end], " \"'")
		if isExternalReference(target) {
			return true
		}
		rest = rest[end+1:]
	}
	return false
}
