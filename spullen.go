package spullen

import "time"

type Database interface {
	IsOpened() bool
	Name() string
	Open(name string, pass []byte, openExisting bool) error
	IsDirty() bool
	Persist() error
	Close() error
}

type ObjectRepository interface {
	GetAll() []*Object
	Get(id string) *Object
	GetDistinctCategories() []string
	GetDistinctTags() []string
	GetDistinctPropertyKeys() []string
	Put(*Object)
	Remove(id string)
	Has(id string) bool
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
