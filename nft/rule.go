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
	"github.com/eddev/go-nft/nft/schema"
)

func NewRule(table *schema.Table, chain *schema.Chain, expr []schema.Statement, handle *int, comment string) *schema.Rule {
	c := &schema.Rule{
		Family:  table.Family,
		Table:   table.Name,
		Chain:   chain.Name,
		Expr:    expr,
		Handle:  handle,
		Comment: comment,
	}

	return c
}

func (c *Config) AddRule(rule *schema.Rule) {
	nftable := schema.Nftable{Rule: rule}
	c.Nftables = append(c.Nftables, nftable)
}

func (c *Config) DeleteRule(rule *schema.Rule) {
	nftable := schema.Nftable{Delete: &schema.Objects{Rule: rule}}
	c.Nftables = append(c.Nftables, nftable)
}
