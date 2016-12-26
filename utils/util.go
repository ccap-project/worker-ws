package utils

import "bufio"
import "bytes"
import "text/template"

func Template (tmpl string, data interface{}) (string) {

  var b bytes.Buffer
  f := bufio.NewWriter(&b)
  //defer f.Close()

  t := template.New("instance")

  t,_ = t.Parse(tmpl)

  t.Execute(f, data)
  f.Flush()

  return(b.String())
}
