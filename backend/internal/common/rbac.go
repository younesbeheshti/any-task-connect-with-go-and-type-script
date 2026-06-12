package common

// Role represents a user's platform role.
type Role string

const (
	RoleAdmin     Role = "ADMIN"
	RoleRequester Role = "REQUESTER"
	RoleAgent     Role = "AGENT"
)

func (r Role) String() string {
	return string(r)
}

// IsValid reports whether the role is a known value.
func (r Role) IsValid() bool {
	switch r {
	case RoleAdmin, RoleRequester, RoleAgent:
		return true
	default:
		return false
	}
}

// APIRole returns the lowercase role string used in JWT/API responses.
func (r Role) APIRole() string {
	switch r {
	case RoleAdmin:
		return "admin"
	case RoleRequester:
		return "requester"
	case RoleAgent:
		return "agent"
	default:
		return ""
	}
}

// RoleFromAPI parses API role strings into domain roles.
func RoleFromAPI(s string) Role {
	switch s {
	case "admin":
		return RoleAdmin
	case "requester":
		return RoleRequester
	case "agent":
		return RoleAgent
	default:
		return ""
	}
}

// Permission represents a granular authorization action.
type Permission string

const (
	PermTaskCreate       Permission = "task:create"
	PermTaskRead         Permission = "task:read"
	PermTaskUpdate       Permission = "task:update"
	PermTaskCancel       Permission = "task:cancel"
	PermTaskVerify       Permission = "task:verify"
	PermApplicationApply Permission = "application:apply"
	PermApplicationRead  Permission = "application:read"
	PermApplicationAccept Permission = "application:accept"
	PermWalletRead       Permission = "wallet:read"
	PermWalletTopup      Permission = "wallet:topup"
	PermWalletWithdraw   Permission = "wallet:withdraw"
	PermChatRead         Permission = "chat:read"
	PermChatSend         Permission = "chat:send"
	PermReviewCreate     Permission = "review:create"
	PermAdminAll         Permission = "admin:*"
)

// RolePermissions maps roles to their granted permissions.
var RolePermissions = map[Role][]Permission{
	RoleAdmin: {PermAdminAll},
	RoleRequester: {
		PermTaskCreate, PermTaskRead, PermTaskUpdate, PermTaskCancel, PermTaskVerify,
		PermApplicationRead, PermApplicationAccept,
		PermWalletRead, PermWalletTopup, PermWalletWithdraw,
		PermChatRead, PermChatSend, PermReviewCreate,
	},
	RoleAgent: {
		PermTaskRead, PermApplicationApply, PermApplicationRead,
		PermWalletRead, PermWalletWithdraw,
		PermChatRead, PermChatSend, PermReviewCreate,
	},
}

// HasPermission checks if a role grants the given permission.
func HasPermission(role Role, perm Permission) bool {
	perms, ok := RolePermissions[role]
	if !ok {
		return false
	}
	for _, p := range perms {
		if p == PermAdminAll || p == perm {
			return true
		}
	}
	return false
}
