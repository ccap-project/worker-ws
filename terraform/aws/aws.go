package aws

import "fmt"
import "../../config"

//type hostgroup config.Hostgroup

func Serializer(config *config.Config) {

	instance(config)
	//config.Hostgroups.Marshall()
	fmt.Println("Here !")
	//fmt.Println(config.HostGroups)
}
