package openstack

import "bytes"

//import "text/template"
//import "os"

import "../../config"
import "../../utils"

const keypair_resource_tmpl = `
resource "openstack_compute_keypair_v2" "{.KeyName}" {
  name = "{.KeyName}"
  public_key = "{.PublicKey}"
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

	k, err := utils.Template(instance_resource_tmpl, config.KeyPair)
	if err != nil {
		return nil, err
	}

	keypair.Write(k.Bytes())

	return &keypair, nil
}
