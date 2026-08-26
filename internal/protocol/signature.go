package protocol

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

func Signature(envelope Envelope) (string, error) {
	raw, err := json.Marshal(envelope)
	if err != nil {
		return "", fmt.Errorf("marshal envelope signature: %w", err)
	}
	digest := sha256.Sum256(raw)
	return hex.EncodeToString(digest[:]), nil
}

func VerifySignature(envelope Envelope, expected string) error {
	actual, err := Signature(envelope)
	if err != nil {
		return err
	}
	if actual != expected {
		return fmt.Errorf("message signature mismatch")
	}
	return nil
}
