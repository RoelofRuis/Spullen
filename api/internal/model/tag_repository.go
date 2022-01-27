package model

import (
	"context"
	"github.com/roelofruis/spullen/internal/db"
	"time"
)

type TagRepository struct {
	DB *db.Proxy
}

func (r TagRepository) Insert(tag *Tag) error {
	if tag.ID == TagID(0) {
		query := `
		INSERT INTO tags(name, description, is_system_tag)
		VALUES (?, ?, ?)
		`

		args := []interface{}{tag.Name, tag.Description, tag.IsSystemTag}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		res, err := r.DB.ExecContext(ctx, query, args...)
		if err != nil {
			return err
		}

		id, err := res.LastInsertId()
		if err != nil {
			return err
		}
		tag.ID = TagID(id)

		return nil
	}

	query := `
	UPDATE objects SET name = ?, description = ?, is_system_tag = ?
	WHERE id = ?
	`

	args := []interface{}{tag.Name, tag.Description, tag.IsSystemTag, tag.ID}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := r.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r TagRepository) GetAll() ([]*Tag, error) {
	query := `
	SELECT id, name, description, is_system_tag
	FROM tags`
	var params []interface{}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := r.DB.QueryContext(ctx, query, params...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	objects := make([]*Tag, 0)

	for rows.Next() {
		var tag Tag

		err := rows.Scan(
			&tag.ID,
			&tag.Name,
			&tag.Description,
			&tag.IsSystemTag,
		)
		if err != nil {
			return nil, err
		}

		objects = append(objects, &tag)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return objects, nil
}
