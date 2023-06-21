#!/usr/bin/env bash
# (c) 2019-2022, Ava Labs, Inc. All rights reserved.
# See the file LICENSE for licensing terms.

set -e

# Set the CGO flags to use the portable version of BLST
#
# We use "export" here instead of just setting a bash variable because we need
# to pass this flag to all child processes spawned by the shell.
export CGO_CFLAGS="-O -D__BLST_PORTABLE__"

# e.g.,
# ./scripts/run.sh 1.7.13
#
# run without e2e tests
# ./scripts/run.sh 1.7.13
#
# to run E2E tests (terminates cluster afterwards)
# E2E=true ./scripts/run.sh 1.7.13
if ! [[ "$0" =~ scripts/run.sh ]]; then
  echo "must be run from repository root"
  exit 255
fi

# TimestampVM root directory
ANTELOPEVM_PATH=$( cd "$( dirname "${BASH_SOURCE[0]}" )"; cd .. && pwd )

# Load the versions
source "$ANTELOPEVM_PATH"/scripts/versions.sh

MODE=${MODE:-run}
E2E=${E2E:-false}
if [[ ${E2E} == true ]]; then
  MODE="test"
fi

METAL_LOG_LEVEL=${METAL_LOG_LEVEL:-INFO}

echo "Running with:"
echo metal_version: ${metal_version}
echo MODE: ${MODE}

############################
# download metalgo
# https://github.com/MetalBlockchain/metalgo/releases
GOARCH=$(go env GOARCH)
GOOS=$(go env GOOS)
METALGO_PATH=/tmp/metalgo-${metal_version}/metalgo
METALGO_PLUGIN_DIR=/tmp/metalgo-${metal_version}/plugins

if [ ! -f "$METALGO_PATH" ]; then
  DOWNLOAD_URL=https://github.com/MetalBlockchain/metalgo/releases/download/${metal_version}/metalgo-linux-${GOARCH}-${metal_version}.tar.gz
  DOWNLOAD_PATH=/tmp/metalgo.tar.gz
  if [[ ${GOOS} == "darwin" ]]; then
    DOWNLOAD_URL=https://github.com/MetalBlockchain/metalgo/releases/download/${metal_version}/metalgo-macos-${metal_version}.zip
    DOWNLOAD_PATH=/tmp/metalgo.zip
  fi

  rm -rf /tmp/metalgo-${metal_version}
  rm -rf /tmp/metalgo-build
  rm -f ${DOWNLOAD_PATH}

  echo "downloading metalgo ${metal_version} at ${DOWNLOAD_URL}"
  curl -L ${DOWNLOAD_URL} -o ${DOWNLOAD_PATH}

  echo "extracting downloaded metalgo"
  if [[ ${GOOS} == "linux" ]]; then
    tar xzvf ${DOWNLOAD_PATH} -C /tmp
  elif [[ ${GOOS} == "darwin" ]]; then
    unzip ${DOWNLOAD_PATH} -d /tmp/metalgo-build
    mv /tmp/metalgo-build/build /tmp/metalgo-${metal_version}
  fi
  find /tmp/metalgo-${metal_version}
fi

############################

############################
echo "building antelopevm"

# delete previous (if exists)
rm -f /tmp/metalgo-${metal_version}/plugins/jukr5oTVE2KfmEGwDRCvQnXcwrzwyzRtwyfSGFGPyuyNc12Fs

go build \
-o /tmp/metalgo-${metal_version}/plugins/jukr5oTVE2KfmEGwDRCvQnXcwrzwyzRtwyfSGFGPyuyNc12Fs \
./main/
find /tmp/metalgo-${metal_version}

############################

############################

echo "creating genesis file"
cp ./chain/genesis_test.json /tmp/.genesis

############################

############################

echo "creating vm config"
echo -n "{}" > /tmp/.config

############################

############################
echo "building e2e.test"
# to install the ginkgo binary (required for test build and run)
go install -v github.com/onsi/ginkgo/v2/ginkgo@v2.1.4
ACK_GINKGO_RC=true ginkgo build ./tests/e2e

#################################
# download metal-network-runner
# https://github.com/MetalBlockchain/metal-network-runner
ANR_REPO_PATH=github.com/MetalBlockchain/metal-network-runner
ANR_VERSION=v1.3.3-rc.1
# version set
go install -v ${ANR_REPO_PATH}@${ANR_VERSION}

#################################
# run "metal-network-runner" server
GOPATH=$(go env GOPATH)
if [[ -z ${GOBIN+x} ]]; then
  # no gobin set
  BIN=${GOPATH}/bin/metal-network-runner
else
  # gobin set
  BIN=${GOBIN}/metal-network-runner
fi

echo "launch metal-network-runner in the background"
$BIN server \
--log-level debug \
--port=":12342" \
--disable-grpc-gateway &
PID=${!}

############################
# By default, it runs all e2e test cases!
# Use "--ginkgo.skip" to skip tests.
# Use "--ginkgo.focus" to select tests.
echo "running e2e tests"
./tests/e2e/e2e.test \
--ginkgo.v \
--network-runner-log-level info \
--network-runner-grpc-endpoint="0.0.0.0:12342" \
--metalgo-path=${METALGO_PATH} \
--metalgo-plugin-dir=${METALGO_PLUGIN_DIR} \
--vm-genesis-path=/tmp/.genesis \
--vm-config-path=/tmp/.config \
--output-path=/tmp/metalgo-${metal_version}/output.yaml \
--mode=${MODE}
STATUS=$?

############################
if [[ -f "/tmp/metalgo-${metal_version}/output.yaml" ]]; then
  echo "cluster is ready!"
  cat /tmp/metalgo-${metal_version}/output.yaml
else
  echo "cluster is not ready in time... terminating ${PID}"
  kill ${PID}
  exit 255
fi

############################
if [[ ${MODE} == "test" ]]; then
  # "e2e.test" already terminates the cluster for "test" mode
  # just in case tests are aborted, manually terminate them again
  echo "network-runner RPC server was running on PID ${PID} as test mode; terminating the process..."
  pkill -P ${PID} || true
  kill -2 ${PID} || true
  pkill -9 -f jukr5oTVE2KfmEGwDRCvQnXcwrzwyzRtwyfSGFGPyuyNc12Fs || true # in case pkill didn't work
  exit ${STATUS}
else
  echo "network-runner RPC server is running on PID ${PID}..."
  echo ""
  echo "use the following command to terminate:"
  echo ""
  echo "pkill -P ${PID} && kill -2 ${PID} && pkill -9 -f jukr5oTVE2KfmEGwDRCvQnXcwrzwyzRtwyfSGFGPyuyNc12Fs"
  echo ""
fi
