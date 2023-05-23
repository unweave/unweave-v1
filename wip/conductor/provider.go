package conductor

type Node interface {
	NodeCreate()
	NodeInit()
	NodeStop()
	NodeTerminate()
	NodeList()
}

type Volume interface {
	VolumeCreate()
	VolumeDelete()
	VolumeList()
	VolumeGet()
	VolumeResize()
}

type Keys interface {
	KeysRegister()
	KeysUnregister()
	KeysList()
}

type Network interface {
	NetworkAddPort()
	NetworkRemovePort()
}

type Provider interface {
	ID() string
	Node
	Volume
}
