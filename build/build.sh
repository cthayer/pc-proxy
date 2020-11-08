#!/bin/sh

# path to the directory this script is in (based on how it was called)
SCRIPT_DIR=$(dirname "$0")

# get the version
export VERSION=$(cat $SCRIPT_DIR/versions.json | jq -r '."pc-proxy"')

mage build

if [ "${BUILD_DOCKER}" = "1" ]; then
  # build docker images (we don't build a docker image for darwin because docker doesn't run there)
  docker build --build-arg "VERSION=${VERSION}" --build-arg "ARCH=amd64" -t cdmnky/pc-proxy:${VERSION}-amd64 .
  docker build --build-arg "VERSION=${VERSION}" --build-arg "ARCH=arm64" -t cdmnky/pc-proxy:${VERSION}-arm64 .
  docker build --build-arg "VERSION=${VERSION}" --build-arg "ARCH=arm" -t cdmnky/pc-proxy:${VERSION}-arm .
fi
