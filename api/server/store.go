package server

import "github.com/unweave/unweave/service/volume"

var store struct {
	VolumeStore volume.Store
}

func GetVolumeStore() volume.Store {
	return store.VolumeStore
}

func InitStore(volumeStore volume.Store) {
	store.VolumeStore = volumeStore
}
