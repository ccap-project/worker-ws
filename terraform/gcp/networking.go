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

package gcp

import (
	"bytes"
	"sort"
	"worker-ws/config"
	"worker-ws/utils"
)

const loadbalancer_http_resource_tmpl = `
#
# HTTP Load Balancer Configuration
#

resource "google_compute_health_check" "{{.}}" {
  name = "{{.}}"
  request_path = "/"
}

resource "google_compute_backend_service" "{{.Name}}" {
  name        = "{{.Name}}"
  port_name   = "{{.Name}}"
  protocol    = "{{.Protocol}}"

{{range $idx, $v := .Members}}
  backend {
    group = "${google_compute_instance_group.{{.}}.self_link}"
  }
{{end}}
  health_checks = ["${google_compute_http_health_check.{{.}}.self_link}"]
}
`

const loadbalancer_resource_tmpl = `
#
# Load Balancer Configuration
#

resource "google_compute_health_check" "{{.}}" {
  name = "{{.}}"

  tcp_health_check {
    port = "{{.Port}}"
  }
}

resource "google_compute_backend_service" "{{.Name}}" {
  name        = "{{.Name}}"
  port_name   = "{{.Name}}"
  protocol    = "{{.Protocol}}"

{{range $idx, $v := .Members}}
  backend {
    group = "${google_compute_instance_group.{{.}}.self_link}"
  }
{{end}}
  health_checks = ["${google_compute_http_health_check.{{.}}.self_link}"]
}
`

const loadbalancer_autoscaling_attachment_resource_tmpl = `
{{- $LbName := .Name -}}
{{- range .Members}}
resource "aws_autoscaling_attachment" "{{.}}" {
  autoscaling_group_name = "${aws_autoscaling_group.{{.}}.id}"
  alb_target_group_arn   = "${aws_lb_target_group.{{$LbName}}.arn}"
}{{end}}
`

const loadbalancer_instance_group_resource_tmpl = `
{{- $LbName := .Name -}}
{{- $LbPort := .Port }}
{{- range .Members}}
resource "google_compute_instance_group" "{{.}}" {
  name  = "{{.}}"
  zone  = "${google_compute_subnetwork.{{ (index .Network 0)}}.region}"
  instances = "${join(google_compute_instance.{{.}}.*.self_link, ",")}"

  named_port {
    name = "$LbName"
    port = "{{$LbPort}}"
  }
}{{end}}
`

const subnet_resource_tmpl = `
#
# Subnet Configuration
#
resource "google_compute_subnetwork" "{{.Name}}" {
  name       = "{{.Name}}"
  network     = "${google_compute_network.{{.Router}}.self_link}"
  ip_cidr_range = "{{.Cidr}}"
  region = "${var.current_region}"
}
`

const secgroup_resource_tmpl = `
#
# SecurityGroup Configuration
#

resource "aws_security_group" "{{.Name}}" {
  name = "{{.Name}}"
  vpc_id = "${aws_vpc.{{.Router}}.id}"
}
{{if .Rules -}}
{{- $SecgroupName := .Name -}}
{{- range .Rules}}
resource "aws_security_group_rule" "{{.SourceSecuritygroup}}_to_{{$SecgroupName}}_on_{{.DestinationPort}}" {
  type         = "ingress"
  protocol     = "{{.Proto}}"
  from_port    = {{.DestinationPort}}
  to_port      = {{.DestinationPort}}
{{if eq .DestinationSecuritygroup ""}}  cidr_blocks  = "{{.DestinationAddr}}"{{else}}  source_security_group_id   = "${aws_security_group.{{.SourceSecuritygroup}}.id}"{{end}}
  security_group_id = "${aws_security_group.{{$SecgroupName}}.id}"
}{{end}}{{end}}
`

const vpc_resource_tmpl = `
#
# Network Configuration
#
resource "google_compute_network" "{{.Name}}" {
  name       = "{{.Name}}"
  auto_create_subnetworks = "false"
}
`

func loadbalancer(config *config.Cell) (*bytes.Buffer, error) {

	var loadbalancer bytes.Buffer

	for _, lb := range config.Loadbalancers {
		sort.Strings(lb.Members)
		n, err := utils.Template(loadbalancer_resource_tmpl, lb)
		if err != nil {
			return nil, err
		}
		loadbalancer.Write(n.Bytes())

		// Check if autoscale is enable on lb members
		for _, h := range config.Hostgroups {

			exists, _ := utils.Grep(lb.Members, h.Name)
			if exists {
				// Autoscale enabled
				if h.MaxSize > 0 {
					n, err := utils.Template(loadbalancer_autoscaling_attachment_resource_tmpl, lb)
					if err != nil {
						return nil, err
					}
					loadbalancer.Write(n.Bytes())

				} else {
					n, err := utils.Template(loadbalancer_instance_group_resource_tmpl, lb)
					if err != nil {
						return nil, err
					}
					loadbalancer.Write(n.Bytes())
				}
			}
		}

	}

	return &loadbalancer, nil
}

func network(config *config.Cell) (*bytes.Buffer, error) {

	var networks bytes.Buffer

	for _, _net := range config.Routers {

		r, err := utils.Template(vpc_resource_tmpl, _net)
		if err != nil {
			return nil, err
		}
		networks.Write(r.Bytes())
	}
	return &networks, nil
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

	for _, net := range config.Networks {
		s, err := utils.Template(subnet_resource_tmpl, net)
		if err != nil {
			return nil, err
		}
		subnets.Write(s.Bytes())
	}

	return &subnets, nil
}
