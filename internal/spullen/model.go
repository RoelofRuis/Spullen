package spullen

import "time"

type DatabaseMode int
const ModeOpenExisting DatabaseMode = 0x1
const ModeUseGzip DatabaseMode = 0x2
const ModeUseEncryption DatabaseMode = 0x4

type Database interface {
	IsOpened() bool
	Name() string
	Open(name string, pass []byte, mode DatabaseMode) (ObjectRepository, error)
	Persist() error
	Close()
}

type ObjectRepositoryFactory interface {
	CreateFromData(data []byte) (ObjectRepository, error)
	CreateNew() ObjectRepository
}

type ObjectRepository interface {
	GetAll() []*Object
	Get(id string) *Object
	Put(*Object) error
	Remove(id string) error
	ToRawData() ([]byte, error)
}

type Object struct {
	Id         string
	Added      time.Time
	Name       string
	Quantity   int
	Categories []string
	Tags       []string
	Properties []*Property
	Hidden     bool
	Notes      string
}

type Property struct {
	Key   string
	Value string
}
