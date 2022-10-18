package samlkeyfields

type SamlKeyField int64

const (
	Certificate SamlKeyField = iota
)

func (s SamlKeyField) String() string {
	switch s {
	case Certificate:
		return "certificate"
	}
	return "unknown"
}
