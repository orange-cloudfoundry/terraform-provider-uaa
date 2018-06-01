package uaaapi

import (
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"code.cloudfoundry.org/cli/cf/configuration"
	"code.cloudfoundry.org/cli/cf/configuration/coreconfig"
	"code.cloudfoundry.org/cli/cf/i18n"
	"code.cloudfoundry.org/cli/cf/net"
)

// Session - wraps the CF CLI session objects
type Session struct {
	Log *Logger

	config     coreconfig.Repository
	refresher  coreconfig.APIConfigRefresher
	uaaGateway net.Gateway

	authManager   *AuthManager
	userManager   *UserManager
	clientManager *ClientManager

	// Used for direct endpoint calls
	httpClient *http.Client
}

// uaaErrorResponse -
type uaaErrorResponse struct {
	Code        string `json:"error"`
	Description string `json:"error_description"`
}

// NewSession -
func NewSession(
	uaaLoginEndpoint, uaaAuthEndpoint,
	uaaClientID, uaaClientSecret, caCert string,
	skipSslValidation bool) (s *Session, err error) {

	s = &Session{
		httpClient: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: skipSslValidation},
			},
		},
	}

	envDialTimeout := os.Getenv("UAA_DIAL_TIMEOUT")

	debug, _ := strconv.ParseBool(os.Getenv("UAA_DEBUG"))
	s.Log = NewLogger(debug, os.Getenv("UAA_TRACE"))

	s.config = coreconfig.NewRepositoryFromPersistor(&noopPersistor{}, func(err error) {
		if err != nil {
			s.Log.UI.Failed(err.Error())
			os.Exit(1)
		}
	})
	if i18n.T == nil {
		i18n.T = i18n.Init(s.config)
	}
	s.config.SetSSLDisabled(skipSslValidation)

	s.config.SetAuthenticationEndpoint(endpointAsURL(uaaLoginEndpoint))
	s.config.SetUaaEndpoint(endpointAsURL(uaaAuthEndpoint))

	s.uaaGateway = net.NewUAAGateway(s.config, s.Log.UI, s.Log.TracePrinter, envDialTimeout)
	s.authManager = NewAuthManager(s.uaaGateway, s.config, net.NewRequestDumper(s.Log.TracePrinter))
	s.uaaGateway.SetTokenRefresher(s.authManager)

	s.userManager, err = newUserManager(s.config, s.uaaGateway, s.Log)
	if err != nil {
		return nil, err
	}

	s.clientManager, err = newClientManager(s.config, s.uaaGateway, s.Log)
	if err != nil {
		return nil, err
	}

	if s.userManager.clientToken, err = s.authManager.GetClientToken(uaaClientID, uaaClientSecret); err == nil {
		err = s.userManager.loadGroups()
	}

	return
}

// UserManager -
func (s *Session) UserManager() *UserManager {
	return s.userManager
}

// UserManager -
func (s *Session) ClientManager() *ClientManager {
	return s.clientManager
}

// AuthManager -
func (s *Session) AuthManager() *AuthManager {
	return s.authManager
}

// noopPersistor - No Op Persistor for CF CLI session
type noopPersistor struct {
}

func newNoopPersistor() configuration.Persistor {
	return &noopPersistor{}
}

func (p *noopPersistor) Delete() {
}

func (p *noopPersistor) Exists() bool {
	return false
}

func (p *noopPersistor) Load(configuration.DataInterface) error {
	return nil
}

func (p *noopPersistor) Save(configuration.DataInterface) error {
	return nil
}

// endpointAsURL
func endpointAsURL(endpoint string) string {

	endpoint = strings.TrimSuffix(endpoint, "/")
	if !strings.HasPrefix(endpoint, "http") {
		endpoint = "https://" + endpoint
	}
	return endpoint
}

// newUUID generates a random UUID according to RFC 4122
func newUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}

	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}
