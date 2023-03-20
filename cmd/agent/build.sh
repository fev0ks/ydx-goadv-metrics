#!/usr/bin/env bash

# STEP 1: Determinate the required values

PACKAGE="github.com/fev0ks/ydx-goadv-metrics/cmd/agent"
VERSION="$(git describe --tags --always --abbrev=0 --match='v[0-9]*.[0-9]*.[0-9]*' 2> /dev/null | sed 's/^.//')"
COMMIT_HASH="$(git rev-parse --short HEAD)"
BUILD_TIMESTAMP=$(date '+%Y-%m-%dT%H:%M:%S')

# STEP 2: Build the ldflags

LDFLAGS=(
  "-X '${PACKAGE}/version.BuildVersion=${VERSION}'"
  "-X '${PACKAGE}/version.BuildDate=${BUILD_TIMESTAMP}'"
  "-X '${PACKAGE}/version.BuildCommit=${COMMIT_HASH}'"
)

# STEP 3: Actual Go build process

go build -ldflags="${LDFLAGS[*]}" github.com/fev0ks/ydx-goadv-metrics/cmd/agent/