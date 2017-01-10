package utils

import "bufio"
import "bytes"
import "text/template"
//import "io/ioutil"

func Template (tmpl string, data interface{}) (*bytes.Buffer) {

  var b bytes.Buffer

  // XXX: Error handling !!
  f := bufio.NewWriter(&b)

  t := template.New("instance")

  t,_ = t.Parse(tmpl)

  t.Execute(f, data)
  f.Flush()

  return(&b)
}

/*
func WriteFile (filename string, file_content []string) (error) {

  err := ioutil.WriteFile(filename,  []byte(file_content), 0644)

  return err
}
*/
