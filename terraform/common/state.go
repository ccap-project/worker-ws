package terraformcommon

import (
	"fmt"
	"os"

	t "github.com/hashicorp/terraform/terraform"
)

func ReadState(file string) (*t.State, error) {

	f, err := os.Open(file)
	defer f.Close()

	if err != nil {
		fmt.Println("err: ", err)
		return nil, err
	}

	state, err := t.ReadState(f)
	if err != nil {
		fmt.Println("err: ", err)
		return nil, err
	}

	return state, nil
}
