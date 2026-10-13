package uaaapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"code.cloudfoundry.org/cli/v9/cf/configuration/coreconfig"
	cfErrors "code.cloudfoundry.org/cli/v9/cf/errors"
	"code.cloudfoundry.org/cli/v9/cf/net"
)

// ClientManager -
type ClientManager struct {
	log        *Logger
	config     coreconfig.Reader
	uaaGateway net.Gateway
}

// ClientSecret         string   `json:"client_secret,omitempty"`

// UAAClient -
type UAAClient struct {
	ClientID             string           `json:"client_id,omitempty"`
	ClientSecret         string           `json:"client_secret,omitempty"`
	AuthorizedGrantTypes *[]string        `json:"authorized_grant_types,omitempty"`
	RedirectURI          *[]string        `json:"redirect_uri,omitempty"`
	Scope                *[]string        `json:"scope,omitempty"`
	ResourceIds          *[]string        `json:"resource_ids,omitempty"`
	Authorities          *[]string        `json:"authorities,omitempty"`
	Autoapprove          *[]string        `json:"autoapprove,omitempty"`
	AccessTokenValidity  int              `json:"access_token_validity,omitempty"`
	RefreshTokenValidity int              `json:"refresh_token_validity,omitempty"`
	Allowedproviders     *[]string        `json:"allowedproviders,omitempty"`
	Name                 string           `json:"name,omitempty"`
	TokenSalt            string           `json:"token_salt,omitempty"`
	CreatedWith          string           `json:"createdwith,omitempty"`
	ApprovalsDeleted     bool             `json:"approvals_deleted,omitempty"`
	RequiredUserGroups   *[]string        `json:"required_user_groups,omitempty"`
	LastModified         int64            `json:"lastModified,omitempty"`
	ClientJwtConfig      *ClientJwtConfig `json:"client_jwt_config,omitempty"`
}

// ClientJwtCredential - federated JWT trust entry (RFC 7523)
type ClientJwtCredential struct {
	Issuer   string `json:"iss"`
	Subject  string `json:"sub"`
	Audience string `json:"aud,omitempty"`
}

// ClientJwtConfig - private_key_jwt trust configuration of a client
type ClientJwtConfig struct {
	JwksURI  string                `json:"jwks_uri,omitempty"`
	Jwks     json.RawMessage       `json:"jwks,omitempty"`
	JwtCreds []ClientJwtCredential `json:"jwt_creds,omitempty"`
}

// IsEmpty -
func (c *ClientJwtConfig) IsEmpty() bool {
	return c == nil || (c.JwksURI == "" && len(c.Jwks) == 0 && len(c.JwtCreds) == 0)
}

// MarshalJSON - UAA expects client_jwt_config as a JSON encoded string
func (c ClientJwtConfig) MarshalJSON() ([]byte, error) {
	type alias ClientJwtConfig
	raw, err := json.Marshal(alias(c))
	if err != nil {
		return nil, err
	}
	return json.Marshal(string(raw))
}

// UnmarshalJSON - accepts client_jwt_config either as a JSON encoded string or as an object
func (c *ClientJwtConfig) UnmarshalJSON(data []byte) error {
	type alias ClientJwtConfig
	if len(data) > 0 && data[0] == '"' {
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		if s == "" {
			return nil
		}
		data = []byte(s)
	}
	return json.Unmarshal(data, (*alias)(c))
}

// ClientJwtChangeRequest - body of PUT /oauth/clients/{client_id}/clientjwt
type ClientJwtChangeRequest struct {
	ClientID   string          `json:"client_id"`
	JwksURI    string          `json:"jwks_uri,omitempty"`
	Jwks       json.RawMessage `json:"jwks,omitempty"`
	Iss        string          `json:"iss,omitempty"`
	Sub        string          `json:"sub,omitempty"`
	Aud        string          `json:"aud,omitempty"`
	ChangeMode string          `json:"changeMode,omitempty"`
}

// UAAClientResourceList -
type UAAClientResourceList struct {
	Resources []UAAClient `json:"resources"`
}

func (c *UAAClient) HasDefaultScope() bool {
	return c.Scope != nil && len(*c.Scope) == 1 && (*c.Scope)[0] == "uaa.none"
}

func (c *UAAClient) HasDefaultAuthorites() bool {
	return c.Authorities != nil && len(*c.Authorities) == 1 && (*c.Authorities)[0] == "uaa.none"
}

func (c *UAAClient) HasDefaultResourceIds() bool {
	return c.ResourceIds != nil && len(*c.ResourceIds) == 1 && (*c.ResourceIds)[0] == "none"
}

// NewClientManager -
func newClientManager(config coreconfig.Reader, uaaGateway net.Gateway, logger *Logger) (cm *ClientManager, err error) {
	cm = &ClientManager{
		log:        logger,
		config:     config,
		uaaGateway: uaaGateway,
	}
	return
}

