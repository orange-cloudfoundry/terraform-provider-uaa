package fields

type IdentityZoneField int64

const (
	AccountChooserEnabled IdentityZoneField = iota
	ClientSecretPolicy
	CorsConfig
	DefaultUserGroups
	Id
	InputPrompts
	IsActive
	IdpDiscoveryEnabled
	IssuerUrl
	LogoutRedirectUrl
	LogoutRedirectParam
	LogoutAllowedRedirectUrls
	HomeRedirectUrl
	MfaEnabled
	MfaIdentityProviders
	Name
	SamlConfig
	SelfServeEnabled
	SelfServeSignupUrl
	SelfServePasswordResetUrl
	SubDomain
	TokenPolicy
)

func (s IdentityZoneField) String() string {
	switch s {
	case AccountChooserEnabled:
		return "account_chooser_enabled"
	case ClientSecretPolicy:
		return "client_secret_policy"
	case CorsConfig:
		return "cors_config"
	case DefaultUserGroups:
		return "default_user_groups"
	case HomeRedirectUrl:
		return "home_redirect_url"
	case Id:
		return "id"
	case InputPrompts:
		return "input_prompt"
	case IdpDiscoveryEnabled:
		return "idp_discovery_enabled"
	case IssuerUrl:
		return "issuer_url"
	case IsActive:
		return "is_active"
	case LogoutRedirectUrl:
		return "logout_redirect_url"
	case LogoutRedirectParam:
		return "logout_redirect_param"
	case LogoutAllowedRedirectUrls:
		return "logout_allowed_redirect_urls"
	case MfaEnabled:
		return "mfa_enabled"
	case MfaIdentityProviders:
		return "mfa_identity_providers"
	case Name:
		return "id"
	case SamlConfig:
		return "saml_config"
	case SelfServeEnabled:
		return "self_serve_enabled"
	case SelfServeSignupUrl:
		return "self_serve_signup_url"
	case SelfServePasswordResetUrl:
		return "self_serve_pw_reset_url"
	case SubDomain:
		return "sub_domain"
	case TokenPolicy:
		return "token_policy"
	}
	return "unknown"
}
