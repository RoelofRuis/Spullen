package main

import "time"

type Storage struct {
	name string
	pass []byte
}

type ObjectRepository interface {
	GetAll() []*Object
	Get(id string) *Object
	PutObject(*Object) error
	RemoveObject(id string) error
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
