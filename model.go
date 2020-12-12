package main

import "time"

type Storage interface {
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
	Description string
	Categories []string
	Tags       []string
	Properties []*Property
	Hidden     bool
}

type Property struct {
	Key   string
	Value string
}
