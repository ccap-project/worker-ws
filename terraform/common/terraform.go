package terraformcommon

import (
	"errors"

	"../../config/"
	"../../utils/"
)

func Apply(SystemConfig *config.SystemConfig, dir string) error {

	cmd, _, stderr := utils.RunCmd(SystemConfig.Commands.Terraform, "apply", dir)

	if err_wait := cmd.Wait(); err_wait != nil {
		return errors.New(stderr.String())
	}

	return nil
}

func Validate(SystemConfig *config.SystemConfig, dir string) error {

	cmd, _, stderr := utils.RunCmd(SystemConfig.Commands.Terraform, "validate", dir)

	if err_wait := cmd.Wait(); err_wait != nil {
		return errors.New(stderr.String())
	}

	return nil
}
