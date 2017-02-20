package utils

import (
	"bufio"
	"bytes"
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
func RunCmd(pwd string, env []string, arg ...string) (*[]byte, error) {

	cmd := exec.Command(arg[0], arg[1:]...)

	if len(pwd) > 0 {
		cmd.Dir = filepath.Dir(pwd)
	}

	if len(env) > 0 {
		cmd.Env = env
	}

	out, err := cmd.CombinedOutput()

	return &out, err
}

func GetULID() (ulid.ULID, error) {
	t := time.Unix(1000000, 0)
	entropy := rand.New(rand.NewSource(t.UnixNano()))

	return ulid.New(ulid.Timestamp(t), entropy)
}
