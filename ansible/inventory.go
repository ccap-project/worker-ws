package ansible

import "bytes"
import "fmt"

import "../config/"

func hosts(config *config.Cell) *bytes.Buffer {

	var hosts bytes.Buffer

	// XXX: Move this loop to template
	for _, host := range config.Hosts {

		fmt.Fprintf(&hosts, "%s", host.Name)

		for _, opt := range host.Options {
			for k, v := range opt {
				fmt.Fprintf(&hosts, " %s=%s", k, v)
			}
		}
		fmt.Fprintf(&hosts, "\n")
	}
	fmt.Fprintf(&hosts, "\n")

	return (&hosts)
}

func hostgroups(config *config.Cell) *bytes.Buffer {

	var hostgroups bytes.Buffer

	// XXX: Move this loop to template
	for _, hostgroup := range config.Hostgroups {
		fmt.Fprintf(&hostgroups, "[%s]\n", hostgroup.Name)
		fmt.Fprintf(&hostgroups, "%s[1:%s]\n\n", hostgroup.Name, hostgroup.Count)
	}

	return (&hostgroups)
}

func group_vars(config *config.Cell) *bytes.Buffer {

	var group_vars bytes.Buffer

	// XXX: Move this loop to template
	for _, hostgroup := range config.Hostgroups {
		fmt.Fprintf(&group_vars, "[%s:vars]\n", hostgroup.Name)

		if len(hostgroup.Username) > 0 {
			fmt.Fprintf(&group_vars, "ansible_ssh_user=%s\n", hostgroup.Username)
		}

		for _, vars := range hostgroup.Vars {
			for k, v := range vars {
				fmt.Fprintf(&group_vars, "%s=%s\n", k, v)
			}
		}
		fmt.Fprintf(&group_vars, "\n")
	}

	return (&group_vars)
}

func GetInventoryFilename(config *config.SystemConfig, cell *config.Cell) string {
	return fmt.Sprintf("%s%s", cell.Environment.Ansible.Dir, config.Files.AnsibleHosts)
}
