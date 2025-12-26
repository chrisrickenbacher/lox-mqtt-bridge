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

The fastest way to get started is using the pre-built Docker image.

### 1. Create a `docker-compose.yml`
```yaml
services:
  bridge:
    image: chric/lox-mqtt-bridge:latest
    container_name: lox-mqtt-bridge
    restart: unless-stopped
    environment:
      - LOXONE_IP=192.168.1.10
      - LOXONE_USER=admin
      - LOXONE_PASS=your-password
      - LOXONE_SNR=504F94D07F5F
      - MQTT_HOST=192.168.1.50
      - MQTT_PORT=1883
```

### 2. Start the bridge
```bash
docker compose up -d
```

Check [Configuration](docs/USER_GUIDE.md#configuration) for all available environment variables.

## Development

### Prerequisites
- Go 1.25+
- Docker & Docker Compose

### Build & Run Locally (for contributors)
1. **Clone the repository:**
   ```bash
   git clone https://github.com/chrisrickenbacher/lox-mqtt-bridge.git
   cd lox-mqtt-bridge
   ```
2. **Run with Docker Compose:**
   ```bash
   docker compose up --build -d
   ```
3. **Run binary only:**
   ```bash
   go mod tidy
   go run cmd/bridge/main.go
   ```

### Running Tests
This project uses `go test` and `testify` for unit testing. To run all tests:

```bash
go test ./...
```

For verbose output:
```bash
go test -v ./...
```

### Releasing a New Version
This project uses GitHub Actions to automate Docker image publishing and GitHub Releases.

1.  **Use Conventional Commits**: Ensure your commit messages follow the [Conventional Commits](https://www.conventionalcommits.org/) specification (e.g., `feat: add light group support`, `fix: reconnect logic`).
2.  **Tag the Release**: Create a semver tag starting with `v`:
    ```bash
    git tag v1.0.0
    git push origin v1.0.0
    ```
3.  **Automatic Workflow**: The push will trigger a workflow that:
    - Generates a categorized changelog based on your commit messages.
    - Creates a GitHub Release.
    - Builds and pushes multi-arch (`amd64`, `arm64`) Docker images to Docker Hub.

## Documentation

[Checkout the documentation](./docs/)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
