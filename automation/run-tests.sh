#!/usr/bin/env bash

set -e

EXEC_PATH=$(dirname "$(realpath "$0")")
PROJECT_PATH="$(dirname $EXEC_PATH)"

CONTAINER_WORKSPACE="/workspace/go-nft"

: "${CONTAINER_CMD:="docker"}"
: "${CONTAINER_IMG:="golang:1.16.4-alpine3.13"}"

: "${DISABLE_IPV6_IN_CONTAINER:=0}"

test -t 1 && USE_TTY="-t"

options=$(getopt --options "" \
    --long build,fmt,unit-test,integration-test,help\
    -- "${@}")
eval set -- "$options"
while true; do
    case "$1" in
    --build)
        OPT_BUILD=1
        ;;
    --fmt)
        OPT_FMT=1
        ;;
    --unit-test)
        OPT_UTEST=1
        ;;
    --integration-test)
        OPT_ITEST=1
        ;;
    --help)
        set +x
        echo "$0 [--build] [--fmt] [--unit-test] [--integration-test]"
        exit
        ;;
    --)
        shift
        break
        ;;
    esac
    shift
done

function run_container {
    ${CONTAINER_CMD} run \
        $USE_TTY \
        -i \
        --rm \
        --cap-add=NET_ADMIN \
        --cap-add=NET_RAW \
        --sysctl net.ipv6.conf.all.disable_ipv6=$DISABLE_IPV6_IN_CONTAINER \
        -v "$PROJECT_PATH":"$CONTAINER_WORKSPACE":Z \
        -w "$CONTAINER_WORKSPACE" \
        "$CONTAINER_IMG" \
    sh -c "$1"
}

if [ -z "${OPT_BUILD}" ] && [ -z "${OPT_FMT}" ] && [ -z "${OPT_UTEST}" ] && [ -z "${OPT_ITEST}" ]; then
    OPT_BUILD=1
    OPT_FMT=1
    OPT_UTEST=1
    OPT_ITEST=1
fi

if [ -n "${OPT_BUILD}" ]; then
    go build -v ./...
fi

if [ -n "${OPT_FMT}" ]; then
        unformatted=$(gofmt -l ./nft ./tests)
        test -z "$unformatted" || (echo "Unformatted: $unformatted" && false)

        unformatted=$(go run golang.org/x/tools/cmd/goimports -l --local "github.com/networkplumbing/go-nft" ./nft ./tests)
        test -z "$unformatted" || (echo "Unformatted imports: $unformatted" && false)
fi

if [ -n "${OPT_UTEST}" ]; then
    go test -v ./nft/...
fi

if [ -n "${OPT_ITEST}" ]; then
    # Manually load `nft_masq` kmod on the host to support NAT definitions.
    # The container is unable to load a kmod (usually done by `nft` automatically).
    sudo modprobe nft_masq
    run_container '
        apk add --no-cache nftables gcc musl-dev
        nft -j list ruleset
        go test -v --tags=exec ./tests/...
        apk add --no-cache nftables-dev
        go test -v ./tests/nftlib
    '
fi
