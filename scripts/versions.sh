#!/usr/bin/env bash

# Set up the versions to be used
antelope_version=${CORETH_VERSION:-'v0.0.1'}
# Don't export them as they're used in the context of other calls
metal_version=${METAL_VERSION:-'v1.9.16'}