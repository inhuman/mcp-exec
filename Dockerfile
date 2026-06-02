# syntax=docker/dockerfile:1

# Build stage — uses the committed vendor tree for reproducible builds.
FROM golang:1.26-bookworm AS build
ARG VERSION=dev
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 go build -mod=vendor -trimpath \
    -ldflags "-X github.com/inhuman/mcp-exec/internal/server.Version=${VERSION}" \
    -o /out/mcp-exec ./cmd/mcp-exec

# Runtime stage — Python sandbox with PyYAML + Jinja2, non-root.
# NOTE (supply chain): pin python by digest in release builds, e.g.
#   FROM python:3.13-slim@sha256:<digest>
FROM python:3.13-slim AS runtime
RUN pip install --no-cache-dir PyYAML Jinja2
RUN useradd --uid 65532 --create-home --shell /usr/sbin/nologin sandbox
COPY --from=build /out/mcp-exec /usr/local/bin/mcp-exec
USER 65532
ENTRYPOINT ["/usr/local/bin/mcp-exec"]
