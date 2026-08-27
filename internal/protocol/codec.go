package protocol

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

type Codec struct {
	Prefix string
}

func NewCodec(prefix string) Codec {
	return Codec{Prefix: strings.TrimSpace(prefix)}
}

func (c Codec) EncodeEnvelope(envelope Envelope) (string, error) {
	raw, err := envelope.JSON()
	if err != nil {
		return "", err
	}
	encoded := base64.RawStdEncoding.EncodeToString(raw)
	if c.Prefix == "" {
		return encoded, nil
	}
	return c.Prefix + "." + encoded, nil
}

func (c Codec) DecodeEnvelope(value string) (Envelope, error) {
	parts := strings.SplitN(value, ".", 2)
	if c.Prefix != "" {
		if len(parts) != 2 || parts[0] != c.Prefix {
			return Envelope{}, fmt.Errorf("invalid protocol prefix")
		}
		value = parts[1]
	}
	raw, err := base64.RawStdEncoding.DecodeString(value)
	if err != nil {
		return Envelope{}, fmt.Errorf("decode base64 envelope: %v", err)
	}
	return ParseEnvelope(raw)
}

func MarshalForLog(envelope Envelope) string {
	raw, err := json.Marshal(struct {
		MessageID string `json:"message_id"`
		Source    string `json:"source"`
		FrameID   string `json:"frame_id"`
	}{envelope.MessageID, envelope.Source, firstFrameID(envelope)})
	if err != nil {
		return "{}"
	}
	return string(raw)
}

func firstFrameID(envelope Envelope) string {
	if len(envelope.Frames) == 0 {
		return ""
	}
	return envelope.Frames[0].ID
}
