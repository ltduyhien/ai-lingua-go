A lightweight Go backend for AI-powered text translation using LangChainGo, Ollama, gRPC, and Redis.

Project layout (container-friendly)

```
ai-lingua-go/
├── compose.yaml       # Runs server + Ollama + Redis (each in its own container)
├── server/            # Go gRPC app (Dockerfile, cmd, internal, api, config)
├── ollama/            # Ollama LLM runtime (Dockerfile → ollama/ollama)
└── redis/             # Redis cache (Dockerfile → redis:7-alpine)
```

Run the full stack: `docker compose up -d`, then `docker compose exec ollama ollama pull qwen2.5:7b`. Call the API on port 50051 (see compose.yaml comments). The default 7B model needs ~4–5 GB RAM; for less memory use a smaller model (e.g. `OLLAMA_MODEL=tinyllama` and pull `tinyllama`).

Frontend (React) is served at http://localhost:3002. To have Docker detect frontend changes and hot-reload, run with the dev profile:

```bash
docker compose -f compose.yaml -f compose.dev.yaml --profile dev up -d
```

The frontend then runs the Vite dev server with your `frontend/` folder mounted; edit files locally and the app updates in the browser. Without the dev profile, the frontend is built once and served by nginx (no auto-reload).

Making translation faster

- Cache: With Redis (default in Docker), repeated translations for the same text and language pair are served from cache and return almost instantly.
- Smaller model: For quicker first-time responses, use a smaller Ollama model: `OLLAMA_MODEL=tinyllama docker compose up -d`, then `docker compose exec ollama ollama pull tinyllama`. Quality is lower than qwen2.5:7b but latency is much lower.
- Backend: The server limits LLM output to 1024 tokens and uses a low temperature (0.2) so the model stops sooner and runs a bit faster.

Technology Choices

1. Go. Backend language. Fast, concurrent, single binary. Good fit for servers.

2. LangChainGo. LLM orchestration and prompt management. Keeps prompt logic clean and makes it easy to swap LLM providers later.

3. Ollama. Local LLM inference. Free, runs on your machine, no API keys. Pull a model like llama2 and use it locally.

4. gRPC. API layer. Request-response, typed contracts via Protocol Buffers, HTTP/2. Best for backend and internal clients.

5. Redis. Translation cache. Fast lookups, TTL support, shareable across instances. Cuts down repeated LLM calls.

6. Protocol Buffers. Message format for gRPC. Schema-first, efficient, tooling support.

Alternatives Considered

gRPC vs REST: REST is simpler and browser-friendly. We chose gRPC for typed API and better performance for backend clients.

gRPC vs WebSocket: WebSocket is for real-time, bidirectional streams. Translation is request-response. gRPC fits better.

LangChainGo vs direct HTTP: Direct HTTP to Ollama is simpler. LangChainGo gives structure and easier extension if we add chains or more providers.

Redis vs in-memory cache: In-memory is fine for single-instance dev. Redis is needed for shared, multi-instance production.

Prerequisites

1. Go 1.22 or later
2. Ollama installed, with a model pulled (e.g. ollama pull llama2)
3. Redis (optional for dev; use in-memory cache instead)

Example to pull a model in Ollama:

    ollama pull llama2

## License

This project is open source under the [MIT License](LICENSE).
