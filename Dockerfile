FROM golang:1.20.9-alpine3.17 AS builder

COPY . /MovieBot/
WORKDIR /MovieBot/

RUN go mod download
RUN go build -o ./bin/bot cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=0 /MovieBot/bin/bot .
COPY --from=0 /MovieBot/configs configs/

EXPOSE 80

CMD ["./bot"]