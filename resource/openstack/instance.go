package openstack

import "fmt"
//import "text/template"
//import "os"

import "../../config/"
import "../../utils"

const instance_resource_tmpl = `
resource "openstack_compute_instance_v2" "{{.Name}}" {
  count = {{.Count}}
  name = "{{.Name}}-${count.index}",
  image_id = "{{.Image}}",
  flavor_name = "{{.Flavor}}",
  #key_pair = "${var.keypair}",
  #floating_ip = "${openstack_compute_floatingip_v2.tf-ds-float-ip.address}",
  #security_groups = ["default"]
  #network {
  #  "uuid" = "${openstack_networking_network_v2.tf-ds-network.id}"
  #}
}
`

func instance (config *config.Config) {

  for i, h := range config.Hostgroups {

    z := utils.Template(instance_resource_tmpl, h)
    fmt.Printf("\n%d %s\n", i, z)
    //fmt.Printf("%d %+v", i, h)

    //t := template.New("instance")
    //t,_ = t.Parse(instance_resource_tmpl)
    //t.Execute(os.Stdout, h)
  }
  //fmt.Println("Here 2 !")
}
