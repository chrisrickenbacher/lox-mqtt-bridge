<div align="center">

# Loxone MQTT Bridge

A high-performance, bidirectional bridge between **Loxone Miniserver** and **MQTT Brokers**. Written in Go, it runs as a lightweight Docker container.

<h3>

[User Guide](./docs/USER_GUIDE.md) | [Reference](./docs/REFERENCE.md) | [Architecture](./docs/ARCHITECTURE.md) 

</h3>

</div>

## Features

- **Bidirectional Bridging**
    - **Loxone to MQTT:** High-performance streaming of state updates via Loxone's binary event stream.
    - **MQTT to Loxone:** Precise control of devices using command topics mapped to Loxone action UUIDs.
    - **Loop Prevention:** State topics only reflect actual confirmations from the Miniserver, ensuring the source of truth is maintained.

- **Automatic Discovery & Mapping**
    - **Smart Registry:** Automatically fetches and parses `LoxAPP3.json` to build a human-readable topic map.
    - **Granular Topics:** Slugs for Rooms and Controls (e.g., `living-room/ceiling-light/switch_active`).
    - **Metadata Publishing:** Publishes detailed metadata (JSON) for the Miniserver, Rooms, and individual Controls to specific `/_info` topics.
    - **Efficient Sync:** Utilizes in-memory caching and `LoxAPPversion3` checks to minimize structure file downloads.

- **Security & Connectivity**
    - **Modern Authentication:** Implements Loxone's Token-Based Authentication (v16.0).
    - **Transport Security:** Exclusive use of **Secure WebSockets (WSS)** via Loxone CloudDNS hostnames for trusted TLS certificates.
    - **App-Layer Encryption:** RSA and AES-256 (CBC) encryption used during the sensitive token acquisition flow.
    - **MQTT Resilience:** Supports both TCP and WebSockets with configurable QoS 1 and Retain flags for persistent state.

## Requirements

- **Loxone Miniserver:** Generation 2 (or Miniserver Compact).
- **Loxone Firmware:** Version 12.0 or higher (Supports TLS/WSS without application-layer encryption).
- **Network:** The Miniserver must be reachable via the local network.

## Quick Start (Docker)

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/chrisrickenbacher/lox-mqtt-bridge.git
    cd lox-mqtt-bridge
    ```

2.  **Configure:**
    Edit `docker-compose.yml` to set environment variables. Check [Configuration](docs/USER_GUIDE.md#configuration) for details.

3.  **Run:**
    ```bash
    docker-compose up --build -d
    ```

4.  **Verify:**
    Check logs: `docker-compose logs -f bridge`

## Development

### Prerequisites
- Go 1.25+
- Docker & Docker Compose

### Build & Run Locally
```bash
go mod tidy
go run cmd/bridge/main.go
```

## Documentation

[Checkout the documentation](./docs/)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
