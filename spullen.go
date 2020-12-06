package main

import "time"

type ObjectList struct {
	Objects []*Object
}

type Object struct {
	Name string
	OwnedSince time.Time
}