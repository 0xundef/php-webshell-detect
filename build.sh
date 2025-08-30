#!/bin/bash
set -e
if [ -z "${BUILD_VERSION}" ];then
  echo 'Please set BUILD_VERSION.'
  exit 1
fi

rm -rf output
mkdir output

#cli for local test
GOARCH=amd64 go build -ldflags="-s -w" -o cli ./cli/cli.go

# Build for Linux
GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o webshell-detect
tar -zcvf ./output/webshell-detect-linux-default-x86_64-${BUILD_VERSION}.tar.gz webshell-detect ./config
rm -rf webshell-detect

# Build for macOS
GOARCH=amd64 GOOS=darwin go build -ldflags="-s -w" -o webshell-detect-macos
tar -zcvf ./output/webshell-detect-macos-default-x86_64-${BUILD_VERSION}.tar.gz webshell-detect-macos ./config
rm -rf webshell-detect-macos
