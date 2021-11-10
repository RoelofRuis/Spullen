package migration

import "database/sql"

func init() {
    migrator.AddMigration(&Migration{
        Version: "20211110154005",
        Up:      mig_20211110154005_structure_up,
    })
}

func mig_20211110154005_structure_up(tx *sql.Tx) error {
    _, err := tx.Exec(`
    CREATE TABLE objects (
        id INTEGER PRIMARY KEY,
        added DATETIME NOT NULL,
        name TEXT NOT NULL,
        quantity INTEGER NOT NULL
    );`)
    if err != nil {
        return err
    }

    return nil
}