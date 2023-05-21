package types

import "time"

type Volume struct {
	ID       string       `json:"id"`
	Name     string       `json:"name"`
	Provider Provider     `json:"provider"`
	State    *VolumeState `json:"state"`
}

type VolumeState struct {
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
