# AI Translate — Frontend

React + TypeScript + Vite app with a chat-style UI for the translation service.

## Run with Docker Compose (recommended)

From the repo root:

```bash
docker compose up -d
```

Then open **http://localhost:3000**. The frontend container serves the built app; the browser calls the REST API at `http://localhost:8080` (default). To override the API URL at build time: `VITE_API_URL=http://your-host:8080 docker compose build frontend`.

## Run locally (dev)

1. Start the backend (REST API on port 8080). From repo root:
   - **Docker:** `docker compose up -d` (server exposes HTTP on 8080)
   - **Local:** from `server/`: `HTTP_PORT=8080 OLLAMA_MODEL=qwen2.5:7b go run ./cmd/server`

2. From this directory:
   ```bash
   npm install
   npm run dev
   ```
   Open http://localhost:5173. The app calls `http://localhost:8080/api/translate` by default.

To use a different API URL: `VITE_API_URL=http://your-host:8080 npm run dev`

## Build

```bash
npm run build
```
Output is in `dist/`. You can serve it with any static host or point the Go server at it later.
