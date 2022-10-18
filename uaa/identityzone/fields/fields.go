package fields

type IdentityZoneField int64

const (
	ClientSecretPolicy IdentityZoneField = iota
	CorsConfig
	Id
	IsActive
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
