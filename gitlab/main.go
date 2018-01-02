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

package gitlab

import (
	"crypto/tls"
	"fmt"
	"net/http"

	g "github.com/xanzy/go-gitlab"
)

type GitGroup_t struct {
	Name  string
	Path  string
	Group string
}

type Role_t struct {
	ID       int
	Name     string
	Url      string
	Versions []string
}

type Project struct {
	*g.Project
}

/**********************************************
 *
 **********************************************/
func Connect(url string, token string, TLSInsecureSkipVerify bool) (*g.Client, error) {

	transport := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: TLSInsecureSkipVerify}}
	http := &http.Client{Transport: transport}

	client := g.NewClient(http, token)

	if err := client.SetBaseURL(url); err != nil {
		return nil, fmt.Errorf("Failure connecting to gitlab(%s), %v", url, err)
	}

	return client, nil
}

func CreateGroup(client *g.Client, name string, path string) (*g.Group, error) {

	Group, _, err := client.Groups.CreateGroup(&g.CreateGroupOptions{Name: g.String(name),
		Path: g.String(path)})
	if err != nil {
		return nil, err
	}

	return Group, nil
}

func CreateProject(client *g.Client, name string, namespace int) (*Project, error) {

	project, _, err := client.Projects.CreateProject(&g.CreateProjectOptions{Name: g.String(name),
		NamespaceID: g.Int(namespace)})

	if err != nil {
		return nil, err
	}

	return &Project{project}, nil
}
