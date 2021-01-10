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

type ObjectId string

type ObjectRepository interface {
	GetAll() []*Object
	Count() int
	Get(id ObjectId) *Object
	GetDistinctCategories(includeHidden bool) []string
	GetDistinctTags(includeHidden bool) []string
	GetDistinctPropertyKeys(includeHidden bool) []string
	Put(*Object)
	Remove(id ObjectId)
	Has(id ObjectId) bool
}

type ObjectDeletionRepository interface {
	Get(id ObjectId) *ObjectDeletion
	Put(*ObjectDeletion)
	Has(id ObjectId) bool
}

type Object struct {
	Id         ObjectId
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

type ObjectDeletion struct {
	Id        ObjectId
	Reason    string
	DeletedAt time.Time
}
