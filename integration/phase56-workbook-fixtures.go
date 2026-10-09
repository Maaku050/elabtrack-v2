package main

import (
	"encoding/json"
	"fmt"
	"github.com/xuri/excelize/v2"
	"os"
	"strings"
)

func main() {
	raw, e := os.ReadFile("/tmp/elabtrack-phase56-items.json")
	if e != nil {
		panic(e)
	}
	var v struct {
		Faculty struct {
			Email string `json:"email"`
		}
	}
	if json.Unmarshal(raw, &v) != nil || len(v.Faculty.Email) < 20 || !strings.HasPrefix(v.Faculty.Email, "faculty-") || !strings.HasSuffix(v.Faculty.Email, "@students.example.invalid") {
		panic("synthetic item fixture missing")
	}
	for _, kind := range []string{"valid", "numeric", "header", "formula", "missing", "domain", "conflict", "duplicate"} {
		f := excelize.NewFile()
		for i, h := range []string{"studentId", "name", "email", "course", "contactNumber"} {
			c, _ := excelize.CoordinatesToCellName(i+1, 1)
			_ = f.SetCellStr("Sheet1", c, h)
		}
		_ = f.SetCellStr("Sheet1", "A2", "0028366-"+v.Faculty.Email[8:20])
		_ = f.SetCellStr("Sheet1", "B2", "TEST roster Student")
		_ = f.SetCellStr("Sheet1", "C2", "student-"+v.Faculty.Email)
		_ = f.SetCellStr("Sheet1", "E2", "00912345678")
		switch kind {
		case "numeric":
			_ = f.SetCellInt("Sheet1", "A2", 28366)
		case "header":
			_ = f.SetCellStr("Sheet1", "A1", "wrong")
		case "formula":
			_ = f.SetCellFormula("Sheet1", "B2", "1+1")
		case "missing":
			_ = f.SetCellStr("Sheet1", "B2", "")
		case "domain":
			_ = f.SetCellStr("Sheet1", "C2", "test@example.invalid")
		case "conflict":
			_ = f.SetCellStr("Sheet1", "C2", v.Faculty.Email)
		case "duplicate":
			_ = f.SetCellStr("Sheet1", "A3", "0028366-"+v.Faculty.Email[8:20])
			_ = f.SetCellStr("Sheet1", "B3", "TEST duplicate Student")
			_ = f.SetCellStr("Sheet1", "C3", "student-"+v.Faculty.Email)
		}
		b, e := f.WriteToBuffer()
		if e != nil {
			panic(e)
		}
		if os.WriteFile("/tmp/elabtrack-phase56-"+kind+".xlsx", b.Bytes(), 0600) != nil {
			panic("fixture write")
		}
		_ = f.Close()
	}
	fmt.Println("Synthetic workbook fixtures generated")
}
