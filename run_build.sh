#!/usr/bin/env sh

set -e

_platform="linux/amd64,linux/arm64,linux/386"

# create builder
#docker buildx create --use --name buildx

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags prod -o wiblog "./cmd/wiblog"

# docker image
docker buildx build --platform "$_platform" \
  -t "iwuxc/wiblog:latest" \
  -t "iwuxc/wiblog:1.0.0" \
  --push .

# clean dir ./bin
rm -rf wiblog
