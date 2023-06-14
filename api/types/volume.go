package types

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/unweave/unweave/tools/random"
)

type Volume struct {
	ID       string       `json:"id"`
	Name     string       `json:"name"`
	Size     int          `json:"size"`
	State    *VolumeState `json:"state"`
	Provider Provider     `json:"provider"`
}

func NewVolume(name string, size int, provider Provider) (Volume, error) {
	rand, err := random.GenerateRandomString(11)
	if err != nil {
		return Volume{}, fmt.Errorf("failed to generate random string, %w", err)
	}
	// not a cryptographically safe identifier, but consistent with execs
	volID := "vol_" + strings.ToLower(rand)

	return Volume{
		ID:   volID,
		Name: name,
		Size: size,
		State: &VolumeState{
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		},
		Provider: provider,
	}, nil
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

func (p *VolumeCreateRequest) Bind(r *http.Request) error {
	if p.Name == "" {
		return &Error{
			Code:    http.StatusBadRequest,
			Message: "Name is required",
		}
	}

	if p.Size <= 0 {
		return &Error{
			Code:    http.StatusBadRequest,
			Message: "Size is required",
		}
	}

	if p.Provider != UnweaveProvider { // Lambda Labs not implemented
		return &Error{
			Code:    http.StatusBadRequest,
			Message: "Provider is required",
		}
	}
	return nil
}

type VolumeDeleteRequest struct {
	IDOrName string `json:"idOrName"`
}

func (p *VolumeDeleteRequest) Bind(r *http.Request) error {
	if p.IDOrName == "" {
		return &Error{
			Code:    http.StatusBadRequest,
			Message: "Name is required",
		}
	}

	return nil
}

type VolumeResizeRequest struct {
	IDOrName string `json:"idOrName"`
	Size     int    `json:"size"`
}

func (p *VolumeResizeRequest) Bind(r *http.Request) error {
	if p.IDOrName == "" {
		return &Error{
			Code:    http.StatusBadRequest,
			Message: "Name is required",
		}
	}

	if p.Size <= 0 {
		return &Error{
			Code:    http.StatusBadRequest,
			Message: "A new size is required",
		}
	}

	return nil
}
