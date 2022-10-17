package util

import (
	"fmt"
	"reflect"
	"strconv"
)

func AssertSame(expected interface{}, actual interface{}) error {
	if actual != expected {
		return fmt.Errorf("expected '%s' found '%s' ", expected, actual)
	}
	return nil
}

func AssertEquals(attributes map[string]string,
	key string, expected interface{}) error {
	v, ok := attributes[key]

	expectedValue := reflect.ValueOf(expected)

	if ok {

		var s string
		if expectedValue.Kind() == reflect.Ptr {

			if expectedValue.IsNil() {
				return fmt.Errorf("expected resource '%s' to not be present but it was '%s'", key, v)
			}

			expectedValueContent := reflect.Indirect(reflect.ValueOf(expected))
			switch expectedValueContent.Kind() {
			case reflect.String:
				s = fmt.Sprintf("%s", expectedValueContent.String())
			case reflect.Int:
				s = fmt.Sprintf("%d", expectedValueContent.Int())
			case reflect.Bool:
				s = fmt.Sprintf("%t", expectedValueContent.Bool())
			default:
				return fmt.Errorf("unable to determine underlying content of expected value: %s", expectedValueContent.Kind())
			}
		} else {
			switch expected.(type) {
			case string:
				s = fmt.Sprintf("%s", expected)
			case int:
				s = fmt.Sprintf("%d", expected)
			case bool:
				s = fmt.Sprintf("%t", expected)
			default:
				s = fmt.Sprintf("%v", expected)
			}
		}
		if v != s {
			return fmt.Errorf("expected resource '%s' to be '%s' but it was '%s'", key, expected, v)
		}
	} else if expectedValue.Kind() == reflect.Ptr && !expectedValue.IsNil() {
		return fmt.Errorf("expected resource '%s' to be '%s' but it was not present", key, reflect.Indirect(reflect.ValueOf(expected)))
	}
	return nil
}

func AssertSetEquals(attributes map[string]string,
	key string, expected []interface{}) (err error) {

	var n int

	num := attributes[key+".#"]
	if len(num) > 0 {
		n, err = strconv.Atoi(num)
		if err != nil {
			return
		}
	} else {
		n = 0
	}

	if len(expected) > 0 && n == 0 {
		return fmt.Errorf(
			"expected resource '%s' to be '%v' but it was empty", key, expected)
	}
	if len(expected) != n {
		return fmt.Errorf(
			"expected resource '%s' to have '%d' elements but it has '%d' elements",
			key, len(expected), n)
	}
	if n > 0 {
		found := 0
		for i := range expected {
			if _, ok := attributes[key+"."+strconv.Itoa(i)]; ok {
				found++
			}
		}
		if n != found {
			return fmt.Errorf(
				"expected set resource '%s' to have elements '%v' but matched only '%d' elements",
				key, expected, found)
		}
	}
	return
}
