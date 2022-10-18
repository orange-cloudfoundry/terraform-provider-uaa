package configfields

type IdentityZoneConfigField int64

const (
	CorsConfig IdentityZoneConfigField = iota
	Saml
	TokenPolicy
)

func (s IdentityZoneConfigField) String() string {
	switch s {
	case CorsConfig:
		return "cors_config"
	case Saml:
		return "saml"
	case TokenPolicy:
		return "token_policy"
	}
	return "unknown"
}
