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

import (
	"encoding/json"
	"fmt"
)

type Root struct {
	Nftables []Nftable `json:"nftables"`
}

type Objects struct {
	Table *Table `json:"table,omitempty"`
	Chain *Chain `json:"chain,omitempty"`
	Rule  *Rule  `json:"rule,omitempty"`
}

type Nftable struct {
	Table *Table `json:"table,omitempty"`
	Chain *Chain `json:"chain,omitempty"`
	Rule  *Rule  `json:"rule,omitempty"`

	Add    *Objects `json:"add,omitempty"`
	Delete *Objects `json:"delete,omitempty"`
	Flush  *Objects `json:"flush,omitempty"`

	Metainfo *Metainfo `json:"metainfo,omitempty"`
}

type Metainfo struct {
	Version           string `json:"version"`
	ReleaseName       string `json:"release_name"`
	JsonSchemaVersion int    `json:"json_schema_version"`
}

// Table Address Families
const (
	FamilyIP     = "ip"     // IPv4 address AddressFamily.
	FamilyIP6    = "ip6"    // IPv6 address AddressFamily.
	FamilyINET   = "inet"   // Internet (IPv4/IPv6) address AddressFamily.
	FamilyARP    = "arp"    // ARP address AddressFamily, handling IPv4 ARP packets.
	FamilyBRIDGE = "bridge" // Bridge address AddressFamily, handling packets which traverse a bridge device.
	FamilyNETDEV = "netdev" // Netdev address AddressFamily, handling packets from ingress.
)

type Table struct {
	Family string `json:"family"`
	Name   string `json:"name"`
}

// Chain Types
const (
	TypeFilter = "filter"
	TypeNAT    = "nat"
	TypeRoute  = "route"
)

// Chain Hooks
const (
	HookPreRouting  = "prerouting"
	HookInput       = "input"
	HookOutput      = "output"
	HookForward     = "forward"
	HookPostRouting = "postrouting"
	HookIngress     = "ingress"
)

// Chain Policies
const (
	PolicyAccept = "accept"
	PolicyDrop   = "drop"
)

type Chain struct {
	Family string `json:"family"`
	Table  string `json:"table"`
	Name   string `json:"name"`
	Type   string `json:"type,omitempty"`
	Hook   string `json:"hook,omitempty"`
	Prio   *int   `json:"prio,omitempty"`
	Policy string `json:"policy,omitempty"`
}

type Rule struct {
	Family  string      `json:"family"`
	Table   string      `json:"table"`
	Chain   string      `json:"chain"`
	Expr    []Statement `json:"expr,omitempty"`
	Handle  *int        `json:"handle,omitempty"`
	Comment string      `json:"comment,omitempty"`
}

type Statement struct {
	Match *Match `json:"match,omitempty"`
	Verdict
}

type Verdict struct {
	SimpleVerdict
	Jump *ToTarget `json:"jump,omitempty"`
	Goto *ToTarget `json:"goto,omitempty"`
}

type SimpleVerdict struct {
	Accept   bool `json:"-"`
	Continue bool `json:"-"`
	Drop     bool `json:"-"`
	Return   bool `json:"-"`
}

type ToTarget struct {
	Target string `json:"target"`
}

type Match struct {
	Op    string     `json:"op"`
	Left  Expression `json:"left"`
	Right Expression `json:"right"`
}

type Expression struct {
	String  *string         `json:"-"`
	Bool    *bool           `json:"-"`
	Int     *int            `json:"-"`
	Payload *Payload        `json:"payload,omitempty"`
	RowData json.RawMessage `json:"-"`
}

type Payload struct {
	Protocol string `json:"protocol"`
	Field    string `json:"field"`
}

// Verdict Operations
const (
	VerdictAccept   = "accept"
	VerdictContinue = "continue"
	VerdictDrop     = "drop"
	VerdictReturn   = "return"
)

// Match Operators
const (
	OperAND = "&"  // Binary AND
	OperOR  = "|"  // Binary OR
	OperXOR = "^"  // Binary XOR
	OperLSH = "<<" // Left shift
	OperRSH = ">>" // Right shift
	OperEQ  = "==" // Equal
	OperNEQ = "!=" // Not equal
	OperLS  = "<"  // Less than
	OperGR  = ">"  // Greater than
	OperLSE = "<=" // Less than or equal to
	OperGRE = ">=" // Greater than or equal to
	OperIN  = "in" // Perform a lookup, i.e. test if bits on RHS are contained in LHS value
)

