package local

import (
	"context"
	"io"
)

// NodeDriver is a local driver that runs commands on the local machine. This should only be
// used for testing.
type NodeDriver struct{}

func (d *NodeDriver) NodeRunCommand(ctx context.Context, id string, command []string, env []string) error {
	return nil
}

func (d *NodeDriver) NodeRunScript(ctx context.Context, id string, script, workingDir string, env []string) error {
	return nil
}

func (d *NodeDriver) NodeTransferFile(ctx context.Context, id string, src io.ReadCloser, dst string) error {
	return nil
}

// VolumeDriver is a local driver that mounts volumes on the local machine.
// This should only be used for testing.
type VolumeDriver struct{}

func (v *VolumeDriver) Mount(ctx context.Context, volumeID, path string) error {
	return nil
}

func (v *VolumeDriver) Unmount(ctx context.Context, volumeID, path string) error {
	return nil
}
