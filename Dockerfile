FROM golang:1.17-bullseye AS builder

ENV GOPROXY=https://mirrors.aliyun.com/goproxy/,direct \
	GO111MODULE=on \
	CGO_ENABLED=0

WORKDIR /workdir/

COPY go.mod go.sum ./

RUN go mod download all

COPY . ./

RUN go build -o /passport-v4-server

FROM scratch
#alpine:latest

COPY --from=builder /passport-v4-server /passport-v4-server

ENTRYPOINT ["/passport-v4-server"]
