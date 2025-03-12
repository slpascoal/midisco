FROM golang:1.24-alpine

WORKDIR /usr/src/app

RUN go mod init midisco-api

COPY . .

RUN go mod tidy

RUN go build -o main ./cmd/server

EXPOSE 8080

CMD ["./main"]