package utils

import "bufio"
import "bytes"
import "fmt"
import "text/template"
import "log"
import "os/exec"

//import "io/ioutil"

func Template(tmpl string, data interface{}) (*bytes.Buffer, error) {

	var err error
	var b bytes.Buffer

	f := bufio.NewWriter(&b)

	t := template.New("instance")

	t, err = t.Parse(tmpl)
	if err != nil {
		return nil, err
	}

	err = t.Execute(f, data)
	if err != nil {
		return nil, err
	}

	f.Flush()

	return &b, nil
}

/*
 * Run command, return stdout scanner and cmd handler
 */
func RunCmd(arg ...string) (*exec.Cmd, *bufio.Scanner, *bytes.Buffer) {

	var stderr bytes.Buffer
	//log.Println(arg)

	cmd := exec.Command(arg[0], arg[1:]...)

	cmd.Stderr = &stderr
	stdout, err_stdout := cmd.StdoutPipe()

	if err_stdout != nil {
		log.Fatalf("Error running(%v), %s", arg, fmt.Errorf("%v", err_stdout))
	}

	if err_start := cmd.Start(); err_start != nil {
		log.Fatalf("Error running(%v), %s", arg, fmt.Errorf("%v", err_start))
	}

	stdout_scanner := bufio.NewScanner(stdout)

	return cmd, stdout_scanner, &stderr
}
