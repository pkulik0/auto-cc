package utility

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

func HashFromStrings(data []string) string {
	joined := strings.Join(data, "")
	hash := sha256.Sum256([]byte(joined))
	return hex.EncodeToString(hash[:])
}
