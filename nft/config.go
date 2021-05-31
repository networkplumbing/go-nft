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
	"encoding/json"

	"github.com/eddev/go-nft/nft/schema"
)

type Config struct {
	schema.Root
}

func NewConfig() *Config {
	c := &Config{}
	c.Nftables = []schema.Nftable{}
	return c
}

func (c *Config) MarshalJSON() ([]byte, error) {
	return json.Marshal(*c)
}

func (c *Config) FromJSON(data []byte) error {
	if err := json.Unmarshal(data, c); err != nil {
		return err
	}
	return nil
}
