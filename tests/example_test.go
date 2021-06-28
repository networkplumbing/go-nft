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
	"sort"
	"testing"

	"github.com/eddev/go-nft/nft"
	"github.com/eddev/go-nft/nft/schema"

	"github.com/stretchr/testify/assert"
)

func TestNoMacSpoofingExample(t *testing.T) {
	runTestWithFlushTable(t, testNoMacSpoofingExample)
}

func testNoMacSpoofingExample(t *testing.T) {
	// The definition of no-mac-spoofing, is to prevent users from modifying an interface mac-address and keep
	// connectivity.
	// A simple implementation is to allow only a pre-defined number of mac-addresses on egress (filtering the source
	// mac) and ingress (filtering the destination mac).

	// User input
	var (
		ifaceName  = "nic0"
		macAddress = "00:00:00:00:00:01"
	)

	desiredConfig := buildNoMacSpoofingConfigImperatively(ifaceName, macAddress)
	assert.Equal(t, desiredConfig, buildNoMacSpoofingConfigDecleratively(ifaceName, macAddress))
	assert.NoError(t, nft.ApplyConfig(desiredConfig))

	actualConfig, err := nft.ReadConfig()
	assert.NoError(t, err)

	expectedNftablesEntries := len(desiredConfig.Nftables) + 1 // +1 for the metainfo entry.
	assert.Len(t, actualConfig.Nftables, expectedNftablesEntries)

	desiredConfig = normalizeConfigForComparison(desiredConfig)
	actualConfig = normalizeConfigForComparison(actualConfig)

	desiredJson, err := desiredConfig.ToJSON()
	assert.NoError(t, err)
	actualJson, err := actualConfig.ToJSON()
	assert.NoError(t, err)
	assert.Equal(t, string(desiredJson), string(actualJson))
}

// normalizeConfigForComparison returns the configuration ready for comparison with another by
// - removing the metainfo entry.
// - removing the handle + index parameters.
// - Sorting the list.
func normalizeConfigForComparison(config *nft.Config) *nft.Config {
	if len(config.Nftables) > 0 && config.Nftables[0].Metainfo != nil {
		config.Nftables = config.Nftables[1:]
	}

	for _, nftable := range config.Nftables {
		if nftable.Rule != nil {
			nftable.Rule.Index = nil
			nftable.Rule.Handle = nil
		}
	}

	sort.Slice(config.Nftables, func(i int, j int) bool {
		s := config.Nftables
		isTableFirst := s[i].Table != nil && (s[j].Chain != nil || s[j].Rule != nil)
		isChainBeforeRule := s[i].Chain != nil && s[j].Rule != nil
		return isTableFirst || isChainBeforeRule
	})
	return config
}

func buildNoMacSpoofingConfigImperatively(ifaceName string, macAddress string) *nft.Config {
	// Configuration Details
	var (
		baseChainName  = "preroute-bridge"
		ifaceChainName = "example-iface-" + ifaceName
		macChainName   = ifaceChainName + "-mac"
	)
	config := nft.NewConfig()

	table := nft.NewTable("example", nft.FamilyBridge)
	config.AddTable(table)

	chainType, chainHook, chainPrio, chainPolicy := nft.TypeFilter, nft.HookPreRouting, -300, nft.PolicyAccept
	baseChain := nft.NewChain(table, baseChainName, &chainType, &chainHook, &chainPrio, &chainPolicy)
	config.AddChain(baseChain)

	ifaceChain := nft.NewRegularChain(table, ifaceChainName)
	config.AddChain(ifaceChain)

	macChain := nft.NewRegularChain(table, macChainName)
	config.AddChain(macChain)

	matchIfaceAndJump := []schema.Statement{
		{Match: &schema.Match{
			Op:    schema.OperEQ,
			Left:  schema.Expression{RowData: []byte(`{"meta":{"key":"iifname"}}`)},
			Right: schema.Expression{String: &ifaceName},
		}},
		{Verdict: schema.Verdict{Jump: &schema.ToTarget{Target: ifaceChainName}}},
	}
	matchIfaceRule := nft.NewRule(table, baseChain, matchIfaceAndJump, nil, nil, "match input interface name")
	config.AddRule(matchIfaceRule)

	jumpToMACChain := []schema.Statement{
		{Verdict: schema.Verdict{Jump: &schema.ToTarget{Target: macChainName}}},
	}
	ifaceRule := nft.NewRule(table, ifaceChain, jumpToMACChain, nil, nil, "redirect to mac-chain")
	config.AddRule(ifaceRule)

	matchSrcMacAndReturn := []schema.Statement{
		{Match: &schema.Match{
			Op: schema.OperEQ,
			Left: schema.Expression{Payload: &schema.Payload{
				Protocol: schema.PayloadProtocolEther,
				Field:    schema.PayloadFieldEtherSAddr,
			}},
			Right: schema.Expression{String: &macAddress},
		}},
		{Verdict: schema.Return()},
	}
	matchSrcMacRule := nft.NewRule(table, macChain, matchSrcMacAndReturn, nil, nil, "match source mac address")
	config.AddRule(matchSrcMacRule)

	drop := []schema.Statement{{Verdict: schema.Drop()}}
	// When multiple rules are added to a chain, index allows to define an order between them.
	macRulesIndex := nft.NewRuleIndex()
	dropRule := nft.NewRule(table, macChain, drop, nil, macRulesIndex.Next(), "drop all the rest")
	config.AddRule(dropRule)

	return config
}

