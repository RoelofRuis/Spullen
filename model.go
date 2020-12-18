package main

import "time"

type App struct {
	Storage *Storage
	PrivateMode bool
}

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
