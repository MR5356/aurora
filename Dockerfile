FROM golang:1.22.2-alpine3.19 AS builder
WORKDIR /build

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk add --no-cache ca-certificates gcc libtool make musl-dev protoc git && \
    go env -w GOPROXY=https://goproxy.cn,direct

COPY Makefile go.mod go.sum ./
RUN make init && go mod download

COPY . .
RUN make build

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /build/_output/bin/aurora /app/

EXPOSE 80
ENTRYPOINT ["/app/aurora"]