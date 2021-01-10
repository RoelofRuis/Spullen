package spullen

import (
	"fmt"
	"time"
)

type DataFlags struct {
	ShowHiddenItems  bool
	ShowDeletedItems bool
}

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

func (p *Property) String() string {
	return fmt.Sprintf("%s=%s", p.Key, p.Value)
}

type DeletionRepository interface {
	Get(id ObjectId) *Deletion
	Put(*Deletion)
	Has(id ObjectId) bool
}

type Deletion struct {
	Id        ObjectId
	Reason    string
	DeletedAt time.Time
}
