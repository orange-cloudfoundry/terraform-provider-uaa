package fields

type ClientField int64

const (
	AccessTokenValidity ClientField = iota
	AllowProviders
	ApprovalsDeleted
	Authorities
	AuthorizedGrantTypes
	AutoApprove
	ClientId
	ClientSecret
	CreatedWith
	Name
	RedirectUri
	RefreshTokenValidity
	RequiredUserGroups
	ResourceIds
	Scope
	TokenSalt
)

func (s ClientField) String() string {
	switch s {
	case AccessTokenValidity:
		return "access_token_validity"
	case AllowProviders:
		return "allowed_providers"
	case ApprovalsDeleted:
		return "approvals_deleted"
	case Authorities:
		return "authorities"
	case AuthorizedGrantTypes:
		return "authorized_grant_types"
	case AutoApprove:
		return "auto_approve"
	case ClientId:
		return "client_id"
	case ClientSecret:
		return "client_secret"
	case CreatedWith:
		return "created_with"
	case Name:
		return "name"
	case RedirectUri:
		return "redirect_uri"
	case RefreshTokenValidity:
		return "refresh_token_validity"
	case RequiredUserGroups:
		return "required_user_groups"
	case ResourceIds:
		return "resource_ids"
	case Scope:
		return "scope"
	case TokenSalt:
		return "token_salt"
	}
	return "unknown"
}
