FROM golang:1.19

WORKDIR /app

# RUN apk update && apk add libc-dev && apk add gcc && apk add make

COPY . .

RUN go mod download && go mod verify

ENV go env -w GO111MODULE=off

RUN go build -o bin/server main.go

CMD ["./bin/server"]
