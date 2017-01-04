package config

import (
  "encoding/json"
  "fmt"
  "log"
  "os"
)

type Provider struct {
  Name        string `json:"name"`
  Region      string `json:"region"`
  Tenantname  string `json:"tenantname"`
  Username    string `json:"username"`
  Password    string `json:"password"`
  AuthUrl     string `json:"auth_url"`
}

type Host struct {
  Name          string        `json:"name"`
  Options       []map[string]string `json:"options"`
}

type Hostgroup struct {
  Name          string `json:"name"`
  Flavor        string `json:"flavor"`
  Image         string `json:"image"`
  Count         string `json:"count"`
  Network       string `json:"network"`
  Vars          []map[string]string  `json:"vars"`
}

type Network struct {
  Name          string `json:"name"`
  AdminState    string `json:"admin_state"`
}

type Router struct {
  Name          string `json:"name"`
  AdminState    string `json:"admin_state"`
}

type RouterInterface struct {
  Name          string `json:"name"`
  Router        string `json:"router"`
  Subnet        string `json:"subnet"`
}

type Subnet struct {
  Name          string `json:"name"`
  Cidr          string `json:"cidr"`
  Network       string `json:"network"`
  IPVersion     string `json:"ip_version"`
  AdminState    string `json:"admin_state"`
}

type Config struct {
  Provider          *Provider           `json:"provider"`
  Hosts             []*Host            `json:"hosts"`
  Hostgroups        []*Hostgroup        `json:"hostgroups"`
  Networks          []*Network          `json:"networks"`
  Subnets           []*Subnet           `json:"subnets"`
  Routers           []*Router           `json:"routers"`
  RoutersInterfaces []*RouterInterface  `json:"routers_interfaces"`
}

//type Resource struct {
//  ResourceType ResourceType
//}

//type TerraForm struct {
//  Resource     Resource
//}

/*
func ReadFile(configFilePath string) *Configuration {

  file, err := os.Open(configFilePath)
  if err != nil {
    log.Fatalf("Can't open config file(%s), %s", configFilePath, err)
  }
  defer file.Close()

  decoder := json.NewDecoder(file)
  config := Configuration{GitlabCfg: GitlabCfg_t{TLSInsecureSkipVerify: true,
                                                  Group: "ansible-roles"}}

  err = decoder.Decode(&config)

  if err != nil {
    log.Fatalf("Can't decode config file(%s), %s", configFilePath, err)
  }

  config.Log = log.New(os.Stderr, "roles-ws: ", log.Llongfile)

  return &config
}
*/

func ReadJson(configFilePath string) *Config {

  var config Config

  file, err := os.Open(configFilePath)
  if err != nil {
    log.Fatalf("Can't open config file(%s), %s", configFilePath, err)
  }
  defer file.Close()

  decoder := json.NewDecoder(file)

  if err := decoder.Decode(&config); err != nil {
    log.Fatalf("Can't decode config file(%s), %s", configFilePath, err)
  }

  //config.Log = log.New(os.Stderr, "roles-ws: ", log.Llongfile)
  fmt.Printf("%+v", config)
  return &config
}
