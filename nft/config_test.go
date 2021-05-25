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
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/eddev/go-nft/nft"
)

func TestDefineEmptyConfig(t *testing.T) {
	config := nft.NewConfig()

	expected := []byte(`{"nftables":[]}`)
	serializedConfig, err := config.MarshalJSON()
	assert.NoError(t, err)
	assert.Equal(t, string(expected), string(serializedConfig))
}
