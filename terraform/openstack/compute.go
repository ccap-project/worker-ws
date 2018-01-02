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

import (
	"bytes"

	"worker-ws/config"
	"worker-ws/utils"
)

const keypair_resource_tmpl = `
resource "openstack_compute_keypair_v2" "{{.Name}}" {
  name = "{{.Name}}"
  public_key = "{{.PublicKey}}"
}
`

const instance_resource_tmpl = `
resource "openstack_compute_instance_v2" "{{.Name}}" {
  count = "{{.Count}}"
  name = "${format("{{.Name}}%d", count.index + 1)}"
  image_id = "{{.Image}}"
  flavor_name = "{{.Flavor}}"
  #key_pair = "${var.keypair}",
  #floating_ip = "${openstack_compute_floatingip_v2.tf-ds-float-ip.address}"
  #security_groups = ["default"]
  network {
    {{if ne .NetworkUUIDByName ""}}uuid = "${openstack_networking_network_v2.{{.Network}}.id}"{{else}}name = "{{.Network}}"{{end}}
  }
  {{if ne .KeyPair "" }}key_pair = "{{.KeyPair}}"{{end}}
  provisioner "remote-exec" {
    inline = [ "ls" ]
    connection {
      agent = "false"
      type  = "ssh"
      user  = "{{.Username}}"
      private_key = "${file("/Users/ale/.ssh/id_rsa")}"
    }
  }
}
`

func instance(config *config.Cell) (*bytes.Buffer, error) {

	var instances bytes.Buffer

	for _, h := range config.Hostgroups {
		h.KeyPair = config.KeyPair.Name
		i, err := utils.Template(instance_resource_tmpl, h)
		if err != nil {
			return nil, err
		}

		instances.Write(i.Bytes())
	}

	return &instances, nil
}

func keypair(config *config.Cell) (*bytes.Buffer, error) {

	var keypair bytes.Buffer

	k, err := utils.Template(keypair_resource_tmpl, config.KeyPair)
	if err != nil {
		return nil, err
	}

	keypair.Write(k.Bytes())

	return &keypair, nil
}
