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

package config_test

import (
	"encoding/json"
	"fmt"
	"testing"

	assert "github.com/stretchr/testify/require"

	"github.com/networkplumbing/go-nft/nft"
	"github.com/networkplumbing/go-nft/nft/schema"
)

type ruleAction string

// Rule Actions
const (
	ruleADD    ruleAction = "add"
	ruleDELETE ruleAction = "delete"
)

func TestRule(t *testing.T) {
	testAddRuleWithMatchAndVerdict(t)
	testDeleteRule(t)

	testAddRuleWithRowExpression(t)
	testAddRuleWithCounter(t)
	testAddRuleWithNamedCounter(t)
	testAddRuleWithNAT(t)

	testRuleLookup(t)

	testReadRuleWithNumericalExpression(t)
}

func testAddRuleWithRowExpression(t *testing.T) {
	const comment = "mycomment"

	table := nft.NewTable(tableName, nft.FamilyIP)
	chain := nft.NewRegularChain(table, chainName)

	t.Run("Add rule with a row expression, check serialization", func(t *testing.T) {
		statements, serializedStatements := matchWithRowExpression()
		rule := nft.NewRule(table, chain, statements, nil, nil, comment)

		config := nft.NewConfig()
		config.AddRule(rule)

		serializedConfig, err := config.ToJSON()
		assert.NoError(t, err)

		expectedConfig := buildSerializedConfig(ruleADD, serializedStatements, nil, nil, comment)
		assert.Equal(t, string(expectedConfig), string(serializedConfig))
	})

	t.Run("Add rule with a row expression, check deserialization", func(t *testing.T) {
		statements, serializedStatements := matchWithRowExpression()

		serializedConfig := buildSerializedConfig(ruleADD, serializedStatements, nil, nil, comment)

		var deserializedConfig nft.Config
		assert.NoError(t, json.Unmarshal(serializedConfig, &deserializedConfig))

		rule := nft.NewRule(table, chain, statements, nil, nil, comment)
		expectedConfig := nft.NewConfig()
		expectedConfig.AddRule(rule)

		assert.Equal(t, expectedConfig, &deserializedConfig)
	})
}

func testAddRuleWithMatchAndVerdict(t *testing.T) {
	const comment = "mycomment"

	table := nft.NewTable(tableName, nft.FamilyIP)
	chain := nft.NewRegularChain(table, chainName)

	t.Run("Add rule with match and verdict, check serialization", func(t *testing.T) {
		statements, serializedStatements := matchSrcIP4withReturnVerdict()
		rule := nft.NewRule(table, chain, statements, nil, nil, comment)

		config := nft.NewConfig()
		config.AddRule(rule)

		serializedConfig, err := config.ToJSON()
		assert.NoError(t, err)

		expectedConfig := buildSerializedConfig(ruleADD, serializedStatements, nil, nil, comment)
		assert.Equal(t, string(expectedConfig), string(serializedConfig))
	})

	t.Run("Add rule with match and verdict, check deserialization", func(t *testing.T) {
		statements, serializedStatements := matchSrcIP4withReturnVerdict()

		serializedConfig := buildSerializedConfig(ruleADD, serializedStatements, nil, nil, comment)

		var deserializedConfig nft.Config
		assert.NoError(t, json.Unmarshal(serializedConfig, &deserializedConfig))

		rule := nft.NewRule(table, chain, statements, nil, nil, comment)
		expectedConfig := nft.NewConfig()
		expectedConfig.AddRule(rule)

		assert.Equal(t, expectedConfig, &deserializedConfig)
	})
}

func testDeleteRule(t *testing.T) {
	table := nft.NewTable(tableName, nft.FamilyIP)
	chain := nft.NewRegularChain(table, chainName)

	t.Run("Delete rule", func(t *testing.T) {
		handleID := 100
		rule := nft.NewRule(table, chain, nil, &handleID, nil, "")

		config := nft.NewConfig()
		config.DeleteRule(rule)

		serializedConfig, err := config.ToJSON()
		assert.NoError(t, err)

		expectedConfig := buildSerializedConfig(ruleDELETE, "", nil, &handleID, "")
		assert.Equal(t, string(expectedConfig), string(serializedConfig))
	})
}

