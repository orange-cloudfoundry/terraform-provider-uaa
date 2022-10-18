package corsconfigfields

type CorsConfigField int64

const (
	AllowedOrigins CorsConfigField = iota
	AllowedOriginPatterns
	AllowedUris
	AllowedUriPatterns
	AllowedHeaders
	AllowedMethods
	AllowedCredentials
	Name
	MaxAge
)

func (s CorsConfigField) String() string {
	switch s {
	case AllowedOrigins:
		return "allowed_origins"
	case AllowedOriginPatterns:
		return "allowed_origin_patterns"
	case AllowedUris:
		return "allowed_uris"
	case AllowedUriPatterns:
		return "allowed_uri_patterns"
	case AllowedHeaders:
		return "allowed_headers"
	case AllowedMethods:
		return "allowed_methods"
	case AllowedCredentials:
		return "allowed_credentials"
	case Name:
		return "name"
	case MaxAge:
		return "max_age"
	}
	return "unknown"
}
