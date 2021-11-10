package main

import (
	"github.com/roelofruis/spullen/internal/migration"
	"log"
)

func main() {
	err := migration.Create("structure")
	if err != nil {
		log.Fatal(err)
	}
}
