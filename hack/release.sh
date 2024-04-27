#!/usr/bin/env bash

OS="linux darwin windows"
ARCHITECTURES="amd64 arm64"

NAME=$1
GO_FLAGS=$2
OUT_DIR=$3

mkdir -p ${OUT_DIR}
for arch in ${ARCHITECTURES}; do
  for os in ${OS}; do
    echo "Building ${os}-${arch}"
    if [ "${os}" == "windows" ]; then
      GOOS=${os} GOARCH=${arch} go build -ldflags "${GO_FLAGS}" -o ${OUT_DIR}/${NAME}-${os}-${arch}.exe ./cmd/aurora
      cd ${OUT_DIR} || exit
      tar zcvf ${NAME}-${os}-${arch}.tar.gz ${NAME}-${os}-${arch}.exe
      rm ${NAME}-${os}-${arch}.exe
      cd - || exit
    else
      GOOS=${os} GOARCH=${arch} go build -ldflags "${GO_FLAGS}" -o ${OUT_DIR}/${NAME}-${os}-${arch} ./cmd/aurora
      cd ${OUT_DIR} || exit
      tar zcvf ${NAME}-${os}-${arch}.tar.gz ${NAME}-${os}-${arch}
      rm ${NAME}-${os}-${arch}
      cd - || exit
    fi
  done
done
