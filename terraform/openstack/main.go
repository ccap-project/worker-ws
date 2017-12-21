package openstack

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"worker-ws/config"
	"worker-ws/utils"
	"worker-ws/terraform/common"
)

type Openstack struct{}

const provider_resource_tmpl = `provider "openstack" {
  user_name  = "{{.Username}}"
  tenant_name = "{{.TenantName}}"
  password  = "{{.Password}}"
  auth_url  = "{{.AuthUrl}}"
	{{if ne .DomainName "" }}domain_name = "{{.DomainName}}"{{end -}}
}
`

func (o *Openstack) Apply(system *config.SystemConfig, cell *config.Cell) (*[]byte, error) {
	return terraformcommon.Apply(system, cell)
}

func (o *Openstack) Validate(system *config.SystemConfig, cell *config.Cell) error {
	return terraformcommon.Validate(system, cell)
}

func (o *Openstack) Serialize(system *config.SystemConfig, cell *config.Cell) error {

	var tf bytes.Buffer
	terraformSite := fmt.Sprintf("%s/%s", cell.Environment.Terraform.Dir, system.Files.TerraformSite)

	provider, err := provider(cell)
	if err != nil {
		return (err)
	}

	router, err := router(cell)
	if err != nil {
		return (err)
	}

	router_interface, err := router_interface(cell)
	if err != nil {
		return (err)
	}

	network, err := network(cell)
	if err != nil {
		return (err)
	}

	subnet, err := subnet(cell)
	if err != nil {
		return (err)
	}

	keypair, err := keypair(cell)
	if err != nil {
		return (err)
	}

	instance, err := instance(cell)
	if err != nil {
		return (err)
	}

	tf.Write(provider.Bytes())
	tf.Write(router.Bytes())
	tf.Write(router_interface.Bytes())
	tf.Write(network.Bytes())
	tf.Write(subnet.Bytes())
	tf.Write(keypair.Bytes())
	tf.Write(instance.Bytes())

	ioutil.WriteFile(terraformSite, tf.Bytes(), 0644)

	return (nil)
}

func provider(cell *config.Cell) (*bytes.Buffer, error) {

	p, err := utils.Template(provider_resource_tmpl, cell.Provider)
	if err != nil {
		return p, err
	}

	return p, nil
}
