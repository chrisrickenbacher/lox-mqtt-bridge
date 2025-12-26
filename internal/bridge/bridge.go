package bridge

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/chrisrickenbacher/lox-mqtt-bridge/internal/config"
	"github.com/chrisrickenbacher/lox-mqtt-bridge/internal/loxone"
	"github.com/chrisrickenbacher/lox-mqtt-bridge/internal/mqtt"
)

// Bridge acts as the middleman between Loxone and MQTT
type Bridge struct {
	cfg      *config.Config
	lox      *loxone.Client
	mqtt     *mqtt.Client
	registry *Registry
	done     chan struct{}
}

// New creates a new Bridge instance
func NewBridge(cfg *config.Config) (*Bridge, error) {
	lox := loxone.NewClient(cfg.Loxone)
	mqttClient := mqtt.NewClient(cfg.MQTT)
	// Registry will be initialized in Start() after fetching structure

	return &Bridge{
		cfg:  cfg,
		lox:  lox,
		mqtt: mqttClient,
		done: make(chan struct{}),
	}, nil
}

// Start begins the bridging process
func (b *Bridge) Start(ctx context.Context) error {
	defer b.Stop()
	slog.Info("Starting Bridge...")

	if err := b.mqtt.Connect(); err != nil {
		return fmt.Errorf("failed to connect to MQTT: %v", err)
	}

	slog.Info("Connecting to Loxone...")
	if err := b.lox.Connect(); err != nil {
		return fmt.Errorf("failed to connect to Loxone: %v", err)
	}

	structure, err := b.lox.GetStructure()
	if err != nil {
		return fmt.Errorf("failed to get structure: %v", err)
	}

	b.registry = NewRegistry(structure)
	slog.Info("Registry initialized", "controls", len(structure.Controls))

	if err := b.lox.EnableStatusUpdates(); err != nil {
		b.lox.Close()
		return fmt.Errorf("failed to enable status updates: %v", err)
	}

	// 1. Publish Miniserver Info
	// Topic: <prefix>/<snr>/_info
	infoTopic := fmt.Sprintf("%s/%s/_info", b.cfg.MQTT.TopicPrefix, b.cfg.Loxone.Snr)
	infoPayload, _ := json.Marshal(structure.MsInfo)
	if err := b.mqtt.Publish(infoTopic, 1, true, infoPayload); err != nil {
		slog.Error("Failed to publish MS info", "error", err)
	}

	// 2. Publish Room Infos
	// Topic: <prefix>/<snr>/<room>/_info
	for _, room := range structure.Rooms {
		roomTopic := fmt.Sprintf("%s/%s/%s/_info", b.cfg.MQTT.TopicPrefix, b.cfg.Loxone.Snr, sanitize(room.Name))
		roomPayload, _ := json.Marshal(room)
		b.mqtt.Publish(roomTopic, 1, true, roomPayload)
	}

	// 3. Publish Control Infos
	// Topic: <prefix>/<snr>/<room>/<control>/_info
	// Iterate through registry's controlLookup to get all controls that are mapped
	for _, ctrl := range b.registry.controlLookup {
		// Lookup room name again to be safe/consistent
		roomName := "unknown"
		if r, ok := structure.Rooms[ctrl.Room]; ok {
			roomName = r.Name
		}

		ctrlTopic := fmt.Sprintf("%s/%s/%s/%s/_info",
			b.cfg.MQTT.TopicPrefix,
			b.cfg.Loxone.Snr,
			sanitize(roomName),
			sanitize(ctrl.Name))

		ctrlPayload, _ := json.Marshal(ctrl)
		b.mqtt.Publish(ctrlTopic, 1, true, ctrlPayload)
	}

	// Subscribe to Commands
	// Format: loxone/<snr>/<room>/<control>/command
	cmdTopic := fmt.Sprintf("%s/%s/+/+/command", b.cfg.MQTT.TopicPrefix, b.cfg.Loxone.Snr)

	if err := b.mqtt.Subscribe(cmdTopic, 1, b.handleMQTTMessage); err != nil {
		return fmt.Errorf("failed to subscribe to MQTT: %v", err)
	}

	return b.runEventLoop(ctx)
}

type Payload struct {
	Value interface{} `json:"value"`
	Ts    string      `json:"ts"`
}

func (b *Bridge) runEventLoop(ctx context.Context) error {
	slog.Info("Starting Event Loop...")
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-b.done:
			return nil
		case event := <-b.lox.Events:
			u, err := ParseUUID(event.UUID)
			if err != nil {
				slog.Error("Invalid UUID in event", "uuid", event.UUID, "error", err)
				continue
			}

			state, found := b.registry.LookupState(u)
			if !found {
				continue
			}

			// Topic: <prefix>/<snr>/<room>/<control>/<type>_<state>
			// Example: loxone/504.../living-room/light-switch/switch_active
			topic := fmt.Sprintf("%s/%s/%s/%s/%s_%s",
				b.cfg.MQTT.TopicPrefix,
				b.cfg.Loxone.Snr,
				sanitize(state.RoomName),
				sanitize(state.Control.Name),
				sanitize(state.Control.Type),
				sanitize(state.Name),
			)

			// Construct Payload
			payload := Payload{
				Value: event.Value,
				Ts:    time.Now().UTC().Format(time.RFC3339),
			}

			// Handle Text events specifically if needed,
			// event.Value is float64, event.Text is string.
			if event.Type == "Text" {
				payload.Value = event.Text
			}

			jsonPayload, err := json.Marshal(payload)
			if err != nil {
				slog.Error("Error marshaling payload", "error", err)
				continue
			}

			err = b.mqtt.Publish(topic, 0, true, jsonPayload)
			if err != nil {
				slog.Error("Failed to publish MQTT message", "error", err)
			}
		}
	}
}

func (b *Bridge) handleMQTTMessage(topic string, payload []byte) {
	// Expected: <prefix>/<snr>/<room>/<control>/command

	// Construct the root path: prefix/snr
	root := fmt.Sprintf("%s/%s", b.cfg.MQTT.TopicPrefix, b.cfg.Loxone.Snr)

	if !strings.HasPrefix(topic, root) {
		slog.Warn("Ignoring topic outside of configured root", "topic", topic)
		return
	}

	// Strip root to get: /room/control/command
	suffix := strings.TrimPrefix(topic, root)
	suffix = strings.TrimPrefix(suffix, "/")

	parts := strings.Split(suffix, "/")
	// Expected: room, control, command (3 parts)
	if len(parts) != 3 {
		slog.Warn("Ignoring malformed topic suffix", "suffix", suffix)
		return
	}

	if parts[2] != "command" {
		return
	}

	room := parts[0]
	control := parts[1]

	// Look up Control
	ctrl, found := b.registry.LookupControlByPath(room, control)
	if !found {
		slog.Warn("Command received for unknown control", "room", room, "control", control)
		return
	}

	// Important: We must send the command to the Control's UUIDAction
	targetUUID := ctrl.UUIDAction
	if targetUUID == "" {
		slog.Warn("Control has no UUIDAction", "control", ctrl.Name)
		return
	}

	val := string(payload)
	slog.Info("Sending command to Loxone", "control", ctrl.Name, "uuid", targetUUID, "type", ctrl.Type, "value", val)

	cmd := fmt.Sprintf("jdev/sps/io/%s/%s", targetUUID, val)
	if err := b.lox.SendCommand(cmd); err != nil {
		slog.Error("Failed to send command to Loxone", "error", err)
	}
}

func (b *Bridge) Stop() {
	select {
	case <-b.done:
	default:
		close(b.done)
	}
	b.lox.Close()
	b.mqtt.Close()
}