#!/usr/bin/env bash
# Build and package the Linux/macOS release archives.
# Usage: BUILD_VERSION=1.2.3 ./build.sh
set -euo pipefail

: "${BUILD_VERSION:?Please set BUILD_VERSION, e.g. BUILD_VERSION=1.2.3 ./build.sh}"

cd "$(dirname "$0")"
make clean
make package-linux package-macos VERSION="$BUILD_VERSION"

echo "Artifacts:"
ls -lh output/*.tar.gz
