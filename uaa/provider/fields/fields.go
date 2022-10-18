package fields

type ProviderField int64

const (
	AuthEndpoint ProviderField = iota
	CaCert
	ClientId
	ClientSecret
	LoginEndpoint
	SkipSslValidation
)

func (s ProviderField) String() string {
	switch s {
	case AuthEndpoint:
		return "auth_endpoint"
	case CaCert:
		return "ca_cert"
	case ClientId:
		return "client_id"
	case ClientSecret:
		return "client_secret"
	case LoginEndpoint:
		return "login_endpoint"
	case SkipSslValidation:
		return "skip_ssl_validation"
	}
	return "unknown"
}
