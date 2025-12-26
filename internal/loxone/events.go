package loxone

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log/slog"
)

// Event represents a parsed Loxone event
type Event struct {
	UUID  string
	Value float64
	Text  string
	Type  string // "Value", "Text", "Daytimer", "Weather"
}

// EnableStatusUpdates tells the Miniserver to start sending events
func (c *Client) EnableStatusUpdates() error {
	slog.Info("Enabling Status Updates...")
	if err := c.SendCommand("jdev/sps/enablebinstatusupdate"); err != nil {
		return err
	}
	_, err := c.WaitForResponse("enablebinstatusupdate")
	return err
}

// parseLoxoneUUID converts Loxone's Little Endian binary UUID to a standard string
// Structure: Data1 (4b LE), Data2 (2b LE), Data3 (2b LE), Data4 (8b)
func parseLoxoneUUID(b []byte) string {
	if len(b) != 16 {
		return ""
	}
	d1 := binary.LittleEndian.Uint32(b[0:4])
	d2 := binary.LittleEndian.Uint16(b[4:6])
	d3 := binary.LittleEndian.Uint16(b[6:8])
	d4 := b[8:]

	return fmt.Sprintf("%08x-%04x-%04x-%02x%02x-%02x%02x%02x%02x%02x%02x",
		d1, d2, d3, d4[0], d4[1], d4[2], d4[3], d4[4], d4[5], d4[6], d4[7])
}

// HandleBinaryMessage parses binary value events (Type 2)
func (c *Client) HandleBinaryMessage(message []byte) {
	reader := bytes.NewReader(message)
	// Value Event: 16 bytes UUID + 8 bytes float64 = 24 bytes
	for reader.Len() >= 24 {
		uuidBytes := make([]byte, 16)
		if _, err := reader.Read(uuidBytes); err != nil {
			slog.Error("Error reading UUID", "error", err)
			break
		}
		var val float64
		if err := binary.Read(reader, binary.LittleEndian, &val); err != nil {
			slog.Error("Error reading value", "error", err)
			break
		}

		uuidStr := parseLoxoneUUID(uuidBytes)

		slog.Debug("Parsed Value Event", "uuid", uuidStr, "value", val)

		select {
		case c.Events <- Event{
			UUID:  uuidStr,
			Value: val,
			Type:  "Value",
		}:
		default:
			slog.Warn("Events channel full, dropping event")
		}
	}
}

// HandleTextMessage parses text events (Type 3)
func (c *Client) HandleTextMessage(message []byte) {
	reader := bytes.NewReader(message)
	// Text Event: 16 bytes UUID + 16 bytes Icon UUID + 4 bytes Length (n) + n bytes Text
	for reader.Len() >= 36 {
		uuidBytes := make([]byte, 16)
		if _, err := reader.Read(uuidBytes); err != nil {
			break
		}
		// Skip Icon UUID (16 bytes)
		if _, err := reader.Seek(16, io.SeekCurrent); err != nil {
			break
		}
		var textLen uint32
		if err := binary.Read(reader, binary.LittleEndian, &textLen); err != nil {
			break
		}

		textBuf := make([]byte, textLen)
		if n, err := reader.Read(textBuf); err != nil || uint32(n) != textLen {
			break
		}

		// Text events are padded to 4-byte boundaries
		padding := (4 - (textLen % 4)) % 4
		if padding > 0 {
			reader.Seek(int64(padding), io.SeekCurrent)
		}

		uuidStr := parseLoxoneUUID(uuidBytes)

		slog.Debug("Parsed Text Event", "uuid", uuidStr, "text", string(textBuf))

		select {
		case c.Events <- Event{
			UUID: uuidStr,
			Text: string(textBuf),
			Type: "Text",
		}:
		default:
		}
	}
}
