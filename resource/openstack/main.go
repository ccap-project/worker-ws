package openstack

import "fmt"
import "../../config/"
import "../../utils"

const provider_resource_tmpl = `
provider "openstack" {
  user_name  = "{{.Username}}"
  tenant_name = "{{.Tenantname}}"
  password  = "{{.Password}}"
  auth_url  = "{{.AuthUrl}}"
}
`

func Serializer (config *config.Config) {

  output := make([]string, 1)

  output = append(output, provider(config))
  output = append(output, router(config)...)
  output = append(output, router_interface(config)...)
  output = append(output, network(config)...)
  output = append(output, subnet(config)...)
  output = append(output, instance(config)...)

  for _,v := range output {
    fmt.Printf("%s\n", v)
  }
}

func provider (config *config.Config) (string) {

  z := utils.Template(provider_resource_tmpl, config.Provider)

  return(z)
}
