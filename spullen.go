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
	Private    bool
}

type Property struct {
	Key   string
	Value string
}
