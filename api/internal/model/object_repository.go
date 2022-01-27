package model

import (
	"context"
	"github.com/roelofruis/spullen/internal/db"
	"time"
)

type ObjectRepository struct {
	DB *db.Proxy
}

func (r ObjectRepository) Insert(obj *Object) error {
	if obj.ID == ObjectID(0) {
		query := `
		INSERT INTO objects(added, name, quantity)
		VALUES (?, ?, ?)
		`

		args := []interface{}{obj.Added, obj.Name, obj.Quantity}

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
		obj.ID = ObjectID(id)

		return nil
	}

	query := `
	UPDATE objects SET added = ?, name = ?, quantity = ?
	WHERE id = ?
	`

	args := []interface{}{obj.Added, obj.Name, obj.Quantity, obj.ID}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := r.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r ObjectRepository) GetAll(name string) ([]*Object, error) {
	query := `
	SELECT id, added, name, quantity
	FROM objects`
	var params []interface{}

	if name != "" {
		query = query + "\nWHERE name LIKE ?"
		params = append(params, name)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := r.DB.QueryContext(ctx, query, params...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	objects := make([]*Object, 0)

	for rows.Next() {
		var object Object

		err := rows.Scan(
			&object.ID,
			&object.Added,
			&object.Name,
			&object.Quantity,
		)
		if err != nil {
			return nil, err
		}

		objects = append(objects, &object)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return objects, nil
}