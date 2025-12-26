package loxone

import (
	"crypto/tls"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/chrisrickenbacher/lox-mqtt-bridge/internal/config"
	"github.com/gorilla/websocket"
)

// Response represents a generic Loxone WebSocket/HTTP response wrapper
type Response struct {
	LL struct {
		Control string          `json:"control"`
		Value   json.RawMessage `json:"value"`
		Code    interface{}     `json:"Code"` // Can be int or string depending on version
	} `json:"LL"`
}

// ApiKeyStruct represents the data from /jdev/cfg/apiKey
type ApiKeyStruct struct {
	Snr     string `json:"snr"`
	Version string `json:"version"`
}

// Token represents the authentication token
type Token struct {
	Token        string `json:"token"`
	ValidUntil   int64  `json:"validUntil"`
	TokenRights  int    `json:"tokenRights"`
	UnsecurePass bool   `json:"unsecurePass"`
	Key          string `json:"key"`
}

type Header struct {
	Type   byte
	Info   byte
	Length uint32
}

// Client handles the connection to the Loxone Miniserver
type Client struct {
	cfg        config.LoxoneConfig
	conn       *websocket.Conn
	httpClient *http.Client

	// Concurrency
	mu          sync.Mutex
	done        chan struct{}
	isConnected bool
	isClosed    bool

	// Channels
	msgChan chan []byte
	Events  chan Event

	// Header state for the read loop
	header *Header

	// Connection info
	host string

	// Structure Cache
	structureCache   *LoxApp3
	structureLastMod string

	// Crypto state
	token *Token
}

// NewClient creates a new Loxone client
func NewClient(cfg config.LoxoneConfig) *Client {
	// Create HTTP client with insecure skip verify for local IP connections
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient := &http.Client{
		Transport: tr,
		Timeout:   10 * time.Second,
	}

	return &Client{
		cfg:        cfg,
		msgChan:    make(chan []byte, 100),
		Events:     make(chan Event, 1000),
		done:       make(chan struct{}),
		httpClient: httpClient,
	}
}

// Connect establishes the connection to the Miniserver
func (c *Client) Connect() error {
	c.mu.Lock()

	if c.cfg.Snr == "" {
		c.mu.Unlock()
		return fmt.Errorf("serial number (snr) is required for TLS connection")
	}

	// Always use Loxone CloudDNS Hostname for Local TLS
	// Format: {cleaned-ip}.{snr}.dyndns.loxonecloud.com
	cleanIP := strings.ReplaceAll(c.cfg.IP, ".", "-")
	cleanIP = strings.ReplaceAll(cleanIP, ":", "-") // Basic IPv6 handling
	c.host = fmt.Sprintf("%s.%s.dyndns.loxonecloud.com", cleanIP, c.cfg.Snr)
	port := 443
	scheme := "wss"

	slog.Info("Connecting to Loxone", "host", c.host, "port", port, "scheme", scheme)

	// 1. Check Reachability via HTTP
	if err := c.checkReachability(scheme, c.host, port); err != nil {
		c.mu.Unlock()
		return fmt.Errorf("miniserver unreachable: %v", err)
	}

	// 2. Dial WebSocket
	// Use hostname without explicit port for standard schemes (wss implies 443) to match browser behavior
	u := url.URL{Scheme: scheme, Host: c.host, Path: "/ws/rfc6455"}
	slog.Debug("Dialing WebSocket", "url", u.String())

	dialer := websocket.DefaultDialer
	dialer.Subprotocols = []string{"remotecontrol"}

	var err error
	c.conn, _, err = dialer.Dial(u.String(), nil)
	if err != nil {
		c.mu.Unlock()
		return fmt.Errorf("websocket dial failed: %v", err)
	}
	c.isConnected = true
	c.mu.Unlock()

	// Start reading messages
	go c.readLoop()
	// Start keepalive loop
	go c.keepAliveLoop()

	// 3. Authenticate
	if err := c.Authenticate(); err != nil {
		return fmt.Errorf("authentication failed: %v", err)
	}

	return nil
}

func (c *Client) checkReachability(scheme, host string, port int) error {
	// Determine HTTP scheme based on WSS scheme
	httpScheme := "http"
	if scheme == "wss" {
		httpScheme = "https"
	}

	url := fmt.Sprintf("%s://%s:%d/jdev/cfg/apiKey", httpScheme, host, port)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("reachability check returned status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("failed to parse apiKey response: %v", err)
	}
	return nil
}

// readLoop handles incoming messages
func (c *Client) readLoop() {
	defer func() {
		if c.conn != nil {
			c.conn.Close()
		}
		c.mu.Lock()
		c.isConnected = false
		c.mu.Unlock()
	}()

	for {
		select {
		case <-c.done:
			return
		default:
			msgType, message, err := c.conn.ReadMessage()
			if err != nil {
				return
			}

			// If we get a Header (8 bytes starting with 0x03)
			if msgType == websocket.BinaryMessage && len(message) == 8 && message[0] == 0x03 {
				newHeader := &Header{
					Type:   message[1],
					Info:   message[2],
					Length: binary.LittleEndian.Uint32(message[4:]),
				}

				c.header = newHeader

				// Type 6 is Keepalive response (no payload follows)
				if c.header.Type == 6 {
					c.header = nil
				}
				continue
			}

			// If we have a header, this message is the payload
			if c.header != nil {
				h := c.header
				c.header = nil

				switch h.Type {
				case 0: // Text-Message
					if msgType == websocket.TextMessage {
						slog.Debug("Loxone Text Message", "message", string(message))
						select {
						case c.msgChan <- message:
						default:
						}
					}
				case 1: // Binary File
					select {
					case c.msgChan <- message:
					default:
					}
				case 2: // Value-States
					c.HandleBinaryMessage(message)
				case 3: // Text-States
					c.HandleTextMessage(message)
				case 5: // Out-Of-Service
					slog.Warn("Miniserver indicates Out-Of-Service")
				}
			} else {
				// No header - usually unsolicited text messages
				if msgType == websocket.TextMessage {
					select {
					case c.msgChan <- message:
					default:
					}
				}
			}
		}
	}
}

// keepAliveLoop sends a keepalive message every 4 minutes
func (c *Client) keepAliveLoop() {
	ticker := time.NewTicker(4 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-c.done:
			return
		case <-ticker.C:
			if err := c.SendCommand("keepalive"); err != nil {
				slog.Error("Error sending keepalive", "error", err)
				return
			}
		}
	}
}

// WaitForResponse waits for a specific control response
func (c *Client) WaitForResponse(controlContains string) (*Response, error) {
	timeout := time.After(5 * time.Second)
	for {
		select {
		case msg := <-c.msgChan:
			var resp Response
			if err := json.Unmarshal(msg, &resp); err != nil {
				continue
			}
			if strings.Contains(resp.LL.Control, controlContains) {
				return &resp, nil
			}
		case <-timeout:
			return nil, fmt.Errorf("timeout waiting for response: %s", controlContains)
		}
	}
}

// SendCommand sends a raw command string
func (c *Client) SendCommand(cmd string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.isConnected {
		return fmt.Errorf("not connected")
	}
	return c.conn.WriteMessage(websocket.TextMessage, []byte(cmd))
}

// Done returns the done channel
func (c *Client) Done() <-chan struct{} {
	return c.done
}

// Close closes the connection
func (c *Client) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.isClosed {
		return
	}
	c.isClosed = true
	close(c.done)
	if c.conn != nil {
		c.conn.Close()
	}
}
