package instance

import (
	"errors"
	"fmt"
	"strings"
)

var actions = []string{"list", "start", "stop"}

// Check command arguments
func CheckArgs(arg string) error {
	for _, v := range actions {
		if v == arg {
			return nil
		}
	}
	str := fmt.Sprintf("Specified action not defined. Must be one of %s.", strings.Join(actions, ","))
	return errors.New(str)
}
