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

package nft

import (
	"github.com/networkplumbing/go-nft/nft/schema"
)

// NewCounter returns a new schema counter structure for a named counter.
func NewCounter(table *schema.Table, name string) *schema.NamedCounter {
	c := &schema.NamedCounter{
		Family: table.Family,
		Table:  table.Name,
		Name:   name,
	}

	return c
}