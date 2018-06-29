package instance

import "testing"

func TestAction(t *testing.T) {
	for _, v := range []string{"list", "start", "stop", "reboot"} {
		err := CheckArgs(v)
		if err != nil {
			t.Errorf("%s action should not error", v)
		}
	}
	a := "foo"
	err := CheckArgs(a)
	if err == nil {
		t.Errorf("%s action should error", a)
	}
}
