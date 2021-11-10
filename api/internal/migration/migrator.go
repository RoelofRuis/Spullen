package migration

import (
	"bytes"
	"context"
	"database/sql"
	"embed"
	_ "embed"
	"errors"
	"fmt"
	"os"
	"text/template"
	"time"
)

//go:embed template.txt
var f embed.FS

type Migration struct {
	Version string
	Up      func(tx *sql.Tx) error

	done bool
}

type Migrator struct {
	db         *sql.DB
	Versions   []string
	Migrations map[string]*Migration
}

var migrator = &Migrator{
	Versions:   []string{},
	Migrations: map[string]*Migration{},
}

func (m *Migrator) AddMigration(mg *Migration) {
	m.Migrations[mg.Version] = mg

	index := 0
	for index < len(m.Versions) {
		if m.Versions[index] > mg.Version {
			break
		}
		index++
	}

	m.Versions = append(m.Versions, mg.Version)
	copy(m.Versions[index+1:], m.Versions[index:])
	m.Versions[index] = mg.Version
}

func (m *Migrator) Up() error {
	tx, err := m.db.BeginTx(context.TODO(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	for _, v := range m.Versions {
		mg := m.Migrations[v]

		if mg.done {
			continue
		}

		fmt.Println("running migration", mg.Version)
		if err := mg.Up(tx); err != nil {
			tx.Rollback()
			return err
		}

		if _, err := tx.Exec("INSERT INTO `schema_migrations` VALUES(?)", mg.Version); err != nil {
			tx.Rollback()
			return err
		}
		fmt.Println("finished running migration", mg.Version)
	}

	tx.Commit()

	return nil
}

func Init(db *sql.DB) (*Migrator, error) {
	migrator.db = db

	if _, err := db.Exec("CREATE TABLE IF NOT EXISTS `schema_migrations` (version varchar(255));"); err != nil {
		return nil, err
	}

	rows, err := db.Query("SELECT version FROM `schema_migrations`;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var version string
		err := rows.Scan(&version)
		if err != nil {
			return nil, err
		}

		if migrator.Migrations[version] != nil {
			migrator.Migrations[version].done = true
		}
	}

	return migrator, nil
}

func Create(name string) error {
	version := time.Now().Format("20060102150405")

	in := struct {
		Version string
		Name    string
	}{
		Version: version,
		Name:    name,
	}

	var out bytes.Buffer

	t := template.Must(template.ParseFS(f, "template.txt"))
	err := t.Execute(&out, in)
	if err != nil {
		return errors.New("unable to execute template: " + err.Error())
	}

	f, err := os.Create(fmt.Sprintf("./internal/migration/%s_%s.go", version, name))
	if err != nil {
		return errors.New("unable to create migration file: " + err.Error())
	}
	defer f.Close()

	if _, err := f.WriteString(out.String()); err != nil {
		return errors.New("unable to write to migration file: " + err.Error())
	}

	fmt.Println("generated new migration files...", f.Name())
	return nil
}
