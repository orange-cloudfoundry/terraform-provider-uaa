package api

import (
	"bytes"
	"code.cloudfoundry.org/cli/cf/configuration/coreconfig"
	"code.cloudfoundry.org/cli/cf/errors"
	"code.cloudfoundry.org/cli/cf/net"
	"encoding/json"
	"fmt"
	apiheaders "github.com/jlpospisil/terraform-provider-uaa/uaa/api/headers"
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

func (gm *GroupManager) CreateGroup(displayName string, description string, zoneId string) (group *UAAGroup, err error) {
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

	request, err := gm.uaaGateway.NewRequest(
		"POST",
		fmt.Sprintf("%s/Groups", uaaEndpoint),
		gm.config.AccessToken(),
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, err
	}
	request.HTTPReq.Header.Set(apiheaders.ZoneId.String(), zoneId)

	group = &UAAGroup{}
	_, err = gm.uaaGateway.PerformRequestForJSONResponse(request, group)

	switch httpErr := err.(type) {
	case errors.HTTPError:
		if httpErr.StatusCode() == http.StatusConflict {
			err = errors.NewModelAlreadyExistsError("group", displayName)
		}
	}
	return
}

func (gm *GroupManager) GetGroup(id, zoneId string) (group *UAAGroup, err error) {
	uaaEndpoint := gm.config.UaaEndpoint()
	if len(uaaEndpoint) == 0 {
		err = errors.New("UAA endpoint missing from config file")
		return
	}

	request, err := gm.uaaGateway.NewRequest(
		"GET",
		fmt.Sprintf("%s/Groups/%s", uaaEndpoint, id),
		gm.config.AccessToken(),
		nil,
	)
	if err != nil {
		return nil, err
	}
	request.HTTPReq.Header.Set(apiheaders.ZoneId.String(), zoneId)

	group = &UAAGroup{}
	_, err = gm.uaaGateway.PerformRequestForJSONResponse(request, group)

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
	}

	body, err := json.Marshal(groupResource)
	if err != nil {
		return
	}

	request, err := gm.uaaGateway.NewRequest(
		"PUT",
		fmt.Sprintf("%s/Groups/%s", uaaEndpoint, id),
		gm.config.AccessToken(),
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, err
	}
	request.HTTPReq.Header.Set(apiheaders.IfMatch.String(), "*")
	request.HTTPReq.Header.Set(apiheaders.ZoneId.String(), zoneId)

	group = &UAAGroup{}
	_, err = gm.uaaGateway.PerformRequestForJSONResponse(request, group)

	return
}

func (gm *GroupManager) DeleteGroup(id, zoneId string) (err error) {

	uaaEndpoint := gm.config.UaaEndpoint()
	if len(uaaEndpoint) == 0 {
		err = errors.New("UAA endpoint missing from config file")
		return
	}

	request, err := gm.uaaGateway.NewRequest(
		"DELETE",
		fmt.Sprintf("%s/Groups/%s", uaaEndpoint, id),
		gm.config.AccessToken(),
		nil,
	)
	if err != nil {
		return err
	}
	request.HTTPReq.Header.Set(apiheaders.ZoneId.String(), zoneId)
	_, err = gm.uaaGateway.PerformRequest(request)

	return
}

func (gm *GroupManager) FindByDisplayName(displayName, zoneId string) (group *UAAGroup, err error) {

	uaaEndpoint := gm.config.UaaEndpoint()
	if len(uaaEndpoint) == 0 {
		err = errors.New("UAA endpoint missing from config file")
		return
	}

	displayNameFilter := url.QueryEscape(fmt.Sprintf(`displayName Eq "%s"`, displayName))
	path := fmt.Sprintf("%s/Groups?filter=%s", uaaEndpoint, displayNameFilter)

	request, err := gm.uaaGateway.NewRequest(
		"GET",
		path,
		gm.config.AccessToken(),
		nil,
	)
	if err != nil {
		return nil, err
	}
	request.HTTPReq.Header.Set(apiheaders.ZoneId.String(), zoneId)

	groupResourceList := &UAAGroupResourceList{}
	_, err = gm.uaaGateway.PerformRequestForJSONResponse(request, groupResourceList)
	if err != nil {
		return nil, err
	}

	if err == nil {
		if len(groupResourceList.Resources) > 0 {
			group = &groupResourceList.Resources[0]
		} else {
			err = errors.NewModelNotFoundError("Group", displayName)
		}
	}
	return
}
