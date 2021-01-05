FROM golang:1.15-alpine as dev

RUN apk add --no-cache git
WORKDIR /work

RUN go get github.com/uudashr/gopkgs/cmd/gopkgs \
    github.com/ramya-rao-a/go-outline \
    github.com/go-delve/delve/cmd/dlv \
    golang.org/x/lint/golint
RUN GO111MODULE=on go get golang.org/x/tools/gopls@master golang.org/x/tools@master

FROM golang:1.15-alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

RUN mkdir /src/
WORKDIR /src/

COPY ./go.mod .
COPY ./go.sum .
RUN go mod download

COPY . .
WORKDIR /src/cmd
RUN go build -o app

FROM alpine as runtime
COPY --from=builder /src/cmd/app /
CMD ./app