# syntax=docker/dockerfile:1
FROM golang:1.24.2-alpine AS build

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=build /app/server .

# If you rely on a .env file, you can COPY it (or better: pass via env)
# COPY .env .

EXPOSE 3000
ENTRYPOINT ["./server"]