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

The bridge acts as a bidirectional gateway. It publishes state changes from Loxone to MQTT and listens for commands on MQTT to control Loxone devices.

### 1. Topic Hierarchy & States (Read-Only)
The bridge follows a strictly normalized topic structure based on the Miniserver's Room and Device names. All names are **sanitized** (lowercased, spaces replaced with hyphens).

*   **Structure:** Please refer to [Architecture > Topic Structure](ARCHITECTURE.md#5-topic-structure) for the complete definition of how topics are constructed (e.g., `lox/504F.../living-room/ceiling-light/...`).
*   **Data Types:** Please refer to [Reference](REFERENCE.md) for a complete list of **Control Types** (like `Switch`, `Dimmer`, `Jalousie`) and exactly which state topics (e.g., `switch_active`, `dimmer_position`) are available for each.

### 2. Controlling Devices (Commands)
To control a device, you publish a message to its specific **command topic**.

**Command Topic Format:**
`<topic-prefix>/<serial-number>/<room>/<control-name>/command`

**Payload:**
The payload is the raw value or command string you want to send to the Loxone control.

#### Examples:

**1. Turn on a Light (Switch)**
*   **Topic:** `lox/504F94A00000/kitchen/ceiling-light/command`
*   **Payload:** `On`
*   *(See [Reference > Switch](REFERENCE.md#switch) for available commands like `On`, `Off`, `Pulse`)*

**2. Set Dimmer to 50%**
*   **Topic:** `lox/504F94A00000/living-room/spots/command`
*   **Payload:** `50`
*   *(See [Reference > Dimmer](REFERENCE.md#dimmer) for details)*

**3. Open Blinds (Jalousie)**
*   **Topic:** `lox/504F94A00000/bedroom/blinds/command`
*   **Payload:** `FullOpen`
*   *(See [Reference > Jalousie](REFERENCE.md#jalousie) for details)*

**Note:** The bridge does not immediately update the state topic upon receiving a command. It sends the command to the Miniserver and waits for the Miniserver to push the new state back. This ensures the MQTT state always reflects the *actual* device state.
