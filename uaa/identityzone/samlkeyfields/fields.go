package samlkeyfields

type SamlKeyField int64

const (
	Certificate SamlKeyField = iota
	Name
)

func (s SamlKeyField) String() string {
	switch s {
	case Certificate:
		return "certificate"
	case Name:
		return "name"
	}
	return "unknown"
}
