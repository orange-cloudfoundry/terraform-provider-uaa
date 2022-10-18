package apiheaders

type ApiHeader int64

const (
	IfMatch ApiHeader = iota
	ZoneId
)

func (s ApiHeader) String() string {
	switch s {
	case IfMatch:
		return "If-Match"
	case ZoneId:
		return "X-Identity-Zone-Id"
	}
	return "unknown"
}
