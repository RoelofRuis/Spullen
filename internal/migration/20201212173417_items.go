package migration

import "database/sql"

func init() {
    migrator.AddMigration(&Migration{
        Version: "20201212173417",
        Up:      mig_20201212173417_items_up,
    })
}

func mig_20201212173417_items_up(tx *sql.Tx) error {
    return nil
}
