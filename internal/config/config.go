package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type LoxoneConfig struct {
	IP   string `envconfig:"LOXONE_IP" required:"true"`
	User string `envconfig:"LOXONE_USER" required:"true"`
	Pass string `envconfig:"LOXONE_PASS" required:"true"`
	Snr  string `envconfig:"LOXONE_SNR" required:"true"`
}

type MQTTConfig struct {
	Host        string `envconfig:"MQTT_HOST" default:"localhost"`
	Port        int    `envconfig:"MQTT_PORT" default:"1883"`
	Protocol    string `envconfig:"MQTT_PROTOCOL" default:"tcp"`
	Path        string `envconfig:"MQTT_PATH"`
	ClientID    string `envconfig:"MQTT_CLIENT_ID" default:"lox-bridge"`
	User        string `envconfig:"MQTT_USER"`
	Pass        string `envconfig:"MQTT_PASS"`
	TopicPrefix string `envconfig:"MQTT_TOPIC_PREFIX" default:"lox"`
}

func (c *MQTTConfig) Validate() error {
	switch c.Protocol {
	case "tcp", "ssl", "ws", "wss":
		// valid
	default:
		return fmt.Errorf("invalid MQTT protocol: %s (must be tcp, ssl, ws, or wss)", c.Protocol)
	}

	if (c.Protocol == "ws" || c.Protocol == "wss") && c.Path == "" {
		c.Path = "/mqtt"
	}

	return nil
}

type SystemConfig struct {
	LogLevel string `envconfig:"LOG_LEVEL" default:"info"`
}

type Config struct {
	Loxone LoxoneConfig
	MQTT   MQTTConfig
	System SystemConfig
}

func Load() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("failed to process env vars: %w", err)
	}

	if err := cfg.MQTT.Validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}
