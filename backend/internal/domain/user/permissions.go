package user

// Only implemented workspace/account boundaries are defined here. Future
// mutations require their own reviewed permission and transactional authority.
type Permission string

const (
	BorrowerWorkspace    Permission = "borrower.workspace"
	StaffWorkspace       Permission = "staff.workspace"
	Administration       Permission = "administration"
	ReadAccountDirectory Permission = "accounts.read"
)

func (r Role) Allows(permission Permission) bool {
	switch permission {
	case BorrowerWorkspace:
		return r == RoleBorrower
	case StaffWorkspace:
		return r == RoleStaff || r == RoleAdmin
	case Administration, ReadAccountDirectory:
		return r == RoleAdmin
	default:
		return false
	}
}
