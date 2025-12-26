package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMQTTConfig_Validate(t *testing.T) {
	tests := []struct {
		name        string
		cfg         MQTTConfig
		expectedErr bool
		expectedPath string // Check if path was defaulted
	}{
		{
			name: "Valid TCP",
			cfg: MQTTConfig{
				Protocol: "tcp",
			},
			expectedErr: false,
		},
		{
			name: "Valid WSS with Path",
			cfg: MQTTConfig{
				Protocol: "wss",
				Path:     "/ws",
			},
			expectedErr: false,
			expectedPath: "/ws",
		},
		{
			name: "Valid WSS Default Path",
			cfg: MQTTConfig{
				Protocol: "wss",
				Path:     "",
			},
			expectedErr: false,
			expectedPath: "/mqtt",
		},
		{
			name: "Invalid Protocol",
			cfg: MQTTConfig{
				Protocol: "http",
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if tt.expectedPath != "" {
				assert.Equal(t, tt.expectedPath, tt.cfg.Path)
			}
		})
	}
}
