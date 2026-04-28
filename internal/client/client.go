package client

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

// SocketPath is the path to the zkeep daemon's Unix socket.
const SocketPath = "/var/run/zkeep.sock"

// Client represents the zkeep client.
type Client struct {
	socketPath string
}

// New creates a new Client instance.
func New(socketPath string) *Client {
	return &Client{
		socketPath: socketPath,
	}
}

// SendCommand sends a command to the daemon and returns the response.
func (c *Client) SendCommand(cmd string) (string, error) {
	// Connect to the Unix socket
	conn, err := net.Dial("unix", c.socketPath)
	if err != nil {
		return "", fmt.Errorf("failed to connect to daemon: %w", err)
	}
	defer conn.Close()

	// Send command
	_, err = fmt.Fprintf(conn, "%s\n", cmd)
	if err != nil {
		return "", fmt.Errorf("failed to send command: %w", err)
	}

	// Close write end to signal we're done sending
	// This allows the daemon's scanner to detect EOF and process the command
	if uc, ok := conn.(*net.UnixConn); ok {
		uc.CloseWrite()
	}

	// Read all response lines until daemon closes connection
	var response strings.Builder
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		response.WriteString(scanner.Text())
		response.WriteString("\n")
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading response: %w", err)
	}

	return strings.TrimSpace(response.String()), nil
}
