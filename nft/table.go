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

import "github.com/eddev/go-nft/nft/schema"

type AddressFamily string

// Address Families
const (
	FamilyIP     AddressFamily = schema.FamilyIP
	FamilyIP6    AddressFamily = schema.FamilyIP6
	FamilyINET   AddressFamily = schema.FamilyINET
	FamilyARP    AddressFamily = schema.FamilyARP
	FamilyBridge AddressFamily = schema.FamilyBridge
	FamilyNETDEV AddressFamily = schema.FamilyNETDEV
)

type TableAction string

// Table Actions
const (
	TableADD    TableAction = "add"
	TableDELETE TableAction = "delete"
	TableFLUSH  TableAction = "flush"
)

func NewTable(name string, family AddressFamily) *schema.Table {
	return &schema.Table{
		Name:   name,
		Family: string(family),
	}
}

func (c *Config) AddTable(table *schema.Table) {
	nftable := schema.Nftable{Table: table}
	c.Nftables = append(c.Nftables, nftable)
}

func (c *Config) DeleteTable(table *schema.Table) {
	nftable := schema.Nftable{Delete: &schema.Objects{Table: table}}
	c.Nftables = append(c.Nftables, nftable)
}

func (c *Config) FlushTable(table *schema.Table) {
	nftable := schema.Nftable{Flush: &schema.Objects{Table: table}}
	c.Nftables = append(c.Nftables, nftable)
}

func (c *Config) LookupTable(toFind *schema.Table) *schema.Table {
	for _, nftable := range c.Nftables {
		if t := nftable.Table; t != nil {
			if t.Name == toFind.Name && t.Family == toFind.Family {
				return t
			}
		}
	}
	return nil
}
