package catalogimage

import (
	"bytes"
	"encoding/binary"
	d "github.com/Maaku050/elabtrack-v2/backend/internal/domain/inventory"
	"image"
	"image/jpeg"
	"image/png"
	"testing"
)

func TestSafeCatalogImage(t *testing.T) {
	im := image.NewRGBA(image.Rect(0, 0, 8, 8))
	for _, jpegInput := range []bool{false, true} {
		var raw bytes.Buffer
		if jpegInput {
			_ = jpeg.Encode(&raw, im, nil)
		} else {
			_ = png.Encode(&raw, im)
		}
		v, e := (Validator{}).Validate(raw.Bytes())
		if e != nil || v.Width != 8 || v.Height != 8 || len(v.Hash) != 64 || !bytes.HasPrefix(v.PNG, []byte("\x89PNG")) {
			t.Fatal("canonical raster")
		}
		raw.WriteString("hidden credential metadata")
		v, e = (Validator{}).Validate(raw.Bytes())
		if e != nil || bytes.Contains(v.PNG, []byte("credential")) {
			t.Fatal("metadata retained")
		}
	}
	for _, raw := range [][]byte{nil, []byte("<svg onload='alert(1)'/>"), []byte("GIF89a"), bytes.Repeat([]byte{1}, d.MaxImageBytes+1), []byte("\x89PNG\r\n\x1a\n")} {
		if _, e := (Validator{}).Validate(raw); e == nil {
			t.Fatal("arbitrary image accepted")
		}
	}
	var b bytes.Buffer
	_ = png.Encode(&b, im)
	bad := append([]byte(nil), b.Bytes()...)
	binary.BigEndian.PutUint32(bad[16:20], 1<<30)
	if _, e := (Validator{}).Validate(bad); e == nil {
		t.Fatal("dimension bomb")
	}
}
