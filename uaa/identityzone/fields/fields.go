package fields

type IdentityZoneField int64

const (
	ClientSecretPolicy IdentityZoneField = iota
	Config
	CorsConfig
	Id
	IsActive
	Name
	SubDomain
)

func (s IdentityZoneField) String() string {
	switch s {
	case ClientSecretPolicy:
		return "client_secret_policy"
	case Config:
		return "config"
	case CorsConfig:
		return "cors_config"
	case Id:
		return "id"
	case IsActive:
		return "is_active"
	case Name:
		return "id"
	case SubDomain:
		return "sub_domain"
	}
	return "unknown"
}
