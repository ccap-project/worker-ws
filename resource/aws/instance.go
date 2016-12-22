package aws

//import "fmt"
import "os"
import "text/template"
import "../../config/"

const instance_templ = `resource "aws_instance" "{{.name}}" {
    instance_type = "{{.flavor}}"
}
`

//func (h *hostgroup) Marshall() {
func Hostgroup(SystemConfig *config.Config) {
  t := template.New("Person template")
  t, _ = t.Parse(instance_templ)
  t.Execute(os.Stdout, SystemConfig.Hostgroup[0])
}
