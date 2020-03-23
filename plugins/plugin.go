package plugins

type Plugin interface {
	NewInstance(instanceFilePath string) (Instance, error)
	DestroyInstance(instance Instance) error
	GetInstancesDir() string
}

type Instance interface {
	Start() error
	Stop() error
}
