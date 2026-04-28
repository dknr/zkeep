package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"zkeep/internal/daemon"
	"zkeep/internal/zfs"
)

const (
	socketPath = "/var/run/zkeep.sock"
	groupName  = "zkeep"
)

var (
	Version   = "dev"
	BuildTime = "unknown"
)

func main() {
	// Ensure run directory exists
	if err := os.MkdirAll("/var/run", 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create run directory: %v\n", err)
		os.Exit(1)
	}

	// Remove existing socket file if it exists
	if _, err := os.Stat(socketPath); err == nil {
		if err := os.Remove(socketPath); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to remove existing socket: %v\n", err)
			os.Exit(1)
		}
	}

	zfsInst := zfs.New()
	d := daemon.New(zfsInst, socketPath, Version, BuildTime)

	// Handle signals for clean shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("Shutting down...")
		os.Remove(socketPath)
		os.Exit(0)
	}()

	fmt.Println("zkeepd started. Listening on", socketPath)
	if err := d.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Daemon error: %v\n", err)
		os.Remove(socketPath)
		os.Exit(1)
	}
}
