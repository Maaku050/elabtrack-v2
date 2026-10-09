package spreadsheet

import (
	"archive/zip"
	"bytes"
	"encoding/binary"
	"errors"
	d "github.com/Maaku050/elabtrack-v2/backend/internal/domain/accounts"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/shared"
	"github.com/xuri/excelize/v2"
	"testing"
)

func workbook(t *testing.T, edit func(*excelize.File)) []byte {
	t.Helper()
	f := excelize.NewFile()
	defer f.Close()
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		_ = f.SetCellStr("Sheet1", cell, h)
	}
	_ = f.SetCellStr("Sheet1", "A2", "0028366")
	_ = f.SetCellStr("Sheet1", "B2", "Synthetic Student")
	_ = f.SetCellStr("Sheet1", "C2", "student@students.example.invalid")
	if edit != nil {
		edit(f)
	}
	b, e := f.WriteToBuffer()
	if e != nil {
		t.Fatal(e)
	}
	return b.Bytes()
}
func TestStudentWorkbook(t *testing.T) {
	p := StudentRoster{}
	v, e := p.Parse(workbook(t, nil))
	if e != nil || len(v) != 1 || v[0].Input.StudentID != "0028366" || v[0].Number != 2 {
		t.Fatal("text identity not preserved")
	}
	v, e = p.Parse(workbook(t, func(f *excelize.File) { _ = f.SetCellInt("Sheet1", "A2", 28366) }))
	if e != nil || v[0].Error == "" {
		t.Fatal("numeric Student ID silently accepted")
	}
	for _, tc := range []struct {
		name string
		edit func(*excelize.File)
	}{{"formula", func(f *excelize.File) { _ = f.SetCellFormula("Sheet1", "B2", "1+1") }}, {"password", func(f *excelize.File) { _ = f.SetCellStr("Sheet1", "F1", "password") }}, {"extra_sheet", func(f *excelize.File) { _, _ = f.NewSheet("Other") }}, {"unknown_header", func(f *excelize.File) { _ = f.SetCellStr("Sheet1", "E1", "role") }}, {"row_bound", func(f *excelize.File) { _ = f.SetCellStr("Sheet1", "A502", "over limit") }}} {
		t.Run(tc.name, func(t *testing.T) {
			if _, e = p.Parse(workbook(t, tc.edit)); e == nil {
				t.Fatal("unsafe workbook accepted")
			}
		})
	}
	if _, e = p.Parse([]byte("not a zip")); e == nil {
		t.Fatal("malformed zip")
	}
	if _, e = p.Parse(make([]byte, MaxUpload+1)); e == nil {
		t.Fatal("oversized upload")
	}
	raw, e := p.Template()
	if e != nil {
		t.Fatal(e)
	}
	f, e := excelize.OpenReader(bytes.NewReader(raw))
	if e != nil {
		t.Fatal(e)
	}
	defer f.Close()
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		value, _ := f.GetCellValue("Sheet1", cell)
		if value != h {
			t.Fatal("template parser mismatch")
		}
	}
}
func TestUnsignedZipBoundsAndExternalRelationships(t *testing.T) {
	var b bytes.Buffer
	z := zip.NewWriter(&b)
	w, _ := z.CreateHeader(&zip.FileHeader{Name: "large.xml", Method: zip.Deflate, UncompressedSize64: 1 << 63})
	_, _ = w.Write(bytes.Repeat([]byte("x"), 4<<20+1))
	_ = z.Close()
	if preflight(b.Bytes()) == nil {
		t.Fatal("oversized archive accepted")
	}
	b.Reset()
	z = zip.NewWriter(&b)
	w, _ = z.Create("_rels/.rels")
	_, _ = w.Write([]byte(`<Relationships><Relationship TargetMode="External" Target="https://example.invalid"/></Relationships>`))
	_ = z.Close()
	if preflight(b.Bytes()) == nil {
		t.Fatal("external relation accepted")
	}
}

