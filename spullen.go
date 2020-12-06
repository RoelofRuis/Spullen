package main

import "time"

type ObjectList struct {
	Objects map[string]*Object
}

type Object struct {
	Id    string
	Name  string
	Added time.Time
	Tags  []string
}
