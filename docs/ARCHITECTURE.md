# Loxone MQTT Bridge: Architecture & Concept

## 1. Project Goal
The primary goal is to create a bidirectional bridge between a Loxone Miniserver and an MQTT broker (specifically VerneMQ). This enables monitoring and control of Loxone devices through MQTT, facilitating integration with other smart home platforms (Node-RED, Home Assistant, etc.).

## 2. Architecture Overview
The system consists of three main components:
1.  **Loxone Miniserver:** The central controller.
2.  **Loxone-MQTT-Bridge:** A standalone Go application running in Docker.
3.  **MQTT Broker:** VerneMQ (supporting both TCP and WebSockets).

## 1. Connection Strategy

The bridge exclusively uses **Secure WebSockets (WSS)** via the official Loxone CloudDNS hostname to ensure a trusted TLS connection, even on local networks.

### Hostname Construction
Instead of connecting to the raw IP, the bridge constructs a hostname:
`wss://{cleaned-ip}.{snr}.dyndns.loxonecloud.com/ws/rfc6455`

*   `{cleaned-ip}`: Local IP with dots `.` replaced by hyphens `-` (e.g., `192-168-1-10`).
*   `{snr}`: Miniserver Serial Number (from config).

### Protocol
*   **Scheme:** `wss` (Port 443 implicit).
*   **Certificate:** Valid Loxone CloudDNS certificate (trusted by default root CAs).

## 2. Authentication Flow

The bridge implements the **Token-Based Authentication** flow (Loxone PDF v16.0).
Even though WSS provides transport security, the Miniserver **requires** Application Layer Encryption (AES-256) for the Token Acquisition command.

1.  **Public Key Acquisition:** `HTTPS GET /jdev/sys/getPublicKey`
    *   Returns the Miniserver's X.509 certificate/public key.
2.  **Session Key Generation:**
    *   Client generates a random AES-256 Key and IV.
    *   Client RSA-encrypts the "Key:IV" string using the Miniserver's Public Key.
3.  **Key Exchange:**
    *   Send via WebSocket: `jdev/sys/keyexchange/{encrypted-session-key}`.
    *   From this point on, commands can be encrypted using the AES Key.
4.  **Salt Acquisition:**
    *   Send: `jdev/sys/getkey2/{user}` (Often sent as `jdev/sys/fenc/...` for privacy).
    *   Response: `{key}`, `{salt}`, `{hashAlg}` (e.g., SHA1 or SHA256).
    *   *Note: This `key` is a one-time key for hashing, distinct from the AES Session Key.*
5.  **Hashing:**
    *   **Password Hash:** `Hash(password:salt)` (using `hashAlg`).
        *   Result must be **Uppercase Hex**.
    *   **Authentication Hash:** `HMAC(key, "user:pwHash")` (using `hashAlg`).
        *   Result must be **Hex**.
6.  **Token Request (Encrypted):**
    *   Command: `jdev/sys/getjwt/{hash}/{user}/{permission}/{uuid}/{info}`.
    *   **Permission:** `4` (App) is used to request a long-lived token (weeks).
    *   **MUST** be encrypted using the AES Session Key.
    *   Send: `jdev/sys/enc/{encrypted-command}`.
    *   Response: `{token}`, `{validUntil}`, `{tokenRights}`, `{key}`.
7. Authentication Complete:
    *   The `token` is stored.
    *   Subsequent commands do not need full re-authentication but must handle token expiration/refresh.
8.  **Structure File:**
    *   Client requests `data/LoxAPP3.json` over the established WebSocket.
    *   Uses in-memory caching: checks `jdev/sps/LoxAPPversion3` before downloading the full file.

## 4. MQTT Integration
-   **Broker:** VerneMQ.
-   **Protocol:** WebSockets (preferred) or TCP.
-   **QoS:** Level 1 (At least once) for state updates to ensure delivery.
-   **Retain:** `true` for state messages. Clients subscribing will immediately receive the last known state.

## 5. Topic Structure
We utilize a **granular topic structure** to allow for precise subscriptions and atomic updates.

### Format

- `<topic-prefix>/<serial-number>/_info`: Miniserver general info.
- `<topic-prefix>/<serial-number>/<room>/_info`: Room-specific info.
- `<topic-prefix>/<serial-number>/<room>/<control-name>/<control-type>_<state>`: Read only state of a specific control.
- `<topic-prefix>/<serial-number>/<room>/<control-name>/_info`: Control metadata/info.
- `<topic-prefix>/<serial-number>/<room>/<control-name>/command`: Command topic for controlling the control.


*   `<topic-prefix>`: Configurable prefix (default: `lox`).
*   `<serial-number>`: The serial number of the Miniserver (e.g., `504F94A00000`).
*   `<room>`: Normalized room name (e.g., `living-room`).
*   `<control-name>`: User-defined control name (e.g., `ceiling-light`).
*   `<control-type>_<state>`: State for the specific control type (e.g., `switch_active`, `pushbutton_active`, `slider_value`)
    * `<control-type>`: Type of control, see [Loxone Control Types](docs/Loxone_Control_types.md)
    * `<state>`: Specific state of the control, see [Loxone Control Types](docs/Loxone_Control_types.md) for details.

> All endpoints are read only except the `command` topics, which accept commands to control the respective Loxone device.


## 6. Data Flow

### 6.1. Loxone to MQTT (Events)
1.  Bridge connects to Loxone WebSocket.
2.  Parses `LoxApp3.json` to build a `UUID <-> Topic` lookup map.
3.  Subscribes to status updates.
4.  On Event (UUID, Value):
    *   Look up the corresponding Topic Parts for the UUID.
    *   Construct the topic: `.../<function-key>/state`.
    *   Publish payload to MQTT (Retained).

### 6.2. MQTT to Loxone (Commands)
1.  Bridge subscribes to `<topic-prefix>/<serial-number>/+/+/command`.
2.  On Message:
    *   Parses the topic to extract the `Device` and `Function`.
    *   Uses the lookup map to find the corresponding Loxone **Action UUID**.
    *   Sends a WebSocket command: `jdev/sps/io/<UUID>/<Value>`.

## 7. Loop Prevention & State Management
*   **Internal State:** The bridge maintains a cache of the last known values.
*   **Command Handling:** When a command arrives via MQTT, it is passed to Loxone. The bridge relies on the subsequent Loxone Event to update the MQTT `state` topic, ensuring the `state` topic always reflects the *actual* confirmation from the Miniserver, not just the *intent* from the command.

## 8. Configuration
The application is configured strictly via **Environment Variables**.
We use `kelseyhightower/envconfig` to map these variables to the internal Go configuration struct.

*   **Loxone:** `LOXONE_IP`, `LOXONE_USER`, `LOXONE_PASS`, `LOXONE_SNR`.
*   **MQTT:** `MQTT_HOST`, `MQTT_PORT`, `MQTT_PROTOCOL`, `MQTT_PATH`, `MQTT_CLIENT_ID`, `MQTT_USER`, `MQTT_PASS`.
    *   `MQTT_PATH`: Optional path for WebSocket connections (default: `/mqtt` if protocol is `ws` or `wss`).
*   **System:** `LOG_LEVEL`.

## 9. Dockerization
*   **Image:** Lightweight (Alpine or Distroless).
*   **Build:** Multi-stage Dockerfile to keep the final image size minimal.
