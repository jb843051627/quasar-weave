package validation

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func PayloadHash(payload []byte) string {
	digest := sha256.Sum256(payload)
	return hex.EncodeToString(digest[:])
}

func VerifyPayload(payload []byte, expected string) error {
	if expected == "" {
		return fmt.Errorf("payload hash is required")
	}
	if PayloadHash(payload) != expected {
		return fmt.Errorf("payload checksum mismatch")
	}
	return nil
}

func CanonicalText(parts ...string) string {
	result := ""
	for _, part := range parts {
		result += fmt.Sprintf("%d:%s|", len(part), part)
	}
	return result
}
