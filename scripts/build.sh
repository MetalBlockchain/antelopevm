#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

# Load the constants
# Set the PATHS
GOPATH="$(go env GOPATH)"

# TimestampVM root directory
TIMESTAMPVM_PATH=$( cd "$( dirname "${BASH_SOURCE[0]}" )"; cd .. && pwd )

# Set default binary directory location
binary_directory="/tmp/plugins"
name="jukr5oTVE2KfmEGwDRCvQnXcwrzwyzRtwyfSGFGPyuyNc12Fs"

if [[ $# -eq 1 ]]; then
    binary_directory=$1
elif [[ $# -eq 2 ]]; then
    binary_directory=$1
    name=$2
elif [[ $# -ne 0 ]]; then
    echo "Invalid arguments to build antelopevm. Requires either no arguments (default) or one arguments to specify binary location."
    exit 1
fi


# Build timestampvm, which is run as a subprocess
echo "Building AntelopeVM in $binary_directory/$name"
go build -o "$binary_directory/$name" "main/"*.go