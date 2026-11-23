# syntax=docker/dockerfile:1.7

FROM golang:1.26.4-bookworm AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    go install github.com/swaggo/swag/cmd/swag@v1.16.6 && \
    swag init \
    -g app/cmd/server/main.go --output app/docs --parseDependency --parseInternal

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/server ./app/cmd/server

FROM gcr.io/distroless/static-debian12
COPY --from=builder --chown=nonroot:nonroot /out/server /server
USER nonroot:nonroot
ENTRYPOINT ["/server"]
