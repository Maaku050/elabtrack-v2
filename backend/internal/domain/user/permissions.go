package user

// Only implemented workspace/account boundaries are defined here. Future
// mutations require their own reviewed permission and transactional authority.
type Permission string

const (
	BorrowerWorkspace    Permission = "borrower.workspace"
	StaffWorkspace       Permission = "staff.workspace"
	Administration       Permission = "administration"
	ReadAccountDirectory Permission = "accounts.read"
	AcceptBorrowerTerms  Permission = "terms.accept"
	ManageInventory      Permission = "inventory.manage"
	ReadBorrowers        Permission = "borrowers.read"
	ManageAccounts       Permission = "accounts.manage"
	PublishTerms         Permission = "terms.publish"
)

func (r Role) Allows(permission Permission) bool {
	switch permission {
	case BorrowerWorkspace, AcceptBorrowerTerms:
		return r == RoleBorrower
	case StaffWorkspace, ReadBorrowers, ManageInventory:
		return r == RoleStaff || r == RoleAdmin
	case Administration, ReadAccountDirectory, PublishTerms, ManageAccounts:
		return r == RoleAdmin
	default:
		return false
	}
}
