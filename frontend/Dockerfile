# Dev stage: run Vite dev server (use with volume mount for live reload).
FROM node:22-alpine AS dev
WORKDIR /app
COPY package.json package-lock.json* ./
RUN npm ci
ARG VITE_API_URL=http://localhost:8081
ENV VITE_API_URL=${VITE_API_URL}
EXPOSE 5173
CMD ["npm", "run", "dev", "--", "--host", "0.0.0.0"]

# Stage 1: build the Vite app (runs in the browser; API is called from the host).
FROM node:22-alpine AS builder

WORKDIR /app

COPY package.json package-lock.json* ./
RUN npm ci

COPY . .
# Browser runs on the host, so the API URL must be reachable from the host (compose exposes REST on 8081).
ARG VITE_API_URL=http://localhost:8081
ENV VITE_API_URL=${VITE_API_URL}
RUN npm run build

# Stage 2: serve static files with nginx.
FROM nginx:alpine

COPY --from=builder /app/dist /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
