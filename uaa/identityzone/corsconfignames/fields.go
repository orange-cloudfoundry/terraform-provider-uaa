package corsconfignames

type CorsConfigName int64

const (
	Default CorsConfigName = iota
	Xhr
)

func (s CorsConfigName) String() string {
	switch s {
	case Default:
		return "default"
	case Xhr:
		return "xhr"
	}
	return "unknown"
}
