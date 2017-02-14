package repo

import (
	"fmt"
	"os"

	"../config/"
	"../git/"
	"../gitlab/"
)

func Build(SystemConfig *config.SystemConfig, cell *config.Cell) error {

	var err error

	cell.Environment.Terraform, err = initialize(SystemConfig, cell.Name, cell.CustomerName, "terraform")

	fmt.Println(cell.Environment.Terraform)

	if err != nil {
		return err
	}

	cell.Environment.Ansible, err = initialize(SystemConfig, cell.Name, cell.CustomerName, "ansible")

	if err != nil {
		return err
	}

	return nil
}

func initialize(SystemConfig *config.SystemConfig, CellName string, CustomerName string, RepoType string) (*config.RepoEnv, error) {

	var Project *gitlab.Project

	SystemConfig.Log.Debugf("CellName(%s) CustomerName(%s) RepoType(%s)", CellName, CustomerName, RepoType)

	RepoEnv := new(config.RepoEnv)

	RepoEnv.Name = fmt.Sprintf("%s-%s", CellName, RepoType)
	projectPath := fmt.Sprintf("%s/%s", CustomerName, RepoEnv.Name)
	RepoEnv.Dir = fmt.Sprintf("%s/%s/%s/%s/", SystemConfig.Files.TempDir, CustomerName, CellName, RepoType)

	Gitlab, err := gitlab.Connect(SystemConfig.Gitlab.Url, SystemConfig.Gitlab.Token, SystemConfig.Gitlab.TLSInsecureSkipVerify)

	if err != nil {
		return nil, fmt.Errorf("connecting to gitlab %s, %v", SystemConfig.Gitlab.Url, err)
	}

	Group, res, err := Gitlab.Groups.GetGroup(CustomerName)

	if err != nil {
		if res == nil {
			return nil, fmt.Errorf("getting gitlab group %s from %s, %v", CustomerName, SystemConfig.Gitlab.Url, err)
		}

		if res.StatusCode != 404 {
			return nil, fmt.Errorf("getting gitlab group %s, %v", CustomerName, err)
		}
	}

	if Group == nil || len(Group.Name) == 0 {
		Group, err = gitlab.CreateGroup(Gitlab, CustomerName, CustomerName)

		if err != nil {
			return nil, fmt.Errorf("creating gitlab group %s, %v", CustomerName, err)
		}
		SystemConfig.Log.Infof("Created group(%s)", CustomerName)
	} else {
		for _, v := range *Group.Projects {
			if v.Name == RepoEnv.Name {
				Project = &gitlab.Project{&v}
				break
			}
		}
	}

	SystemConfig.Log.Debugf("Project (%s)", RepoEnv.Name)
	if Project == nil {
		Project, err = gitlab.CreateProject(Gitlab, RepoEnv.Name, Group.ID)

		if err != nil {
			return nil, fmt.Errorf("creating gitlab project %s, %v", projectPath, err)
		}
		SystemConfig.Log.Infof("Created project(%s)", projectPath)
	}

	d, err := os.Stat(RepoEnv.Dir)

	if err == nil && d.IsDir() {
		SystemConfig.Log.Infof("Pulling %s on %s", Project.HTTPURLToRepo, RepoEnv.Dir)

		git.Pull(RepoEnv.Dir, Project.HTTPURLToRepo)

		if err != nil {
			return nil, fmt.Errorf("pulling gitlab project %s, %v", projectPath, err)
		}
	} else {
		os.Remove(RepoEnv.Dir)

		if err := os.MkdirAll(RepoEnv.Dir, 0750); err != nil {
			return nil, fmt.Errorf("Creating %s, %v", RepoEnv.Dir, err)
		}

		SystemConfig.Log.Infof("Cloning %s on %s", Project.HTTPURLToRepo, RepoEnv.Dir)

		err = git.Clone(RepoEnv.Dir, Project.HTTPURLToRepo, SystemConfig.Gitlab.TLSInsecureSkipVerify)
		if err != nil {
			return nil, fmt.Errorf("cloning gitlab project %s, %v", projectPath, err)
		}
	}

	return RepoEnv, nil
}
