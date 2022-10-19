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
	AccountChooserEnabled bool                           `json:"accountChooserEnabled"`
	ClientSecretPolicy    IdentityZoneClientSecretPolicy `json:"clientSecretPolicy,omitempty"`
	CorsPolicy            IdentityZoneCorsPolicy         `json:"corsPolicy,omitempty"`
	IdpDiscoveryEnabled   bool                           `json:"idpDiscoveryEnabled"`
	InputPrompts          []InputPrompt                  `json:"prompts,omitempty"'`
	IssuerUrl             string                         `json:"issuer,omitempty"`
	Links                 IdentityZoneLinks              `json:"links,omitempty"`
	MfaConfig             MfaConfig                      `json:"MfaConfig,omitempty"`
	TokenPolicy           IdentityZoneTokenPolicy        `json:"tokenPolicy,omitempty"`
	Saml                  IdentityZoneSamlConfig         `json:"samlConfig,omitempty"`
	UserConfig            UserConfig                     `json:"userConfig,omitempty"`
}

type IdentityZoneClientSecretPolicy struct {
	MaxLength             int64 `json:"maxLength,omitempty"`
	MinLength             int64 `json:"minLength,omitempty"`
	MinUpperCaseCharacter int64 `json:"requireUpperCaseCharacter,omitempty"`
	MinLowerCaseCharacter int64 `json:"requireLowerCaseCharacter,omitempty"`
	MinDigit              int64 `json:"requireDigit,omitempty"`
	MinSpecialCharacter   int64 `json:"requireSpecialCharacter,omitempty"`
}

type IdentityZoneCorsPolicy struct {
	DefaultConfiguration IdentityZoneCorsConfig `json:"defaultConfiguration,omitempty"`
	XhrConfiguration     IdentityZoneCorsConfig `json:"xhrConfiguration,omitempty"`
}

type IdentityZoneCorsConfig struct {
	AllowedOrigins        []string `json:"allowedOrigins,omitempty"`
	AllowedOriginPatterns []string `json:"allowedOriginPatterns,omitempty"`
	AllowedUris           []string `json:"allowedUris,omitempty"`
	AllowedUriPatterns    []string `json:"allowedUriPatterns,omitempty"`
	AllowedHeaders        []string `json:"allowedHeaders,omitempty"`
	AllowedMethods        []string `json:"allowedMethods,omitempty"`
	AllowedCredentials    bool     `json:"allowedCredentials"`
	MaxAge                int64    `json:"maxAge,omitempty"`
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

type IdentityZoneLinks struct {
	HomeRedirect string                  `json:"homeRedirect,omitempty"`
	Logout       IdentityZoneLogoutLinks `json:"logout,omitempty"`
	SelfService  SelfServiceLinks        `json:"selfService,omitempty"`
}

type IdentityZoneLogoutLinks struct {
	RedirectUrl           string   `json:"redirectUrl,omitempty"`
	RedirectParameterName string   `json:"redirectParameterName,omitempty"`
	AllowedRedirectUrls   []string `json:"whitelist"`
}

type SelfServiceLinks struct {
	Enabled          bool   `json:"selfServiceLinksEnabled"`
	SignupUrl        string `json:"signup,omitempty"`
	PasswordResetUrl string `json:"signup,passwd"`
}

type InputPrompt struct {
	Name  string `json:"name,omitempty"`
	Type  string `json:"type,omitempty"`
	Value string `json:"text,omitempty"`
}

type UserConfig struct {
	DefaultGroups []string `json:"defaultGroups,omitempty"`
}

type MfaConfig struct {
	IsEnabled         bool     `json:"enabled"`
	IdentityProviders []string `json:"identityProviders,omitempty"`
}
