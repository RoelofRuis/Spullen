package spullen

import "time"

type ObjectRepository interface {
	GetAll() []*Object
	Get(id string) *Object
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
