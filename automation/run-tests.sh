#!/usr/bin/env bash

set -xe

EXEC_PATH=$(dirname "$(realpath "$0")")
PROJECT_PATH="$(dirname $EXEC_PATH)"

CONTAINER_WORKSPACE="/workspace/go-nft"

: "${CONTAINER_CMD:="docker"}"
: "${CONTAINER_IMG:="golang:alpine"}"

test -t 1 && USE_TTY="-t"

function run_container {
    ${CONTAINER_CMD} run $USE_TTY -i --rm  --cap-add=NET_ADMIN --cap-add=NET_RAW -v "$PROJECT_PATH":"$CONTAINER_WORKSPACE":Z -w "$CONTAINER_WORKSPACE" "$CONTAINER_IMG" \
    sh -c "$1"
}

run_container '
    apk add --no-cache nftables gcc musl-dev
    nft -j list ruleset
    go test -v ./tests
'