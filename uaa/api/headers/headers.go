package apiheaders

type ApiHeader int64

const (
	ZoneId ApiHeader = iota
)

func (s ApiHeader) String() string {
	switch s {
	case ZoneId:
		return "X-Identity-Zone-Id"
	}
	return "unknown"
}
