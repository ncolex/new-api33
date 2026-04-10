# VPS Deployment Troubleshooting (Node.js expectation vs actual backend)

This repository is **not a pure Node.js API backend**. The API server runtime is Go (`main.go`), and the Node projects under `web/` and `electron/` are frontend/desktop layers.

## Deployment failure checklist

### 1) Entry point mismatch (most common)

- VPS startup commands like `npm start`, `node server.js`, or `pm2 start app.js` at repository root fail because there is no root Node API entrypoint.
- The real API process starts with Go: `go build -o new-api main.go && ./new-api`.

### 2) Missing dependencies (runtime/build)

- Backend requires a Go toolchain (Go 1.22+).
- Frontend build uses Bun in `web/` (`bun install`, `bun run build`) to produce `web/dist` assets embedded/served by the backend.
- Using npm-only flow for root backend deploy is incorrect for this repo.

### 3) Environment variables that can hard-fail startup

- `SESSION_SECRET` **must not** be `random_string` (startup exits fatally if set to that value).
- `PORT` should be set by VPS/process manager if non-default is required.
- `SQL_DSN` controls external DB usage; if unset, app falls back to SQLite.
- `SQLITE_PATH` (when used) must be writable by the service account.

### 4) Port binding

- Server binds to `PORT` when present.
- If `PORT` is unset, it falls back to CLI/default port (`3000`).
- Reverse proxy and process manager must target the same port.

### 5) Async/unhandled promises

- For API deployment on VPS, Node async/unhandled promise handling is **not the primary failure vector**, because backend runtime is Go.
- Node promise issues would mostly affect `electron/` or frontend tooling, not root API startup.

### 6) Build/start script mismatch

- `web/package.json` scripts (`dev`, `build`, `preview`) are for Vite frontend only.
- `electron/package.json` scripts are for desktop app packaging/runtime.
- Neither replaces the Go API start command in production VPS deployment.

## Exact fix sequence (production)

```bash
cd /workspace/new-api33 && \
  cd web && bun install && bun run build && cd .. && \
  go build -o new-api main.go && \
  SESSION_SECRET='replace-with-random-secret' PORT=${PORT:-3000} ./new-api
```

## Corrected start command (minimal)

```bash
PORT=${PORT:-3000} ./new-api
```

> Use this minimal command only after `new-api` is already built and frontend assets are prepared.
