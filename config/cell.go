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

type Host struct {
	Name    string              `json:"name"`
	Options []map[string]string `json:"options"`
}

type File struct {
	Key      string `json:"key"`
	Filename string `json:"filename"`
	DontCopy string `json:"dont_copy"`
}

type RoleParam struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Role struct {
	Name    string `json:"name"`
	Source  string `json:"url"`
	Version string `json:"version"`
	Files   []*File
	Params  []*RoleParam `json:"params"`
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
	Roles             []*Role     `json:"roles"`
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
	Networks          []*Network         `json:"networks"`
	Subnets           []*Subnet          `json:"subnets"`
	Routers           []*Router          `json:"routers"`
	RoutersInterfaces []*RouterInterface `json:"routers_interfaces"`
}

func DecodeJson(r io.Reader) (*Cell, error) {

	var cell Cell

	decoder := json.NewDecoder(r)

	if err := decoder.Decode(&cell); err != nil {
		return nil, fmt.Errorf("Can't decode request, %s", err)
	}

	return &cell, nil
}
