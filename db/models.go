// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0

package db

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type UnweaveBuildStatus string

const (
	UnweaveBuildStatusInitializing    UnweaveBuildStatus = "initializing"
	UnweaveBuildStatusBuilding        UnweaveBuildStatus = "building"
	UnweaveBuildStatusSuccess         UnweaveBuildStatus = "success"
	UnweaveBuildStatusFailed          UnweaveBuildStatus = "failed"
	UnweaveBuildStatusError           UnweaveBuildStatus = "error"
	UnweaveBuildStatusCanceled        UnweaveBuildStatus = "canceled"
	UnweaveBuildStatusSyncingSnapshot UnweaveBuildStatus = "syncing_snapshot"
)

func (e *UnweaveBuildStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = UnweaveBuildStatus(s)
	case string:
		*e = UnweaveBuildStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for UnweaveBuildStatus: %T", src)
	}
	return nil
}

type NullUnweaveBuildStatus struct {
	UnweaveBuildStatus UnweaveBuildStatus
	Valid              bool // Valid is true if UnweaveBuildStatus is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullUnweaveBuildStatus) Scan(value interface{}) error {
	if value == nil {
		ns.UnweaveBuildStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.UnweaveBuildStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullUnweaveBuildStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.UnweaveBuildStatus), nil
}

type UnweaveExecStatus string

const (
	UnweaveExecStatusInitializing UnweaveExecStatus = "initializing"
	UnweaveExecStatusRunning      UnweaveExecStatus = "running"
	UnweaveExecStatusTerminated   UnweaveExecStatus = "terminated"
	UnweaveExecStatusError        UnweaveExecStatus = "error"
	UnweaveExecStatusSnapshotting UnweaveExecStatus = "snapshotting"
	UnweaveExecStatusPending      UnweaveExecStatus = "pending"
)

func (e *UnweaveExecStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = UnweaveExecStatus(s)
	case string:
		*e = UnweaveExecStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for UnweaveExecStatus: %T", src)
	}
	return nil
}

type NullUnweaveExecStatus struct {
	UnweaveExecStatus UnweaveExecStatus
	Valid             bool // Valid is true if UnweaveExecStatus is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullUnweaveExecStatus) Scan(value interface{}) error {
	if value == nil {
		ns.UnweaveExecStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.UnweaveExecStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullUnweaveExecStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.UnweaveExecStatus), nil
}

type UnweaveAccount struct {
	ID string `json:"id"`
}

type UnweaveBuild struct {
	ID          string             `json:"id"`
	Name        string             `json:"name"`
	ProjectID   string             `json:"projectID"`
	BuilderType string             `json:"builderType"`
	Status      UnweaveBuildStatus `json:"status"`
	CreatedBy   string             `json:"createdBy"`
	CreatedAt   time.Time          `json:"createdAt"`
	StartedAt   sql.NullTime       `json:"startedAt"`
	FinishedAt  sql.NullTime       `json:"finishedAt"`
	UpdatedAt   time.Time          `json:"updatedAt"`
	MetaData    json.RawMessage    `json:"metaData"`
}

type UnweaveExec struct {
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	Region       string            `json:"region"`
	CreatedBy    string            `json:"createdBy"`
	CreatedAt    time.Time         `json:"createdAt"`
	ReadyAt      sql.NullTime      `json:"readyAt"`
	ExitedAt     sql.NullTime      `json:"exitedAt"`
	Status       UnweaveExecStatus `json:"status"`
	ProjectID    string            `json:"projectID"`
	Error        sql.NullString    `json:"error"`
	BuildID      sql.NullString    `json:"buildID"`
	Spec         json.RawMessage   `json:"spec"`
	CommitID     sql.NullString    `json:"commitID"`
	GitRemoteUrl sql.NullString    `json:"gitRemoteUrl"`
	Command      []string          `json:"command"`
	Metadata     json.RawMessage   `json:"metadata"`
	Image        string            `json:"image"`
	Provider     string            `json:"provider"`
}

type UnweaveExecSshKey struct {
	ExecID   string `json:"execID"`
	SshKeyID string `json:"sshKeyID"`
}

type UnweaveExecVolume struct {
	ExecID    string `json:"execID"`
	VolumeID  string `json:"volumeID"`
	MountPath string `json:"mountPath"`
}

type UnweaveNode struct {
	ID           string          `json:"id"`
	Provider     string          `json:"provider"`
	Region       string          `json:"region"`
	Metadata     json.RawMessage `json:"metadata"`
	Status       string          `json:"status"`
	CreatedAt    time.Time       `json:"createdAt"`
	ReadyAt      sql.NullTime    `json:"readyAt"`
	OwnerID      string          `json:"ownerID"`
	TerminatedAt sql.NullTime    `json:"terminatedAt"`
}

type UnweaveNodeSshKey struct {
	NodeID    string    `json:"nodeID"`
	SshKeyID  string    `json:"sshKeyID"`
	CreatedAt time.Time `json:"createdAt"`
}

type UnweaveProject struct {
	ID             string         `json:"id"`
	DefaultBuildID sql.NullString `json:"defaultBuildID"`
}

type UnweaveSshKey struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	OwnerID   string    `json:"ownerID"`
	CreatedAt time.Time `json:"createdAt"`
	PublicKey string    `json:"publicKey"`
	IsActive  bool      `json:"isActive"`
}

type UnweaveVolume struct {
	ID        string       `json:"id"`
	Name      string       `json:"name"`
	ProjectID string       `json:"projectID"`
	Provider  string       `json:"provider"`
	CreatedAt time.Time    `json:"createdAt"`
	UpdatedAt time.Time    `json:"updatedAt"`
	Size      int32        `json:"size"`
	DeletedAt sql.NullTime `json:"deletedAt"`
}
