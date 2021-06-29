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

package nft

import (
	"bytes"
	"encoding/json"

	"github.com/eddev/go-nft/nft/schema"
)

// NewRule returns a new schema rule structure.
func NewRule(table *schema.Table, chain *schema.Chain, expr []schema.Statement, handle *int, index *int, comment string) *schema.Rule {
	c := &schema.Rule{
		Family:  table.Family,
		Table:   table.Name,
		Chain:   chain.Name,
		Expr:    expr,
		Handle:  handle,
		Index:   index,
		Comment: comment,
	}

	return c
}

// AddRule appends the given rule to the nftable config.
// The rule is added without an explicit action (`add`).
// Adding multiple times the same rule will result in multiple identical rules when applied.
func (c *Config) AddRule(rule *schema.Rule) {
	nftable := schema.Nftable{Rule: rule}
	c.Nftables = append(c.Nftables, nftable)
}

// DeleteRule appends a given rule to the nftable config
// with the `delete` action.
// A rule is identified by its handle ID and it must be present in the given rule.
// Attempting to delete a non-existing rule, results with a failure when the config is applied.
// A common usage is to use LookupRule() and then to pass the result to DeleteRule.
func (c *Config) DeleteRule(rule *schema.Rule) {
	nftable := schema.Nftable{Delete: &schema.Objects{Rule: rule}}
	c.Nftables = append(c.Nftables, nftable)
}

// LookupRule searches the configuration for a matching rule and returns it.
// The rule is matched first by the table and chain.
// Other matching fields are optional (nil or an empty string arguments imply no-matching).
// Mutating the returned chain will result in mutating the configuration.
func (c *Config) LookupRule(toFind *schema.Rule) []*schema.Rule {
	var rules []*schema.Rule

	for _, nftable := range c.Nftables {
		if r := nftable.Rule; r != nil {
			match := r.Table == toFind.Table && r.Family == toFind.Family && r.Chain == toFind.Chain
			if match {
				if h := toFind.Handle; h != nil {
					match = match && r.Handle != nil && *r.Handle == *h
				}
				if i := toFind.Index; i != nil {
					match = match && r.Index != nil && *r.Index == *i
				}
				if co := toFind.Comment; co != "" {
					match = match && r.Comment == co
				}
				if toFindStatements := toFind.Expr; toFindStatements != nil {
					if match = match && len(toFindStatements) == len(r.Expr); match {
						for i, toFindStatement := range toFindStatements {
							equal, err := areStatementsEqual(toFindStatement, r.Expr[i])
							match = match && err == nil && equal
						}
					}
				}
				if match {
					rules = append(rules, r)
				}
			}
		}
	}
	return rules
}

func areStatementsEqual(statementA, statementB schema.Statement) (bool, error) {
	statementARow, err := json.Marshal(statementA)
	if err != nil {
		return false, err
	}
	statementBRow, err := json.Marshal(statementB)
	if err != nil {
		return false, err
	}
	return bytes.Equal(statementARow, statementBRow), nil
}

type RuleIndex int

// NewRuleIndex returns a rule index object which acts as an iterator.
// When multiple rules are added to a chain, index allows to define an order between them.
// The first rule which is added to a chain should have no index (it is assigned index 0),
// following rules should have the index set, referencing after/before which rule the new one is to be added/inserted.
func NewRuleIndex() *RuleIndex {
	var index RuleIndex = -1
	return &index
}

// Next returns the next iteration value as an integer pointer.
// When first time called, it returns the value 0.
func (i *RuleIndex) Next() *int {
	*i++
	var index = int(*i)
	return &index
}
