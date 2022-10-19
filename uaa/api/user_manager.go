package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	apiheaders "github.com/jlpospisil/terraform-provider-uaa/uaa/api/headers"
	"net/http"
	"net/url"

	"code.cloudfoundry.org/cli/cf/configuration/coreconfig"
	"code.cloudfoundry.org/cli/cf/errors"
	"code.cloudfoundry.org/cli/cf/net"
)

// UserManager -
type UserManager struct {
	log *Logger

	config     coreconfig.Reader
	uaaGateway net.Gateway

	clientToken string

	groupMap      map[string]string
	defaultGroups map[string]byte
}

// UAAUser -
type UAAUser struct {
	ID       string `json:"id,omitempty"`
	Username string `json:"userName,omitempty"`
	Password string `json:"password,omitempty"`
	Origin   string `json:"origin,omitempty"`

	Name   UAAUserName    `json:"name,omitempty"`
	Emails []UAAUserEmail `json:"emails,omitempty"`
	Groups []UAAUserGroup `json:"groups,omitempty"`
}

// UAAUserResourceList -
type UAAUserResourceList struct {
	Resources []UAAUser `json:"resources"`
}

// UAAUserEmail -
type UAAUserEmail struct {
	Value string `json:"value"`
}

// UAAUserName -
type UAAUserName struct {
	GivenName  string `json:"givenName"`
	FamilyName string `json:"familyName"`
}

// UAAUserGroup -
type UAAUserGroup struct {
	Value   string `json:"value"`
	Display string `json:"display"`
	Type    string `json:"type"`
}

// NewUserManager -
func newUserManager(config coreconfig.Reader, uaaGateway net.Gateway, logger *Logger) (um *UserManager, err error) {

	um = &UserManager{
		log: logger,

		config:        config,
		uaaGateway:    uaaGateway,
		groupMap:      make(map[string]string),
		defaultGroups: make(map[string]byte),
	}
	return
}

func (um *UserManager) loadGroups() (err error) {

	uaaEndpoint := um.config.UaaEndpoint()
	if len(uaaEndpoint) == 0 {
		err = errors.New("UAA endpoint missing from config file")
		return
	}

	// Retrieve all groups
	groupList := &UAAGroupResourceList{}
	err = um.uaaGateway.GetResource(
		fmt.Sprintf("%s/Groups", uaaEndpoint),
		groupList)
	if err != nil {
		return
	}
	for _, r := range groupList.Resources {
		um.groupMap[r.DisplayName] = r.ID
	}

	// Retrieve default scope/groups for a new user by creating
	// a dummy user and extracting the default scope of that user
	username, err := newUUID()
	if err != nil {
		return
	}
	userResource := UAAUser{
		Username: username,
		Password: "password",
		Origin:   "uaa",
		Emails:   []UAAUserEmail{{Value: "email@domain.com"}},
	}
	body, err := json.Marshal(userResource)
	if err != nil {
		return
	}
	user := &UAAUser{}
	err = um.uaaGateway.CreateResource(uaaEndpoint, "/Users", bytes.NewReader(body), user)
	if err != nil {
		return err
	}
	err = um.uaaGateway.DeleteResource(uaaEndpoint, fmt.Sprintf("/Users/%s", user.ID))
	if err != nil {
		return err
	}
	for _, g := range user.Groups {
		um.defaultGroups[g.Display] = 1
	}

	return
}

// IsDefaultGroup -
func (um *UserManager) IsDefaultGroup(group string) (ok bool) {
	_, ok = um.defaultGroups[group]
	return
}

// GetUser -
func (um *UserManager) GetUser(id string) (user *UAAUser, err error) {

	uaaEndpoint := um.config.UaaEndpoint()
	if len(uaaEndpoint) == 0 {
		err = errors.New("UAA endpoint missing from config file")
		return
	}

	user = &UAAUser{}
	err = um.uaaGateway.GetResource(
		fmt.Sprintf("%s/Users/%s", uaaEndpoint, id),
		user)

	return
}

// CreateUser -
func (um *UserManager) CreateUser(
	username, password, origin, givenName, familyName, email string) (user UAAUser, err error) {

	uaaEndpoint := um.config.UaaEndpoint()
	if len(uaaEndpoint) == 0 {
		err = errors.New("UAA endpoint missing from config file")
		return
	}

	userResource := UAAUser{
		Username: username,
		Password: password,
		Origin:   origin,
		Name: UAAUserName{
			GivenName:  givenName,
			FamilyName: familyName,
		},
	}
	if len(email) > 0 {
		userResource.Emails = append(userResource.Emails, UAAUserEmail{email})
	} else {
		userResource.Emails = append(userResource.Emails, UAAUserEmail{username})
	}

	body, err := json.Marshal(userResource)
	if err != nil {
		return
	}

	user = UAAUser{}
	err = um.uaaGateway.CreateResource(uaaEndpoint, "/Users", bytes.NewReader(body), &user)
	switch httpErr := err.(type) {
	case errors.HTTPError:
		if httpErr.StatusCode() == http.StatusConflict {
			err = errors.NewModelAlreadyExistsError("user", username)
		}
	}
	return
}

