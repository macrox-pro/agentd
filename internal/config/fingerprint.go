package config

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

// Fingerprint returns sha256 hex of canonical JSON for the merged fileConfig.
func Fingerprint(merged *fileConfig) (string, error) {
	if merged == nil {
		merged = &fileConfig{}
	}
	b, err := json.Marshal(merged)
	if err != nil {
		return "", fmt.Errorf("canonical json: %w", err)
	}
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:]), nil
}
