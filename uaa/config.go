package uaa

import "github.com/orange-cloudfoundry/terraform-provider-uaa/uaa/uaaapi"

// Config -
type Config struct {
	loginEndpoint     string
	authEndpoint      string
	clientID          string
	clientSecret      string
	caCert            string
	skipSslValidation bool
}

// Client - Terraform providor client initialization
func (c *Config) Client() (*uaaapi.Session, error) {
	return uaaapi.NewSession(c.loginEndpoint, c.authEndpoint, c.clientID, c.clientSecret, c.caCert, c.skipSslValidation)
}
