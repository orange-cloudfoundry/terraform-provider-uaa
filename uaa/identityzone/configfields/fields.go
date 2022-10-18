package configfields

type IdentityZoneConfigField int64

const (
	ClientSecretPolicy IdentityZoneConfigField = iota
	CorsConfig
	Saml
	TokenPolicy
)

func (s IdentityZoneConfigField) String() string {
	switch s {
	case ClientSecretPolicy:
		return "client_secret_policy"
	case CorsConfig:
		return "cors_config"
	case Saml:
		return "saml"
	case TokenPolicy:
		return "token_policy"
	}
	return "unknown"
}
