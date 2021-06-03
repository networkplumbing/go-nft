#!/usr/bin/env bash

set -e

gofmt -l -w ./nft ./tests
go run golang.org/x/tools/cmd/goimports -l --local "github.com/eddev/go-nft" ./nft ./tests
