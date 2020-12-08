package main

import "time"

type ObjectList struct {
	Objects map[string]*Object
}

type Object struct {
	Id    string
	Name  string
	Added time.Time
	Categories []string
	Tags []string
	Properties []*Property
	Private bool
}

type Property struct {
	Key string
	Value string
}