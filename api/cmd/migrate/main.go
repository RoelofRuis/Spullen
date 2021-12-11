package main

import (
	"flag"
	"fmt"
	"github.com/roelofruis/spullen/internal/migration"
	"log"
	"os"
)

func main() {
	var migrationName string

	flag.StringVar(&migrationName, "name", "", "Specify a name for the migration")
	flag.Parse()

	if migrationName == "" {
		fmt.Printf("Migration name is required. Set with --name <name>\n")
		os.Exit(1)
	}

	err := migration.Create(migrationName)
	if err != nil {
		log.Fatal(err)
	}
}
