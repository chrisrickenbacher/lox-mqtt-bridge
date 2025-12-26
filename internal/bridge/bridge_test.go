package bridge

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/chrisrickenbacher/lox-mqtt-bridge/internal/config"
	"github.com/chrisrickenbacher/lox-mqtt-bridge/internal/loxone"
	"github.com/stretchr/testify/mock"
)

func TestBridge_Start_And_Events(t *testing.T) {
	// Setup Mocks
	mockLox := new(MockLoxoneProvider)
	mockMQTT := new(MockMQTTProvider)

	cfg := &config.Config{
		Loxone: config.LoxoneConfig{Snr: "504F94A00000"},
		MQTT:   config.MQTTConfig{TopicPrefix: "loxone"},
	}

	// Prepare Structure Data for Registry
	uuidActive := "10000000-0000-0000-0000-000000000001"
	structure := &loxone.LoxApp3{
		MsInfo: map[string]interface{}{"serial": "504F94A00000"},
		Rooms: map[string]*loxone.Room{
			"r1": {Name: "Living Room", UUID: "r1"},
		},
		Controls: map[string]*loxone.Control{
			"c1": {
				Name: "Light",
				Type: "Switch",
				Room: "r1",
				States: map[string]interface{}{
					"active": uuidActive,
				},
			},
		},
	}

	// --- Expectations ---

	// 1. Connect
	mockMQTT.On("Connect").Return(nil)
	mockLox.On("Connect").Return(nil)

	// 2. Get Structure & Enable Updates
	mockLox.On("GetStructure").Return(structure, nil)
	mockLox.On("EnableStatusUpdates").Return(nil)

	// 3. Info Publishing (Startup)
	// We expect multiple calls. We use loose matching for the "info" topics.
	mockMQTT.On("Publish", mock.MatchedBy(func(topic string) bool {
		return strings.HasSuffix(topic, "/_info")
	}), mock.Anything, mock.Anything, mock.Anything).Return(nil)

	// 4. Subscribe
	mockMQTT.On("Subscribe", "loxone/504F94A00000/+/+/command", mock.Anything, mock.Anything).Return(nil)

	// 5. Event Loop Setup
	events := make(chan loxone.Event, 1)
	mockLox.On("GetEvents").Return((<-chan loxone.Event)(events))

	// 6. Event Processing Expectation
	// When we send the event, we expect a SPECIFIC publish
	expectedTopic := "loxone/504F94A00000/living-room/light/switch_active"
	mockMQTT.On("Publish", expectedTopic, byte(0), true, mock.MatchedBy(func(payload []byte) bool {
		// Verify payload contains value
		var p map[string]interface{}
		json.Unmarshal(payload, &p)
		return p["value"] == 1.0
	})).Return(nil)

	// 7. Cleanup
	mockLox.On("Close").Return()
	mockMQTT.On("Close").Return()

	// --- Execution ---

	b := &Bridge{
		cfg:  cfg,
		lox:  mockLox,
		mqtt: mockMQTT,
		done: make(chan struct{}),
	}

	ctx, cancel := context.WithCancel(context.Background())

	// Run Start in background
	go func() {
		b.Start(ctx)
	}()

	// Allow startup to proceed
	time.Sleep(50 * time.Millisecond)

	// Inject Event
	events <- loxone.Event{
		UUID:  uuidActive,
		Value: 1.0,
		Type:  "Value",
	}

	// Allow event processing
	time.Sleep(50 * time.Millisecond)

	// Stop
	cancel()
	time.Sleep(20 * time.Millisecond) // Allow shutdown

	// Verify
	mockMQTT.AssertExpectations(t)
	mockLox.AssertExpectations(t)
}

func TestBridge_CommandHandling(t *testing.T) {
	// Setup
	mockLox := new(MockLoxoneProvider)
	mockMQTT := new(MockMQTTProvider)
	cfg := &config.Config{
		Loxone: config.LoxoneConfig{Snr: "504F94A00000"},
		MQTT:   config.MQTTConfig{TopicPrefix: "loxone"},
	}

	// Mock Structure with Action UUID
	uuidAction := "20000000-0000-0000-0000-000000000001"
	structure := &loxone.LoxApp3{
		Rooms: map[string]*loxone.Room{"r1": {Name: "Living Room"}},
		Controls: map[string]*loxone.Control{
			"c1": {
				Name:       "Light",
				Room:       "r1",
				Type:       "Switch",
				UUIDAction: uuidAction,
			},
		},
	}

	b := &Bridge{
		cfg:      cfg,
		lox:      mockLox,
		mqtt:     mockMQTT,
		registry: NewRegistry(structure), // Manually inject registry
	}

	// Expectation: SendCommand
	// Note: Validation of payload "On" depends on how handleMQTTMessage sends it.
	// It sends "On" as string.
	cmdStr := fmt.Sprintf("jdev/sps/io/%s/On", uuidAction)
	mockLox.On("SendCommand", cmdStr).Return(nil)

	// Trigger Command
	// Topic: loxone/504F94A00000/living-room/light/command
	topic := "loxone/504F94A00000/living-room/light/command"
	payload := []byte("On")

	b.handleMQTTMessage(topic, payload)

	mockLox.AssertExpectations(t)
}

// Test case for ignoring invalid topics
func TestBridge_CommandHandling_Ignored(t *testing.T) {
	mockLox := new(MockLoxoneProvider)
	mockMQTT := new(MockMQTTProvider)
	cfg := &config.Config{
		Loxone: config.LoxoneConfig{Snr: "504F94A00000"},
		MQTT:   config.MQTTConfig{TopicPrefix: "loxone"},
	}
	b := &Bridge{cfg: cfg, lox: mockLox, mqtt: mockMQTT}

	// No expectations on mockLox because it should NOT be called

	// Invalid Topic
	b.handleMQTTMessage("loxone/other/command", []byte("On"))
	// Not a command topic
	b.handleMQTTMessage("loxone/504F94A00000/room/light/state", []byte("On"))

	mockLox.AssertExpectations(t)
}
