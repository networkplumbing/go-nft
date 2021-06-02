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

package nft_test

import (
	"fmt"
	"testing"

	"github.com/eddev/go-nft/nft/schema"

	"github.com/stretchr/testify/assert"

	"github.com/eddev/go-nft/nft"
)

func TestDefineEmptyConfig(t *testing.T) {
	config := nft.NewConfig()

	expected := []byte(`{"nftables":[]}`)
	serializedConfig, err := config.ToJSON()
	assert.NoError(t, err)
	assert.Equal(t, string(expected), string(serializedConfig))
}

func TestReadEmptyConfigWithMetaInfo(t *testing.T) {
	const version = "0.9.3"
	const releaseName = "Topsy"
	const schemaVersion = 1
	serializedConfig := []byte(fmt.Sprintf(
		`{"nftables":[{"metainfo":{"version":%q,"release_name":%q,"json_schema_version":%d}}]}`,
		version, releaseName, schemaVersion,
	))

	config := nft.NewConfig()
	assert.NoError(t, config.FromJSON(serializedConfig))

	expectedConfig := nft.NewConfig()
	expectedConfig.Nftables = append(expectedConfig.Nftables, schema.Nftable{Metainfo: &schema.Metainfo{
		Version:           version,
		ReleaseName:       releaseName,
		JsonSchemaVersion: schemaVersion,
	}})

	assert.Equal(t, expectedConfig, config)
}

func TestFlushRuleset(t *testing.T) {
	config := nft.NewConfig()
	config.FlushRuleset()

	expected := []byte(`{"nftables":[{"flush":{"ruleset":null}}]}`)
	serializedConfig, err := config.ToJSON()
	assert.NoError(t, err)
	assert.Equal(t, string(expected), string(serializedConfig))
}
