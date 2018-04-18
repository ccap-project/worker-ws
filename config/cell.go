/*
 *
 * Copyright (c) 2016, 2017, 2018 Alexandre Biancalana <ale@biancalanas.net>.
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *     * Redistributions of source code must retain the above copyright
 *       notice, this list of conditions and the following disclaimer.
 *     * Redistributions in binary form must reproduce the above copyright
 *       notice, this list of conditions and the following disclaimer in the
 *       documentation and/or other materials provided with the distribution.
 *     * Neither the name of the <organization> nor the
 *       names of its contributors may be used to endorse or promote products
 *       derived from this software without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
 * ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
 * WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL <COPYRIGHT HOLDER> BE LIABLE FOR ANY
 * DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
 * (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
 * LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
 * ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
 * SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 */

package config

import (
	"encoding/json"
	"fmt"
	"io"
)

type RepoEnv struct {
	Name string
	Dir  string
	Url  string
	Env  []string
}

type CustomerEnv struct {
	Ansible   *RepoEnv
	Terraform *RepoEnv
}

/*
 * Bellow json received data
 */
type Provider struct {
	Name       string `json:"name"`
	DomainName string `json:"domain_name"`
	Region     string `json:"region"`
	TenantName string `json:"tenant_name"`
	Type       string `json:"type"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	AuthUrl    string `json:"auth_url"`
}

type KeyPair struct {
	Name      string `json:"name"`
	PublicKey string `json:"public_key"`
}

type Param struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Host struct {
	Name    string   `json:"name"`
	Options []*Param `json:"options"`
}

type File struct {
	Key      string `json:"key"`
	Filename string `json:"filename"`
	DontCopy string `json:"dont_copy"`
}

type Role struct {
	Name    string `json:"name"`
	Source  string `json:"url"`
	Version string `json:"version"`
	Files   []*File
	Params  []*Param `json:"params"`
}

type Hostgroup struct {
	Name              string      `json:"name"`
	Flavor            string      `json:"flavor"`
	Image             string      `json:"image"`
	KeyPair           string      `json:"key_pair"`
	Count             json.Number `json:"count,Number"`
	Network           string      `json:"network"`
	NetworkUUIDByName string      `json:"network_uuid_by_name"`
	Username          string      `json:"username"`
	Component         string      `json:"component"`
	BootstrapCommand  string      `json:"bootstrap_command"`
	Order             json.Number `json:"order"`
	Roles             []*Role     `json:"roles"`
	Securitygroups    []string    `json:"securitygroups"`
}

type Hostgroups []*Hostgroup

func (h Hostgroups) Len() int      { return len(h) }
func (h Hostgroups) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

type HostgroupByName struct{ Hostgroups }

func (s HostgroupByName) Less(i, j int) bool {
	return s.Hostgroups[i].BootstrapCommand < s.Hostgroups[j].BootstrapCommand
}

type HostgroupByOrder struct{ Hostgroups }

func (s HostgroupByOrder) Less(i, j int) bool {
	return s.Hostgroups[i].Order < s.Hostgroups[j].Order
}

type Loadbalancer struct {
	Algorithm             *string `json:"algorithm"`
	ConnectionDrain       string  `json:"connection_drain,omitempty"`
	ConnectionIDLETimeout int64   `json:"connection_idle_timeout,omitempty"`
	Members               string  `json:"members"`
	Name                  *string `json:"name"`
	Port                  *int64  `json:"port"`
	Protocol              *string `json:"protocol"`
	Type                  string  `json:"type,omitempty"`
}

type Network struct {
	Name       string `json:"name"`
	AdminState string `json:"admin_state"`
}

type Router struct {
	Name       string `json:"name"`
	AdminState string `json:"admin_state"`
}

type RouterInterface struct {
	Name   string `json:"name"`
	Router string `json:"router"`
	Subnet string `json:"subnet"`
}

type Securitygroup struct {
	Name  string                `json:"name"`
	Rules []*SecuritygroupRules `json:"rules"`
}

type SecuritygroupRules struct {
	DestinationAddr          string `json:"destination_addr,omitempty"`
	DestinationPort          string `json:"destination_port,omitempty"`
	DestinationSecuritygroup string `json:"destination_securitygroup,omitempty"`
	Ethertype                string `json:"ethertype,omitempty"`
	Proto                    string `json:"proto,omitempty"`
	SourceAddr               string `json:"source_addr,omitempty"`
	SourcePort               string `json:"source_port,omitempty"`
	SourceSecuritygroup      string `json:"source_securitygroup,omitempty"`
}

type Subnet struct {
	Name       string `json:"name"`
	Cidr       string `json:"cidr"`
	Network    string `json:"network"`
	IPVersion  string `json:"ip_version"`
	AdminState string `json:"admin_state"`
}

type Cell struct {
	CustomerName      string `json:"customer_name"`
	Name              string `json:"name"`
	Environment       CustomerEnv
	Provider          *Provider          `json:"provider"`
	KeyPair           *KeyPair           `json:"keypair"`
	Hosts             []*Host            `json:"hosts"`
	Hostgroups        []*Hostgroup       `json:"hostgroups"`
	Loadbalancers     []*Loadbalancer    `json:"loadbalancers"`
	Networks          []*Network         `json:"networks"`
	Subnets           []*Subnet          `json:"subnets"`
	Routers           []*Router          `json:"routers"`
	RoutersInterfaces []*RouterInterface `json:"routers_interfaces"`
	Securitygroups    []*Securitygroup   `json:"securitygroups"`
}

func DecodeJson(r io.Reader) (*Cell, error) {

	var cell Cell

	decoder := json.NewDecoder(r)

	if err := decoder.Decode(&cell); err != nil {
		return nil, fmt.Errorf("Can't decode request, %s", err)
	}

	return &cell, nil
}
