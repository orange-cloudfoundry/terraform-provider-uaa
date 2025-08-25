package uaa

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider
var testAccProvidersFactories map[string]func() (*schema.Provider, error)

func init() {

	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"uaa": testAccProvider,
	}
	testAccProvidersFactories = map[string]func() (*schema.Provider, error){
		"uaa": func() (*schema.Provider, error) {
			return testAccProvider, nil
		},
	}

}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ = Provider()
}

func testAccPreCheck(t *testing.T) {

	if !testAccEnvironmentSet() {
		t.Fatal("Acceptance environment has not been set.")
	}
}

func testAccEnvironmentSet() bool {

	loginEndpoint := os.Getenv("UAA_LOGIN_URL")
	authEndpoint := os.Getenv("UAA_AUTH_URL")
	clientID := os.Getenv("UAA_CLIENT_ID")
	clientSecret := os.Getenv("UAA_CLIENT_SECRET")
	skipSslValidation := strings.ToLower(os.Getenv("UAA_SKIP_SSL_VALIDATION"))

	if len(loginEndpoint) == 0 ||
		len(authEndpoint) == 0 ||
		len(clientID) == 0 ||
		len(clientSecret) == 0 ||
		len(skipSslValidation) == 0 {

		fmt.Println("UAA_LOGIN_URL, UAA_AUTH_URL, UAA_CLIENT_ID, UAA_CLIENT_SECRET " +
			"and UAA_SKIP_SSL_VALIDATION must be set for acceptance tests to work.")
		return false
	}
	return true
}

func assertSame(actual interface{}, expected interface{}) error {
	if actual != expected {
		return fmt.Errorf("expected '%s' found '%s' ", expected, actual)
	}
	return nil
}

func assertEquals(attributes map[string]string,
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
				s = expectedValueContent.String()
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

func assertSetEquals(attributes map[string]string,
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
