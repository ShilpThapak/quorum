# syntax=docker/dockerfile:1
# Quorum grader — single Render web service (Docker).
# Build context = repo ROOT (paths below are root-relative).
# Runtime MUST keep the full Go toolchain + gcc: the grader spawns `go vet` /
# `go test -race` subprocesses, and -race links with cgo (needs gcc).

FROM golang:1.24 AS gobuild
WORKDIR /src
COPY server/go.mod server/go.mod
COPY server/*.go ./
RUN CGO_ENABLED=0 go build -o /out/quorum-server .

FROM golang:1.24
RUN apt-get update && apt-get install -y --no-install-recommends gcc \
    && rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY --from=gobuild /out/quorum-server /app/quorum-server
COPY questions-src/ /app/questions-src/
ENV QUORUM_QUESTIONS=/app/questions-src
ENV PORT=8090
EXPOSE 8090
CMD ["/app/quorum-server"]
