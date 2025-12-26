package loxone

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"strings"
)

// HashUserPassword hashes the password with the salt using the specified algorithm (SHA1 or SHA256)
// Returns Uppercase Hex string.
// Format: Hash("password:salt")
func HashUserPassword(password, salt, alg string) string {
	payload := fmt.Sprintf("%s:%s", password, salt)
	var h hash.Hash

	switch strings.ToUpper(alg) {
	case "SHA1":
		h = sha1.New()
	case "SHA256":
		h = sha256.New()
	default:
		// Default to SHA1 if unknown
		h = sha1.New()
	}

	h.Write([]byte(payload))
	sum := h.Sum(nil)
	return strings.ToUpper(hex.EncodeToString(sum))
}

// ComputeHMAC computes the HMAC of the message using the provided key and algorithm.
// The key is provided as a Hex string (from getkey2 response).
// Returns Hex string.
func ComputeHMAC(keyHex, message, alg string) string {
	keyBytes, err := hex.DecodeString(keyHex)
	if err != nil {
		return ""
	}

	var h hash.Hash

	switch strings.ToUpper(alg) {
	case "SHA1":
		h = hmac.New(sha1.New, keyBytes)
	case "SHA256":
		h = hmac.New(sha256.New, keyBytes)
	default:
		h = hmac.New(sha1.New, keyBytes)
	}

	h.Write([]byte(message))
	sum := h.Sum(nil)
	return hex.EncodeToString(sum)
}
