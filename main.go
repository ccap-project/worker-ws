/*
 *
 * Copyright (c) 2016, 2017, 2018 Alexandre Biancalana <ale@biancalanas.net>.
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *     * Redistributions of source code must retain the above copyright
 *       notice, this list of conditions and the following disclaimer.
 *     * Redistributions in binary form must reproduce the above copyright
 *       notice, this list of conditions and the following disclaimer in the
 *       documentation and/or other materials provided with the distribution.
 *     * Neither the name of the <organization> nor the
 *       names of its contributors may be used to endorse or promote products
 *       derived from this software without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
 * ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
 * WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL <COPYRIGHT HOLDER> BE LIABLE FOR ANY
 * DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
 * (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
 * LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
 * ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
 * SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 */

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
