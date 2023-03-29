package uaa

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/orange-cloudfoundry/terraform-provider-uaa/uaa/uaaapi"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider
var testAccProvidersFactories map[string]func() (*schema.Provider, error)

var tstSession *uaaapi.Session

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

func testSession() *uaaapi.Session {

	if !testAccEnvironmentSet() {
		panic(fmt.Errorf("ERROR! test UAA_* environment variables have not been set"))
	}

	if tstSession == nil {
		c := Config{
			loginEndpoint: os.Getenv("UAA_LOGIN_URL"),
			authEndpoint:  os.Getenv("UAA_AUTH_URL"),
			clientID:      os.Getenv("UAA_CLIENT_ID"),
			clientSecret:  os.Getenv("UAA_CLIENT_SECRET"),
		}
		c.skipSslValidation, _ = strconv.ParseBool(os.Getenv("UAA_SKIP_SSL_VALIDATION"))

		var (
			err     error
			session *uaaapi.Session
		)

		if session, err = c.Client(); err != nil {
			fmt.Printf("ERROR! Error creating a new session: %s\n", err.Error())
			panic(err.Error())
		}
		tstSession = session
	}
	return tstSession
}

func assertContains(str string, list []string) bool {
	for _, s := range list {
		if str == s {
			return true
		}
	}
	return false
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

func assertListEquals(attributes map[string]string,
	key string, actualLen int,
	match func(map[string]string, int) bool) (err error) {

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

	if actualLen > 0 && n == 0 {
		return fmt.Errorf(
			"expected resource '%s' to be empty but it has '%d' elements", key, actualLen)
	}
	if actualLen != n {
		return fmt.Errorf(
			"expected resource '%s' to have '%d' elements but it has '%d' elements",
			key, n, actualLen)
	}
	if n > 0 {
		found := 0

		var (
			values map[string]string
			ok     bool
		)

		keyValues := make(map[string]map[string]string)
		for k, v := range attributes {
			keyParts := strings.Split(k, ".")
			if key == keyParts[0] && keyParts[1] != "#" {
				i := keyParts[1]
				if values, ok = keyValues[i]; !ok {
					values = make(map[string]string)
					keyValues[i] = values
				}
				if len(keyParts) == 2 {
					values["value"] = v
				} else {
					values[strings.Join(keyParts[2:], ".")] = v
				}
			}
		}

		for _, values := range keyValues {
			for j := 0; j < actualLen; j++ {
				if match(values, j) {
					found++
					break
				}
			}
		}
		if n != found {
			return fmt.Errorf(
				"expected list resource '%s' to match '%d' elements but matched only '%d' elements",
				key, n, found)
		}
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

func assertMapEquals(key string, attributes map[string]string, actual map[string]interface{}) (err error) {

	expected := make(map[string]interface{})
	for k, v := range attributes {
		keyParts := strings.Split(k, ".")
		if keyParts[0] == key && keyParts[1] != "%" {

			l := len(keyParts)
			m := expected
			for _, kk := range keyParts[1 : l-1] {
				if _, ok := m[kk]; !ok {
					m[kk] = make(map[string]interface{})
				}
				m = m[kk].(map[string]interface{})
			}
			m[keyParts[l-1]] = v
		}
	}
	if !reflect.DeepEqual(expected, actual) {
		err = fmt.Errorf("map with key '%s' expected to be %#v but was %#v", key, expected, actual)
	}
	return nil
}

func assertHTTPResponse(url string, expectedStatusCode int, expectedResponses *[]string) (err error) {

	var resp *http.Response
	if resp, err = http.Get(url); err != nil {
		return
	}
	if expectedStatusCode != resp.StatusCode {
		err = fmt.Errorf(
			"expected response status code from url '%s' to be '%d', but actual was: %s",
			url, expectedStatusCode, resp.Status)
		return
	}
	if expectedResponses != nil {
		in := resp.Body
		out := bytes.NewBuffer(nil)
		if _, err = io.Copy(out, in); err != nil {
			return
		}
		content := out.String()

		found := false
		for _, e := range *expectedResponses {
			if e == content {
				found = true
				break
			}
		}
		if !found {
			err = fmt.Errorf(
				"expected response from url '%s' to be one of '%v', but actual was '%s'",
				url, *expectedResponses, content)
		}
	}
	return
}
