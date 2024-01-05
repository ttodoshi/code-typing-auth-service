FROM golang:1.21-alpine AS builder

WORKDIR /usr/local/src/app

COPY go.mod go.sum ./
COPY internal internal
COPY pkg pkg
COPY cmd cmd

RUN go mod tidy

RUN go build -o ./bin/app ./cmd/main/main.go

FROM alpine:latest AS final

WORKDIR /usr/local/src/app

COPY --from=builder /usr/local/src/app/bin/app .

EXPOSE 8090

CMD ["./app"]
