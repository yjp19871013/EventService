package plugins

type Plugin interface {
	NewInstance(instanceFilePath string) (Instance, error)
	DestroyInstance(instance Instance) error
}

type Instance interface {
	Start() error
	Stop() error
}
