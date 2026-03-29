FROM golang:1.25

RUN go install github.com/air-verse/air@latest

ENV CGO_ENABLED=0 GO111MODULE=on GOOS=linux

WORKDIR /go/src/app

COPY . .

RUN go mod download

CMD air -c .air.users.toml