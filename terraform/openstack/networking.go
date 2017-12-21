package openstack

import "bytes"

import "worker-ws/config"
import "worker-ws/utils"

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
    i,err := utils.Template(router_interface_resource_tmpl, router_interface)
    if err != nil {
      return nil, err
    }
    routers_interfaces.Write(i.Bytes())
  }

  return &routers_interfaces, nil
}

func subnet(config *config.Cell) (*bytes.Buffer, error) {

  var subnets bytes.Buffer

  for _, subnet := range config.Subnets {
    s,err := utils.Template(subnet_resource_tmpl, subnet)
    if err != nil {
      return nil, err
    }
    subnets.Write(s.Bytes())
  }

  return &subnets, nil
}