// GetClient -
func (cm *ClientManager) GetClient(id string) (client *UAAClient, err error) {
	uaaEndpoint := cm.config.UaaEndpoint()
	if len(uaaEndpoint) == 0 {
		err = cfErrors.New("UAA endpoint missing from config file")
		return
	}

	client = &UAAClient{}
	err = cm.uaaGateway.GetResource(
		fmt.Sprintf("%s/oauth/clients/%s", uaaEndpoint, id),
		client)
	return
}

func (cm *ClientManager) UaaEndPoint() string {
	return cm.config.UaaEndpoint()
}

// Create - Create a client
func (cm *ClientManager) Create(nCli UAAClient) (client UAAClient, err error) {
	uaaEndpoint := cm.config.UaaEndpoint()
	if len(uaaEndpoint) == 0 {
		err = cfErrors.New("UAA endpoint missing from config file")
		return
	}

	body, err := json.Marshal(nCli)
	if err != nil {
		return
	}

	client = nCli
	err = cm.uaaGateway.CreateResource(uaaEndpoint, "/oauth/clients", bytes.NewReader(body), &client)
	var httpErr cfErrors.HTTPError
	if errors.As(err, &httpErr) && httpErr.StatusCode() == http.StatusConflict {
		err = cfErrors.NewModelAlreadyExistsError("client", nCli.ClientID)
	}
	return
}

// UpdateClient -
func (cm *ClientManager) UpdateClient(nCli *UAAClient) (client UAAClient, err error) {
	uaaEndpoint := cm.config.UaaEndpoint()
	if len(uaaEndpoint) == 0 {
		err = cfErrors.New("UAA endpoint missing from config file")
		return
	}

	body, err := json.Marshal(nCli)
	if err != nil {
		return
	}

	request, err := cm.uaaGateway.NewRequest("PUT",
		fmt.Sprintf("%s/oauth/clients/%s", uaaEndpoint, nCli.ClientID),
		cm.config.AccessToken(), bytes.NewReader(body))
	if err != nil {
		return
	}

	client = *nCli
	_, err = cm.uaaGateway.PerformRequestForJSONResponse(request, client)
	return
}

// DeleteClient -
func (cm *ClientManager) DeleteClient(id string) (err error) {
	uaaEndpoint := cm.config.UaaEndpoint()
	if len(uaaEndpoint) == 0 {
		err = cfErrors.New("UAA endpoint missing from config file")
		return
	}
	err = cm.uaaGateway.DeleteResource(uaaEndpoint, fmt.Sprintf("/oauth/clients/%s", id))
	return
}

// ChangeSecret - Change the password/secret
func (cm *ClientManager) ChangeSecret(id, oldSecret, newSecret string) (err error) {
	uaaEndpoint := cm.config.UaaEndpoint()
	if len(uaaEndpoint) == 0 {
		err = cfErrors.New("UAA endpoint missing from config file")
		return
	}

	data := map[string]string{
		"secret": newSecret,
	}

	if len(oldSecret) != 0 {
		data["oldSecret"] = oldSecret
	}

	body, err := json.Marshal(data)
	if err != nil {
		return
	}

	request, err := cm.uaaGateway.NewRequest("PUT",
		uaaEndpoint+fmt.Sprintf("/oauth/clients/%s/secret", id),
		cm.config.AccessToken(), bytes.NewReader(body))
	if err != nil {
		return err
	}

	response := make(map[string]interface{})
	_, err = cm.uaaGateway.PerformRequestForJSONResponse(request, response)
	if err != nil {
		return err
	}
	return
}

// ChangeClientJwt - Add, update or delete the private_key_jwt trust configuration of a client
func (cm *ClientManager) ChangeClientJwt(req ClientJwtChangeRequest) (err error) {
	uaaEndpoint := cm.config.UaaEndpoint()
	if len(uaaEndpoint) == 0 {
		err = cfErrors.New("UAA endpoint missing from config file")
		return
	}

	body, err := json.Marshal(req)
	if err != nil {
		return
	}

	request, err := cm.uaaGateway.NewRequest("PUT",
		uaaEndpoint+fmt.Sprintf("/oauth/clients/%s/clientjwt", req.ClientID),
		cm.config.AccessToken(), bytes.NewReader(body))
	if err != nil {
		return err
	}

	response := make(map[string]interface{})
	_, err = cm.uaaGateway.PerformRequestForJSONResponse(request, &response)
	return
}

// FindByClientID -
func (cm *ClientManager) FindByClientID(clientID string) (client UAAClient, err error) {
	uaaEndpoint := cm.config.UaaEndpoint()
	if len(uaaEndpoint) == 0 {
		err = cfErrors.New("UAA endpoint missing from config file")
		return
	}

	filter := url.QueryEscape(fmt.Sprintf(`client_id Eq "%s"`, clientID))
	path := fmt.Sprintf("%s/oauth/clients?filter=%s", uaaEndpoint, filter)

	clientResourceList := &UAAClientResourceList{}
	err = cm.uaaGateway.GetResource(path, clientResourceList)

	if err == nil {
		if len(clientResourceList.Resources) > 0 {
			client = clientResourceList.Resources[0]
		} else {
			err = cfErrors.NewModelNotFoundError("Client", clientID)
		}
	}
	return
}