func buildSerializedConfig(
	action ruleAction,
	serializedRuleStatements string,
	namedCounter *schema.NamedCounter,
	handle *int,
	comment string,
) []byte {
	ruleArgs := fmt.Sprintf(`"family":%q,"table":%q,"chain":%q`, nft.FamilyIP, tableName, chainName)
	if serializedRuleStatements != "" {
		ruleArgs += "," + serializedRuleStatements
	}
	if handle != nil {
		ruleArgs += fmt.Sprintf(`,"handle":%d`, *handle)
	}
	if comment != "" {
		ruleArgs += fmt.Sprintf(`,"comment":%q`, comment)
	}

	serialzedCounter := ""
	if namedCounter != nil {
		counterArgs := fmt.Sprintf(`"family":%q,"table":%q,"name":%q`, nft.FamilyIP, tableName, namedCounter.Name)
		serialzedCounter = fmt.Sprintf(`{"counter":{%s}},`, counterArgs)
	}

	var config string
	if action == ruleADD {
		config = fmt.Sprintf(`{"nftables":[%s{"rule":{%s}}]}`, serialzedCounter, ruleArgs)
	} else {
		if namedCounter != nil {
			serialzedCounter = fmt.Sprintf(`{%q:{%s}},`, action, serialzedCounter)
		}
		config = fmt.Sprintf(`{"nftables":[%s{%q:{"rule":{%s}}}]}`, serialzedCounter, action, ruleArgs)
	}
	return []byte(config)
}

func matchSrcIP4withReturnVerdict() ([]schema.Statement, string) {
	ipAddress := "10.10.10.10"
	matchSrcIP4 := schema.Statement{
		Match: &schema.Match{
			Op: schema.OperEQ,
			Left: schema.Expression{
				Payload: &schema.Payload{
					Protocol: schema.PayloadProtocolIP4,
					Field:    schema.PayloadFieldIPSAddr,
				},
			},
			Right: schema.Expression{String: &ipAddress},
		},
	}

	verdict := schema.Statement{}
	verdict.Return = true

	statements := []schema.Statement{matchSrcIP4, verdict}

	expectedMatch := fmt.Sprintf(
		`"match":{"op":"==","left":{"payload":{"protocol":"ip","field":"saddr"}},"right":%q}`, ipAddress,
	)
	expectedVerdict := `"return":null`
	serializedStatements := fmt.Sprintf(`"expr":[{%s},{%s}]`, expectedMatch, expectedVerdict)

	return statements, serializedStatements
}

func matchWithRowExpression() ([]schema.Statement, string) {
	stringExpression := "string-expression"
	rowExpression := `{"foo":"boo"}`
	match := schema.Statement{
		Match: &schema.Match{
			Op:    schema.OperEQ,
			Left:  schema.Expression{RowData: json.RawMessage(rowExpression)},
			Right: schema.Expression{String: &stringExpression},
		},
	}

	statements := []schema.Statement{match}

	expectedMatch := fmt.Sprintf(`"match":{"op":"==","left":%s,"right":%q}`, rowExpression, stringExpression)
	serializedStatements := fmt.Sprintf(`"expr":[{%s}]`, expectedMatch)

	return statements, serializedStatements
}

