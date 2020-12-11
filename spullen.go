package main

import "time"

type Object struct {
	Id         string
	Name       string
	Quantity   int
	Added      time.Time
	Categories []string
	Tags       []string
	Properties []*Property
	Hidden     bool
}

type Property struct {
	Key   string
	Value string
}
