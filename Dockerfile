FROM golang:1.17-bullseye AS builder

ENV \
	GOPROXY=https://mirrors.aliyun.com/goproxy/,direct \
	GO111MODULE=on \
	CGO_ENABLED=0 \
	GIN_MODE=release

WORKDIR /workdir/

COPY go.mod go.sum ./

RUN go mod tidy

COPY . ./

RUN go build -o /passport-v4-server


# need SSL certification
# so not from "scratch"
FROM debian:bullseye

# install SSL ca
RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates curl

COPY --from=builder /passport-v4-server /passport-v4-server
COPY config.yml /
COPY handler/redirect.html /handler/

EXPOSE 8000
ENTRYPOINT ["/passport-v4-server"]
