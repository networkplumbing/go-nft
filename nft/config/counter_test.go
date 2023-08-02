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

package config_test

import (
	"fmt"
	"testing"

	assert "github.com/stretchr/testify/require"

	"github.com/networkplumbing/go-nft/nft"
	"github.com/networkplumbing/go-nft/nft/schema"
)

type CounterAction string

// Counter Actions
const (
	CounterADD    CounterAction = "add"
	CounterDELETE CounterAction = "delete"
)

type counterActionFunc func(*nft.Config, *schema.NamedCounter)

const counterName = "test-counter"

func TestCounter(t *testing.T) {
	testCounterActions(t)
	testCounterLookup(t)
}

func testCounterActions(t *testing.T) {
	actions := map[CounterAction]counterActionFunc{
		CounterADD:    func(c *nft.Config, t *schema.NamedCounter) { c.AddCounter(t) },
		CounterDELETE: func(c *nft.Config, t *schema.NamedCounter) { c.DeleteCounter(t) },
	}
	table := nft.NewTable(tableName, nft.FamilyIP)
	counter := nft.NewCounter(table, counterName)

	for action, actionFunc := range actions {
		testName := fmt.Sprintf("%s counter", string(action))

		t.Run(testName, func(t *testing.T) {
			config := nft.NewConfig()
			actionFunc(config, counter)

			serializedConfig, err := config.ToJSON()
			assert.NoError(t, err)

			counterArgs := fmt.Sprintf(`"family":%q,"table":%q,"name":%q`, table.Family, table.Name, counterName)
			var expected []byte
			if action == CounterADD {
				expected = []byte(fmt.Sprintf(`{"nftables":[{"counter":{%s}}]}`, counterArgs))
			} else {
				expected = []byte(fmt.Sprintf(`{"nftables":[{%q:{"counter":{%s}}}]}`, action, counterArgs))
			}

			assert.Equal(t, string(expected), string(serializedConfig))
		})
	}
}

func testCounterLookup(t *testing.T) {
	config := nft.NewConfig()
	table_inet := nft.NewTable("table-inet", nft.FamilyINET)
	config.AddTable(table_inet)

	counter := nft.NewCounter(table_inet, "test-counter")
	config.AddCounter(counter)

	t.Run("Lookup an existing named counter", func(t *testing.T) {
		c := config.LookupCounter(counter)
		assert.Equal(t, counter, c)
	})

	t.Run("Lookup a missing named counter", func(t *testing.T) {
		c := nft.NewCounter(table_inet, "counter-na")
		assert.Nil(t, config.LookupCounter(c))
	})
}
