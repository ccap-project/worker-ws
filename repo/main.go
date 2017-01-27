package repo

import (
	"fmt"
	"os"

	"../config/"
	"../gitlab/"
)

func BuildInfrastructureEnv(SystemConfig *config.SystemConfig, customer string, cell string) (string, error) {

	var Project *gitlab.Project

	projectName := fmt.Sprintf("%s-terraform", cell)
	projectPath := fmt.Sprintf("%s/%s", customer, projectName)
	tempPath := fmt.Sprintf("%s/%s/%s/terraform", SystemConfig.Files.TempDir, customer, cell)

	Gitlab, err := gitlab.Connect(SystemConfig.Gitlab.Url, SystemConfig.Gitlab.Token, SystemConfig.Gitlab.TLSInsecureSkipVerify)

	if err != nil {
		return tempPath, fmt.Errorf("connecting to gitlab %s, %v", SystemConfig.Gitlab.Url, err)
	}

	Group, res, err := Gitlab.Groups.GetGroup(customer)

	if err != nil && res.StatusCode != 404 {
		return tempPath, fmt.Errorf("getting gitlab group %s, %v", customer, err)
	}

	if Group == nil || len(Group.Name) == 0 {
		Group, err = gitlab.CreateGroup(Gitlab, customer, customer)

		if err != nil {
			return tempPath, fmt.Errorf("creating gitlab group %s, %v", customer, err)
		}
		SystemConfig.Log.Infof("Created group(%s)", customer)
	} else {
		for _, v := range *Group.Projects {
			if v.Name == projectName {
				Project = &gitlab.Project{&v}
				break
			}
		}
	}

	SystemConfig.Log.Debugf("Project (%s)", Project.Name)

	if Project == nil {
		Project, err = gitlab.CreateProject(Gitlab, projectName, Group.ID)

		if err != nil {
			return tempPath, fmt.Errorf("creating gitlab project %s, %v", projectPath, err)
		}
		SystemConfig.Log.Infof("Created project(%s)", projectPath)
	}

	d, err := os.Stat(tempPath)

	if err == nil && d.IsDir() {
		SystemConfig.Log.Infof("Pulling %s on %s", Project.HTTPURLToRepo, tempPath)

		//* XXX: While go-git isn't capable of repo clone with worktree...
		//git.Pull(tempPath, Project.HTTPURLToRepo)

		if err != nil {
			return tempPath, fmt.Errorf("pulling gitlab project %s, %v", projectPath, err)
		}
	} else {
		os.Remove(tempPath)

		if err := os.MkdirAll(tempPath, 0750); err != nil {
			return tempPath, fmt.Errorf("Creating %s, %v", tempPath, err)
		}

		SystemConfig.Log.Infof("Cloning %s on %s", Project.HTTPURLToRepo, tempPath)

		// XXX: While go-git isn't capable of repo clone with worktree...
		//err = git.Clone(tempPath, Project.HTTPURLToRepo, SystemConfig.Gitlab.TLSInsecureSkipVerify)
		//if err != nil {
		//	return tempPath, fmt.Errorf("cloning gitlab project %s, %v", projectPath, err)
		//}
	}

	return tempPath, nil
}
