package main

import (
	"github.com/roelofruis/spullen/internal/migration"
	"log"
)

func main() {
	err := migration.Create("items")
	if err != nil {
		log.Fatal(err)
	}
}