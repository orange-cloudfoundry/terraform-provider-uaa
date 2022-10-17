package fields

type GroupField int64

const (
	Description GroupField = iota
	DisplayName
	ZoneId
)

func (s GroupField) String() string {
	switch s {
	case Description:
		return "description"
	case DisplayName:
		return "display_name"
	case ZoneId:
		return "zone_id"
	}
	return "unknown"
}
