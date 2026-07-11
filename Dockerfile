FROM golang:alpine AS builder

WORKDIR /app
RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -a -o drop-server .

FROM alpine:latest
WORKDIR /app
RUN apk add --no-cache sqlite-libs tzdata

COPY --from=builder /app/drop-server .
COPY templates ./templates
COPY static ./static

ENV PORT=80
ENV DATA_DIR=/app/data

EXPOSE 80

CMD ["./drop-server"]
