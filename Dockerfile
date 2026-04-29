# syntax=docker/dockerfile:1.7

FROM golang:1.26-bookworm AS builder

WORKDIR /src

ENV CGO_ENABLED=0 \
    GOOS=linux

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .

ARG VERSION=dev
ARG COMMIT=unknown
ARG BUILD_DATE=unknown

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -trimpath \
        -ldflags "-s -w \
            -X main.version=${VERSION} \
            -X main.commit=${COMMIT} \
            -X main.date=${BUILD_DATE}" \
        -o /out/ingest ./cmd/ingest


FROM debian:bookworm-slim AS runtime

RUN apt-get update \
    && apt-get install -y --no-install-recommends \
        ca-certificates \
        libffi8 \
        tini \
    && rm -rf /var/lib/apt/lists/*

RUN groupadd --system --gid 1000 ingest \
    && useradd --system --uid 1000 --gid ingest --home-dir /app --shell /usr/sbin/nologin ingest \
    && mkdir -p /app /data/catalog /data/staging \
    && chown -R ingest:ingest /app /data

WORKDIR /app
COPY --from=builder --chown=ingest:ingest /out/ingest /app/ingest

USER ingest

EXPOSE 8080
VOLUME ["/data"]

ENTRYPOINT ["/usr/bin/tini", "--", "/app/ingest"]