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
 * Copyright 2023 Anapaya Systems
 *
 */

package config

import (
	"github.com/networkplumbing/go-nft/nft/schema"
)

// AddCounter appends the given named counter to the nftable config.
// The rule is added without an explicit action (`add`).
// Adding multiple times the same named counter has no affect when the config is applied.
func (c *Config) AddCounter(counter *schema.NamedCounter) {
	nftable := schema.Nftable{Counter: counter}
	c.Nftables = append(c.Nftables, nftable)
}

// DeleteCounter appends a given rule to the nftable config with the `delete` action.
// Attempting to delete a non-existing named counter, results with a failure when the config is applied.
func (c *Config) DeleteCounter(counter *schema.NamedCounter) {
	nftable := schema.Nftable{Delete: &schema.Objects{Counter: counter}}
	c.Nftables = append(c.Nftables, nftable)
}

// LookupCounter searches the configuration for a matching counter and returns it.
// The counter is matched first by the table.
// Mutating the returned counter will result in mutating the configuration.
func (c *Config) LookupCounter(toFind *schema.NamedCounter) *schema.NamedCounter {
	for _, nftable := range c.Nftables {
		if counter := nftable.Counter; counter != nil {
			match := counter.Table == toFind.Table && counter.Family == toFind.Family && counter.Name == toFind.Name
			if match {
				return counter
			}
		}
	}
	return nil
}
