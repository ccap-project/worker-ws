package config

import (
  "encoding/json"
  "fmt"
  "log"
  "os"
)

type Provider struct {
  Name        string `json: "name"`
  Region      string `json: "region"`
}

type Hostgroup struct {
  Name          string `json: "name"`
  Flavor        string `json: "flavor"`
  Image         string `json: "image"`
  Count         string `json: "count"`
}

type Config struct {
  Provider *Provider      `json: "provider"`
  Hostgroups []*Hostgroup  `json: "hostgroups"`
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
