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
