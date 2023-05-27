package types

import "time"

type Volume struct {
	ID       string       `json:"id"`
	Name     string       `json:"name"`
	Size     int          `json:"size"`
	State    *VolumeState `json:"state"`
	Provider Provider     `json:"provider"`
}

type VolumeState struct {
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type VolumeCreateRequest struct {
	Size     int      `json:"size"`
	Name     string   `json:"name"`
	Provider Provider `json:"provider"`
}
