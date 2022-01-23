package data

import (
	"context"
	"github.com/roelofruis/spullen/internal/validator"
	"time"
)

type Object struct {
	ID       int64     `json:"id"`
	Added    time.Time `json:"added"`
	Name     string    `json:"name"`
	Quantity int       `json:"quantity"`
}

func ValidateObject(v *validator.Validator, obj *Object) {
	v.Check(obj.Name != "", "name", "must not be empty")
	v.Check(obj.Quantity > 0, "quantity", "quantity must be a positive integer")
}

type ObjectModel struct {
	DB *DBProxy
}

func (r ObjectModel) Insert(obj *Object) error {
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
	obj.ID = id

	return nil
}

func (r ObjectModel) GetAll(name string) ([]*Object, error) {
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
