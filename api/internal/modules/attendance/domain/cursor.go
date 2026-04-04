package domain

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

// CursorToken is the decoded opaque cursor (id + HMAC signature).
type CursorToken struct {
	ID  int64  `json:"id"`
	Sig string `json:"sig"`
}

const attendanceCursorMaxLimit = 100

// CapAttendanceLimit clamps limit to a safe range (default 20, max 100).
func CapAttendanceLimit(limit int) int {
	if limit <= 0 {
		return 20
	}
	if limit > attendanceCursorMaxLimit {
		return attendanceCursorMaxLimit
	}
	return limit
}

// EncodeAttendanceCursor builds a Base64 JSON cursor with HMAC over the id.
func EncodeAttendanceCursor(id int64, secret string) (string, error) {
	if secret == "" {
		return "", fmt.Errorf("cursor secret required")
	}
	sig := hmacSignAttendance(id, secret)
	tok := CursorToken{ID: id, Sig: sig}
	raw, err := json.Marshal(tok)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(raw), nil
}

// DecodeAttendanceCursor validates Base64, JSON, and HMAC; returns cursor id (0 means start).
func DecodeAttendanceCursor(cursor, secret string) (int64, error) {
	if strings.TrimSpace(cursor) == "" {
		return 0, nil
	}
	raw, err := base64.URLEncoding.DecodeString(cursor)
	if err != nil {
		return 0, ErrInvalidCursor
	}
	var tok CursorToken
	if err := json.Unmarshal(raw, &tok); err != nil {
		return 0, ErrInvalidCursor
	}
	if tok.ID < 0 {
		return 0, ErrInvalidCursor
	}
	expected := hmacSignAttendance(tok.ID, secret)
	if !hmac.Equal([]byte(expected), []byte(tok.Sig)) {
		return 0, ErrInvalidCursor
	}
	return tok.ID, nil
}

func hmacSignAttendance(id int64, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(fmt.Sprintf("%d", id)))
	return base64.URLEncoding.EncodeToString(mac.Sum(nil))
}
