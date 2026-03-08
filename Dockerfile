FROM --platform=$BUILDPLATFORM golang:1.24-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o excalidraw-backend .

FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/excalidraw-backend /app/excalidraw-backend
COPY --from=builder /app/config /app/config
COPY --from=builder /app/.env.example.dex /app/.env.example.dex

EXPOSE 3002

CMD ["/app/excalidraw-backend"]
