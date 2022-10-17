package util

import (
	"fmt"
)

// Assertions

func AssertSame(expected interface{}, actual interface{}) error {
	if actual != expected {
		return fmt.Errorf("expected '%s' found '%s' ", expected, actual)
	}
	return nil
}
