# Deployment Checklist

## Environment Variables

### Backend (Leapcell)
```bash
# Application
APP_ENV=production
LOG_MODE=stdout

# Database
POSTGRES_URL=your-postgres-connection-string

# JWT
JWT_SECRET=your-long-random-secret-key
ACCESS_TOKEN_TTL_MIN=15
REFRESH_TOKEN_TTL_DAY=7

# CORS - Multiple origins separated by comma
FRONTEND_ORIGIN=https://tany-ai.vercel.app,https://your-custom-domain.com

# AI Provider
AI_PROVIDER=gemini
GOOGLE_GENAI_API_KEY=your-gemini-api-key
GEMINI_MODEL=gemini-2.5-flash

# Storage
STORAGE_DRIVER=supabase
SUPABASE_URL=https://your-project.supabase.co
SUPABASE_SERVICE_ROLE=your-service-role-key
SUPABASE_BUCKET=public-assets

# Rate Limiting
KB_CACHE_TTL_SECONDS=60
KB_RATE_LIMIT_PER_5MIN=30
CHAT_RATE_LIMIT_PER_5MIN=30
```

### Frontend (Vercel)
```bash
# API Endpoint
NEXT_PUBLIC_API_URL=https://your-backend.apn.leapcell.dev

# JWT Secret (must match backend)
JWT_SECRET=your-long-random-secret-key
```

## Deployment Steps

### 1. Backend (Leapcell)

1. Push code ke repository
2. Set semua environment variables di Leapcell dashboard
3. Deploy backend
4. Test health endpoint: `https://your-backend.apn.leapcell.dev/healthz`

### 2. Frontend (Vercel)

1. Set environment variables di Vercel dashboard
2. Trigger redeploy atau push ke main branch
3. Verify CSP tidak blocking requests di browser console

## Post-Deployment Verification

### Backend Health Check
```bash
curl https://your-backend.apn.leapcell.dev/healthz
```

Expected: `{"status":"ok"}`

### CORS Check
```bash
curl -X OPTIONS https://your-backend.apn.leapcell.dev/api/v1/chat \
  -H "Origin: https://tany-ai.vercel.app" \
  -H "Access-Control-Request-Method: POST" \
  -H "Access-Control-Request-Headers: content-type,authorization" \
  -v
```

Expected headers in response:
- `Access-Control-Allow-Origin: https://tany-ai.vercel.app`
- `Access-Control-Allow-Methods: GET, POST, PUT, PATCH, DELETE, OPTIONS`
- `Access-Control-Allow-Headers: Content-Type, Authorization, X-Requested-With`

### Frontend Check

1. Open https://tany-ai.vercel.app
2. Open browser console (F12)
3. Check Network tab for any CORS errors
4. Verify images from Supabase load correctly
5. Test chat functionality

## Common Issues

### Issue: CORS Error
**Solution**: Ensure `FRONTEND_ORIGIN` in backend matches your frontend URL exactly (no trailing slash)

### Issue: Images not loading (400 Bad Request)
**Solution**: Verify Supabase domain is in `next.config.ts` `images.remotePatterns`

### Issue: CSP blocking requests
**Solution**: Check `middleware.ts` includes backend domain in `connect-src`

### Issue: 500 Error on frontend
**Solution**: Check backend logs and ensure all environment variables are set correctly
