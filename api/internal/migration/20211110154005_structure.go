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
        name TEXT NOT NULL,
        description TEXT
    );`)
    if err != nil {
        return err
    }

    _, err = tx.Exec(`
    CREATE TABLE deletions (
        object_id INTEGER NOT NULL,
        deleted_at DATETIME NOT NULL,
        description TEXT,
        FOREIGN KEY (object_id) REFERENCES objects(id)
    );`)
    if err != nil {
        return err
    }

    _, err = tx.Exec(`
    CREATE TABLE tags (
        id INTEGER PRIMARY KEY,
        name TEXT NOT NULL,
        description TEXT,
        is_system_tag BOOLEAN NOT NULL CHECK (is_system_tag IN (0, 1))
    );`)
    if err != nil {
        return err
    }

    _, err = tx.Exec(`
    CREATE TABLE object_tags (
        object_id INTEGER NOT NULL,
        tag_id INTEGER NOT NULL,
        PRIMARY KEY (object_id, tag_id),
        FOREIGN KEY (object_id) REFERENCES objects(id),
        FOREIGN KEY (tag_id) REFERENCES tags(id)
    );`)

    return nil
}