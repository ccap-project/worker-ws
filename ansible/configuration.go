package ansible

import (
	"bytes"

	"../config"
	"../utils"
)

const cfg_tmpl = `[defaults]
inventory = {{.Environment.Ansible.Dir}}hosts
roles_path = {{.Environment.Ansible.Dir}}roles
log_path = {{.Environment.Ansible.Dir}}log
host_key_checking = False
[ssh_connection]
control_path = %(directory)s/%%C
`

func configuration(cell *config.Cell) (*bytes.Buffer, error) {

	var cfg bytes.Buffer

	c, err := utils.Template(cfg_tmpl, cell)
	if err != nil {
		return nil, err
	}

	cfg.Write(c.Bytes())

	return &cfg, nil
}