// UpdateUser -
func (um *UserManager) UpdateUser(
	id, username, givenName, familyName, email string) (user *UAAUser, err error) {

	uaaEndpoint := um.config.UaaEndpoint()
	if len(uaaEndpoint) == 0 {
		err = errors.New("UAA endpoint missing from config file")
		return
	}

	userResource := UAAUser{
		Username: username,
		Name: UAAUserName{
			GivenName:  givenName,
			FamilyName: familyName,
		},
	}
	if len(email) > 0 {
		userResource.Emails = append(userResource.Emails, UAAUserEmail{email})
	} else {
		userResource.Emails = append(userResource.Emails, UAAUserEmail{username})
	}

	body, err := json.Marshal(userResource)
	if err != nil {
		return
	}

	request, err := um.uaaGateway.NewRequest("PUT",
		fmt.Sprintf("%s/Users/%s", uaaEndpoint, id),
		um.config.AccessToken(), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	request.HTTPReq.Header.Set(apiheaders.IfMatch.String(), "*")

	user = &UAAUser{}
	_, err = um.uaaGateway.PerformRequestForJSONResponse(request, user)

	return
}

// DeleteUser -
func (um *UserManager) DeleteUser(id string) (err error) {

	uaaEndpoint := um.config.UaaEndpoint()
	if len(uaaEndpoint) == 0 {
		err = errors.New("UAA endpoint missing from config file")
		return
	}
	err = um.uaaGateway.DeleteResource(uaaEndpoint, fmt.Sprintf("/Users/%s", id))
	return
}

// ChangePassword -
func (um *UserManager) ChangePassword(
	id, oldPassword, newPassword string) (err error) {

	uaaEndpoint := um.config.UaaEndpoint()
	if len(uaaEndpoint) == 0 {
		err = errors.New("UAA endpoint missing from config file")
		return
	}

	body, err := json.Marshal(map[string]string{
		"oldPassword": oldPassword,
		"password":    newPassword,
	})
	if err != nil {
		return
	}

	request, err := um.uaaGateway.NewRequest("PUT",
		uaaEndpoint+fmt.Sprintf("/Users/%s/password", id),
		um.config.AccessToken(), bytes.NewReader(body))
	if err != nil {
		return err
	}
	request.HTTPReq.Header.Set("Authorization", um.clientToken)

	response := make(map[string]interface{})
	_, err = um.uaaGateway.PerformRequestForJSONResponse(request, response)
	if err != nil {
		return err
	}
	return
}

// UpdateRoles -
func (um *UserManager) UpdateRoles(
	id string, scopesToDelete, scopesToAdd []string, origin string) (err error) {

	uaaEndpoint := um.config.UaaEndpoint()
	if len(uaaEndpoint) == 0 {
		err = errors.New("UAA endpoint missing from config file")
		return
	}

	for _, s := range scopesToDelete {
		roleID := um.groupMap[s]
		err = um.uaaGateway.DeleteResource(uaaEndpoint,
			fmt.Sprintf("/Groups/%s/members/%s", roleID, id))
	}
	for _, s := range scopesToAdd {
		roleID, exists := um.groupMap[s]
		if !exists {
			err = fmt.Errorf("Group '%s' was not found", s)
			return
		}

		var body []byte
		body, err = json.Marshal(map[string]string{
			"origin": origin,
			"type":   "USER",
			"value":  id,
		})
		if err != nil {
			return
		}

		response := make(map[string]interface{})
		err = um.uaaGateway.CreateResource(uaaEndpoint,
			fmt.Sprintf("/Groups/%s/members", roleID),
			bytes.NewReader(body), &response)
		if err != nil {
			return
		}
	}

	return
}

// FindByUsername -
func (um *UserManager) FindByUsername(username string) (user UAAUser, err error) {

	uaaEndpoint := um.config.UaaEndpoint()
	if len(uaaEndpoint) == 0 {
		err = errors.New("UAA endpoint missing from config file")
		return
	}

	usernameFilter := url.QueryEscape(fmt.Sprintf(`userName Eq "%s"`, username))
	path := fmt.Sprintf("%s/Users?filter=%s", uaaEndpoint, usernameFilter)

	userResourceList := &UAAUserResourceList{}
	err = um.uaaGateway.GetResource(path, userResourceList)

	if err == nil {
		if len(userResourceList.Resources) > 0 {
			user = userResourceList.Resources[0]
		} else {
			err = errors.NewModelNotFoundError("User", username)
		}
	}
	return
}
