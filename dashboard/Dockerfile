# ── Build Stage ──────────────────────────────────────────────────────────────
FROM node:20-alpine AS build

WORKDIR /app

COPY package*.json ./
RUN npm ci

COPY . .

# API calls happen in the browser (client-side), so point to the host-exposed gateway port
ENV VITE_API_URL=http://localhost:8080
RUN npm run build

# ── Runtime Stage ────────────────────────────────────────────────────────────
FROM node:20-alpine

WORKDIR /app

COPY --from=build /app/build ./build
COPY --from=build /app/package*.json ./
RUN npm ci --omit=dev

ENV PORT=3000
ENV ORIGIN=http://localhost:3000
EXPOSE 3000

CMD ["node", "build"]
