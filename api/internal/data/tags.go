package data

import (
	"context"
	"github.com/roelofruis/spullen/internal/model"
	"github.com/roelofruis/spullen/internal/validator"
	"time"
)

func ValidateTag(v *validator.Validator, tag *model.Tag) {
	v.Check(tag.Name != "", "name" , "must not be empty")
}

type TagModel struct {
	DB *DBProxy
}

func (r TagModel) Insert(tag *model.Tag) error {
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
	tag.ID = id

	return nil
}

func (r TagModel) GetAll() ([]*model.Tag, error) {
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

	objects := make([]*model.Tag, 0)

	for rows.Next() {
		var tag model.Tag

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
