# User Guide: Loxone MQTT Bridge

## Configuration

The bridge is configured exclusively via Environment Variables.

### Loxone Config

| Variable | Description | Example |
|---|---|---|
| `LOXONE_IP` | Local IP of your Miniserver | `192.168.1.10` |
| `LOXONE_USER` | User with Web/App access | `admin` |
| `LOXONE_PASS` | Password | `password` |
| `LOXONE_SNR` | Serial Number (**MANDATORY** for TLS certificate generation) | `504F94D0F02C` |

**Note:** The bridge automatically constructs the secure local hostname (e.g., `192-168-1-10.snr.dyndns.loxonecloud.com`) to enable TLS (WSS) connections. This avoids certificate errors.

### MQTT Configuration
| Variable | Description | Default |
|---|---|---|
| `MQTT_HOST` | MQTT Broker Host | `localhost` |
| `MQTT_PORT` | MQTT Broker Port | `1883` |
| `MQTT_PROTOCOL` | MQTT Protocol (`tcp`, `ssl`, `ws`, `wss`) | `tcp` |
| `MQTT_PATH` | Path for WebSocket connections | `/mqtt` (if ws/wss) |
| `MQTT_CLIENT_ID` | MQTT Client Identifier | `lox-bridge` |
| `MQTT_USER` | MQTT Username | *(Empty)* |
| `MQTT_PASS` | MQTT Password | *(Empty)* |
| `MQTT_TOPIC_PREFIX` | Base topic for bridge messages | `lox` |

### System Configuration
| Variable | Description | Default |
|---|---|---|
| `LOG_LEVEL` | Logging verbosity (`debug`, `info`, `warn`, `error`) | `info` |

### Example `docker-compose.yml`
```yaml
services:
  bridge:
    image: lox-mqtt-bridge
    environment:
      - LOXONE_IP=192.168.1.10
      - LOXONE_USER=admin
      - LOXONE_PASS=securepassword
      - LOXONE_SNR=504F94D0F02C
      - MQTT_HOST=vernemq
      - MQTT_PORT=1883
```

## Usage & Topic Structure

The bridge follows a strictly normalized topic structure based on the Miniserver's Room and Device names. All names are **sanitized** (lowercased, spaces replaced with hyphens).

Check [Architecture](docs/ARCHITECTURE.md) for detailed topic structure. 

Check [Reference](docs/REFERENCE.md) for supported Loxone Control Types and their commands/states.
