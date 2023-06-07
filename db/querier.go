// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0

package db

import (
	"context"
)

type Querier interface {
	BuildCreate(ctx context.Context, arg BuildCreateParams) (string, error)
	BuildGet(ctx context.Context, id string) (UnweaveBuild, error)
	BuildGetUsedBy(ctx context.Context, id string) ([]BuildGetUsedByRow, error)
	BuildUpdate(ctx context.Context, arg BuildUpdateParams) error
	ExecCreate(ctx context.Context, arg ExecCreateParams) error
	ExecGet(ctx context.Context, idOrName string) (UnweaveExec, error)
	ExecGetAllActive(ctx context.Context) ([]UnweaveExec, error)
	ExecListActiveByProvider(ctx context.Context, provider string) ([]UnweaveExec, error)
	ExecListByProvider(ctx context.Context, provider string) ([]UnweaveExec, error)
	ExecSetError(ctx context.Context, arg ExecSetErrorParams) error
	ExecStatusUpdate(ctx context.Context, arg ExecStatusUpdateParams) error
	ExecUpdateConnectionInfo(ctx context.Context, arg ExecUpdateConnectionInfoParams) error
	ExecsGet(ctx context.Context, arg ExecsGetParams) ([]UnweaveExec, error)
	//-----------------------------------------------------------------
	// The queries below return data in the format expected by the API.
	//-----------------------------------------------------------------
	MxExecGet(ctx context.Context, id string) (MxExecGetRow, error)
	MxExecsGet(ctx context.Context, projectID string) ([]MxExecsGetRow, error)
	NodeCreate(ctx context.Context, arg NodeCreateParams) error
	NodeStatusUpdate(ctx context.Context, arg NodeStatusUpdateParams) error
	ProjectGet(ctx context.Context, id string) (UnweaveProject, error)
	SSHKeyAdd(ctx context.Context, arg SSHKeyAddParams) error
	SSHKeyGetByName(ctx context.Context, arg SSHKeyGetByNameParams) (UnweaveSshKey, error)
	SSHKeyGetByPublicKey(ctx context.Context, arg SSHKeyGetByPublicKeyParams) (UnweaveSshKey, error)
	SSHKeysGet(ctx context.Context, ownerID string) ([]UnweaveSshKey, error)
	VolumeCreate(ctx context.Context, arg VolumeCreateParams) (UnweaveVolume, error)
	VolumeDelete(ctx context.Context, id string) error
	VolumeGet(ctx context.Context, id string) (UnweaveVolume, error)
	VolumeList(ctx context.Context, projectID string) ([]UnweaveVolume, error)
}

var _ Querier = (*Queries)(nil)
