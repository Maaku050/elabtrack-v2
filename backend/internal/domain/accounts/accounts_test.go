package accounts

import (
	"errors"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/shared"
	"testing"
)

func TestProvisioningPolicy(t *testing.T) {
	p := Policy{StudentDomains: []string{"students.example.invalid"}}
	student := Input{Name: "Student Example", Email: "S@students.example.invalid", BorrowerType: "STUDENT", StudentID: "0028366"}
	v, e := p.Validate(student, false)
	if e != nil || v.StudentID != "0028366" || v.Email != "s@students.example.invalid" {
		t.Fatal("identity normalization must retain string ID")
	}
	student.Email = "outside@example.invalid"
	if _, e = p.Validate(student, false); !errors.Is(e, ErrStudentDomain) {
		t.Fatal("Student external domain accepted")
	}
	student.Email = "s@students.example.invalid"
	if _, e = (Policy{}).Validate(student, false); !errors.Is(e, ErrDomainsMissing) {
		t.Fatal("missing domains must fail closed")
	}
	faculty := Input{Name: "Faculty Example", Email: "faculty@example.invalid", BorrowerType: "FACULTY"}
	if _, e = (Policy{}).Validate(faculty, false); e != nil {
		t.Fatal("Faculty independent of Student domains")
	}
	faculty.StudentID = "123"
	if _, e = p.Validate(faculty, false); !errors.Is(e, shared.ErrInvalidInput) {
		t.Fatal("Faculty must not have Student ID")
	}
	if _, e = p.Validate(student, true); e == nil {
		t.Fatal("Staff payload accepted borrower attributes")
	}
	faculty.StudentID = ""
	faculty.BorrowerType = "ADMIN"
	if _, e = p.Validate(faculty, false); e == nil {
		t.Fatal("role escalation accepted")
	}
}
