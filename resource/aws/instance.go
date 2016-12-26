package aws

import "fmt"
import "text/template"
import "os"

import "../../config/"

const instance_resource_tmpl = `
resource "aws_instance" "{{.Name}}" {
  name = "{{.Name}}-${count.index}",
  ami = "{{.Image}}"
  instance_type = "{{.Flavor}}"
  count = "{{.Count}}"
}
`

func instance (config *config.Config) {

  for i, h := range config.Hostgroups {
    fmt.Printf("%d %+v", i, h)

    t := template.New("instance")
    t,_ = t.Parse(instance_resource_tmpl)
    t.Execute(os.Stdout, h)
  }
  //fmt.Println("Here 2 !")
}
