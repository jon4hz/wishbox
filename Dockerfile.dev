FROM golang:1.17 as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/wishbox ./
ENTRYPOINT [ "/app/wishbox" ]
