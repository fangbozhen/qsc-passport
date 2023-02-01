FROM golang:1.19

MAINTAINER sync

ENV GO111MODULE=on \
    CGO_ENABLE=0 \
    GOPROXY="https://goproxy.io,direct" \
    GIN_MODE=release

WORKDIR /go/src/qscpassport

COPY . /go/src/qscpassport

RUN go build -o /go/bin/passport .

COPY config.yml .

EXPOSE 3000

CMD ["/go/bin/passport"]