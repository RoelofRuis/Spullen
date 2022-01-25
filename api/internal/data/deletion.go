package data

import (
	"context"
	"github.com/roelofruis/spullen/internal/model"
	"time"
)

type DeletionModel struct {
	DB *DBProxy
}

func (r DeletionModel) Insert(objectId uint64, deletion *model.Deletion) error {
	query := `
	INSERT INTO deletions (object_id, deleted_at, description)
	VALUES (?, ?, ?)
	`

	args := []interface{}{objectId, deletion.DeletedAt, deletion.Description}

	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()

	_, err := r.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}