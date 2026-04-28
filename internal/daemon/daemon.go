package daemon

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"zkeep/internal/zfs"
)

const groupName = "zkeep"

// Daemon represents the zkeep daemon.
type Daemon struct {
	zfs        *zfs.ZFS
	socketPath string
	version    string
	buildTime  string
}

// New creates a new Daemon instance.
func New(zfsInst *zfs.ZFS, socketPath, version, buildTime string) *Daemon {
	return &Daemon{
		zfs:        zfsInst,
		socketPath: socketPath,
		version:    version,
		buildTime:  buildTime,
	}
}

// Run starts the daemon and listens for connections on the Unix socket.
func (d *Daemon) Run() error {
	// Remove existing socket file if it exists
	if _, err := net.DialTimeout("unix", d.socketPath, 0); err == nil {
		if err := os.Remove(d.socketPath); err != nil {
			return fmt.Errorf("failed to remove existing socket: %w", err)
		}
	}

	// Listen on Unix socket
	listener, err := net.Listen("unix", d.socketPath)
	if err != nil {
		return fmt.Errorf("failed to listen on socket: %w", err)
	}

	// Set socket permissions
	if err := d.setSocketPermissions(); err != nil {
		listener.Close()
		return err
	}

	fmt.Printf("Daemon listening on %s\n", d.socketPath)

	for {
		conn, err := listener.Accept()
		if err != nil {
			return fmt.Errorf("failed to accept connection: %w", err)
		}

		go d.handleConnection(conn)
	}
}

func (d *Daemon) setSocketPermissions() error {
	gid := getGroupID(groupName)
	if gid == -1 {
		fmt.Fprintf(os.Stderr, "Warning: group '%s' not found, skipping socket ownership\n", groupName)
		return nil
	}

	if err := os.Chown(d.socketPath, 0, gid); err != nil {
		return fmt.Errorf("failed to set socket group ownership: %w", err)
	}

	if err := os.Chmod(d.socketPath, 0660); err != nil {
		return fmt.Errorf("failed to set socket permissions: %w", err)
	}

	return nil
}

func (d *Daemon) handleConnection(conn net.Conn) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		resp, err := d.handleCommand(line)
		if err != nil {
			resp = fmt.Sprintf("Error: %v", err)
		}

		_, err = conn.Write([]byte(resp + "\n"))
		if err != nil {
			return
		}

		// Close the connection after sending the response
		conn.Close()
		return
	}
}

func (d *Daemon) handleCommand(cmd string) (string, error) {
	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return "", fmt.Errorf("empty command")
	}

	switch parts[0] {
	case "snapshot":
		if len(parts) < 2 {
			return "", fmt.Errorf("usage: snapshot <dataset@tag>")
		}
		err := d.zfs.Snapshot(parts[1])
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Snapshot %s created", parts[1]), nil
	case "list":
		if len(parts) < 2 {
			return "", fmt.Errorf("usage: list <dataset>")
		}
		snapshots, err := d.zfs.ListSnapshots(parts[1])
		if err != nil {
			return "", err
		}
		return snapshots, nil
	case "destroy":
		if len(parts) < 2 {
			return "", fmt.Errorf("usage: destroy <snapshot>")
		}
		err := d.zfs.DestroySnapshot(parts[1])
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Snapshot %s destroyed", parts[1]), nil
	case "rollback":
		if len(parts) < 2 {
			return "", fmt.Errorf("usage: rollback <snapshot>")
		}
		err := d.zfs.Rollback(parts[1])
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Rolled back to %s", parts[1]), nil
	case "status":
		return "Status: OK", nil
	case "version":
		return fmt.Sprintf("%s (built %s)", d.version, d.buildTime), nil
	default:
		return "", fmt.Errorf("unknown command: %s", parts[0])
	}
}

func getGroupID(name string) int {
	content, err := os.ReadFile("/etc/group")
	if err != nil {
		return -1
	}

	for _, line := range splitLines(string(content)) {
		if len(line) > 0 && line[0] != '#' {
			fields := splitFields(line)
			if len(fields) >= 3 && fields[0] == name {
				var id int
				fmt.Sscanf(fields[2], "%d", &id)
				return id
			}
		}
	}
	return -1
}

func splitLines(s string) []string {
	var lines []string
	var current string
	for _, c := range s {
		if c == '\n' {
			lines = append(lines, current)
			current = ""
		} else {
			current += string(c)
		}
	}
	if len(current) > 0 {
		lines = append(lines, current)
	}
	return lines
}

func splitFields(s string) []string {
	var fields []string
	var current string
	for _, c := range s {
		if c == ':' {
			fields = append(fields, current)
			current = ""
		} else {
			current += string(c)
		}
	}
	if len(current) > 0 {
		fields = append(fields, current)
	}
	return fields
}
