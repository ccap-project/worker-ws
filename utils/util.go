package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"os/exec"
	"path/filepath"
	"text/template"
	"time"

	"github.com/oklog/ulid"
)

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
func RunCmd(pwd string, env []string, arg ...string) (*exec.Cmd, *bufio.Scanner, *bytes.Buffer) {

	var stderr bytes.Buffer
	//log.Println(arg)

	cmd := exec.Command(arg[0], arg[1:]...)

	cmd.Dir = filepath.Dir(pwd)
	cmd.Env = env

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

func GetULID() (ulid.ULID, error) {
	t := time.Unix(1000000, 0)
	entropy := rand.New(rand.NewSource(t.UnixNano()))

	return ulid.New(ulid.Timestamp(t), entropy)
}
