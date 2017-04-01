package repo

import (
	"fmt"
	"os"

	"../config/"
	"../git/"
	"../gitlab/"
)

func Build(ctx *config.RequestContext) error {

	var err error

	ctx.Cell.Environment.Terraform, err = initialize(ctx, "terraform")

	if err != nil {
		return err
	}

	ctx.Cell.Environment.Ansible, err = initialize(ctx, "ansible")

	if err != nil {
		return err
	}

	return nil
}

func Persist(ctx *config.RequestContext, repoEnv *config.RepoEnv, needTag bool) error {

	commit, err := git.Commit(repoEnv.Dir, repoEnv.Dir)
	if err != nil {
		return err
	}

	//ctx.Log = ctx.SystemConfig.Log.WithFields(log.Fields{"commit_id": commit.String()})

	if needTag {
		err = git.Tag(repoEnv.Dir, commit, ctx.RunID)
		if err != nil {
			return err
		}
		err = git.Push(repoEnv.Dir, ctx.RunID)
	} else {
		err = git.Push(repoEnv.Dir, "")
	}

	if err != nil {
		return err
	}

	return nil
}

func initialize(ctx *config.RequestContext, RepoType string) (*config.RepoEnv, error) {

	var Project *gitlab.Project

	ctx.Log.Debugf("RepoType(%s)", RepoType)

	RepoEnv := new(config.RepoEnv)

	RepoEnv.Name = fmt.Sprintf("%s-%s", ctx.Cell.Name, RepoType)
	projectPath := fmt.Sprintf("%s/%s", ctx.Cell.CustomerName, RepoEnv.Name)
	RepoEnv.Dir = fmt.Sprintf("%s/%s/%s/%s/", ctx.SystemConfig.Files.TempDir, ctx.Cell.CustomerName, ctx.Cell.Name, RepoType)

	switch RepoType {
	case "ansible":
		RepoEnv.Env = append(RepoEnv.Env, fmt.Sprintf("ANSIBLE_PRIVATE_KEY_FILE=%s", "/Users/ale/.ssh/id_rsa"))
		RepoEnv.Env = append(RepoEnv.Env, "ANSIBLE_GALAXY_IGNORE=true")
		RepoEnv.Env = append(RepoEnv.Env, "GIT_SSL_NO_VERIFY=true")
		RepoEnv.Env = append(RepoEnv.Env, "HOST_KEY_CHECKING=true")

	case "terraform":
		RepoEnv.Env = append(RepoEnv.Env, fmt.Sprintf("TF_LOG=%s", "INFO"))
		RepoEnv.Env = append(RepoEnv.Env, fmt.Sprintf("TF_LOG_PATH=%s/log", RepoEnv.Dir))
	}

	ctx.Log.Debugf(fmt.Sprintf("Env(%s)", RepoEnv.Env))

	Gitlab, err := gitlab.Connect(ctx.SystemConfig.Gitlab.Url, ctx.SystemConfig.Gitlab.Token, ctx.SystemConfig.Gitlab.TLSInsecureSkipVerify)

	if err != nil {
		return nil, fmt.Errorf("connecting to gitlab %s, %v", ctx.SystemConfig.Gitlab.Url, err)
	}

	Group, res, err := Gitlab.Groups.GetGroup(ctx.Cell.CustomerName)

	if err != nil {
		if res == nil {
			return nil, fmt.Errorf("getting gitlab group %s from %s, %v", ctx.Cell.CustomerName, ctx.SystemConfig.Gitlab.Url, err)
		}

		if res.StatusCode != 404 {
			return nil, fmt.Errorf("getting gitlab group %s, %v", ctx.Cell.CustomerName, err)
		}
	}

	if Group == nil || len(Group.Name) == 0 {
		Group, err = gitlab.CreateGroup(Gitlab, ctx.Cell.CustomerName, ctx.Cell.CustomerName)

		if err != nil {
			return nil, fmt.Errorf("creating gitlab group %s, %v", ctx.Cell.CustomerName, err)
		}
		ctx.Log.Infof("Created group(%s)", ctx.Cell.CustomerName)
	} else {
		for _, v := range *Group.Projects {
			if v.Name == RepoEnv.Name {
				Project = &gitlab.Project{&v}
				break
			}
		}
	}

	ctx.Log.Debugf("Project (%s)", RepoEnv.Name)
	if Project == nil {
		Project, err = gitlab.CreateProject(Gitlab, RepoEnv.Name, Group.ID)

		if err != nil {
			return nil, fmt.Errorf("creating gitlab project %s, %v", projectPath, err)
		}
		//projectCreated = true
		ctx.Log.Infof("Created project(%s)", projectPath)
	}

	/*
	 * Local repo handling
	 */
	d, err := os.Stat(RepoEnv.Dir)

	/*
	 * Use existant directory structure
	 */
	if err == nil && d.IsDir() {
		ctx.Log.Infof("Checkout %s on %s", Project.HTTPURLToRepo, RepoEnv.Dir)

		err = git.Checkout(RepoEnv.Dir, ctx.SystemConfig.Gitlab.TLSInsecureSkipVerify, "")

		if err != nil {
			return nil, fmt.Errorf("Checkout repo %s, %v", projectPath, err)
		}

	} else {
		/*
		 * Clone remote directory structure
		 */
		os.Remove(RepoEnv.Dir)

		if err := os.MkdirAll(RepoEnv.Dir, 0750); err != nil {
			return nil, fmt.Errorf("Creating %s, %v", RepoEnv.Dir, err)
		}

		ctx.Log.Infof("Cloning %s on %s", Project.HTTPURLToRepo, RepoEnv.Dir)

		err = git.Clone(RepoEnv.Dir, Project.SSHURLToRepo, ctx.SystemConfig.Gitlab.TLSInsecureSkipVerify)
		if err != nil {
			return nil, fmt.Errorf("cloning gitlab project %s, %v", projectPath, err)
		}

		//if projectCreated {
		//	git.CreateBranch(RepoEnv.Dir, "main")
		//}

		os.MkdirAll(RepoEnv.Dir+"/roles", 0750)
	}

	// XXX: checkout defined tag if exists
	if len(ctx.TagID) > 0 {
		ctx.Log.Infof("Checkout %s on %s tag %s", Project.HTTPURLToRepo, RepoEnv.Dir, ctx.TagID)

		err = git.Checkout(RepoEnv.Dir, ctx.SystemConfig.Gitlab.TLSInsecureSkipVerify, ctx.TagID)

		if err != nil {
			return nil, fmt.Errorf("Checkout gitlab project %s, tag %s, %v", projectPath, ctx.RunID, err)
		}
	}

	return RepoEnv, nil
}
