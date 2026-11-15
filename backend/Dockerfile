# Build the static binaries first, then ship only what the process needs.
FROM golang:1.23-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/service ./cmd/service \
 && CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/migrate ./cmd/migrate

# Alpine keeps a shell for the pre-deploy command and bundles CA roots for Neon.
FROM alpine:3.20

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=build /out/service /app/service
COPY --from=build /out/migrate /app/migrate
COPY --from=build /src/migrations /app/migrations

ENV MIGRATIONS_DIR=/app/migrations

EXPOSE 8080

USER nobody
ENTRYPOINT ["/app/service"]
