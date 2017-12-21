package main

import "worker-ws/config"
import "worker-ws/webservice"

func main() {

	//var err error
	//var terraform terraform.Serialize

	SystemConfig := config.ReadFile("etc/system.conf")

	SystemConfig.Log.Debugf("GitlabUrl(%s) GitlabToken(%s)", SystemConfig.Gitlab.Url, SystemConfig.Gitlab.Token)

	webservice.Start(SystemConfig)

	/*
		terraformSerializer := reflect.ValueOf(&terraform).MethodByName(SystemConfig.Provider.Name)

		if !terraformSerializer.IsValid() {
			fmt.Printf("Terraform serializer for provider(%s) is not supported !\n", SystemConfig.Provider.Name)
			os.Exit(-1)
		}

		err = terraformSerializer.Call([]reflect.Value{reflect.ValueOf(SystemConfig)})
		if err != nil {
			fmt.Println("Failure serializing Terraform %s file, %v", SystemConfig.Provider.Name, err)
			os.Exit(-1)
		}

		err = ansible.Serializer(SystemConfig)
		if err != nil {
			fmt.Println("Failure serializing Ansible Openstack file, ", err)
		}
	*/
}
