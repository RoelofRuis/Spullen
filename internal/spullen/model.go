package spullen

import "time"

type Database interface {
	IsOpened() bool
	Name() string
	Open(name string, pass []byte, isExisting bool) (ObjectRepository, error)
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
	PutObject(*Object) error
	RemoveObject(id string) error
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
