FROM --platform=${BUILDPLATFORM} node:20.12.2-bullseye AS node-builder
WORKDIR /build

COPY frontend/package.json frontend/yarn.lock ./

RUN yarn config set registry 'https://registry.npmmirror.com' && \
    yarn install

COPY frontend .
RUN yarn build-only

FROM golang:1.22.2-alpine3.19 AS builder
WORKDIR /build

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk add --no-cache ca-certificates gcc libtool make musl-dev protoc git && \
    go env -w GOPROXY=https://goproxy.cn,direct

RUN apk --no-cache add bash

COPY Makefile go.mod go.sum ./
RUN make init && go mod download

COPY . .
COPY --from=node-builder /build/dist ./pkg/server/static
RUN make deps && chmod +x hack/release.sh && ./hack/release.sh aurora "-s -w -X 'github.com/MR5356/aurora/pkg/version.Version=v0.0.1'" _output

FROM scratch
COPY --from=builder /build/_output /
