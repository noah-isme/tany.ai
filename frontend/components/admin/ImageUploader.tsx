"use client";

import { useCallback, useEffect, useMemo, useRef, useState, type ChangeEvent } from "react";
import { UploadCloud, RotateCcw } from "lucide-react";

import { useToast } from "@/components/admin/ToastProvider";

import { Button } from "./ui/button";
import { Input } from "./ui/input";

const ACCEPTED_TYPES = "image/png,image/jpeg,image/webp,image/svg+xml";

type UploadResponse = {
  url: string;
  key: string;
  contentType: string;
  size: number;
};

type ImageUploaderProps = {
  value?: string;
  onChange: (value: string) => void;
  onBlur?: () => void;
  id?: string;
  name?: string;
  accept?: string;
  disabled?: boolean;
};

type UploadState = {
  isUploading: boolean;
  progress: number | null;
  error: string | null;
};

const initialState: UploadState = {
  isUploading: false,
  progress: null,
  error: null,
};

export function ImageUploader({
  value,
  onChange,
  onBlur,
  id,
  name,
  accept = ACCEPTED_TYPES,
  disabled = false,
}: ImageUploaderProps) {
  const pushToast = useToast();
  const [state, setState] = useState<UploadState>(initialState);
  const [objectURL, setObjectURL] = useState<string | null>(null);
  const [lastFile, setLastFile] = useState<File | null>(null);
  const inputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    return () => {
      if (objectURL && typeof URL.revokeObjectURL === "function") {
        URL.revokeObjectURL(objectURL);
      }
    };
  }, [objectURL]);

  const updateObjectURL = useCallback((next: string | null) => {
    setObjectURL((current) => {
      if (current && current !== next && typeof URL.revokeObjectURL === "function") {
        URL.revokeObjectURL(current);
      }
      return next;
    });
  }, []);

  const handleManualChange = (event: ChangeEvent<HTMLInputElement>) => {
    const next = event.target.value;
    updateObjectURL(null);
    onChange(next);
  };

  const resetInput = () => {
    if (inputRef.current) {
      inputRef.current.value = "";
    }
  };

  const extractErrorMessage = (xhr: XMLHttpRequest) => {
    if (xhr.responseType === "json" && xhr.response) {
      const errorBody = xhr.response.error;
      if (errorBody && typeof errorBody.message === "string") {
        return errorBody.message;
      }
    }
    try {
      const raw = typeof xhr.response === "string" && xhr.response ? xhr.response : typeof xhr.responseText === "string" ? xhr.responseText : "";
      if (raw) {
        const payload = JSON.parse(raw);
        if (payload?.error?.message) {
          return String(payload.error.message);
        }
      }
    } catch {
      /* ignore */
    }
    if (xhr.status === 0) {
      return "Unggahan gagal: koneksi jaringan tidak tersedia.";
    }
    return `Unggahan gagal (${xhr.status}).`;
  };

  const handleUpload = useCallback(
    (file: File) => {
      const objectUrl = URL.createObjectURL(file);
      updateObjectURL(objectUrl);
      setLastFile(file);
      setState({ isUploading: true, progress: 0, error: null });

      const formData = new FormData();
      formData.append("file", file);

      const xhr = new XMLHttpRequest();
      xhr.open("POST", "/api/admin/uploads");
      xhr.withCredentials = true;
      xhr.responseType = "json";

      xhr.upload.onprogress = (event) => {
        if (event.lengthComputable) {
          const percent = Math.round((event.loaded / event.total) * 100);
          setState((prev) => ({ ...prev, progress: percent }));
        } else {
          setState((prev) => ({ ...prev, progress: null }));
        }
      };

      xhr.onerror = () => {
        const message = extractErrorMessage(xhr);
        setState({ isUploading: false, progress: null, error: message });
        pushToast({ type: "error", message });
      };

      xhr.onload = () => {
        if (xhr.status >= 200 && xhr.status < 300) {
          const payload = xhr.response ?? {};
          const data: UploadResponse | undefined = payload?.data ?? payload;
          if (data && typeof data.url === "string") {
            updateObjectURL(null);
            onChange(data.url);
            pushToast({ type: "success", message: "Gambar berhasil diunggah." });
            setState({ isUploading: false, progress: null, error: null });
            resetInput();
            return;
          }
          const message = "Respons unggahan tidak valid.";
          setState({ isUploading: false, progress: null, error: message });
          pushToast({ type: "error", message });
        } else {
          const message = extractErrorMessage(xhr);
          setState({ isUploading: false, progress: null, error: message });
          pushToast({ type: "error", message });
        }
      };

      xhr.send(formData);
    },
    [onChange, pushToast, updateObjectURL],
  );

  const onFileSelect = (event: ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (!file) {
      return;
    }
    handleUpload(file);
    resetInput();
  };

  const triggerFilePicker = () => {
    inputRef.current?.click();
  };

  const retryUpload = () => {
    if (lastFile) {
      handleUpload(lastFile);
    }
  };

  const acceptLabel = useMemo(() => accept.replace(/,/g, ", "), [accept]);
  const preview = objectURL ?? value ?? "";

  return (
    <div className="space-y-2">
      <div className="flex gap-2">
        <Input
          id={id}
          name={name}
          value={preview}
          onChange={handleManualChange}
          onBlur={onBlur}
          placeholder="https://cdn.example.com/image.png"
          disabled={state.isUploading || disabled}
        />
        <Button
          type="button"
          variant="secondary"
          onClick={triggerFilePicker}
          disabled={state.isUploading || disabled}
        >
          <UploadCloud className="mr-2 h-4 w-4" /> Unggah
        </Button>
      </div>
      <input
        ref={inputRef}
        type="file"
        accept={accept}
        className="hidden"
        onChange={onFileSelect}
        data-testid="image-uploader-input"
      />
      <p className="text-xs text-muted-foreground">Tipe diperbolehkan: {acceptLabel}.</p>
      {state.isUploading ? (
        <div className="space-y-1">
          <div className="flex items-center justify-between text-xs text-muted-foreground">
            <span>Mengunggah…</span>
            {typeof state.progress === "number" ? <span>{state.progress}%</span> : null}
          </div>
          <div className="h-2 w-full overflow-hidden rounded-full bg-muted">
            <div
              className="h-full rounded-full bg-primary transition-all"
              style={{ width: `${state.progress ?? 100}%` }}
            />
          </div>
        </div>
      ) : null}
      {state.error ? (
        <div className="flex items-center justify-between rounded-md border border-rose-400/40 bg-rose-500/10 px-3 py-2 text-xs text-rose-300">
          <span>{state.error}</span>
          {lastFile ? (
            <button
              type="button"
              className="inline-flex items-center gap-1 text-rose-200 underline"
              onClick={retryUpload}
            >
              <RotateCcw className="h-3 w-3" /> Coba lagi
            </button>
          ) : null}
        </div>
      ) : null}
      <div className="flex items-center justify-center overflow-hidden rounded-lg border border-dashed border-border bg-muted/50 p-2 text-xs text-muted-foreground supports-[backdrop-filter]:bg-muted/40">
        {preview ? (
          // eslint-disable-next-line @next/next/no-img-element
          <img
            src={preview}
            alt="Pratinjau gambar"
            className="max-h-40 w-full rounded-md object-cover"
          />
        ) : (
          <span>Tidak ada gambar</span>
        )}
      </div>
    </div>
  );
}
