package terraformcommon

import (
	"fmt"

	"../../config/"
	"../../utils/"
)

func Apply(SystemConfig *config.SystemConfig, dir string) error {

	out, err := utils.RunCmd(dir, []string{}, SystemConfig.Commands.Terraform, "apply", dir)

	if err != nil {
		return fmt.Errorf("%v, %s", err, out)
	}

	return nil
}

func Validate(SystemConfig *config.SystemConfig, dir string) error {

	out, err := utils.RunCmd(dir, []string{}, SystemConfig.Commands.Terraform, "validate", dir)

	if err != nil {
		return fmt.Errorf("%v, %s", err, out)
	}

	return nil
}
