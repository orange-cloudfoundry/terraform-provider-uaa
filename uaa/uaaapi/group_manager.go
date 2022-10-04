package uaaapi

import (
	"code.cloudfoundry.org/cli/cf/configuration/coreconfig"
	"code.cloudfoundry.org/cli/cf/errors"
	"code.cloudfoundry.org/cli/cf/net"
	"fmt"
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

func (gm *GroupManager) CreateGroup(displayName string, zoneId string) (group UAAGroup, err error) {

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
