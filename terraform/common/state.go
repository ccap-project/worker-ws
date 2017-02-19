package terraformcommon

import (
	"os"

	t "github.com/hashicorp/terraform/terraform"
)

func ReadState(file string) (*t.State, error) {

	f, err := os.Open(file)
	defer f.Close()

	if err != nil {
		return nil, err
	}

	state, err := t.ReadState(f)
	if err != nil {
		return nil, err
	}

	return state, nil
}
