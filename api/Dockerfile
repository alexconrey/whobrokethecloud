FROM golang:1.18-alpine

WORKDIR /app

ADD . .

RUN go build 

ENTRYPOINT ["/app/api"]