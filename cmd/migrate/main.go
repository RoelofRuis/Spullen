package main

import (
	"github.com/roelofruis/spullen/internal_/migration"
	"log"
)

func main() {
	err := migration.Create("structure")
	if err != nil {
		log.Fatal(err)
	}
}
