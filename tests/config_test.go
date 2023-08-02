/*
 * This file is part of the go-nft project
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * Copyright 2021 Red Hat, Inc.
 *
 */

package tests

import (
	"testing"

	assert "github.com/stretchr/testify/require"

	"github.com/networkplumbing/go-nft/tests/testlib"

	"github.com/networkplumbing/go-nft/nft"
	"github.com/networkplumbing/go-nft/nft/schema"

	"context"
	"time"
)

func TestConfig(t *testing.T) {
	testlib.RunTestWithFlushTable(t, testReadEmptyConfig)
	testlib.RunTestWithFlushTable(t, testApplyConfigWithAnEmptyTable)
	testlib.RunTestWithFlushTable(t, testReadFilteredConfig)
	testlib.RunTestWithFlushTable(t, testApplyConfigWithSampleStatements)
}

func testReadEmptyConfig(t *testing.T) {
	config, err := nft.ReadConfig()
	assert.NoError(t, err)
	assert.Len(t, config.Nftables, 1, "Expecting just the metainfo entry")

	expectedConfig := nft.NewConfig()
	expectedConfig.Nftables = append(expectedConfig.Nftables, schema.Nftable{Metainfo: &schema.Metainfo{JsonSchemaVersion: 1}})

	// The underling nftable userspace version depends on where it is ran, therefore remove it from the diff.
	expectedConfig.Nftables[0].Metainfo.Version = config.Nftables[0].Metainfo.Version
	expectedConfig.Nftables[0].Metainfo.ReleaseName = config.Nftables[0].Metainfo.ReleaseName
	assert.Equal(t, expectedConfig, config)
}

func testApplyConfigWithAnEmptyTable(t *testing.T) {
	config := nft.NewConfig()
	config.AddTable(nft.NewTable("mytable", nft.FamilyIP))

	assert.NoError(t, nft.ApplyConfig(config))

	newConfig, err := nft.ReadConfig()
	assert.NoError(t, err)

	assert.Len(t, newConfig.Nftables, 2, "Expecting the metainfo and an empty table entry")
	newConfig = testlib.NormalizeConfigForComparison(newConfig)
	assert.Equal(t, config.Nftables[0], newConfig.Nftables[0])
}

func testReadFilteredConfig(t *testing.T) {
	const (
		tableName1 = "mytable1"
		tableName2 = "mytable2"
		chainName1 = "mychain1"
		chainName2 = "mychain2"
	)
	config := nft.NewConfig()
	table1 := nft.NewTable(tableName1, nft.FamilyIP)
	table2 := nft.NewTable(tableName2, nft.FamilyIP)
	config.AddTable(table1)
	config.AddTable(table2)

	chain1 := nft.NewChain(table1, chainName1, nil, nil, nil, nil)
	chain2 := nft.NewChain(table1, chainName2, nil, nil, nil, nil)
	config.AddChain(chain1)
	config.AddChain(chain2)
	assert.NoError(t, nft.ApplyConfig(config))

	singleTableConfig, err := nft.ReadConfig("table", table2.Family, table2.Name)
	assert.NoError(t, err)
	assert.Len(t, singleTableConfig.Nftables, 2, "Expecting the metainfo and an empty table entry")
	singleTableConfig = testlib.NormalizeConfigForComparison(singleTableConfig)
	assert.Equal(t, config.Nftables[1], singleTableConfig.Nftables[0])

	singleChainConfig, err := nft.ReadConfig("chain", chain1.Family, chain1.Table, chain1.Name)
	assert.NoError(t, err)
	assert.Len(t, singleChainConfig.Nftables, 2, "Expecting the metainfo and an empty chain entry")
	singleChainConfig = testlib.NormalizeConfigForComparison(singleChainConfig)
	assert.Equal(t, config.Nftables[2], singleChainConfig.Nftables[0])
}

func testApplyConfigWithSampleStatements(t *testing.T) {
	testApplyConfigWithStatements(t,
		schema.Statement{Counter: &schema.Counter{}},
	)
}

func testApplyConfigWithStatements(t *testing.T, statements ...schema.Statement) {
	const tableName = "mytable"
	config := nft.NewConfig()
	table := nft.NewTable(tableName, nft.FamilyIP)
	config.AddTable(table)

	const chainName = "mychain"
	chain := nft.NewChain(table, chainName, nil, nil, nil, nil)
	config.AddChain(chain)

	rule := nft.NewRule(table, chain, statements, nil, nil, "test")
	config.AddRule(rule)

	assert.NoError(t, nft.ApplyConfig(config))

	newConfig, err := nft.ReadConfig()
	assert.NoError(t, err)

	config = testlib.NormalizeConfigForComparison(config)
	newConfig = testlib.NormalizeConfigForComparison(newConfig)
	assert.Equal(t, config.Nftables, newConfig.Nftables)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	newConfig, err = nft.ApplyConfigEcho(ctx, config)
	assert.NoError(t, err)

	newConfig = testlib.NormalizeConfigForComparison(newConfig)
	assert.Equal(t, config.Nftables, newConfig.Nftables)
}
