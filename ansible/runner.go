package ansible

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"../config"
	"../utils"
)

func RolesInstall(system *config.SystemConfig, repoPath string) error {

	rolePath := fmt.Sprintf("%s/%s", repoPath, "roles")
	cmd, _, stderr := utils.RunCmd(system.Commands.AnsibleGalaxy, "install", "--roles-path", rolePath, "-r", system.Files.AnsibleRequirements)

	if err_wait := cmd.Wait(); err_wait != nil {
		return errors.New(stderr.String())
	}

	return nil
}

func Run(system *config.SystemConfig, repoPath string) error {

	cmd, _, stderr := utils.RunCmd(system.Commands.Ansible, "-i", "hosts", system.Files.AnsiblePlaybook)

	if !strings.HasSuffix(repoPath, "/") {
		repoPath = repoPath + "/"
	}

	cmd.Env = []string{}
	cmd.Dir = filepath.Dir(repoPath)

	if err_wait := cmd.Wait(); err_wait != nil {
		return errors.New(stderr.String())
	}

	return nil
}

func SyntaxCheck(system *config.SystemConfig, rolePath string) error {

	cmd, _, stderr := utils.RunCmd(system.Commands.Ansible, "-i", "hosts", "--syntax-check", system.Files.AnsiblePlaybook)

	if err_wait := cmd.Wait(); err_wait != nil {
		return errors.New(stderr.String())
	}

	return nil
}
