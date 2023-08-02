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

package example

import (
	"testing"

	assert "github.com/stretchr/testify/require"

	"github.com/networkplumbing/go-nft/tests/testlib"

	"github.com/networkplumbing/go-nft/nft"
	"github.com/networkplumbing/go-nft/nft/schema"
)

// TestNATExamples are based on the examples given at [1].
// 1. https://wiki.nftables.org/wiki-nftables/index.php/Performing_Network_Address_Translation_(NAT)
//
// The inet family is mentioned in [1]:
// "Since Linux kernel 5.2, there is support for performing stateful NAT in inet family chains.".
// Therefore, inet family is not used in this example.
func TestNATExamples(t *testing.T) {
	testlib.RunTestWithFlushTable(t, testMasqueradeExample)
}

func testMasqueradeExample(t *testing.T) {
	desiredJson, actualJson, err := setupExample(buildMasqueradeConfig)
	assert.NoError(t, err)
	assert.JSONEq(t, desiredJson, actualJson)
}

func buildMasqueradeConfig() *nft.Config {
	const (
		ip4TableName = "ip4-masquerade-example"
		ip6TableName = "ip6-masquerade-example"

		baseChainName = "postrouting"
	)

	chainPriority := 100

	return &nft.Config{Root: schema.Root{Nftables: []schema.Nftable{
		{Table: &schema.Table{Family: schema.FamilyIP, Name: ip4TableName}},
		{Table: &schema.Table{Family: schema.FamilyIP6, Name: ip6TableName}},

		{Chain: &schema.Chain{
			Family: schema.FamilyIP,
			Table:  ip4TableName,
			Name:   baseChainName,
			Type:   schema.TypeNAT,
			Hook:   schema.HookPostRouting,
			Prio:   &chainPriority,
			Policy: schema.PolicyAccept,
		}},

		{Chain: &schema.Chain{
			Family: schema.FamilyIP6,
			Table:  ip6TableName,
			Name:   baseChainName,
			Type:   schema.TypeNAT,
			Hook:   schema.HookPostRouting,
			Prio:   &chainPriority,
			Policy: schema.PolicyAccept,
		}},

		{Rule: &schema.Rule{
			Family: schema.FamilyIP,
			Table:  ip4TableName,
			Chain:  baseChainName,
			Expr: []schema.Statement{
				{Nat: schema.Nat{Masquerade: &schema.Masquerade{Enabled: true}}},
			},
			Comment: "generic masquerade",
		}},

		{Rule: &schema.Rule{
			Family: schema.FamilyIP6,
			Table:  ip6TableName,
			Chain:  baseChainName,
			Expr: []schema.Statement{
				{Nat: schema.Nat{Masquerade: &schema.Masquerade{Enabled: true}}},
			},
			Comment: "generic masquerade",
		}},
	}}}
}
