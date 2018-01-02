package instance

import (
	"errors"
)

var actions = []string{"list", "start", "stop"}

// Check command arguments
func CheckArgs(arg string) error {
	for _, v := range actions {
		if v == arg {
			return nil
		}
	}
	return errors.New("Specified action not defined")
}
