package data

import (
	"context"
	"github.com/roelofruis/spullen/internal/model"
	"github.com/roelofruis/spullen/internal/validator"
	"time"
)

func ValidateObject(v *validator.Validator, obj *model.Object) {
	v.Check(obj.Name != "", "name", "must not be empty")
	v.Check(obj.Quantity > 0, "quantity", "quantity must be a positive integer")
}

type ObjectModel struct {
	DB *DBProxy
}

func (r ObjectModel) Insert(obj *model.Object) error {
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

func (r ObjectModel) GetAll(name string) ([]*model.Object, error) {
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

	objects := make([]*model.Object, 0)

	for rows.Next() {
		var object model.Object

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
