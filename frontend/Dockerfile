FROM node:22-alpine AS dev
WORKDIR /app
COPY package.json package-lock.json* ./
RUN npm ci
ARG VITE_API_URL=http://localhost:8081
ENV VITE_API_URL=${VITE_API_URL}
EXPOSE 5173
CMD ["npm", "run", "dev", "--", "--host", "0.0.0.0"]

FROM node:22-alpine AS builder

WORKDIR /app

COPY package.json package-lock.json* ./
RUN npm ci

COPY . .
ARG VITE_API_URL=http://localhost:8081
ENV VITE_API_URL=${VITE_API_URL}
RUN npm run build

FROM nginx:alpine

COPY --from=builder /app/dist /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
