package user

import "testing"

func TestProductPermissionMatrix(t *testing.T) {
	for _, c := range []struct {
		role                   Role
		borrower, staff, admin bool
	}{
		{RoleBorrower, true, false, false}, {RoleStaff, false, true, false}, {RoleAdmin, false, true, true},
		{"user", false, false, false}, {"admin", false, false, false}, {"Student", false, false, false}, {"Faculty", false, false, false}, {"", false, false, false},
	} {
		t.Run(string(c.role), func(t *testing.T) {
			for _, p := range []struct {
				permission Permission
				allowed    bool
			}{{BorrowerWorkspace, c.borrower}, {StaffWorkspace, c.staff}, {Administration, c.admin}, {ReadAccountDirectory, c.admin}, {"unknown", false}} {
				if c.role.Allows(p.permission) != p.allowed {
					t.Fatalf("incorrect permission %s", p.permission)
				}
			}
		})
	}
}
