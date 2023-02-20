# go-nft

[![Licensed under Apache License version 2.0](https://img.shields.io/github/license/kubevirt/kubevirt.svg)](https://www.apache.org/licenses/LICENSE-2.0)
[![Build Workflow](https://github.com/networkplumbing/go-nft/actions/workflows/main.yml/badge.svg)](https://github.com/networkplumbing/go-nft/actions/workflows/main.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/networkplumbing/go-nft)](https://goreportcard.com/report/github.com/networkplumbing/go-nft)

Go bindings for nft utility.

go-nft wraps invocation of the nft utility with functions to append and delete
rules; create, clear and delete tables and chains.

## To start using go-nft

go-nft is a library that provides a structured API to nftables.

go-nft uses the [libnftables-json specification](https://www.mankier.com/5/libnftables-json)
and exposes a subset of its structures.

- Apply the configuration:
```golang
config := nft.NewConfig()
config.AddTable(nft.NewTable("mytable", nft.FamilyIP))
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
err := nft.ApplyConfigContext(ctx, config)
```

- Read the configuration:
```golang
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
config, err := nft.ReadConfigContext(ctx)
nftVersion := config.Nftables[0].Metainfo.Version
```

For full setup example, see the integration test [examples](tests/example).

## Contribution

We welcome contribution of any kind!
Read [CONTRIBUTING](CONTRIBUTING.md) to learn how to contribute to the project.

## Changelog

Please refer to [CHANGELOG](CHANGELOG)
