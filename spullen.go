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

type VersionRepository interface {
	GetAppVersion() int
	GetStoredVersion() int
}

type ObjectRepository interface {
	GetAll() []*Object
	Count() int
	Get(id string) *Object
	GetDistinctCategories(includeHidden bool) []string
	GetDistinctTags(includeHidden bool) []string
	GetDistinctPropertyKeys(includeHidden bool) []string
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
