# --- Build stage ---
FROM golang:1.25.3-bookworm AS build

# 1. Install bun
RUN apt-get update && apt-get install -y curl unzip
RUN curl -fsSL https://bun.sh/install | bash
ENV PATH="/root/.bun/bin:${PATH}"

# 2. Set working directory
WORKDIR /app

# 3. Install dependencies
COPY go.mod go.sum ./
COPY package.json bun.lock ./
RUN go mod download
RUN bun install

# 4. Copy source code
COPY . .

# 5. Build the project
RUN make build/prod/static

# --- Run stage ---
FROM gcr.io/distroless/static AS run

# 1. Set working directory
WORKDIR /app

# 2. Copy the built binary from the build stage
COPY --from=build /app/bin/web /app/bin/web

# 3. Expose port - just for documentation purposes
EXPOSE 8080

# 4. Set the entrypoint
ENTRYPOINT [ "/app/bin/web" ]
