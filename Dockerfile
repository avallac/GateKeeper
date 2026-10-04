# syntax=docker/dockerfile:1

ARG BUILD_FROM=debian:bookworm-slim

FROM golang:1.25-bookworm AS build

RUN apt-get update && apt-get install -y --no-install-recommends \
    pkg-config \
    libzmq3-dev \
    libczmq-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 go build -o /out/gatekeeper ./cmd/app

FROM ${BUILD_FROM}

RUN apt-get update && apt-get install -y --no-install-recommends \
    libzmq5 \
    libczmq4 \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=build /out/gatekeeper ./gatekeeper
COPY web ./web
COPY configs/app.yaml ./configs/app.yaml

EXPOSE 8080

ENTRYPOINT ["./gatekeeper"]
