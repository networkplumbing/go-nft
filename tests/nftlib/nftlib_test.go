//go:build !exec
// +build !exec

/* This file is part of the go-nft project
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
package main

import (
	"testing"

	assert "github.com/stretchr/testify/require"

	"github.com/networkplumbing/go-nft/nft"
	nftlib "github.com/networkplumbing/go-nft/nft/lib"

	"github.com/networkplumbing/go-nft/tests/testlib"
)

func TestNftlib(t *testing.T) {
	testlib.RunTestWithFlushTable(t, func(t *testing.T) {
		config := nft.NewConfig()
		config.AddTable(nft.NewTable("mytable", nft.FamilyIP))

		assert.NoError(t, nftlib.ApplyConfig(config))

		newConfig, err := nftlib.ReadConfig()
		assert.NoError(t, err)

		assert.Len(t, newConfig.Nftables, 2, "Expecting the metainfo and an empty table entry")
		newConfig = testlib.NormalizeConfigForComparison(newConfig)
		assert.Equal(t, config.Nftables[0], newConfig.Nftables[0])
		_, err = nftlib.ReadConfig()
		assert.NoError(t, err)

		newConfig, err = nftlib.ApplyConfigEcho(config)
		assert.NoError(t, err)
		newConfig = testlib.NormalizeConfigForComparison(newConfig)
		assert.Len(t, newConfig.Nftables, 1, "Expecting just the empty table entry")
		assert.Equal(t, config.Nftables, newConfig.Nftables)
	})
}
