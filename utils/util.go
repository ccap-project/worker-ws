/*
 *
 * Copyright (c) 2016, 2017, 2018 Alexandre Biancalana <ale@biancalanas.net>.
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *     * Redistributions of source code must retain the above copyright
 *       notice, this list of conditions and the following disclaimer.
 *     * Redistributions in binary form must reproduce the above copyright
 *       notice, this list of conditions and the following disclaimer in the
 *       documentation and/or other materials provided with the distribution.
 *     * Neither the name of the <organization> nor the
 *       names of its contributors may be used to endorse or promote products
 *       derived from this software without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
 * ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
 * WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL <COPYRIGHT HOLDER> BE LIABLE FOR ANY
 * DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
 * (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
 * LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
 * ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
 * SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 */

package utils

import (
	"bufio"
	"bytes"
	"math/rand"
	"os/exec"
	"path/filepath"
	"reflect"
	"text/template"
	"time"

	"github.com/oklog/ulid"
)

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
	t := time.Now()
	entropy := rand.New(rand.NewSource(t.UnixNano()))

	return ulid.New(ulid.Timestamp(t), entropy)
}

func Grep(val interface{}, array interface{}) (exists bool, index int) {
	exists = false
	index = -1

	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
				index = i
				exists = true
				return
			}
		}
	}

	return
}
