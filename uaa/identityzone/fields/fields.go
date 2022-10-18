package fields

type IdentityZoneField int64

const (
	Config IdentityZoneField = iota
	ClientSecretPolicy
	Id
	IsActive
	Name
	SubDomain
)

func (s IdentityZoneField) String() string {
	switch s {
	case Config:
		return "config"
	case ClientSecretPolicy:
		return "client_secret_policy"
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
