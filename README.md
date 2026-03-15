A lightweight Go backend for AI-powered text translation using LangChainGo, Ollama, gRPC, and Redis.

Project layout (container-friendly)

```
ai-lingua-go/
├── compose.yaml       # Runs server + Ollama + Redis (each in its own container)
├── server/            # Go gRPC app (Dockerfile, cmd, internal, api, config)
├── ollama/            # Ollama LLM runtime (Dockerfile → ollama/ollama)
└── redis/             # Redis cache (Dockerfile → redis:7-alpine)
```

Run the full stack: `docker compose up -d`, then `docker compose exec ollama ollama pull qwen2.5:7b`. Call the API on port 50051 (see compose.yaml comments).

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
