package profile

import (
	"context"
	"github.com/google/uuid"
)

const MaxBytes = 512 * 1024

// Photos are optional account presentation, never return evidence.
type Metadata struct {
	AccountID uuid.UUID  `json:"account_id"`
	ImageID   *uuid.UUID `json:"image_id"`
	Version   int64      `json:"version"`
}
type Image struct {
	Metadata
	PNG           []byte
	Hash          string
	Width, Height int
}
type Repository interface {
	Metadata(context.Context, uuid.UUID) (Metadata, error)
	Image(context.Context, uuid.UUID, uuid.UUID) (Image, error)
	Save(context.Context, Image) error
}
type Validator interface{ Validate([]byte) (Image, error) }
