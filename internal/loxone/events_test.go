package loxone

import (
	"bytes"
	"encoding/binary"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseLoxoneUUID(t *testing.T) {
	// Construct a known UUID in Little Endian Loxone format
	// Target: 504f94a0-1234-5678-0102030405060708
	
	// Data1: 504F94A0 -> LE: A0 94 4F 50
	// Data2: 1234 -> LE: 34 12
	// Data3: 5678 -> LE: 78 56
	// Data4: 01 02 03 04 05 06 07 08 (Array, not reversed usually in Loxone UUIDs for the last part? 
	// The implementation takes b[8:] and prints them in order. So standard big endian for the node part?)
	
	input := []byte{
		0xA0, 0x94, 0x4F, 0x50, // d1
		0x34, 0x12,             // d2
		0x78, 0x56,             // d3
		0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, // d4
	}

	expected := "504f94a0-1234-5678-0102-030405060708"
	
	result := parseLoxoneUUID(input)
	assert.Equal(t, expected, result)

	// Test invalid length
	assert.Equal(t, "", parseLoxoneUUID([]byte{0x00}))
}

func TestClient_HandleBinaryMessage(t *testing.T) {
	c := &Client{
		Events: make(chan Event, 10),
	}

	// Construct message
	// UUID (16 bytes) + Value (8 bytes float64)
	buf := new(bytes.Buffer)
	
	// UUID: 504f94a0-1234-5678-0102-030405060708
	uuidBytes := []byte{
		0xA0, 0x94, 0x4F, 0x50,
		0x34, 0x12,
		0x78, 0x56,
		0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
	}
	buf.Write(uuidBytes)
	
	// Value: 123.456
	binary.Write(buf, binary.LittleEndian, float64(123.456))

	c.HandleBinaryMessage(buf.Bytes())

	select {
	case e := <-c.Events:
		assert.Equal(t, "504f94a0-1234-5678-0102-030405060708", e.UUID)
		assert.Equal(t, 123.456, e.Value)
		assert.Equal(t, "Value", e.Type)
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Timeout waiting for event")
	}
}

func TestClient_HandleTextMessage(t *testing.T) {
	c := &Client{
		Events: make(chan Event, 10),
	}

	// UUID (16) + Icon (16) + Len (4) + Text (n) + Padding
	buf := new(bytes.Buffer)

	// UUID
	uuidBytes := make([]byte, 16)
	uuidBytes[0] = 0xFF // Just a marker
	buf.Write(uuidBytes)

	// Icon (ignored)
	buf.Write(make([]byte, 16))

	text := "Hello World"
	textLen := uint32(len(text)) // 11

	binary.Write(buf, binary.LittleEndian, textLen)
	buf.WriteString(text)

	// Padding: 11 % 4 = 3. 4 - 3 = 1 byte padding needed.
	buf.WriteByte(0x00)

	c.HandleTextMessage(buf.Bytes())

	select {
	case e := <-c.Events:
		expectedUUID := parseLoxoneUUID(uuidBytes)
		assert.Equal(t, expectedUUID, e.UUID)
		assert.Equal(t, text, e.Text)
		assert.Equal(t, "Text", e.Type)
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Timeout waiting for event")
	}
}