func TestHostileZip64UncompressedSize(t *testing.T) {
	var buffer bytes.Buffer
	writer := zip.NewWriter(&buffer)
	file, _ := writer.Create("tiny.xml")
	_, _ = file.Write([]byte("x"))
	_ = writer.Close()
	data := buffer.Bytes()
	central := bytes.Index(data, []byte{'P', 'K', 1, 2})
	if central < 0 {
		t.Fatal("fixture central directory")
	}
	nameLen := int(binary.LittleEndian.Uint16(data[central+28 : central+30]))
	extraLen := int(binary.LittleEndian.Uint16(data[central+30 : central+32]))
	if extraLen != 0 {
		t.Fatal("unexpected fixture extra")
	}
	binary.LittleEndian.PutUint32(data[central+24:central+28], 0xffffffff)
	binary.LittleEndian.PutUint16(data[central+30:central+32], 12)
	extra := make([]byte, 12)
	binary.LittleEndian.PutUint16(extra, 1)
	binary.LittleEndian.PutUint16(extra[2:], 8)
	binary.LittleEndian.PutUint64(extra[4:], 1<<63)
	at := central + 46 + nameLen
	hostile := append(append(append([]byte{}, data[:at]...), extra...), data[at:]...)
	end := bytes.LastIndex(hostile, []byte{'P', 'K', 5, 6})
	size := binary.LittleEndian.Uint32(hostile[end+12 : end+16])
	binary.LittleEndian.PutUint32(hostile[end+12:end+16], size+12)
	archive, err := zip.NewReader(bytes.NewReader(hostile), int64(len(hostile)))
	if err != nil || archive.File[0].UncompressedSize64 != 1<<63 {
		t.Fatal("fixture must carry actual unsigned ZIP64 size")
	}
	if preflight(hostile) == nil {
		t.Fatal("huge unsigned size must be rejected before parser conversion/decompression")
	}
}

func TestLocatedStudentWorkbookFailuresAndBlankID(t *testing.T) {
	p := StudentRoster{}
	for _, tc := range []struct {
		name        string
		edit        func(*excelize.File)
		kind        d.RosterIssueKind
		row, column int
	}{
		{"formula", func(f *excelize.File) { _ = f.SetCellFormula("Sheet1", "B2", "1+1") }, d.RosterFormula, 2, 2},
		{"header", func(f *excelize.File) { _ = f.SetCellStr("Sheet1", "A1", "wrong") }, d.RosterHeaders, 1, 1},
		{"extra column", func(f *excelize.File) { _ = f.SetCellStr("Sheet1", "F1", "password") }, d.RosterColumns, 1, 0},
		{"extra worksheet", func(f *excelize.File) { _, _ = f.NewSheet("Other") }, d.RosterWorksheets, 0, 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := p.Parse(workbook(t, tc.edit))
			var issue *d.RosterIssue
			if !errors.As(err, &issue) || issue.Kind != tc.kind || issue.Row != tc.row || issue.Column != tc.column || !errors.Is(err, shared.ErrInvalidInput) {
				t.Fatalf("classification: %v", err)
			}
		})
	}
	rows, err := p.Parse(workbook(t, func(f *excelize.File) { _ = f.SetCellStr("Sheet1", "A2", "") }))
	if err != nil || len(rows) != 1 || rows[0].Error != "" || rows[0].Input.StudentID != "" {
		t.Fatal("blank ID must reach required identity validation")
	}
	raw, err := p.Template()
	if err != nil {
		t.Fatal(err)
	}
	f, err := excelize.OpenReader(bytes.NewReader(raw))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	for _, cell := range []string{"A2", "A501", "E2", "E501"} {
		id, _ := f.GetCellStyle("Sheet1", cell)
		style, err := f.GetStyle(id)
		if err != nil || style.NumFmt != 49 {
			t.Fatal("input cell must use Text format", cell)
		}
	}
	_ = f.SetCellStr("Sheet1", "A2", "0028366")
	_ = f.SetCellStr("Sheet1", "B2", "Synthetic Student")
	_ = f.SetCellStr("Sheet1", "C2", "student@students.example.invalid")
	_ = f.SetCellStr("Sheet1", "E2", "00987654321")
	b, err := f.WriteToBuffer()
	if err != nil {
		t.Fatal(err)
	}
	rows, err = p.Parse(b.Bytes())
	if err != nil || len(rows) != 1 || rows[0].Input.StudentID != "0028366" || rows[0].Input.ContactNumber != "00987654321" {
		t.Fatal("text input template round trip")
	}
}
