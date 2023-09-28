#!/usr/bin/env bash

# Set up the versions to be used
antelope_version=${CORETH_VERSION:-'v0.0.1'}
# Don't export them as they're used in the context of other calls
metal_version=${METAL_VERSION:-'v1.10.0'}
metal_network_runner_version=${METAL_NETWORK_RUNNER_VERSION:-'v1.6.0'}