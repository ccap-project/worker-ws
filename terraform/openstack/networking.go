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

package openstack

import "bytes"

import "worker-ws/config"
import "worker-ws/utils"

const loadbalancer_resource_tmpl = `
resource "openstack_lb_monitor_v1" "monitor_{{.Name}}" {
  type = "{{.Protocol}}"
  delay = 30
  timeout = 5
  max_retries = 3
  admin_state_up = "true"
}

resource "openstack_lb_pool_v1" "pool_{{.Name}}" {
  name = "pool_{{.Name}}"
  protocol = "{{.Protocol}}"
  #subnet_id = "${openstack_networking_subnet_v2.subnet-web.id}"
  lb_method = "{{.Algorithm}}"
  monitor_ids = ["${openstack_lb_monitor_v1.monitor_{{.Name}}.id}"]
}

resource "openstack_lb_member_v1" "members_{{.Name}}" {
  count   = "${length(join(",", openstack_compute_instance_v2.{{.Members}}.*.id))}"
  pool_id = "${openstack_lb_pool_v1.pool_{{.Name}}.id}"
  address = "${element(openstack_compute_instance_v2.{{.Members}}.*.network.0.fixed_ip_v4, count.index)}"
  port    = {{.Port}}
}

resource "openstack_lb_vip_v1" "vip_{{.Name}}" {
  name        = "vip_1"
  #subnet_id   = "${openstack_networking_subnet_v2.subnet-web.id}"
  protocol    = "{{.Protocol}}"
  port        = {{.Port}}
  pool_id     = "${openstack_lb_pool_v1.pool_{{.Name}}.id}"
  #floating_ip = "${openstack_compute_floatingip_v2.web_public_ip.address}"
}
`

const network_resource_tmpl = `
resource "openstack_networking_network_v2" "{{.Name}}" {
  name = "{{.Name}}"
{{if ne .AdminState "" }}  admin_state_up = {{if eq .AdminState "up"}}"true" {{else}} "false"{{end}}{{end -}}
}
`
const router_resource_tmpl = `
resource "openstack_networking_router_v2" "{{.Name}}" {
  name = "{{.Name}}"
{{if ne .AdminState "" }}  admin_state_up = {{if eq .AdminState "up"}}"true" {{else}} "false"{{end}}{{end -}}
}
`

const router_interface_resource_tmpl = `
resource "openstack_networking_router_interface_v2" "{{.Name}}" {
  router_id = "${openstack_networking_router_v2.{{.Router}}.id}"
  subnet_id = "${openstack_networking_subnet_v2.{{.Subnet}}.id}"
}
`

const subnet_resource_tmpl = `
resource "openstack_networking_subnet_v2" "{{.Name}}" {
  name = "{{.Name}}"
  network_id = "${openstack_networking_network_v2.{{.Network}}.id}"
{{if ne .Cidr "" }}  cidr = "{{.Cidr}}"{{end -}}
{{if ne .IPVersion "" }}  ip_version = "{{.IPVersion}}" {{end}}
}
`

const secgroup_resource_tmpl = `
resource "openstack_networking_secgroup_v2" "{{.Name}}" {
  name = "{{.Name}}"
}
{{if .Rules -}}
{{- $SecgroupName := .Name -}}
{{- range .Rules}}
resource "openstack_networking_secgroup_rule_v2" "{{.SourceSecuritygroup}}_to_{{$SecgroupName}}_on_{{.DestinationPort}}" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "{{.Proto}}"
  port_range_min    = {{.DestinationPort}}
  port_range_max    = {{.DestinationPort}}
{{if eq .DestinationSecuritygroup ""}}  remote_ip_prefix  = "{{.DestinationAddr}}"{{else}}  remote_group_id   = "{{.SourceSecuritygroup}}"{{end}}
  security_group_id = "${openstack_networking_secgroup_v2.{{$SecgroupName}}.id}"
}{{end}}{{end}}
`

func loadbalancer(config *config.Cell) (*bytes.Buffer, error) {

	var loadbalancer bytes.Buffer

	for _, lb := range config.Loadbalancers {
		n, err := utils.Template(loadbalancer_resource_tmpl, lb)
		if err != nil {
			return nil, err
		}
		loadbalancer.Write(n.Bytes())
	}

	return &loadbalancer, nil
}

func network(config *config.Cell) (*bytes.Buffer, error) {

	var networks bytes.Buffer

	for _, net := range config.Networks {
		n, err := utils.Template(network_resource_tmpl, net)
		if err != nil {
			return nil, err
		}
		networks.Write(n.Bytes())
	}

	return &networks, nil
}

func router(config *config.Cell) (*bytes.Buffer, error) {

	var routers bytes.Buffer

	for _, router := range config.Routers {

		r, err := utils.Template(router_resource_tmpl, router)
		if err != nil {
			return nil, err
		}
		routers.Write(r.Bytes())
	}
	return &routers, nil
}

func router_interface(config *config.Cell) (*bytes.Buffer, error) {

	var routers_interfaces bytes.Buffer

	for _, router_interface := range config.RoutersInterfaces {
		i, err := utils.Template(router_interface_resource_tmpl, router_interface)
		if err != nil {
			return nil, err
		}
		routers_interfaces.Write(i.Bytes())
	}

	return &routers_interfaces, nil
}

func securitygroup(config *config.Cell) (*bytes.Buffer, error) {
	var securitygroups bytes.Buffer

	for _, secgroup := range config.Securitygroups {
		s, err := utils.Template(secgroup_resource_tmpl, secgroup)
		if err != nil {
			return nil, err
		}
		securitygroups.Write(s.Bytes())

	}
	return &securitygroups, nil
}

func subnet(config *config.Cell) (*bytes.Buffer, error) {

	var subnets bytes.Buffer

	for _, subnet := range config.Subnets {
		s, err := utils.Template(subnet_resource_tmpl, subnet)
		if err != nil {
			return nil, err
		}
		subnets.Write(s.Bytes())
	}

	return &subnets, nil
}
