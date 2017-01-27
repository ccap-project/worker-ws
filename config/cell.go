package config

import (
	"encoding/json"
	"fmt"
	"io"
)

type Provider struct {
	Name       string `json:"name"`
	Region     string `json:"region"`
	Tenantname string `json:"tenantname"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	AuthUrl    string `json:"auth_url"`
}

type Host struct {
	Name    string              `json:"name"`
	Options []map[string]string `json:"options"`
}

type Role struct {
	Name    string `json:"name"`
	Source  string `json:"src"`
	Version string `json:"version"`
}

type Hostgroup struct {
	Name    string              `json:"name"`
	Flavor  string              `json:"flavor"`
	Image   string              `json:"image"`
	Count   string              `json:"count"`
	Network string              `json:"network"`
	Vars    []map[string]string `json:"vars"`
	Roles   []*Role             `json:"roles"`
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
	CustomerName      string             `json:"customer_name"`
	Name              string             `json:"name"`
	Provider          *Provider          `json:"provider"`
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
