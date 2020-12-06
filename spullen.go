package main

import "time"

type ObjectList struct {
	Objects []*Object
}

type Object struct {
	Id    string
	Name  string
	Added time.Time
	Tags  []string
}
