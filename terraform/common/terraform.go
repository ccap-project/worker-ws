package terraformcommon

import (
	"fmt"

	"../../config"
	"../../utils"
)

func Apply(SystemConfig *config.SystemConfig, cell *config.Cell) (*[]byte, error) {

	out, err := utils.RunCmd(cell.Environment.Terraform.Dir, cell.Environment.Terraform.Env, SystemConfig.Commands.Terraform, "apply", cell.Environment.Terraform.Dir)

	if err != nil {
		return nil, fmt.Errorf("%v, %s", err, out)
	}

	return out, nil
}

func Validate(SystemConfig *config.SystemConfig, cell *config.Cell) error {

	out, err := utils.RunCmd(cell.Environment.Terraform.Dir, cell.Environment.Terraform.Env, SystemConfig.Commands.Terraform, "validate", cell.Environment.Terraform.Dir)

	if err != nil {
		return fmt.Errorf("%v, %s", err, out)
	}

	return nil
}
