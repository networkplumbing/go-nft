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

package config_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	nftconfig "github.com/networkplumbing/go-nft/nft/config"
	"github.com/networkplumbing/go-nft/nft/schema"
)

func TestDefineEmptyConfig(t *testing.T) {
	config := nftconfig.New()

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

	config := nftconfig.New()
	assert.NoError(t, config.FromJSON(serializedConfig))

	expectedConfig := nftconfig.New()
	expectedConfig.Nftables = append(expectedConfig.Nftables, schema.Nftable{Metainfo: &schema.Metainfo{
		Version:           version,
		ReleaseName:       releaseName,
		JsonSchemaVersion: schemaVersion,
	}})

	assert.Equal(t, expectedConfig, config)
}

func TestFlushRuleset(t *testing.T) {
	config := nftconfig.New()
	config.FlushRuleset()

	expected := []byte(`{"nftables":[{"flush":{"ruleset":null}}]}`)
	serializedConfig, err := config.ToJSON()
	assert.NoError(t, err)
	assert.Equal(t, string(expected), string(serializedConfig))
}

func TestFoo(t *testing.T) {
	normalizeConfig := func(in []byte) []byte {
		var foo interface{}
		err := json.Unmarshal(in, &foo)
		require.NoError(t, err)

		fmted, err := json.MarshalIndent(foo, "", "  ")
		require.NoError(t, err)
		return fmted
	}
	var cfg nftconfig.Config
	files := []string{
		"config1.json",
		"config2.json",
	}
	for _, f := range files {
		f := f
		name := fmt.Sprintf("file %s", f)
		t.Run(name, func(t *testing.T) {
			expected, err := os.ReadFile(path.Join("testdata", f))
			require.NoError(t, err)

			expected = normalizeConfig(expected)
			assert.NoError(t, cfg.FromJSON(expected))

			raw, err := json.MarshalIndent(cfg, "", "  ")
			require.NoError(t, err)
			raw = normalizeConfig(raw)
			if !assert.Equal(t, string(expected), string(raw)) {
				err = os.WriteFile(path.Join("testdata", "actual_"+f), raw, 0644)
				require.NoError(t, err)
			}
		})
	}
}
