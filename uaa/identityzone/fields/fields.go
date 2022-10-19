package fields

type IdentityZoneField int64

const (
	ClientSecretPolicy IdentityZoneField = iota
	CorsConfig
	Id
	IsActive
	LogoutRedirectUrl
	LogoutRedirectParam
	LogoutAllowedRedirectUrls
	Name
	SamlConfig
	SubDomain
	TokenPolicy
)

func (s IdentityZoneField) String() string {
	switch s {
	case ClientSecretPolicy:
		return "client_secret_policy"
	case CorsConfig:
		return "cors_config"
	case Id:
		return "id"
	case IsActive:
		return "is_active"
	case LogoutRedirectUrl:
		return "logout_redirect_url"
	case LogoutRedirectParam:
		return "logout_redirect_param"
	case LogoutAllowedRedirectUrls:
		return "logout_allowed_redirect_urls"
	case Name:
		return "id"
	case SamlConfig:
		return "saml_config"
	case SubDomain:
		return "sub_domain"
	case TokenPolicy:
		return "token_policy"
	}
	return "unknown"
}
