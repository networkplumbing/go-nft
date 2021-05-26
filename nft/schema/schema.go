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

package schema

type Root struct {
	Nftables []Nftable `json:"nftables"`
}

type Objects struct {
	Table *Table `json:"table,omitempty"`
	Chain *Chain `json:"chain,omitempty"`
}

type Nftable struct {
	Table *Table `json:"table,omitempty"`
	Chain *Chain `json:"chain,omitempty"`

	Add    *Objects `json:"add,omitempty"`
	Delete *Objects `json:"delete,omitempty"`
	Flush  *Objects `json:"flush,omitempty"`
}

type Table struct {
	Family string `json:"family"`
	Name   string `json:"name"`
}

type Chain struct {
	Family string `json:"family"`
	Table  string `json:"table"`
	Name   string `json:"name"`
	Type   string `json:"type,omitempty"`
	Hook   string `json:"hook,omitempty"`
	Prio   *int   `json:"prio,omitempty"`
	Policy string `json:"policy,omitempty"`
}
