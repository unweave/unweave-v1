// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0

package db

import (
	"context"
)

type Querier interface {
	BuildCreate(ctx context.Context, arg BuildCreateParams) (string, error)
	BuildGet(ctx context.Context, id string) (UnweaveBuild, error)
	BuildGetUsedBy(ctx context.Context, id string) ([]BuildGetUsedByRow, error)
	BuildUpdate(ctx context.Context, arg BuildUpdateParams) error
	EndpointCreate(ctx context.Context, arg EndpointCreateParams) error
	EndpointDelete(ctx context.Context, id string) error
	EndpointEval(ctx context.Context, endpointID string) ([]UnweaveEndpointEval, error)
	EndpointEvalAttach(ctx context.Context, arg EndpointEvalAttachParams) error
	EndpointGet(ctx context.Context, id string) (UnweaveEndpoint, error)
	EndpointsForProject(ctx context.Context, projectID string) ([]UnweaveEndpoint, error)
	EvalCreate(ctx context.Context, arg EvalCreateParams) error
	EvalDelete(ctx context.Context, id string) error
	EvalGet(ctx context.Context, id string) (EvalGetRow, error)
	EvalList(ctx context.Context, dollar_1 []string) ([]EvalListRow, error)
	ExecCreate(ctx context.Context, arg ExecCreateParams) error
	ExecGet(ctx context.Context, idOrName string) (UnweaveExec, error)
	ExecGetAllActive(ctx context.Context) ([]UnweaveExec, error)
	ExecList(ctx context.Context, arg ExecListParams) ([]UnweaveExec, error)
	ExecListActiveByProvider(ctx context.Context, provider string) ([]UnweaveExec, error)
	ExecListByProvider(ctx context.Context, provider string) ([]UnweaveExec, error)
	ExecSSHKeyDelete(ctx context.Context, arg ExecSSHKeyDeleteParams) error
	ExecSSHKeyGet(ctx context.Context, arg ExecSSHKeyGetParams) (UnweaveExecSshKey, error)
	ExecSSHKeyInsert(ctx context.Context, arg ExecSSHKeyInsertParams) error
	ExecSSHKeysGetByExecID(ctx context.Context, execID string) ([]UnweaveExecSshKey, error)
	ExecSetError(ctx context.Context, arg ExecSetErrorParams) error
	ExecSetFailed(ctx context.Context, arg ExecSetFailedParams) error
	ExecStatusUpdate(ctx context.Context, arg ExecStatusUpdateParams) error
	ExecUpdateConnectionInfo(ctx context.Context, arg ExecUpdateConnectionInfoParams) error
	ExecUpdateNetwork(ctx context.Context, arg ExecUpdateNetworkParams) error
	ExecVolumeCreate(ctx context.Context, arg ExecVolumeCreateParams) error
	ExecVolumeDelete(ctx context.Context, execID string) error
	ExecVolumeGet(ctx context.Context, execID string) ([]UnweaveExecVolume, error)
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
	SSHKeysGetByIDs(ctx context.Context, ids []string) ([]UnweaveSshKey, error)
	VolumeCreate(ctx context.Context, arg VolumeCreateParams) (UnweaveVolume, error)
	VolumeDelete(ctx context.Context, id string) error
	VolumeGet(ctx context.Context, arg VolumeGetParams) (UnweaveVolume, error)
	VolumeList(ctx context.Context, projectID string) ([]UnweaveVolume, error)
	VolumeUpdate(ctx context.Context, arg VolumeUpdateParams) error
}

var _ Querier = (*Queries)(nil)