func testRuleLookup(t *testing.T) {
	config := nft.NewConfig()
	table_br := nft.NewTable("table-br", nft.FamilyBridge)
	config.AddTable(table_br)

	chainRegular := nft.NewRegularChain(table_br, "chain-regular")
	config.AddChain(chainRegular)

	ruleSimple := nft.NewRule(table_br, chainRegular, nil, nil, nil, "comment123")
	config.AddRule(ruleSimple)

	ruleWithStatement := nft.NewRule(table_br, chainRegular, []schema.Statement{{}}, nil, nil, "comment456")
	ruleWithStatement.Expr[0].Drop = true
	config.AddRule(ruleWithStatement)

	handle := 10
	index := 1
	ruleWithAllParams := nft.NewRule(table_br, chainRegular, []schema.Statement{{}, {}}, &handle, &index, "comment789")
	config.AddRule(ruleWithAllParams)

	t.Run("Lookup an existing rule by table, chain and comment", func(t *testing.T) {
		rules := config.LookupRule(ruleSimple)
		assert.Len(t, rules, 1)
		assert.Equal(t, ruleSimple, rules[0])
	})

	t.Run("Lookup an existing rule by table, chain, statement and comment", func(t *testing.T) {
		rules := config.LookupRule(ruleWithStatement)
		assert.Len(t, rules, 1)
		assert.Equal(t, ruleWithStatement, rules[0])
	})

	t.Run("Lookup an existing rule by all (root) parameters", func(t *testing.T) {
		rules := config.LookupRule(ruleWithAllParams)
		assert.Len(t, rules, 1)
		assert.Equal(t, ruleWithAllParams, rules[0])
	})

	t.Run("Lookup a missing rule (comment not matching)", func(t *testing.T) {
		rule := nft.NewRule(table_br, chainRegular, nil, nil, nil, "comment-missing")
		assert.Empty(t, config.LookupRule(rule))
	})

	t.Run("Lookup a missing rule (statement content not matching)", func(t *testing.T) {
		rule := nft.NewRule(table_br, chainRegular, []schema.Statement{{}}, nil, nil, "comment456")
		rule.Expr[0].Drop = false
		rule.Expr[0].Return = true
		assert.Empty(t, config.LookupRule(rule))
	})

	t.Run("Lookup a missing rule (statements count not matching)", func(t *testing.T) {
		rule := nft.NewRule(table_br, chainRegular, []schema.Statement{{}, {}}, nil, nil, "comment456")
		rule.Expr[0].Drop = true
		assert.Empty(t, config.LookupRule(rule))
	})

	t.Run("Lookup a missing rule (handle not matching)", func(t *testing.T) {
		changedHandle := 99
		rule := nft.NewRule(table_br, chainRegular, []schema.Statement{{}, {}}, &changedHandle, &index, "comment789")
		assert.Empty(t, config.LookupRule(rule))
	})
}

func testReadRuleWithNumericalExpression(t *testing.T) {
	t.Run("Read rule with numerical expression", func(t *testing.T) {
		c := nft.NewConfig()
		assert.NoError(t, c.FromJSON([]byte(`
		{"nftables":[{"rule":{
		   "expr":[{"match":{"op":"==","left":"foo","right":12345}}]
		}}]}
		`)))
	})
}

func testAddRuleWithCounter(t *testing.T) {
	const comment = "mycomment"

	table := nft.NewTable(tableName, nft.FamilyIP)
	chain := nft.NewRegularChain(table, chainName)

	statements, serializedStatements := counterStatements()
	rule := nft.NewRule(table, chain, statements, nil, nil, comment)

	t.Run("Add rule with counter, check serialization", func(t *testing.T) {
		config := nft.NewConfig()
		config.AddRule(rule)

		serializedConfig, err := config.ToJSON()
		assert.NoError(t, err)

		expectedConfig := buildSerializedConfig(ruleADD, serializedStatements, nil, nil, comment)
		assert.JSONEq(t, string(expectedConfig), string(serializedConfig))
	})

	t.Run("Add rule with counter, check deserialization", func(t *testing.T) {
		serializedConfig := buildSerializedConfig(ruleADD, serializedStatements, nil, nil, comment)

		var deserializedConfig nft.Config
		assert.NoError(t, json.Unmarshal(serializedConfig, &deserializedConfig))

		expectedConfig := nft.NewConfig()
		expectedConfig.AddRule(rule)

		assert.Equal(t, expectedConfig, &deserializedConfig)
	})
}

func counterStatements() ([]schema.Statement, string) {
	statements := []schema.Statement{{
		Counter: &schema.Counter{
			Packets: 0,
			Bytes:   0,
		},
	}}

	expectedCounter := `"counter":{"packets":0,"bytes":0}`
	serializedStatements := fmt.Sprintf(`"expr":[{%s}]`, expectedCounter)

	return statements, serializedStatements
}

func testAddRuleWithNamedCounter(t *testing.T) {
	const comment = "mycomment"
	const counterName = "mycounter"

	table := nft.NewTable(tableName, nft.FamilyIP)
	chain := nft.NewRegularChain(table, chainName)
	counter := nft.NewCounter(table, counterName)

	statements, serializedStatements := namedCounterStatements(counterName)
	rule := nft.NewRule(table, chain, statements, nil, nil, comment)

	t.Run("Add rule with named counter, check serialization", func(t *testing.T) {
		config := nft.NewConfig()
		config.AddCounter(counter)
		config.AddRule(rule)

		serializedConfig, err := config.ToJSON()
		assert.NoError(t, err)

		expectedConfig := buildSerializedConfig(ruleADD, serializedStatements, counter, nil, comment)
		assert.JSONEq(t, string(expectedConfig), string(serializedConfig))
	})

	t.Run("Add rule with named counter, check deserialization", func(t *testing.T) {
		serializedConfig := buildSerializedConfig(ruleADD, serializedStatements, counter, nil, comment)

		var deserializedConfig nft.Config
		assert.NoError(t, json.Unmarshal(serializedConfig, &deserializedConfig))

		expectedConfig := nft.NewConfig()
		expectedConfig.AddCounter(counter)
		expectedConfig.AddRule(rule)

		assert.Equal(t, expectedConfig, &deserializedConfig)
	})
}

