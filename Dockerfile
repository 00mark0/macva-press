# =========================
# 1. Base Stage
# =========================
FROM golang:1.24-alpine AS base
# Enable CGO and install required dependencies
ENV CGO_ENABLED=1
RUN apk add --no-cache \
    gcc \
    g++ \
    make \
    libwebp-dev \
    ffmpeg

# Install dependencies required for both dev and build stages
RUN apk add --no-cache nodejs npm make
WORKDIR /app
# Copy module files first to leverage caching
COPY package.json package-lock.json ./
RUN npm install
COPY go.mod go.sum ./
RUN go mod tidy
# Copy the rest of the source code
COPY . .

# =========================
# 2. Development Stage
# =========================
FROM base AS dev
# Install dev tools for live reloading, templating, etc.
RUN go install github.com/a-h/templ/cmd/templ@latest
RUN go install github.com/air-verse/air@latest
RUN apk add --no-cache libwebp-dev
RUN npm install
EXPOSE 3000
# Mounting the project in the dev compose file will allow live editing.
CMD ["make", "live"]

# =========================
# 3. Builder Stage for Production
# =========================
FROM base AS builder
# Run Tailwind to generate production-ready CSS before building the Go binary
RUN apk add --no-cache libwebp-dev
RUN make prod/tailwind
# 🔧 Install Templ CLI and regenerate templates
RUN go install github.com/a-h/templ/cmd/templ@latest
RUN templ generate
# Build the Go binary (add any production build flags if desired)
RUN go build -o app
RUN npm run build
RUN apk add --no-cache curl
# Install golang-migrate by downloading its prebuilt binary (adjust version if needed)
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.2/migrate.linux-amd64.tar.gz | tar xvz

# =========================
# 4. Minimal Production Stage
# =========================
FROM alpine:latest AS prod
WORKDIR /app
# (Optional) Install ca-certificates if your app needs to make HTTPS requests
RUN apk --no-cache add ca-certificates tzdata
# Copy the compiled Go binary from the builder stage
COPY --from=builder /app/app .
# Copy the .env file from the builder stage
COPY --from=builder /app/.env .env
# Copy the static files (adjust the source path if needed)
COPY --from=builder /app/static /app/static
# Copy the migrations binary from the builder stage
COPY --from=builder /app/migrate ./migrate
COPY db/migration ./migration
COPY start.sh /app/start.sh
RUN chmod +x /app/start.sh
EXPOSE 3000
# Run the binary directly; this final image is minimal as it doesn't include Go, Node.js, npm, or build tools.
CMD ["./app"]
ENTRYPOINT [ "/app/start.sh" ]
