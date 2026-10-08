package spreadsheet

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	d "github.com/Maaku050/elabtrack-v2/backend/internal/domain/accounts"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/shared"
	"github.com/xuri/excelize/v2"
	"io"
	"strings"
)

const MaxUpload = 768 << 10

var headers = []string{"studentId", "name", "email", "course", "contactNumber"}

type StudentRoster struct{}

func (StudentRoster) Template() ([]byte, error) {
	f := excelize.NewFile()
	defer f.Close()
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		if e := f.SetCellStr("Sheet1", cell, h); e != nil {
			return nil, e
		}
	}
	style, e := f.NewStyle(&excelize.Style{NumFmt: 49})
	if e != nil {
		return nil, e
	}
	if e = f.SetColStyle("Sheet1", "A:E", style); e != nil {
		return nil, e
	}
	_ = f.SetColWidth("Sheet1", "A", "E", 24)
	b, e := f.WriteToBuffer()
	if e != nil {
		return nil, e
	}
	return b.Bytes(), nil
}
func preflight(raw []byte) error {
	if len(raw) == 0 || len(raw) > MaxUpload {
		return shared.ErrInvalidInput
	}
	z, e := zip.NewReader(bytes.NewReader(raw), int64(len(raw)))
	if e != nil || len(z.File) > 64 {
		return shared.ErrInvalidInput
	}
	seen := map[string]bool{}
	var total uint64
	for _, f := range z.File {
		if seen[f.Name] || strings.Contains(f.Name, "..") || strings.HasPrefix(f.Name, "/") || strings.Contains(f.Name, "\\") || f.UncompressedSize64 > 4<<20 {
			return shared.ErrInvalidInput
		}
		seen[f.Name] = true
		total += f.UncompressedSize64
		if total > 16<<20 {
			return shared.ErrInvalidInput
		}
		lower := strings.ToLower(f.Name)
		if strings.Contains(lower, "vbaproject") || strings.Contains(lower, "externallinks") || strings.Contains(lower, "embeddings") {
			return shared.ErrInvalidInput
		}
		r, e := f.Open()
		if e != nil {
			return shared.ErrInvalidInput
		}
		body, e := io.ReadAll(io.LimitReader(r, 4<<20+1))
		_ = r.Close()
		if e != nil || len(body) > 4<<20 || uint64(len(body)) != f.UncompressedSize64 {
			return shared.ErrInvalidInput
		}
		if strings.HasSuffix(lower, ".rels") || strings.HasPrefix(lower, "xl/worksheets/") && strings.HasSuffix(lower, ".xml") {
			dec := xml.NewDecoder(bytes.NewReader(body))
			for {
				t, e := dec.Token()
				if e == io.EOF {
					break
				}
				if e != nil {
					return shared.ErrInvalidInput
				}
				if start, ok := t.(xml.StartElement); ok {
					if start.Name.Local == "f" {
						return shared.ErrInvalidInput
					}
					for _, a := range start.Attr {
						if a.Name.Local == "TargetMode" && strings.EqualFold(a.Value, "External") {
							return shared.ErrInvalidInput
						}
					}
				}
			}
		}
	}
	return nil
}
func (StudentRoster) Parse(raw []byte) ([]d.Row, error) {
	if e := preflight(raw); e != nil {
		return nil, e
	}
	f, e := excelize.OpenReader(bytes.NewReader(raw), excelize.Options{UnzipSizeLimit: 16 << 20, UnzipXMLSizeLimit: 4 << 20, RawCellValue: true})
	if e != nil {
		return nil, shared.ErrInvalidInput
	}
	defer f.Close()
	sheets := f.GetSheetList()
	if len(sheets) != 1 {
		return nil, shared.ErrInvalidInput
	}
	sheet := sheets[0]
	iterator, e := f.Rows(sheet)
	if e != nil {
		return nil, shared.ErrInvalidInput
	}
	defer iterator.Close()
	out := []d.Row{}
	number := 0
	for iterator.Next() {
		number++
		if number > 501 {
			return nil, shared.ErrInvalidInput
		}
		cells, e := iterator.Columns(excelize.Options{RawCellValue: true})
		if e != nil || len(cells) > 5 {
			return nil, shared.ErrInvalidInput
		}
		values := make([]string, 5)
		copy(values, cells)
		for i := range values {
			cell, _ := excelize.CoordinatesToCellName(i+1, number)
			formula, e := f.GetCellFormula(sheet, cell)
			if e != nil || formula != "" {
				return nil, shared.ErrInvalidInput
			}
			if len(values[i]) > 512 || strings.HasPrefix(strings.TrimSpace(values[i]), "=") {
				return nil, shared.ErrInvalidInput
			}
		}
		if number == 1 {
			for i, h := range headers {
				if values[i] != h {
					return nil, shared.ErrInvalidInput
				}
			}
			continue
		}
		if strings.TrimSpace(strings.Join(values, "")) == "" {
			continue
		}
		cell, _ := excelize.CoordinatesToCellName(1, number)
		kind, e := f.GetCellType(sheet, cell)
		if e != nil {
			return nil, shared.ErrInvalidInput
		}
		row := d.Row{Number: number, Input: d.Input{StudentID: values[0], Name: values[1], Email: values[2], Course: values[3], ContactNumber: values[4], BorrowerType: "STUDENT"}}
		if kind != excelize.CellTypeSharedString && kind != excelize.CellTypeInlineString {
			row.Error = "Student ID must be a text cell; numeric cells can lose leading zeroes"
		}
		out = append(out, row)
	}
	if iterator.Error() != nil || len(out) == 0 {
		return nil, shared.ErrInvalidInput
	}
	return out, nil
}
