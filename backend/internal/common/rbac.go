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

// Permission represents a granular authorization action (Phase 3 model).
type Permission string

const (
	PermTaskCreate        Permission = "task:create"
	PermTaskView          Permission = "task:view"
	PermTaskEdit          Permission = "task:edit"
	PermTaskDelete        Permission = "task:delete"
	PermApplicationCreate Permission = "application:create"
	PermApplicationAccept Permission = "application:accept"
	PermWalletView        Permission = "wallet:view"
	PermWalletWithdraw    Permission = "wallet:withdraw"
	PermUserManage        Permission = "user:manage"
	PermAdminDashboard    Permission = "admin:dashboard"
	PermReviewCreate      Permission = "review:create"
	PermChatSend          Permission = "chat:send"
	PermNotificationView  Permission = "notification:view"
)

// AllPermissions lists every known permission.
var AllPermissions = []Permission{
	PermTaskCreate, PermTaskView, PermTaskEdit, PermTaskDelete,
	PermApplicationCreate, PermApplicationAccept,
	PermWalletView, PermWalletWithdraw,
	PermUserManage, PermAdminDashboard,
	PermReviewCreate, PermChatSend, PermNotificationView,
}

// RolePermissions maps roles to their granted permissions.
var RolePermissions = map[Role][]Permission{
	RoleAdmin: AllPermissions,
	RoleRequester: {
		PermTaskCreate, PermTaskView, PermTaskEdit, PermTaskDelete,
		PermApplicationAccept,
		PermWalletView, PermWalletWithdraw,
		PermReviewCreate, PermChatSend, PermNotificationView,
	},
	RoleAgent: {
		PermTaskView, PermApplicationCreate,
		PermWalletView, PermWalletWithdraw,
		PermReviewCreate, PermChatSend, PermNotificationView,
	},
}

// HasPermission checks if a role grants the given permission.
func HasPermission(role Role, perm Permission) bool {
	perms, ok := RolePermissions[role]
	if !ok {
		return false
	}
	for _, p := range perms {
		if p == perm {
			return true
		}
	}
	return false
}

// PermissionsForRole returns permission strings for API responses.
func PermissionsForRole(role Role) []string {
	perms := RolePermissions[role]
	out := make([]string, len(perms))
	for i, p := range perms {
		out[i] = string(p)
	}
	return out
}
