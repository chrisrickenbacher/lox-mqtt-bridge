package loxone

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"
)

// LoxApp3 represents the root of the structure file
type LoxApp3 struct {
	LastModified string                 `json:"lastModified"`
	MsInfo       map[string]interface{} `json:"msInfo"`
	Global       map[string]interface{} `json:"global"`
	Rooms        map[string]*Room       `json:"rooms"`
	Cats         map[string]*Cat        `json:"cats"`
	Controls     map[string]*Control    `json:"controls"`
}

type Room struct {
	Name string `json:"name"`
	UUID string `json:"uuid"`
}

type Cat struct {
	Name string `json:"name"`
	Type string `json:"type"`
	UUID string `json:"uuid"`
}

// Control represents a generic device/block
type Control struct {
	Name        string                 `json:"name"`
	Type        string                 `json:"type"`
	UUIDAction  string                 `json:"uuidAction"`
	Room        string                 `json:"room"`   // UUID of the room
	Cat         string                 `json:"cat"`    // UUID of the category
	States      map[string]interface{} `json:"states"` // Map of state-name -> UUID (or array/object)
	Details     map[string]interface{} `json:"details"`
	SubControls map[string]*Control    `json:"subControls"`
}

// GetStructure fetches the structure file, using in-memory caching and version checks
func (c *Client) GetStructure() (*LoxApp3, error) {
	// 1. Check Version
	if err := c.SendCommand("jdev/sps/LoxAPPversion3"); err != nil {
		return nil, fmt.Errorf("failed to request structure version: %v", err)
	}

	resp, err := c.WaitForResponse("LoxAPPversion3")
	if err != nil {
		return nil, fmt.Errorf("failed to get structure version: %v", err)
	}

	// Parse version from response value
	// Value should be a string like "2023-10-25 10:00:00"
	var currentVersion string
	if err := json.Unmarshal(resp.LL.Value, &currentVersion); err != nil {
		// Fallback: sometimes it might be a raw string or different format?
		// But usually it is a JSON string.
		slog.Warn("Could not unmarshal version value", "error", err, "raw", string(resp.LL.Value))
		// Proceed to fetch anyway if we can't parse
	}

	// 2. Return Cache if valid
	if currentVersion != "" && currentVersion == c.structureLastMod && c.structureCache != nil {
		return c.structureCache, nil
	}

	// 3. Fetch New Structure
	if err := c.SendCommand("data/LoxAPP3.json"); err != nil {
		return nil, err
	}

	// 4. Wait for the file
	// The structure file is sent as a raw text message (Type 0), not wrapped in { "LL": ... }
	// We need to filter incoming messages until we find it.
	timeout := time.After(30 * time.Second)
	for {
		select {
		case msg := <-c.msgChan:
			// Check if this is the structure file
			// It starts with { "lastModified": ...
			if strings.Contains(string(msg), "lastModified") && strings.Contains(string(msg), "msInfo") {
				var structure LoxApp3
				if err := json.Unmarshal(msg, &structure); err != nil {
					return nil, fmt.Errorf("failed to parse structure JSON: %v", err)
				}

				// Update Cache
				c.structureCache = &structure
				c.structureLastMod = structure.LastModified
				slog.Info("Structure updated", "controls", len(structure.Controls))

				return &structure, nil
			}
			// Ignore other messages (e.g. keepalives or other control responses)
			// Ideally we would put them back or handle them, but for now we consume them.
			// TODO: A better multiplexer might be needed if this loop eats important packets.

		case <-timeout:
			return nil, fmt.Errorf("timeout waiting for structure file")
		}
	}
}