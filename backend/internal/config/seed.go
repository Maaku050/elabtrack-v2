package config

import "errors"

// CheckDevelopmentSeed must pass before filesystem or database access.
func CheckDevelopmentSeed(env Environment, requested bool) error {
	if env != Development {
		return errors.New("seed: permitted only in development")
	}
	if !requested {
		return errors.New("seed: explicit --seed request required")
	}
	return nil
}