func buildNoMacSpoofingConfigDecleratively(ifaceName string, macAddress string) *nft.Config {
	// Configuration Details
	const tableName = "example"
	var (
		baseChainName  = "preroute-bridge"
		ifaceChainName = "example-iface-" + ifaceName
		macChainName   = ifaceChainName + "-mac"

		chainPriority = -300

		macRulesIndex = nft.NewRuleIndex()
	)

	return &nft.Config{schema.Root{Nftables: []schema.Nftable{
		{Table: &schema.Table{Family: schema.FamilyBridge, Name: tableName}},

		{Chain: &schema.Chain{
			Family: schema.FamilyBridge,
			Table:  tableName,
			Name:   baseChainName,
			Type:   schema.TypeFilter,
			Hook:   schema.HookPreRouting,
			Prio:   &chainPriority,
			Policy: schema.PolicyAccept,
		}},
		{Chain: &schema.Chain{
			Family: schema.FamilyBridge,
			Table:  tableName,
			Name:   ifaceChainName,
		}},
		{Chain: &schema.Chain{
			Family: schema.FamilyBridge,
			Table:  tableName,
			Name:   macChainName,
		}},

		{Rule: &schema.Rule{
			Family: schema.FamilyBridge,
			Table:  tableName,
			Chain:  baseChainName,
			Expr: []schema.Statement{
				{Match: &schema.Match{
					Op:    schema.OperEQ,
					Left:  schema.Expression{RowData: []byte(`{"meta":{"key":"iifname"}}`)},
					Right: schema.Expression{String: &ifaceName},
				}},
				{Verdict: schema.Verdict{Jump: &schema.ToTarget{Target: ifaceChainName}}},
			},
			Comment: "match input interface name",
		}},
		{Rule: &schema.Rule{
			Family: schema.FamilyBridge,
			Table:  tableName,
			Chain:  ifaceChainName,
			Expr: []schema.Statement{
				{Verdict: schema.Verdict{Jump: &schema.ToTarget{Target: macChainName}}},
			},
			Comment: "redirect to mac-chain",
		}},
		{Rule: &schema.Rule{
			Family: schema.FamilyBridge,
			Table:  tableName,
			Chain:  macChainName,
			Expr: []schema.Statement{
				{Match: &schema.Match{
					Op: schema.OperEQ,
					Left: schema.Expression{Payload: &schema.Payload{
						Protocol: schema.PayloadProtocolEther,
						Field:    schema.PayloadFieldEtherSAddr,
					}},
					Right: schema.Expression{String: &macAddress},
				}},
				{Verdict: schema.Return()},
			},
			Comment: "match source mac address",
		}},
		{Rule: &schema.Rule{
			Family:  schema.FamilyBridge,
			Table:   tableName,
			Chain:   macChainName,
			Index:   macRulesIndex.Next(),
			Expr:    []schema.Statement{{Verdict: schema.Drop()}},
			Comment: "drop all the rest",
		}},
	}}}
}
