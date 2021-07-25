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

package example

import (
	"fmt"

	"github.com/networkplumbing/go-nft/tests/testlib"

	"github.com/networkplumbing/go-nft/nft"
)

type buildConfigT func() *nft.Config

func setupExample(buildConfig buildConfigT) (string, string, error) {
	desiredConfig := buildConfig()
	if err := nft.ApplyConfig(desiredConfig); err != nil {
		return "", "", fmt.Errorf("failed to setup example: %v", err)
	}

	actualConfig, err := nft.ReadConfig()
	if err != nil {
		return "", "", fmt.Errorf("failed to setup example: %v", err)
	}

	desiredNftablesEntries := len(desiredConfig.Nftables) + 1 // +1 for the metainfo entry.
	actualNftablesEntries := len(actualConfig.Nftables)
	if actualNftablesEntries != desiredNftablesEntries {
		desiredJson, _ := desiredConfig.ToJSON()
		actualJson, _ := actualConfig.ToJSON()
		return "", "", fmt.Errorf(
			"failed to setup example, unexpected entries in post-setup resuts, desired(%d): %s, actual(%d): %s",
			desiredNftablesEntries, desiredJson, actualNftablesEntries, actualJson,
		)
	}

	desiredConfig = testlib.NormalizeConfigForComparison(desiredConfig)
	actualConfig = testlib.NormalizeConfigForComparison(actualConfig)

	desiredJson, err := desiredConfig.ToJSON()
	if err != nil {
		return "", "", fmt.Errorf("failed to setup example, error json encoding: %v", err)
	}
	actualJson, err := actualConfig.ToJSON()
	if err != nil {
		return "", "", fmt.Errorf("failed to setup example, error json encoding: %v", err)
	}

	return string(desiredJson), string(actualJson), nil
}
