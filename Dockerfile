# syntax=docker/dockerfile:1

# Build stage — uses the committed vendor tree for reproducible builds.
# Pinned by multi-arch manifest-list digest (constitution X).
FROM golang:1.26-bookworm@sha256:5d2b868674b57c9e48cdd39e891acce4196b6926ca6d11e9c270a8f85106203d AS build
ARG VERSION=dev
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 go build -mod=vendor -trimpath \
    -ldflags "-X github.com/inhuman/mcp-exec/internal/server.Version=${VERSION}" \
    -o /out/mcp-exec ./cmd/mcp-exec

# Runtime stage — Python sandbox with PyYAML + Jinja2, non-root.
# Pinned by multi-arch manifest-list digest (constitution X). Update the digest
# intentionally when bumping the base; refresh via:
#   docker buildx imagetools inspect python:3.13-slim
FROM python:3.13-slim@sha256:b04b5d7233d2ad9c379e22ea8927cd1378cd15c60d4ef876c065b25ea8fb3bf3 AS runtime
# Pinned exact versions so the sandbox toolset can't drift on rebuild (MarkupSafe
# is Jinja2's runtime dep — pinned too so it doesn't float).
RUN pip install --no-cache-dir PyYAML==6.0.3 Jinja2==3.1.6 MarkupSafe==3.0.3
# Scientific/data + imaging stack for analysis & visualization tasks (matrices,
# dataframes, charts, image generation). Heavy (~200 MB) but the sandbox is the
# place this work happens. Separate layer so the small deps above stay cached.
RUN pip install --no-cache-dir \
    Pillow==12.2.0 \
    numpy==2.4.6 \
    pandas==3.0.3 \
    matplotlib==3.10.9
RUN useradd --uid 65532 --create-home --shell /usr/sbin/nologin sandbox
COPY --from=build /out/mcp-exec /usr/local/bin/mcp-exec
USER 65532
ENTRYPOINT ["/usr/local/bin/mcp-exec"]
