package fields

type UserField int64

const (
	Email UserField = iota
	FamilyName
	GivenName
	Groups
	Name
	Origin
	Password
)

func (s UserField) String() string {
	switch s {
	case Email:
		return "email"
	case FamilyName:
		return "family_name"
	case GivenName:
		return "given_name"
	case Groups:
		return "groups"
	case Name:
		return "name"
	case Origin:
		return "origin"
	case Password:
		return "password"
	}
	return "unknown"
}