func namedCounterStatements(name string) ([]schema.Statement, string) {
	statements := []schema.Statement{{
		Counter: &schema.Counter{
			Name: name,
		},
	}}

	expectedCounter := fmt.Sprintf(`"counter": "%s"`, name)
	serializedStatements := fmt.Sprintf(`"expr":[{%s}]`, expectedCounter)

	return statements, serializedStatements
}

func testAddRuleWithNAT(t *testing.T) {
	tableTests := []struct {
		typeName         string
		createStatements func() ([]schema.Statement, string)
	}{
		{"dnat", dNATStatements},
		{"snat", sNATStatements},
		{"masquerade", masqueradeStatements},
		{"redirect", redirectStatements},
	}
	for _, tt := range tableTests {
		t.Run(fmt.Sprintf("Add rule with %s, check serialization", tt.typeName), func(t *testing.T) {
			testSerializationWith(t, dNATStatements)
		})
		t.Run(fmt.Sprintf("Add rule with %s, check deserialization", tt.typeName), func(t *testing.T) {
			testDeserializationWith(t, dNATStatements)
		})
	}
}

func testSerializationWith(t *testing.T, createStatements func() ([]schema.Statement, string)) {
	const comment = "mycomment"

	table := nft.NewTable(tableName, nft.FamilyIP)
	chain := nft.NewRegularChain(table, chainName)

	statements, serializedStatements := createStatements()
	rule := nft.NewRule(table, chain, statements, nil, nil, comment)

	config := nft.NewConfig()
	config.AddRule(rule)

	serializedConfig, err := config.ToJSON()
	assert.NoError(t, err)

	expectedConfig := buildSerializedConfig(ruleADD, serializedStatements, nil, nil, comment)
	assert.Equal(t, string(expectedConfig), string(serializedConfig))
}

func testDeserializationWith(t *testing.T, createStatements func() ([]schema.Statement, string)) {
	const comment = "mycomment"

	table := nft.NewTable(tableName, nft.FamilyIP)
	chain := nft.NewRegularChain(table, chainName)

	statements, serializedStatements := createStatements()

	serializedConfig := buildSerializedConfig(ruleADD, serializedStatements, nil, nil, comment)

	var deserializedConfig nft.Config
	assert.NoError(t, json.Unmarshal(serializedConfig, &deserializedConfig))

	rule := nft.NewRule(table, chain, statements, nil, nil, comment)
	expectedConfig := nft.NewConfig()
	expectedConfig.AddRule(rule)

	assert.Equal(t, expectedConfig, &deserializedConfig)
}

func dNATStatements() ([]schema.Statement, string) {
	address0 := "1.2.3.4"
	addressWithFamily := schema.Statement{}
	familyIP4 := schema.FamilyIP
	addressWithFamily.Dnat = &schema.Dnat{
		Addr:   &schema.Expression{String: &address0},
		Family: &familyIP4,
	}

	portList := schema.Statement{}
	portList.Dnat = &schema.Dnat{
		Port: &schema.Expression{RowData: json.RawMessage(`[80,8080]`)},
	}

	var port float64 = 12345
	address1 := "feed::c0fe"
	fullHouse := schema.Statement{}
	familyIP6 := schema.FamilyIP6
	fullHouse.Dnat = &schema.Dnat{
		Addr:   &schema.Expression{String: &address1},
		Family: &familyIP6,
		Port:   &schema.Expression{Float64: &port},
		Flags:  &schema.Flags{Flags: []string{schema.NATFlagRandom, schema.NATFlagPersistent}},
	}

	statements := []schema.Statement{addressWithFamily, portList, fullHouse}

	expectedDNATIP4 := `"dnat":{"addr":"1.2.3.4","family":"ip"}`
	expectedDNATMultiPorts := `"dnat":{"port":[80,8080]}`
	expectedDNATIP6PortAndFlags := `"dnat":{"addr":"feed::c0fe","family":"ip6","port":12345,"flags":["random","persistent"]}`
	serializedStatements := fmt.Sprintf(
		`"expr":[{%s},{%s},{%s}]`, expectedDNATIP4, expectedDNATMultiPorts, expectedDNATIP6PortAndFlags,
	)

	return statements, serializedStatements
}

