package openstack

import (
	"../../config"
	"../common"
)

func (o *Openstack) ReadState(cell *config.Cell, file string) error {

	state, err := terraformcommon.ReadState(file)

	if err != nil {
		return err
	}

	for m := range state.Modules {
		for _, rv := range state.Modules[m].Resources {

			switch rv.Type {
			case "openstack_compute_instance_v2":
				host := new(config.Host)
				option := new(config.Param)

				host.Name = rv.Primary.Attributes["name"]
				option.Name = "ansible_host"
				option.Value = rv.Primary.Attributes["access_ip_v4"]

				host.Options = append(host.Options, option)
				cell.Hosts = append(cell.Hosts, host)

			default:

			}
		}
	}

	return nil
}
