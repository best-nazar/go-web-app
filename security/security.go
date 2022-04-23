package security

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

const secret = "126c11b2ecf89c9e145616c506c2e767"

func ComputeHmac256(message string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}