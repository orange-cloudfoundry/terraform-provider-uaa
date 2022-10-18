package tokenpolicyfields

type TokenPolicyField int64

const (
	AccessTokenTtl TokenPolicyField = iota
	ActiveKeyId
	IsJwtRevocable
	IsRefreshTokenUnique
	RefreshTokenFormat
	RefreshTokenTtl
)

func (s TokenPolicyField) String() string {
	switch s {
	case AccessTokenTtl:
		return "access_token_ttl"
	case ActiveKeyId:
		return "active_key_id"
	case IsJwtRevocable:
		return "is_jwt_revocable"
	case IsRefreshTokenUnique:
		return "is_refresh_token_unique"
	case RefreshTokenFormat:
		return "refresh_token_format"
	case RefreshTokenTtl:
		return "refresh_token_ttl"
	}
	return "unknown"
}
