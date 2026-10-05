package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime"

	"github.com/gofiber/fiber/v3"
)

// bindStrictJSON prevents role/status/identity/security fields from hiding in
// create/self-update payloads. Callers return a generic 400, never decode details.
func bindStrictJSON(c fiber.Ctx, target any) error {
	media, _, err := mime.ParseMediaType(c.Get("Content-Type"))
	if err != nil || media != "application/json" {
		return errors.New("JSON content type required")
	}
	decoder := json.NewDecoder(bytes.NewReader(c.Body()))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return err
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		return errors.New("single JSON body required")
	}
	return nil
}
