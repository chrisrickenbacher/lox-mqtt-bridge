package bridge

import (
	"testing"

	"github.com/chrisrickenbacher/lox-mqtt-bridge/internal/loxone"
	"github.com/stretchr/testify/assert"
)

func TestRegistry_ProcessControls(t *testing.T) {
	// Define UUIDs for test
	uuidSwitch := "10000000-0000-0000-0000000000000001"
	uuidDimmerPos := "10000000-0000-0000-0000000000000002"
	uuidDimmerAct := "10000000-0000-0000-0000000000000003"
	uuidBlindPos := "10000000-0000-0000-0000000000000004"
	uuidMeterActive := "10000000-0000-0000-0000000000000005"
	uuidInfoAnalog := "10000000-0000-0000-0000000000000006"
	uuidColorRGB := "10000000-0000-0000-0000000000000007"
	uuidRoom := "20000000-0000-0000-0000000000000001"
	uuidCat := "30000000-0000-0000-0000000000000001"

	// Construct a mock Structure
	structure := &loxone.LoxApp3{
		LastModified: "2024-01-01 12:00:00",
		Rooms: map[string]*loxone.Room{
			uuidRoom: {Name: "Living Room", UUID: uuidRoom},
		},
		Cats: map[string]*loxone.Cat{
			uuidCat: {Name: "Lighting", Type: "Lighting", UUID: uuidCat},
		},
		Controls: map[string]*loxone.Control{
			"Ctrl1": {
				Name:       "Ceiling Light",
				Type:       "Switch",
				UUIDAction: "10000000-0000-0000-0000-0000000000000000",
				Room:       uuidRoom,
				Cat:        uuidCat,
				States: map[string]interface{}{
					"active": uuidSwitch,
				},
			},
			"Ctrl2": {
				Name:       "Dimmer Spot",
				Type:       "Dimmer",
				UUIDAction: "10000000-0000-0000-0000000000000010",
				Room:       uuidRoom,
				Cat:        uuidCat,
				States: map[string]interface{}{
					"position": uuidDimmerPos,
					"active":   uuidDimmerAct,
				},
			},
			"Ctrl3": {
				Name:       "Living Blind",
				Type:       "Jalousie",
				UUIDAction: "10000000-0000-0000-0000000000000020",
				Room:       uuidRoom,
				States: map[string]interface{}{
					"position": []interface{}{uuidBlindPos, "some-other-uuid"},
				},
			},
			"Ctrl4": {
				Name:       "Power Meter",
				Type:       "Meter",
				Room:       uuidRoom,
				States: map[string]interface{}{
					"active": uuidMeterActive,
				},
			},
			"Ctrl5": {
				Name:       "Temp Sensor",
				Type:       "InfoOnlyAnalog",
				Room:       uuidRoom,
				States: map[string]interface{}{
					"value": uuidInfoAnalog,
				},
			},
			"Ctrl6": {
				Name:       "LED Strip",
				Type:       "ColorPicker",
				Room:       uuidRoom,
				States: map[string]interface{}{
					"color": uuidColorRGB,
				},
			},
		},
	}

	// Initialize Registry
	registry := NewRegistry(structure)

	tests := []struct {
		name         string
		uuidStr      string
		expectedFound bool
		expectedRoom string
		expectedCtrl string
		expectedState string
	}{
		{
			name:          "Switch Active",
			uuidStr:       uuidSwitch,
			expectedFound: true,
			expectedRoom:  "Living Room",
			expectedCtrl:  "Ceiling Light",
			expectedState: "active",
		},
		{
			name:          "Dimmer Position",
			uuidStr:       uuidDimmerPos,
			expectedFound: true,
			expectedRoom:  "Living Room",
			expectedCtrl:  "Dimmer Spot",
			expectedState: "position",
		},
		{
			name:          "Dimmer Active",
			uuidStr:       uuidDimmerAct,
			expectedFound: true,
			expectedRoom:  "Living Room",
			expectedCtrl:  "Dimmer Spot",
			expectedState: "active",
		},
		{
			name:          "Blind Position (Array)",
			uuidStr:       uuidBlindPos,
			expectedFound: true,
			expectedRoom:  "Living Room",
			expectedCtrl:  "Living Blind",
			expectedState: "position",
		},
		{
			name:          "Meter Active",
			uuidStr:       uuidMeterActive,
			expectedFound: true,
			expectedRoom:  "Living Room",
			expectedCtrl:  "Power Meter",
			expectedState: "active",
		},
		{
			name:          "InfoAnalog Value",
			uuidStr:       uuidInfoAnalog,
			expectedFound: true,
			expectedRoom:  "Living Room",
			expectedCtrl:  "Temp Sensor",
			expectedState: "value",
		},
		{
			name:          "ColorPicker Color",
			uuidStr:       uuidColorRGB,
			expectedFound: true,
			expectedRoom:  "Living Room",
			expectedCtrl:  "LED Strip",
			expectedState: "color",
		},
		{
			name:          "Unknown UUID",
			uuidStr:       "99999999-9999-9999-9999-999999999999",
			expectedFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use the same parsing logic as the registry to ensure consistency
			// or use a standard UUID string for the test input.
			// Since we want to test Loxone-style UUIDs, let's use the helper or standard Parse if format allows.
			u, err := ParseUUID(tt.uuidStr)
			if err != nil {
				// Fallback if the test data was meant to be invalid
				if tt.expectedFound {
					t.Fatalf("Failed to parse UUID in test setup: %v", err)
				}
			}
			
			state, found := registry.LookupState(u)

			assert.Equal(t, tt.expectedFound, found)
			if found {
				assert.Equal(t, tt.expectedRoom, state.RoomName)
				assert.Equal(t, tt.expectedCtrl, state.Control.Name)
				assert.Equal(t, tt.expectedState, state.Name)
			}
		})
	}
}

func TestRegistry_LookupByPath(t *testing.T) {
	// Re-use logic or construct minimal registry
	// ... (Setup similar to above)
	// For brevity, skipping full re-setup, assuming NewRegistry works as tested above.
	// But let's test the path lookup specifically with a fresh struct.

	uuidSwitch := "10000000-0000-0000-0000000000000001"
	structure := &loxone.LoxApp3{
		Rooms: map[string]*loxone.Room{
			"room1": {Name: "Living Room", UUID: "room1"},
		},
		Controls: map[string]*loxone.Control{
			"c1": {
				Name: "Light",
				Room: "room1",
				States: map[string]interface{}{
					"active": uuidSwitch,
				},
			},
		},
	}
	registry := NewRegistry(structure)

	// Registry sanitizes names: "Living Room" -> "living-room", "Light" -> "light"
	state, found := registry.LookupStateByPath("Living Room", "Light", "active")
	assert.True(t, found)
	assert.NotNil(t, state)
	if state != nil {
		// Parse the expected UUID using the same helper to get the normalized string
		expectedUUID, _ := ParseUUID(uuidSwitch)
		assert.Equal(t, expectedUUID.String(), state.UUID.String())
	}

	// Test case insensitive / sanitized lookup
	state2, found2 := registry.LookupStateByPath("living room", "light", "active")
	assert.True(t, found2)
	assert.Equal(t, state, state2)
}
