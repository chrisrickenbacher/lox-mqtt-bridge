package bridge

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/chrisrickenbacher/lox-mqtt-bridge/internal/loxone"
	"github.com/google/uuid"
)

// Registry holds the mapping between UUIDs and Controls
type Registry struct {
	states        map[uuid.UUID]State
	rooms         map[string]*loxone.Room
	lookup        map[string]uuid.UUID       // Key: "room/control/function"
	controlLookup map[string]*loxone.Control // Key: "room/control"
}

// State represents a specific state of a control (e.g. "value", "temp", "active")
type State struct {
	Control  *loxone.Control
	Name     string
	UUID     uuid.UUID
	RoomName string
}

// NewRegistry creates a new Registry and processes the structure
func NewRegistry(structure *loxone.LoxApp3) *Registry {
	r := &Registry{
		states:        make(map[uuid.UUID]State),
		rooms:         make(map[string]*loxone.Room),
		lookup:        make(map[string]uuid.UUID),
		controlLookup: make(map[string]*loxone.Control),
	}

	if structure != nil {
		r.rooms = structure.Rooms
		r.processControls(structure.Controls)
	}

	// Debug logging
	slog.Debug("Registry loaded", "states", len(r.states), "lookup_keys", len(r.lookup))
	// Log first 5 for verification
	i := 0
	for u, s := range r.states {
		if i >= 5 {
			break
		}
		slog.Debug("Registered State", "uuid", u, "room", s.RoomName, "control", s.Control.Name, "state", s.Name)
		i++
	}

	return r
}

func (r *Registry) processControls(controls map[string]*loxone.Control) {
	for _, ctrl := range controls {
		roomName := "unknown"
		if room, ok := r.rooms[ctrl.Room]; ok {
			roomName = room.Name
		}

		// Populate control lookup
		// Key format: sanitized room/control
		ctrlKey := fmt.Sprintf("%s/%s", sanitize(roomName), sanitize(ctrl.Name))
		r.controlLookup[ctrlKey] = ctrl

		for stateName, uuidVal := range ctrl.States {
			switch v := uuidVal.(type) {
			case string:
				u, err := ParseUUID(v)
				if err == nil {
					state := State{
						Control:  ctrl,
						Name:     stateName,
						UUID:     u,
						RoomName: roomName,
					}
					r.states[u] = state

					// Populate lookup map
					// Key format: sanitized room/control/function
					key := fmt.Sprintf("%s/%s/%s", sanitize(roomName), sanitize(ctrl.Name), stateName)
					r.lookup[key] = u
				} else {
					slog.Warn("Failed to parse UUID", "control", ctrl.Name, "state", stateName, "value", v, "error", err)
				}
			case []interface{}:
				for _, item := range v {
					if s, ok := item.(string); ok {
						u, err := ParseUUID(s)
						if err == nil {
							state := State{
								Control:  ctrl,
								Name:     stateName,
								UUID:     u,
								RoomName: roomName,
							}
							r.states[u] = state

							// For arrays, we might need a strategy. For now, we overwrite or ignore.
							// Usually these are alternative UUIDs or specific sub-functions.
							// Let's index the first one or all if possible.
							// Given the 1:1 map, last one wins.
							key := fmt.Sprintf("%s/%s/%s", sanitize(roomName), sanitize(ctrl.Name), stateName)
							r.lookup[key] = u
						}
					}
				}
			}
		}

		if ctrl.SubControls != nil {
			r.processControls(ctrl.SubControls)
		}
	}
}

func ParseUUID(s string) (uuid.UUID, error) {
	// Loxone UUIDs in LoxAPP3.json often use 8-4-4-16 format (35 chars)
	// or 8-4-4-4-12 (36 chars).
	// The most robust way is to strip hyphens and parse the 32-char hex string.
	clean := ""
	for _, r := range s {
		if (r >= '0' && r <= '9') || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F') {
			clean += string(r)
		}
	}
	return uuid.Parse(clean)
}

// LookupState returns the state for a given UUID
func (r *Registry) LookupState(u uuid.UUID) (*State, bool) {
	s, ok := r.states[u]
	if !ok {
		return nil, false
	}
	return &s, true
}

// LookupStateByPath finds a State by room, control, and function name
func (r *Registry) LookupStateByPath(room, control, function string) (*State, bool) {
	// keys are stored sanitized
	key := fmt.Sprintf("%s/%s/%s", sanitize(room), sanitize(control), function)
	u, ok := r.lookup[key]
	if !ok {
		return nil, false
	}
	return r.LookupState(u)
}

// LookupControlByPath finds a Control by room and control name
func (r *Registry) LookupControlByPath(room, control string) (*loxone.Control, bool) {
	key := fmt.Sprintf("%s/%s", sanitize(room), sanitize(control))
	c, ok := r.controlLookup[key]
	return c, ok
}

// sanitize replaces spaces with hyphens and lowercases the string
func sanitize(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	return s
}