func sNATStatements() ([]schema.Statement, string) {
	address0 := "1.2.3.4"
	addressWithFamily := schema.Statement{}
	familyIP4 := schema.FamilyIP
	addressWithFamily.Snat = &schema.Snat{
		Addr:   &schema.Expression{String: &address0},
		Family: &familyIP4,
	}

	portList := schema.Statement{}
	portList.Snat = &schema.Snat{
		Port: &schema.Expression{RowData: json.RawMessage(`[80,8080]`)},
	}

	var port float64 = 12345
	address1 := "feed::c0fe"
	fullHouse := schema.Statement{}
	familyIP6 := schema.FamilyIP6
	fullHouse.Snat = &schema.Snat{
		Addr:   &schema.Expression{String: &address1},
		Family: &familyIP6,
		Port:   &schema.Expression{Float64: &port},
		Flags:  &schema.Flags{Flags: []string{schema.NATFlagFullyRandom}},
	}

	statements := []schema.Statement{addressWithFamily, portList, fullHouse}

	expectedDNATIP4 := `"snat":{"addr":"1.2.3.4","family":"ip"}`
	expectedDNATMultiPorts := `"snat":{"port":[80,8080]}`
	expectedDNATIP6PortAndFlag := `"snat":{"addr":"feed::c0fe","family":"ip6","port":12345,"flags":"fully-random"}`
	serializedStatements := fmt.Sprintf(
		`"expr":[{%s},{%s},{%s}]`, expectedDNATIP4, expectedDNATMultiPorts, expectedDNATIP6PortAndFlag,
	)

	return statements, serializedStatements
}

func masqueradeStatements() ([]schema.Statement, string) {
	basic := schema.Statement{}
	basic.Masquerade = &schema.Masquerade{Enabled: true}

	portList := schema.Statement{}
	portList.Masquerade = &schema.Masquerade{
		Port: &schema.Expression{RowData: json.RawMessage(`[80,8080]`)},
	}

	var port float64 = 12345
	portAndFlags := schema.Statement{}
	portAndFlags.Masquerade = &schema.Masquerade{
		Port:  &schema.Expression{Float64: &port},
		Flags: &schema.Flags{Flags: []string{schema.NATFlagFullyRandom}},
	}

	statements := []schema.Statement{basic, portList, portAndFlags}

	expectedMasqueradeNoValues := `"masquerade":null`
	expectedMasqueradeMultiPorts := `"masquerade":{"port":[80,8080]}`
	expectedMasqueradePortAndFlag := `"masquerade":{"port":12345,"flags":"fully-random"}`
	serializedStatements := fmt.Sprintf(
		`"expr":[{%s},{%s},{%s}]`,
		expectedMasqueradeNoValues, expectedMasqueradeMultiPorts, expectedMasqueradePortAndFlag,
	)

	return statements, serializedStatements
}

func redirectStatements() ([]schema.Statement, string) {
	basic := schema.Statement{}
	basic.Redirect = &schema.Redirect{Enabled: true}

	portList := schema.Statement{}
	portList.Redirect = &schema.Redirect{
		Port: &schema.Expression{RowData: json.RawMessage(`[80,8080]`)},
	}

	var port float64 = 12345
	portAndFlags := schema.Statement{}
	portAndFlags.Redirect = &schema.Redirect{
		Port:  &schema.Expression{Float64: &port},
		Flags: &schema.Flags{Flags: []string{schema.NATFlagFullyRandom}},
	}

	statements := []schema.Statement{basic, portList, portAndFlags}

	expectedRedirectNoValues := `"redirect":null`
	expectedRedirectMultiPorts := `"redirect":{"port":[80,8080]}`
	expectedRedirectPortAndFlag := `"redirect":{"port":12345,"flags":"fully-random"}`
	serializedStatements := fmt.Sprintf(
		`"expr":[{%s},{%s},{%s}]`,
		expectedRedirectNoValues, expectedRedirectMultiPorts, expectedRedirectPortAndFlag,
	)

	return statements, serializedStatements
}
