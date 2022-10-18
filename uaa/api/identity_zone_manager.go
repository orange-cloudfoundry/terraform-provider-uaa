package api

import (
	"code.cloudfoundry.org/cli/cf/configuration/coreconfig"
	"code.cloudfoundry.org/cli/cf/net"
	"errors"
	"fmt"
)

type IdentityZoneManager struct {
	log        *Logger
	config     coreconfig.Reader
	uaaGateway net.Gateway
}

func newIdentityZoneManager(config coreconfig.Reader, uaaGateway net.Gateway, logger *Logger) (izm *IdentityZoneManager, err error) {

	izm = &IdentityZoneManager{
		log:        logger,
		config:     config,
		uaaGateway: uaaGateway,
	}
	return
}

// CRUD methods

func (api *IdentityZoneManager) FindById(id string) (*IdentityZone, error) {

	uaaEndpoint := api.config.UaaEndpoint()
	if len(uaaEndpoint) == 0 {
		return nil, errors.New("UAA endpoint missing from config file")
	}

	path := fmt.Sprintf("%s/identity-zones/%s", uaaEndpoint, id)
	identityZone := &IdentityZone{}
	err := api.uaaGateway.GetResource(path, identityZone)
	if err != nil {
		return nil, err
	}

	return identityZone, nil
}

// DTOs

type IdentityZone struct {
	Id        string             `json:"id"`
	IsActive  bool               `json:"active"`
	Name      string             `json:"name,omitempty"`
	SubDomain string             `json:"subdomain,omitempty"`
	Config    IdentityZoneConfig `json:"config,omitempty"`
}

type IdentityZoneConfig struct {
	ClientSecretPolicy IdentityZoneClientSecretPolicy `json:"clientSecretPolicy"`
	TokenPolicy        IdentityZoneTokenPolicy        `json:"tokenPolicy"`
	Saml               IdentityZoneSamlConfig         `json:"samlConfig"`
}

type IdentityZoneClientSecretPolicy struct {
	MaxLength             int64 `json:"maxLength,omitempty"`
	MinLength             int64 `json:"minLength,omitempty"`
	MinUpperCaseCharacter int64 `json:"requireUpperCaseCharacter,omitempty"`
	MinLowerCaseCharacter int64 `json:"requireLowerCaseCharacter,omitempty"`
	MinDigit              int64 `json:"requireDigit,omitempty"`
	MinSpecialCharacter   int64 `json:"requireSpecialCharacter,omitempty"`
}

type IdentityZoneTokenPolicy struct {
	AccessTokenTtl       int64  `json:"accessTokenValidity,omitempty"`
	RefreshTokenTtl      int64  `json:"refreshTokenValidity,omitempty"`
	IsJwtRevocable       bool   `json:"jwtRevocable"`
	IsRefreshTokenUnique bool   `json:"refreshTokenUnique"`
	RefreshTokenFormat   string `json:"refreshTokenFormat,omitempty"`
	ActiveKeyId          string `json:"activeKeyId,omitempty"`
}

type IdentityZoneSamlConfig struct {
	ActiveKeyId              string                         `json:"activeKeyId,omitempty"`
	AssertionTtlSeconds      int64                          `json:"assertionTimeToLiveSeconds,omitempty"`
	Certificate              string                         `json:"certificate,omitempty"`
	DisableInResponseToCheck bool                           `json:"disableInResponseToCheck"`
	EntityId                 string                         `json:"entityID,omitempty"`
	IsAssertionSigned        bool                           `json:"assertionSigned"`
	IsRequestSigned          bool                           `json:"requestSigned"`
	Keys                     map[string]IdentityZoneSamlKey `json:"keys,omitempty"`
	WantAssertionSigned      bool                           `json:"wantAssertionSigned"`
	WantAuthnRequestSigned   bool                           `json:"wantAuthnRequestSigned"`
}

type IdentityZoneSamlKey struct {
	Certificate string `json:"certificate,omitempty"`
}
