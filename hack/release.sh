#!/usr/bin/env bash

OS="linux darwin windows"
ARCHITECTURES="amd64 arm64"

NAME=$1
OUT_DIR=$2
MODULE_NAME="github.com/MR5356/aurora"

# Get the version from git tags
VERSION=$(git describe --tags 2>/dev/null)

# If git describe returned an error (i.e., VERSION is empty)
if [ -z "$VERSION" ]; then
  # Get the current branch name
  BRANCH=$(git rev-parse --abbrev-ref HEAD 2>/dev/null)
  # Get the current commit hash
  COMMIT=$(git rev-parse --short HEAD 2>/dev/null)
  # Construct the version from the branch and commit
  VERSION="${BRANCH}_${COMMIT}"
fi

# Set GO_FLAGS with the version information
GO_FLAGS="-s -w -X '${MODULE_NAME}/pkg/version.Version=${VERSION}'"

echo "Building ${NAME} with GO_FLAGS=${GO_FLAGS}"

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
