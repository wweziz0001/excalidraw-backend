# Excalidraw Backend

A backend-only fork of `excalidraw-full`, focused on identity, canvas management, document access, flexible storage, and AI proxying.

This repository removes the embedded frontend and keeps only the backend services needed to integrate with an existing Excalidraw frontend stack.

## Purpose

This project is designed for deployments where the frontend, scene storage backend, and collaboration room may already exist as separate services, and you want to add a dedicated backend that provides:

- **Authentication via Dex (OIDC)**
- **Multi-canvas management**
- **Backend storage** with support for **SQLite**, **Filesystem**, and **S3**
- **Sharing / document access**
- **OpenAI proxy**

## Scope

This repository is intended to provide only the following backend capabilities:

1. **Authentication**
   - Login through **Dex** using OpenID Connect
   - Suitable for environments where Dex is connected to **LDAP** or another upstream identity provider
   - Session and token handling for authenticated users

2. **Multi-Canvas Management**
   - Create, list, retrieve, update, and delete canvases
   - Associate canvases with authenticated users
   - Support user-owned workspaces and persistent saved boards

3. **Backend Storage**
   - `sqlite`
   - `filesystem`
   - `s3`

4. **Sharing / Document Access**
   - Create shareable document or canvas links
   - Retrieve shared canvases/documents through controlled backend endpoints

5. **OpenAI Proxy**
   - Optional backend proxy for OpenAI-compatible APIs
   - Keeps API keys on the server side instead of exposing them in the browser

## What was intentionally removed

Compared with the original upstream project structure, this repository is now **backend-only**.

The following categories were intentionally removed or excluded from the deployment path:

- Embedded frontend serving
- Patched frontend build pipeline
- Excalidraw frontend submodule dependency for all-in-one deployment
- Cloudflare Worker / BYOC frontend storage helpers
- Built-in realtime collaboration room service
- All-in-one Docker packaging for frontend + backend together

## Repository layout

```text
.
├── config/
├── core/
├── handlers/
├── middleware/
├── stores/
├── .env.example
├── .env.example.dex
├── go.mod
├── go.sum
└── main.go
```

## Authentication model

Authentication is handled through **Dex** using **OIDC**.

Typical flow:

```text
User -> Excalidraw Frontend -> Excalidraw Backend -> Dex -> LDAP
```

Recommended use cases:

- Internal enterprise deployments
- LDAP-backed identity environments
- Centralized SSO through Dex

## Storage backends

The backend is designed to support multiple storage engines.

### 1. SQLite
Recommended as the default starting point.

Use when:
- You want a simple single-node deployment
- You want easy backup and restore
- You want minimal operational overhead

### 2. Filesystem
Useful when documents should be stored directly on mounted volumes.

Use when:
- You want direct file persistence on disk
- You prefer simple host-mounted storage

### 3. S3
Recommended for scalable object storage deployments.

Use when:
- You want external object storage
- You need scalability beyond local disk
- You are using MinIO, AWS S3, or another compatible endpoint

## Configuration

Create a `.env` file from `.env.example.dex` when using Dex.

### Core settings

```env
PORT=3002
APP_BASE_URL=https://backend.example.com
JWT_SECRET=change-me
```

### Dex / OIDC settings

```env
DEX_ISSUER=https://dex.example.com/dex
DEX_CLIENT_ID=excalidraw-backend
DEX_CLIENT_SECRET=change-me
DEX_REDIRECT_URI=https://backend.example.com/auth/callback
DEX_SCOPES=openid profile email groups
```

### Storage settings

#### SQLite
```env
STORAGE_TYPE=sqlite
DATA_SOURCE_NAME=/data/excalidraw.db
```

#### Filesystem
```env
STORAGE_TYPE=filesystem
LOCAL_STORAGE_PATH=/data/storage
```

#### S3
```env
STORAGE_TYPE=s3
S3_BUCKET_NAME=excalidraw
AWS_REGION=us-east-1
AWS_ACCESS_KEY_ID=change-me
AWS_SECRET_ACCESS_KEY=change-me
S3_ENDPOINT=
S3_FORCE_PATH_STYLE=false
```

### OpenAI proxy

```env
OPENAI_API_KEY=
OPENAI_BASE_URL=https://api.openai.com/v1
```

## Expected API areas

This backend is intended to expose endpoints in these areas:

- `/auth/*`
- `/api/canvases/*`
- `/api/share/*`
- `/api/documents/*`
- `/api/ai/*`
- `/healthz`
- `/readyz`

## Suggested integration with an existing stack

This backend is designed to be added alongside an existing Excalidraw deployment, for example:

- `app` for the frontend
- `storage` for existing scene/image storage
- `room` for realtime collaboration
- `mongodb` for the current storage backend
- `backend_full` for authentication, canvas management, sharing, and AI

This lets you introduce platform features incrementally without replacing your current stack immediately.

## Development goals

The current target architecture is:

```text
backend_full =
  Dex authentication
+ multi-canvas management
+ sqlite/filesystem/s3 storage
+ sharing/document access
+ OpenAI proxy
- embedded frontend
- collaboration room
- Cloudflare worker
- GitHub OAuth dependency
```

## Notes

- Use **Dex** when you want LDAP-backed authentication.
- Start with **SQLite** unless you already need filesystem or S3.
- Keep OpenAI proxy optional so the backend can run without AI features.
- This repository is intended to be packaged as a **backend-only Docker image**.

## Status

This repository is currently being refactored from the original `excalidraw-full` layout into a dedicated backend service.

The recommended order of work is:

1. Finalize backend-only cleanup
2. Verify routes and handlers
3. Verify storage implementations
4. Add Dockerfile
5. Build and publish image
6. Add the service into your existing `compose.yaml`
