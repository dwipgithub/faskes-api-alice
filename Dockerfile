FROM golang:1.26.5-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/app \
    .

FROM alpine:3.22

WORKDIR /app

RUN apk add --no-cache ca-certificates
COPY --from=builder /out/app /app/app

EXPOSE 8081

ENTRYPOINT ["/app/app"]
