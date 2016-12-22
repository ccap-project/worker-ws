package openstack

import "../../config/"
import "../../utils"
//import "fmt"

const network_resource_tmpl = `
resource "openstack_networking_network_v2" "{{.Name}}" {
  name = "{{.Name}}"
{{if ne .AdminState "" }}  admin_state_up = {{if eq .AdminState "up"}}"true" {{else}} "false"{{end}}{{end -}}
}
`
const router_resource_tmpl = `
resource "openstack_networking_router_v2" "{{.Name}}" {
  name = "{{.Name}}"
{{if ne .AdminState "" }}  admin_state_up = {{if eq .AdminState "up"}}"true" {{else}} "false"{{end}}{{- end}}
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

func network(config *config.Config) ([]string) {

  var networks []string

  for _, net := range config.Networks {
    z := utils.Template(network_resource_tmpl, net)
    networks = append(networks, z)
  }

  return(networks)
}

func router(config *config.Config) ([]string) {

  var routers []string

  for _, router := range config.Routers {
    z := utils.Template(router_resource_tmpl, router)
    routers = append(routers, z)
  }

  return(routers)
}

func router_interface(config *config.Config) ([]string) {

  var routers_interfaces []string

  for _, router_interface := range config.RoutersInterfaces {
    z := utils.Template(router_interface_resource_tmpl, router_interface)
    routers_interfaces = append(routers_interfaces, z)
  }

  return(routers_interfaces)
}

func subnet(config *config.Config) ([]string) {

  var subnets []string

  for _, subnet := range config.Subnets {
    z := utils.Template(subnet_resource_tmpl, subnet)
    subnets = append(subnets, z)
  }

  return(subnets)
}
