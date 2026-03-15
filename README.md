A lightweight Go backend for AI-powered text translation using LangChainGo, Ollama, gRPC, and Redis.

Project layout

```
ai-lingua-go/
├── compose.yaml       # Full stack: server + Ollama + Redis + frontend
├── compose.dev.yaml   # Dev overrides (frontend hot-reload)
├── server/            # Go backend (gRPC, REST, Dockerfile, cmd, internal, api)
├── frontend/          # React + Vite app (src, public)
├── ollama/            # Ollama LLM (Dockerfile)
└── redis/             # Redis cache (Dockerfile)
```

Run the full stack:

```bash
docker compose up -d
docker compose exec ollama ollama pull qwen2.5:7b
```

Call the API on port 50051 (see compose.yaml comments). The default 7B model needs ~4–5 GB RAM; for less memory use a smaller model (e.g. `OLLAMA_MODEL=tinyllama` and pull `tinyllama`).

Frontend (React) is served at http://localhost:3002. To have Docker detect frontend changes and hot-reload, run with the dev profile:

```bash
docker compose -f compose.yaml -f compose.dev.yaml --profile dev up -d
```

The frontend then runs the Vite dev server with your `frontend/` folder mounted; edit files locally and the app updates in the browser. Without the dev profile, the frontend is built once and served by nginx (no auto-reload).

Making translation faster

- Cache: With Redis (default in Docker), repeated translations for the same text and language pair are served from cache and return almost instantly.
- Smaller model: For quicker first-time responses, use a smaller Ollama model: `OLLAMA_MODEL=tinyllama docker compose up -d`, then `docker compose exec ollama ollama pull tinyllama`. Quality is lower than qwen2.5:7b but latency is much lower.
- Backend: The server limits LLM output to 1024 tokens and uses a low temperature (0.2) so the model stops sooner and runs a bit faster.

Tech stack

- Go, gRPC, REST, LangChainGo, Ollama, Redis, Protocol Buffers
- Frontend: React, TypeScript, Vite

Prerequisites

1. Go 1.22 or later
2. Ollama installed, with a model pulled (e.g. ollama pull llama2)
3. Redis (optional for dev; use in-memory cache instead)

Example to pull a model in Ollama:

    ollama pull llama2

## License

This project is open source under the [MIT License](LICENSE).
