package uaaapi

import (
	"bytes"
	"code.cloudfoundry.org/cli/cf/configuration/coreconfig"
	"code.cloudfoundry.org/cli/cf/errors"
	"code.cloudfoundry.org/cli/cf/net"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type GroupManager struct {
	log        *Logger
	config     coreconfig.Reader
	uaaGateway net.Gateway
}

type UAAGroup struct {
	ID          string `json:"id,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
	Description string `json:"description,omitempty"`
	ZoneId      string `json:"zoneId,omitempty"`
}

type UAAGroupResourceList struct {
	Resources []UAAGroup `json:"resources"`
}

func newGroupManager(config coreconfig.Reader, uaaGateway net.Gateway, logger *Logger) (gm *GroupManager, err error) {
	gm = &GroupManager{
		log:        logger,
		config:     config,
		uaaGateway: uaaGateway,
	}
	return
}

func (gm *GroupManager) CreateGroup(displayName string, description string, zoneId string) (group UAAGroup, err error) {
	uaaEndpoint := gm.config.UaaEndpoint()
	if len(uaaEndpoint) == 0 {
		err = errors.New("UAA endpoint missing from config file")
		return
	}

	groupResource := UAAGroup{
		DisplayName: displayName,
		Description: description,
		ZoneId:      zoneId,
	}

	body, err := json.Marshal(groupResource)
	if err != nil {
		return
	}

	group = UAAGroup{}
	err = gm.uaaGateway.CreateResource(uaaEndpoint, "/Groups", bytes.NewReader(body), &group)
	switch httpErr := err.(type) {
	case errors.HTTPError:
		if httpErr.StatusCode() == http.StatusConflict {
			err = errors.NewModelAlreadyExistsError("group", displayName)
		}
	}
	return
}

func (gm *GroupManager) GetGroup(id string) (group *UAAGroup, err error) {
	uaaEndpoint := gm.config.UaaEndpoint()
	if len(uaaEndpoint) == 0 {
		err = errors.New("UAA endpoint missing from config file")
		return
	}

	group = &UAAGroup{}
	err = gm.uaaGateway.GetResource(
		fmt.Sprintf("%s/Groups/%s", uaaEndpoint, id),
		group)

	return
}

func (gm *GroupManager) FindByDisplayName(displayName string) (group UAAGroup, err error) {

	uaaEndpoint := gm.config.UaaEndpoint()
	if len(uaaEndpoint) == 0 {
		err = errors.New("UAA endpoint missing from config file")
		return
	}

	displayNameFilter := url.QueryEscape(fmt.Sprintf(`displayName Eq "%s"`, displayName))
	path := fmt.Sprintf("%s/Groups?filter=%s", uaaEndpoint, displayNameFilter)

	groupResourceList := &UAAGroupResourceList{}
	err = gm.uaaGateway.GetResource(path, groupResourceList)

	if err == nil {
		if len(groupResourceList.Resources) > 0 {
			group = groupResourceList.Resources[0]
		} else {
			err = errors.NewModelNotFoundError("Group", displayName)
		}
	}
	return
}

func (gm *GroupManager) UpdateGroup(id, displayName, description, zoneId string) (group *UAAGroup, err error) {

	uaaEndpoint := gm.config.UaaEndpoint()
	if len(uaaEndpoint) == 0 {
		err = errors.New("UAA endpoint missing from config file")
		return
	}

	groupResource := UAAGroup{
		DisplayName: displayName,
		Description: description,
		ZoneId:      zoneId,
	}

	body, err := json.Marshal(groupResource)
	if err != nil {
		return
	}

	request, err := gm.uaaGateway.NewRequest("PUT",
		fmt.Sprintf("%s/Groups/%s", uaaEndpoint, id),
		gm.config.AccessToken(), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	request.HTTPReq.Header.Set("If-Match", "*")

	group = &UAAGroup{}
	_, err = gm.uaaGateway.PerformRequestForJSONResponse(request, group)

	return
}

func (gm *GroupManager) DeleteGroup(id string) (err error) {

	uaaEndpoint := gm.config.UaaEndpoint()
	if len(uaaEndpoint) == 0 {
		err = errors.New("UAA endpoint missing from config file")
		return
	}
	err = gm.uaaGateway.DeleteResource(uaaEndpoint, fmt.Sprintf("/Groups/%s", id))
	return
}
