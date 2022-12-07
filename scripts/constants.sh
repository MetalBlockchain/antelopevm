#!/usr/bin/env bash

# Set the PATHS
GOPATH="$(go env GOPATH)"

# Set binary location
binary_path=${ANTELOPE_BINARY_PATH:-"$GOPATH/src/github.com/!metal!blockchain/metalgo/build/plugins/jukr5oTVE2KfmEGwDRCvQnXcwrzwyzRtwyfSGFGPyuyNc12Fs"}

# Avalabs docker hub
dockerhub_repo="metalblockchain/metalgo"

# Current branch
current_branch=${CURRENT_BRANCH:-$(git describe --tags --exact-match 2> /dev/null || git symbolic-ref -q --short HEAD || git rev-parse --short HEAD)}
echo "Using branch: ${current_branch}"

# Image build id
# Use an abbreviated version of the full commit to tag the image.

# WARNING: this will use the most recent commit even if there are un-committed changes present
antelope_commit="$(git --git-dir="$ANTELOPE_PATH/.git" rev-parse HEAD)"
antelope_commit_id="${antelope_commit::8}"

build_image_id=${BUILD_IMAGE_ID:-"$metal_version-$antelope_commit_id"}