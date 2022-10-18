package samlconfigfields

type SamlConfigField int64

const (
	ActiveKeyId SamlConfigField = iota
	AssertionTtlSeconds
	Certificate
	DisableInResponseToCheck
	EntityId
	IsAssertionSigned
	IsRequestSigned
	Key
	WantAssertionSigned
	WantAuthRequestSigned
)

func (s SamlConfigField) String() string {
	switch s {
	case ActiveKeyId:
		return "active_key_id"
	case AssertionTtlSeconds:
		return "assertion_ttl_seconds"
	case Certificate:
		return "certificate"
	case DisableInResponseToCheck:
		return "disable_in_response_to_check"
	case EntityId:
		return "entity_id"

	case IsAssertionSigned:
		return "is_assertion_signed"
	case IsRequestSigned:
		return "is_request_signed"
	case Key:
		return "key"
	case WantAssertionSigned:
		return "want_assertion_signed"
	case WantAuthRequestSigned:
		return "want_authn_request_signed"
	}
	return "unknown"
}
