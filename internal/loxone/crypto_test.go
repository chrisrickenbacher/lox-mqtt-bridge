package loxone

import (
	"encoding/hex"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashUserPassword(t *testing.T) {
	password := "password"
	salt := "salt"

	// Expected values calculated using:
	// echo -n "password:salt" | shasum
	// echo -n "password:salt" | shasum -a 256

	tests := []struct {
		name     string
		alg      string
		expected string // Uppercase Hex
	}{
		{
			name:     "SHA1",
			alg:      "SHA1",
			expected: strings.ToUpper("676f03a8c8530384eb7551b1685f2828546625d9"),
		},
		{
			name:     "SHA256",
			alg:      "SHA256",
			expected: strings.ToUpper("f64671af1dd46e4a00a48a2c7c6a3658d107507391b6eb0d9111b2b3d326512b"),
		},
		{
			name:     "Default (Invalid Alg)",
			alg:      "UNKNOWN",
			expected: strings.ToUpper("676f03a8c8530384eb7551b1685f2828546625d9"), // Should default to SHA1
		},
		{
			name:     "Lowercase Input",
			alg:      "sha1",
			expected: strings.ToUpper("676f03a8c8530384eb7551b1685f2828546625d9"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HashUserPassword(password, salt, tt.alg)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestComputeHMAC(t *testing.T) {
	keyRaw := "00112233445566778899aabbccddeeff" // 16 bytes key
	message := "user:pwHash"

	// Expected values calculated using:
	// echo -n "user:pwHash" | openssl dgst -sha1 -mac HMAC -macopt hexkey:00112233445566778899aabbccddeeff
	// echo -n "user:pwHash" | openssl dgst -sha256 -mac HMAC -macopt hexkey:00112233445566778899aabbccddeeff

	tests := []struct {
		name     string
		alg      string
		expected string // Lowercase Hex
	}{
		{
			name:     "SHA1",
			alg:      "SHA1",
			expected: "0d433c1134e632fc2fed17eed3d12e1af5ef0b07",
		},
		{
			name:     "SHA256",
			alg:      "SHA256",
			expected: "b2076826cbda86ddaaed4183a4a585bf7d183028d68ad41ddf45907a0a27b3f1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ComputeHMAC(keyRaw, message, tt.alg)
			assert.Equal(t, tt.expected, result)
		})
	}

	t.Run("Invalid Key Hex", func(t *testing.T) {
		result := ComputeHMAC("ZZZZ", message, "SHA1")
		assert.Equal(t, "", result)
	})
}

// Simple integration check to ensure key decoding works as expected
func TestKeyDecoding(t *testing.T) {
	keyHex := "deadbeef"
	keyBytes, err := hex.DecodeString(keyHex)
	assert.NoError(t, err)
	assert.Equal(t, []byte{0xde, 0xad, 0xbe, 0xef}, keyBytes)
}