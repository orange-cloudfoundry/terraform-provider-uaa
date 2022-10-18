package configfields

type IdentityZoneConfigField int64

const (
	TokenPolicy IdentityZoneConfigField = iota
)

func (s IdentityZoneConfigField) String() string {
	switch s {
	case TokenPolicy:
		return "token_policy"
	}
	return "unknown"
}
