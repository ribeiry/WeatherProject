FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git

WORKDIR /src

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main . 

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /src .

EXPOSE 8080

CMD ["./main"]