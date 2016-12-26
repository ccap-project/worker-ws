package openstack

import "fmt"
import "../../config/"

const provider_resource_tmpl = `
provider "openstack" {
    user_name  = "{{.Username}}"
    tenant_name = "{{.Tenant_name}}"
    password  = "{{.Password}}"
    auth_url  = "{{.Auth_url}}"
}
`

func Serializer (config *config.Config) {

  instance(config)
  //config.Hostgroups.Marshall()
  fmt.Println("Here !")
  //fmt.Println(config.HostGroups)
}

func provider (config *config.Config) {

}
