package uaaapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"code.cloudfoundry.org/cli/cf/configuration/coreconfig"
	"code.cloudfoundry.org/cli/cf/errors"
	"code.cloudfoundry.org/cli/cf/net"
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
	ClientID             string   `json:"client_id,omitempty"`
	ClientSecret         string   `json:"client_secret,omitempty"`
	AuthorizedGrantTypes []string `json:"authorized_grant_types,omitempty"`
	RedirectURI          []string `json:"redirect_uri,omitempty"`
	Scope                []string `json:"scope,omitempty"`
	ResourceIds          []string `json:"resource_ids,omitempty"`
	Authorities          []string `json:"authorities,omitempty"`
	Autoapprove          []string `json:"autoapprove,omitempty"`
	AccessTokenValidity  int      `json:"access_token_validity,omitempty"`
	RefreshTokenValidity int      `json:"refresh_token_validity,omitempty"`
	Allowedproviders     []string `json:"allowedproviders,omitempty"`
	Name                 string   `json:"name,omitempty"`
	TokenSalt            string   `json:"token_salt,omitempty"`
	CreatedWith          string   `json:"createdwith,omitempty"`
	ApprovalsDeleted     bool     `json:"approvals_deleted,omitempty"`
	RequiredUserGroups   []string `json:"required_user_groups,omitempty"`
	LastModified         int64    `json:"lastModified,omitempty"`
}

// UAAClientResourceList -
type UAAClientResourceList struct {
	Resources []UAAClient `json:"resources"`
}

func (c *UAAClient) HasDefaultScope() bool {
	return len(c.Scope) == 1 && c.Scope[0] == "uaa.none"
}

func (c *UAAClient) HasDefaultAuthorites() bool {
	return len(c.Authorities) == 1 && c.Authorities[0] == "uaa.none"
}

func (c *UAAClient) HasDefaultResourceIds() bool {
	return len(c.ResourceIds) == 1 && c.ResourceIds[0] == "none"
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
		err = errors.New("UAA endpoint missing from config file")
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

// CreateClient -
func (cm *ClientManager) Create(nCli UAAClient) (client UAAClient, err error) {
	uaaEndpoint := cm.config.UaaEndpoint()
	if len(uaaEndpoint) == 0 {
		err = errors.New("UAA endpoint missing from config file")
		return
	}

	body, err := json.Marshal(nCli)
	if err != nil {
		return
	}

	client = nCli
	err = cm.uaaGateway.CreateResource(uaaEndpoint, "/oauth/clients", bytes.NewReader(body), &client)
	switch httpErr := err.(type) {
	case errors.HTTPError:
		if httpErr.StatusCode() == http.StatusConflict {
			err = errors.NewModelAlreadyExistsError("client", nCli.ClientID)
		}
	}
	return
}

// UpdateClient -
func (cm *ClientManager) UpdateClient(nCli *UAAClient) (client UAAClient, err error) {
	uaaEndpoint := cm.config.UaaEndpoint()
	if len(uaaEndpoint) == 0 {
		err = errors.New("UAA endpoint missing from config file")
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
		err = errors.New("UAA endpoint missing from config file")
		return
	}
	err = cm.uaaGateway.DeleteResource(uaaEndpoint, fmt.Sprintf("/oauth/clients/%s", id))
	return
}

// ChangePassword -
func (cm *ClientManager) ChangeSecret(id, oldSecret, newSecret string) (err error) {
	uaaEndpoint := cm.config.UaaEndpoint()
	if len(uaaEndpoint) == 0 {
		err = errors.New("UAA endpoint missing from config file")
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

// FindByClientID -
func (cm *ClientManager) FindByClientID(clientID string) (client UAAClient, err error) {
	uaaEndpoint := cm.config.UaaEndpoint()
	if len(uaaEndpoint) == 0 {
		err = errors.New("UAA endpoint missing from config file")
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
			err = errors.NewModelNotFoundError("Client", clientID)
		}
	}
	return
}
