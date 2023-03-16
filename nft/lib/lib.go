//go:build cgo
// +build cgo

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
package lib

// #cgo CFLAGS: -g -Wall
// #cgo LDFLAGS: -lnftables
// #include <nftables/libnftables.h>
// #include <stdlib.h>
// #include <string.h>
import "C"
import (
	"fmt"
	"strings"
	"unsafe"

	"github.com/networkplumbing/go-nft/nft"
)

const (
	cmdList    = "list"
	cmdRuleset = "ruleset"
)

// ReadConfig loads the nftables configuration from the system and
// returns it as a nftables config structure.
// The system is expected to have the `nft` executable deployed and nftables enabled in the kernel.
func ReadConfig(filterCommands ...string) (*nft.Config, error) {
	whatToList := cmdRuleset
	if len(filterCommands) > 0 {
		whatToList = strings.Join(filterCommands, " ")
	}
	stdout, err := libNftablesRunCmd(fmt.Sprintf("%s %s", cmdList, whatToList))
	if err != nil {
		return nil, err
	}

	config := nft.NewConfig()
	if err := config.FromJSON(stdout); err != nil {
		return nil, fmt.Errorf("failed to list ruleset: %v", err)
	}

	return config, nil
}

// ApplyConfig applies the given nftables config on the system.
// The system is expected to have the `nft` executable deployed and nftables enabled in the kernel.
func ApplyConfig(c *nft.Config) error {
	data, err := c.ToJSON()
	if err != nil {
		return err
	}

	if _, err = libNftablesRunCmd(string(data)); err != nil {
		return err
	}

	return nil
}

func libNftablesRunCmd(cmd string) ([]byte, error) {
	nft := C.nft_ctx_new(C.NFT_CTX_DEFAULT)
	defer C.nft_ctx_free(nft)

	C.nft_ctx_output_set_flags(nft, C.NFT_CTX_OUTPUT_JSON)

	buf := C.CString(cmd)
	defer C.free(unsafe.Pointer(buf))

	rc := C.nft_ctx_buffer_output(nft)
	if rc != C.EXIT_SUCCESS {
		return nil, fmt.Errorf("failed enabling output buffering (rc=%d)", rc)
	}

	rc = C.nft_ctx_buffer_error(nft)
	if rc != C.EXIT_SUCCESS {
		return nil, fmt.Errorf("failed enabling error buffering (rc=%d)", rc)
	}

	rc = C.nft_run_cmd_from_buffer(nft, buf)
	if rc != C.EXIT_SUCCESS {
		errMsg := C.nft_ctx_get_error_buffer(nft)
		return nil, fmt.Errorf("failed running cmd (rc=%d): %s", rc, C.GoString(errMsg))
	}

	config := C.nft_ctx_get_output_buffer(nft)
	configLen := C.int(C.strlen(config))
	return C.GoBytes(unsafe.Pointer(config), configLen), nil
}
