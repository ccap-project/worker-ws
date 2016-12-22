package openstack

//import "fmt"
//import "text/template"
//import "os"

import "../../config/"
import "../../utils"

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
    "uuid" = "${openstack_networking_network_v2.{{.Network}}.id}"
  }
}
`

func instance (config *config.Config) ([]string) {

  var instances []string

  for _, h := range config.Hostgroups {
    z := utils.Template(instance_resource_tmpl, h)
    instances = append(instances, z)
  }

  return(instances)
}
