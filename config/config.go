package config

import (
  "encoding/json"
  "fmt"
  "log"
  "os"
)

type CellMap map[string]interface{}

//type CellMap struct {
//  Resource map[string]interface{}
//}


type Provider struct {
  Name        string
  Region      string
}

type Hostgroup struct {
  Flavor       string
  counter      int
}

type Config struct {
  Provider *Provider
  Hostgroup []*Hostgroup
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
