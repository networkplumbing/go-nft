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

type ChainType string
type ChainHook string
type ChainPolicy string

// Chain Types
const (
	TypeFilter ChainType = "filter"
	TypeNAT    ChainType = "nat"
	TypeRoute  ChainType = "route"
)

// Chain Hooks
const (
	HookPreRouting  ChainHook = "prerouting"
	HookInput       ChainHook = "input"
	HookOutput      ChainHook = "output"
	HookForward     ChainHook = "forward"
	HookPostRouting ChainHook = "postrouting"
	HookIngress     ChainHook = "ingress"
)

// Chain Policies
const (
	PolicyAccept ChainPolicy = "accept"
	PolicyDrop   ChainPolicy = "drop"
)

func NewRegularChain(table *schema.Table, name string) *schema.Chain {
	return NewChain(table, name, nil, nil, nil, nil)
}
func NewChain(table *schema.Table, name string, ctype *ChainType, hook *ChainHook, prio *int, policy *ChainPolicy) *schema.Chain {
	c := &schema.Chain{
		Family: table.Family,
		Table:  table.Name,
		Name:   name,
	}

	if ctype != nil {
		c.Type = string(*ctype)
	}
	if hook != nil {
		c.Hook = string(*hook)
	}
	if prio != nil {
		c.Prio = prio
	}
	if policy != nil {
		c.Policy = string(*policy)
	}

	return c
}

func (c *Config) AddChain(chain *schema.Chain) {
	nftable := schema.Nftable{Chain: chain}
	c.Nftables = append(c.Nftables, nftable)
}

func (c *Config) DeleteChain(chain *schema.Chain) {
	nftable := schema.Nftable{Delete: &schema.Objects{Chain: chain}}
	c.Nftables = append(c.Nftables, nftable)
}

func (c *Config) FlushChain(chain *schema.Chain) {
	nftable := schema.Nftable{Flush: &schema.Objects{Chain: chain}}
	c.Nftables = append(c.Nftables, nftable)
}
