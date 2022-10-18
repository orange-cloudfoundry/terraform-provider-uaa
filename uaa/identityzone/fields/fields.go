package fields

type IdentityZoneField int64

const (
	Config IdentityZoneField = iota
	Id
	IsActive
	Name
	SubDomain
)

func (s IdentityZoneField) String() string {
	switch s {
	case Config:
		return "config"
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
