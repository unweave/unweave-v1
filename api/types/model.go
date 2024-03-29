package types

import (
	"fmt"
	"io"
	"time"

	"github.com/rs/zerolog"
)

type Status string

const (
	StatusPending      Status = "pending"
	StatusInitializing Status = "initializing"
	StatusRunning      Status = "running"
	StatusTerminated   Status = "terminated"
	StatusError        Status = "error"
	StatusFailed       Status = "failed"
	StatusSuccess      Status = "success"
	StatusUnknown      Status = "unknown"
)

func (s Status) IsTerminal() bool {
	switch s {
	case StatusTerminated, StatusError, StatusFailed, StatusSuccess:
		return true
	default:
		return false
	}
}

type NoOpLogHook struct{}

func (d NoOpLogHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {}

var NewErrLogHook = func() zerolog.Hook { return NoOpLogHook{} }

// Provider is the platform that the node is spawned on. This is where the user
// code runs
type Provider string

func (r Provider) String() string {
	return string(r)
}

const (
	LambdaLabsProvider Provider = "lambdalabs"
	UnweaveProvider    Provider = "unweave"
	AWSProvider        Provider = "aws"
)

func (r Provider) DisplayName() string {
	switch r {
	case LambdaLabsProvider:
		return "LambdaLabs"
	case UnweaveProvider:
		return "Unweave"
	case AWSProvider:
		return "AWS"
	default:
		return "Unknown"
	}
}

type Build struct {
	BuildID     string     `json:"buildID"`
	Name        string     `json:"name"`
	ProjectID   string     `json:"projectID"`
	Status      string     `json:"status"`
	BuilderType string     `json:"builderType"`
	CreatedAt   time.Time  `json:"createdAt"`
	StartedAt   *time.Time `json:"startedAt,omitempty"`
	FinishedAt  *time.Time `json:"finishedAt,omitempty"`
}

type LogEntry struct {
	TimeStamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
	Level     string    `json:"level"`
}

func (l LogEntry) String() string {
	return fmt.Sprintf("%s %s %s", l.TimeStamp.Format(time.RFC3339), l.Level, l.Message)
}

type NodeType struct {
	Type     string       `json:"type"`
	ID       string       `json:"id"`
	Name     *string      `json:"name"`
	Price    *int         `json:"price"`
	Regions  []string     `json:"regions"`
	Provider Provider     `json:"provider"`
	Specs    HardwareSpec `json:"specs"`
}

type Node struct {
	ID       string       `json:"id"`
	TypeID   string       `json:"typeID"`
	OwnerID  string       `json:"ownerID"`
	Price    int          `json:"price"`
	Region   string       `json:"region"`
	KeyPair  SSHKey       `json:"sshKeyPair"`
	Status   Status       `json:"status"`
	Provider Provider     `json:"provider"`
	Specs    HardwareSpec `json:"specs"`
	Host     string       `json:"host"`
	User     string       `json:"user"`
	Port     int          `json:"port"`
}

type Project struct {
	ID         string    `json:"id"`
	Owner      string    `json:"owner"`
	Name       string    `json:"name"`
	Icon       string    `json:"icon"`
	CreatedAt  time.Time `json:"createdAt"`
	Visibility string    `json:"visibility"`
	GithubURL  *string   `json:"githubURL"`
}

type SSHKey struct {
	Name       string     `json:"name"`
	PublicKey  *string    `json:"publicKey,omitempty"`
	PrivateKey *string    `json:"privateKey,omitempty"`
	CreatedAt  *time.Time `json:"createdAt,omitempty"`
}

type ConnectionInfo struct {
	Host string `json:"host"`
	Port int    `json:"port"`
	User string `json:"user"`
}

type HTTPService struct {
	InternalPort int32 `json:"internalPort"`
}

type ExecNetwork struct {
	Host        string       `json:"host"`
	Port        int          `json:"port"`
	User        string       `json:"user"`
	HTTPService *HTTPService `json:"httpService,omitempty"`
}

type ExecVolume struct {
	VolumeID  string `json:"volumeID"`
	MountPath string `json:"mountPath"`
}

type Exec struct {
	ID        string       `json:"id"`
	Name      string       `json:"name"`
	CreatedAt time.Time    `json:"createdAt,omitempty"`
	ExitedAt  *time.Time   `json:"exitedAt,omitempty"`
	CreatedBy string       `json:"createdBy,omitempty"`
	Image     string       `json:"image,omitempty"`
	BuildID   *string      `json:"buildID,omitempty"`
	Status    Status       `json:"status"`
	Command   []string     `json:"command"`
	Keys      []SSHKey     `json:"keys"`
	Volumes   []ExecVolume `json:"volumes"`
	Network   ExecNetwork  `json:"network"`
	Spec      HardwareSpec `json:"spec"`
	CommitID  *string      `json:"commitID,omitempty"`
	GitURL    *string      `json:"gitURL,omitempty"`
	Region    string       `json:"region"`
	Provider  Provider     `json:"provider"`
}

type ExecConfig struct {
	Image   string         `json:"image"`
	Command []string       `json:"command"`
	Keys    []SSHKey       `json:"keys"`
	Volumes []Volume       `json:"volumes"`
	Src     *SourceContext `json:"src,omitempty"`
}

type GitConfig struct {
	CommitID *string `json:"commitID"`
	GitURL   *string `json:"gitURL"`
}

type SourceContext struct {
	MountPath string        `json:"mountPath"`
	Context   io.ReadCloser `json:"-"`
}

// UserAccessToken is an internal type to safely pass user tokens, never exposes token or hash
type UserAccessToken struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	DisplayText string    `json:"displayText"`
	ExpiresAt   time.Time `json:"expiresAt"`
	UserID      string    `json:"userId"`
}

func NewUserAccessToken(userID, tokenID, tokenName, displayText string, expiresAt time.Time) UserAccessToken {
	return UserAccessToken{
		UserID:      userID,
		ID:          tokenID,
		Name:        tokenName,
		DisplayText: displayText,
		ExpiresAt:   expiresAt,
	}
}

type Volume struct {
	ID        string      `json:"id"`
	Name      string      `json:"name"`
	Size      int         `json:"size"`
	State     VolumeState `json:"state"`
	Provider  Provider    `json:"provider"`
	ProjectID string      `json:"projectID"`
}

type VolumeState struct {
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type CheckStatus string

var (
	CheckPending    CheckStatus = "pending"
	CheckInProgress CheckStatus = "in_progress"
	CheckCompleted  CheckStatus = "completed"
)

type CheckConclusion string

var (
	CheckSuccess CheckConclusion = "success"
	CheckFailure CheckConclusion = "failure"
	CheckError   CheckConclusion = "error"
)

func (c CheckConclusion) String() string {
	return string(c)
}

type EvalManifest struct {
	DatasetURL string `json:"datasetURL,omitempty"`
	AssertURL  string `json:"assertURL,omitempty"`
	RunURL     string `json:"runURL,omitempty"`
}

func DefaultEvalManifest() EvalManifest {
	return EvalManifest{
		DatasetURL: "/dataset",
		AssertURL:  "/assert",
		RunURL:     "",
	}
}