// Payload Expressions
const (
	PayloadKey = "payload"
	// Ethernet
	PayloadProtocolEther   = "ether"
	PayloadFieldEtherDAddr = "daddr"
	PayloadFieldEtherSAddr = "saddr"
	PayloadFieldEtherType  = "type"

	// IP (common)
	PayloadFieldIPVer   = "version"
	PayloadFieldIPDscp  = "dscp"
	PayloadFieldIPEcn   = "ecn"
	PayloadFieldIPLen   = "length"
	PayloadFieldIPSAddr = "saddr"
	PayloadFieldIPDAddr = "daddr"

	// IPv4
	PayloadProtocolIP4      = "ip"
	PayloadFieldIP4HdrLen   = "hdrlength"
	PayloadFieldIP4Id       = "id"
	PayloadFieldIP4FragOff  = "frag-off"
	PayloadFieldIP4Ttl      = "ttl"
	PayloadFieldIP4Protocol = "protocol"
	PayloadFieldIP4Chksum   = "checksum"

	// IPv6
	PayloadProtocolIP6       = "ip6"
	PayloadFieldIP6FlowLabel = "flowlabel"
	PayloadFieldIP6NextHdr   = "nexthdr"
	PayloadFieldIP6HopLimit  = "hoplimit"
)

func (s Statement) MarshalJSON() ([]byte, error) {
	type _Statement Statement
	statement := _Statement(s)

	// Convert to a dynamic structure
	data, err := json.Marshal(statement)
	if err != nil {
		return nil, err
	}
	dynamicStructure := make(map[string]json.RawMessage)
	if err := json.Unmarshal(data, &dynamicStructure); err != nil {
		return nil, err
	}

	switch {
	case s.Accept:
		dynamicStructure[VerdictAccept] = nil
	case s.Continue:
		dynamicStructure[VerdictContinue] = nil
	case s.Drop:
		dynamicStructure[VerdictDrop] = nil
	case s.Return:
		dynamicStructure[VerdictReturn] = nil
	}

	data, err = json.Marshal(dynamicStructure)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Statement) UnmarshalJSON(data []byte) error {
	type _Statement Statement
	statement := _Statement{}

	if err := json.Unmarshal(data, &statement); err != nil {
		return err
	}
	*s = Statement(statement)

	dynamicStructure := make(map[string]json.RawMessage)
	if err := json.Unmarshal(data, &dynamicStructure); err != nil {
		return err
	}
	_, s.Accept = dynamicStructure[VerdictAccept]
	_, s.Continue = dynamicStructure[VerdictContinue]
	_, s.Drop = dynamicStructure[VerdictDrop]
	_, s.Return = dynamicStructure[VerdictReturn]

	return nil
}

func (e Expression) MarshalJSON() ([]byte, error) {
	var dynamicStruct interface{}

	switch {
	case e.RowData != nil:
		return e.RowData, nil
	case e.String != nil:
		dynamicStruct = *e.String
	case e.Int != nil:
		dynamicStruct = *e.Int
	case e.Bool != nil:
		dynamicStruct = *e.Bool
	default:
		type _Expression Expression
		dynamicStruct = _Expression(e)
	}

	return json.Marshal(dynamicStruct)
}

func (e *Expression) UnmarshalJSON(data []byte) error {
	var dynamicStruct interface{}
	if err := json.Unmarshal(data, &dynamicStruct); err != nil {
		return err
	}

	switch dynamicStruct.(type) {
	case string:
		d := dynamicStruct.(string)
		e.String = &d
	case int:
		d := dynamicStruct.(int)
		e.Int = &d
	case bool:
		d := dynamicStruct.(bool)
		e.Bool = &d
	case map[string]interface{}:
		type _Expression Expression
		expression := _Expression(*e)
		if err := json.Unmarshal(data, &expression); err != nil {
			return err
		}
		*e = Expression(expression)
	default:
		return fmt.Errorf("unsupported field type in expression")
	}

	if e.String == nil && e.Int == nil && e.Bool == nil && e.Payload == nil {
		e.RowData = data
	}

	return nil
}
