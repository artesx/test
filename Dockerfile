FROM golang:1.22.4-alpine3.19 AS builder

RUN go version

RUN apk add --no-cache git

COPY . /webhook-test
WORKDIR /webhook-test

RUN go mod download
RUN GOOS=linux go build -o ./.bin/main ./api/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=0 /webhook-test/.bin/main .
COPY --from=0 /webhook-test/.env .


RUN apk add tzdata ca-certificates curl

EXPOSE 5004

CMD ["./main"]