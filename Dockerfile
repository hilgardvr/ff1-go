FROM golang:1.19.5-alpine3.17

RUN apk add build-base

EXPOSE 3000

WORKDIR /usr/src/app

COPY . .

RUN rm Dockerfile

RUN go build

RUN go test ./...

CMD ["./ff1-go"]
