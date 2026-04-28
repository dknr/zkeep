package zfs

import (
	"fmt"
	"os/exec"
	"strings"
)

// ZFS represents a ZFS interface.
type ZFS struct{}

// New creates a new ZFS instance.
func New() *ZFS {
	return &ZFS{}
}

// Snapshot creates a ZFS snapshot.
func (z *ZFS) Snapshot(snapshot string) error {
	cmd := exec.Command("zfs", "snapshot", snapshot)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("zfs snapshot failed: %w, output: %s", err, string(output))
	}
	return nil
}

// ListSnapshots lists ZFS snapshots for a dataset.
func (z *ZFS) ListSnapshots(dataset string) (string, error) {
	cmd := exec.Command("zfs", "list", "-t", "snapshot", dataset)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("zfs list failed: %w, output: %s", err, string(output))
	}
	return strings.TrimSpace(string(output)), nil
}

// DestroySnapshot destroys a ZFS snapshot.
func (z *ZFS) DestroySnapshot(snapshot string) error {
	cmd := exec.Command("zfs", "destroy", snapshot)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("zfs destroy failed: %w, output: %s", err, string(output))
	}
	return nil
}

// Rollback rolls back to a ZFS snapshot.
func (z *ZFS) Rollback(snapshot string) error {
	cmd := exec.Command("zfs", "rollback", snapshot)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("zfs rollback failed: %w, output: %s", err, string(output))
	}
	return nil
}
