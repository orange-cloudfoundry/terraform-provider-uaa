package util

import "fmt"

func AssertSame(expected interface{}, actual interface{}) error {
	if actual != expected {
		return fmt.Errorf("expected '%s' found '%s' ", expected, actual)
	}
	return nil
}
