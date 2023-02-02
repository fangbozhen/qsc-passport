FROM golang:1.19

MAINTAINER sync

ENV GO111MODULE=on \
    CGO_ENABLE=0 \
    GINMODE=release

WORKDIR C:\Users\gratian\Desktop\QSCpassport

RUN go mod init passport \
    && go mod tidy \
    && go build -o passport .

COPY config.yml .

EXPOSE 8333

CMD ["./passport"]