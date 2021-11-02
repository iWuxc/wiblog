#!/usr/bin/env sh

set -e

_registry="$1"
_tag="$2"
#_platform="linux/amd64,linux/arm64,linux/386"

if [ -z "$_registry" ] || [ -z "$_tag" ]; then
  echo "Please specify image repository and tag."
  exit 0
fi

# create build dir
mkdir -p ../bin

# create builder
# docker buildx create --use --name builder

for file in ../pkg/core/*; do
  app="$(basename $file)"
  CGO_ENABLED=0 go build -tags prod -o ../bin/wiblog "../cmd/$app"
  #docker image
  # docker build --platform "$_platform" \
  docker build \
  -f "../build/package/$app.Dockerfile" \
  -t "$_registry/$app:latest" \
  -t "$_registry/$app:$_tag" .
#  --push .
done

# clean dir ./bin
rm -rf ../bin
