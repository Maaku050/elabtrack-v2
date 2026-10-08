package catalogimage

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	d "github.com/Maaku050/elabtrack-v2/backend/internal/domain/inventory"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/shared"
	"image"
	_ "image/jpeg"
	"image/png"
)

type Validator struct{}

func (Validator) Validate(raw []byte) (d.Image, error) {
	var result d.Image
	if len(raw) == 0 || len(raw) > d.MaxImageBytes {
		return result, shared.ErrInvalidInput
	}
	cfg, format, e := image.DecodeConfig(bytes.NewReader(raw))
	if e != nil || (format != "png" && format != "jpeg") || cfg.Width < 1 || cfg.Height < 1 || cfg.Width > 2048 || cfg.Height > 2048 || int64(cfg.Width)*int64(cfg.Height) > 4194304 {
		return result, shared.ErrInvalidInput
	}
	im, _, e := image.Decode(bytes.NewReader(raw))
	if e != nil {
		return result, shared.ErrInvalidInput
	}
	var b bytes.Buffer
	if e = png.Encode(&b, im); e != nil || b.Len() > d.MaxImageBytes {
		return result, shared.ErrInvalidInput
	}
	sum := sha256.Sum256(b.Bytes())
	result.PNG = b.Bytes()
	result.Hash = hex.EncodeToString(sum[:])
	result.Width = cfg.Width
	result.Height = cfg.Height
	return result, nil
}
