package profileimage

import (
	d "github.com/Maaku050/elabtrack-v2/backend/internal/domain/profile"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/catalogimage"
)

type Validator struct{}

func (Validator) Validate(raw []byte) (d.Image, error) {
	v, e := (catalogimage.Validator{}).Validate(raw)
	return d.Image{PNG: v.PNG, Hash: v.Hash, Width: v.Width, Height: v.Height}, e
}
