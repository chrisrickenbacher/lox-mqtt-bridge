package loxone

import (
	"encoding/json"
	"fmt"
	"log/slog"
)

// Authenticate performs the Token-Based Authentication
func (c *Client) Authenticate() error {
	// 1. Get Key and Salt
	// Request: jdev/sys/getkey2/{user}
	cmd := fmt.Sprintf("jdev/sys/getkey2/%s", c.cfg.User)
	if err := c.SendCommand(cmd); err != nil {
		return fmt.Errorf("failed to send getkey2: %v", err)
	}

	resp, err := c.WaitForResponse("getkey2")
	if err != nil {
		return fmt.Errorf("getkey2 failed: %v", err)
	}

	// Parse response
	// The response value should be a JSON object with key, salt, hashAlg
	var keyInfo struct {
		Key     string `json:"key"`
		Salt    string `json:"salt"`
		HashAlg string `json:"hashAlg"`
	}
	if err := json.Unmarshal(resp.LL.Value, &keyInfo); err != nil {
		return fmt.Errorf("failed to parse getkey2 response: %v", err)
	}

	// 2. Hash Password
	// Password Hash: Hash(password:salt)
	// Result must be Uppercase Hex
	pwHash := HashUserPassword(c.cfg.Pass, keyInfo.Salt, keyInfo.HashAlg)

	// 3. Create Authentication Hash
	// Authentication Hash: HMAC(key, "user:pwHash")
	// Result must be Hex
	userPwHash := fmt.Sprintf("%s:%s", c.cfg.User, pwHash)
	authHash := ComputeHMAC(keyInfo.Key, userPwHash, keyInfo.HashAlg)

	// 4. Request Token
	// Command: jdev/sys/getjwt/{hash}/{user}/{permission}/{uuid}/{info}
	// Permission 2 = Web

	// UUID: We use a static UUID for the bridge to ensure token persistence if needed,
	// or we could generate one. For now, we use a fixed bridge UUID.
	uuid := "50325345-5200-0000-0000-000000000000"
	info := "LoxoneMQTTBridge"

	tokenCmd := fmt.Sprintf("jdev/sys/getjwt/%s/%s/%d/%s/%s", authHash, c.cfg.User, 2, uuid, info)

	if err := c.SendCommand(tokenCmd); err != nil {
		return fmt.Errorf("failed to send getjwt: %v", err)
	}

	resp, err = c.WaitForResponse("getjwt")
	if err != nil {
		return fmt.Errorf("getjwt failed: %v", err)
	}

	// 5. Store Token
	var token Token
	if err := json.Unmarshal(resp.LL.Value, &token); err != nil {
		return fmt.Errorf("failed to parse token: %v", err)
	}
	c.token = &token
	slog.Info("Authentication successful")

	return nil
}
