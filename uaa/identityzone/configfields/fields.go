package configfields

type IdentityZoneConfigField int64

const (
	Saml IdentityZoneConfigField = iota
	TokenPolicy
)

func (s IdentityZoneConfigField) String() string {
	switch s {
	case Saml:
		return "saml"
	case TokenPolicy:
		return "token_policy"
	}
	return "unknown"
}
