package utils

import "bufio"
import "bytes"
import "text/template"
//import "io/ioutil"

func Template (tmpl string, data interface{}) (*bytes.Buffer, error) {

  var err error
  var b bytes.Buffer

  f := bufio.NewWriter(&b)

  t := template.New("instance")

  t,err = t.Parse(tmpl)
  if err != nil {
    return &b, err
  }

  err = t.Execute(f, data)
  if err != nil {
    return &b, err
  }

  f.Flush()

  return &b, nil
}
