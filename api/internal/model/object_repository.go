package model

import (
	"context"
	"fmt"
	"github.com/roelofruis/spullen/internal/db"
	"strconv"
	"strings"
	"time"
)

type ObjectRepository struct {
	DB *db.Proxy
}

func (r ObjectRepository) Insert(obj *Object) error {
	query := db.Insert("objects", map[string]interface{}{"name": obj.Name, "description": obj.Description})

	if obj.ID != ObjectID(0) {
		query.Update("id", obj.ID)
	}

	res, err := r.DB.Exec(query)
	if err != nil {
		return err
	}

	if obj.ID == ObjectID(0) {
		id, err := res.LastInsertId()
		if err != nil {
			return err
		}
		obj.ID = ObjectID(id)
	}

	// TODO: insert tags

	for _, qChange := range obj.QuantityChanges {
		if qChange.ID == QuantityChangeID(0) {
			query := db.Insert("quantity_changes", map[string]interface{}{
				"object_id":   obj.ID,
				"at":          qChange.At,
				"quantity":    qChange.Quantity,
				"description": qChange.Description,
			})

			res, err := r.DB.Exec(query)
			if err != nil {
				return err
			}

			id, err := res.LastInsertId()
			if err != nil {
				return err
			}
			qChange.ID = QuantityChangeID(id)
		}
	}

	return nil
}

func (r ObjectRepository) GetOne(id ObjectID) (*Object, error) {
	objects, err := r.GetAll("", id)
	if err != nil {
		return nil, err
	}

	if len(objects) == 0 {
		return nil, db.ErrNoSuchRecord
	}

	return objects[0], nil
}

// TODO: configurable query
func (r ObjectRepository) GetAll(name string, id ObjectID) ([]*Object, error) {
	query := `
	SELECT id, name, description, COALESCE(tag_list, ""), COALESCE(quantity_change_list, "")
	FROM objects
	LEFT JOIN (
		SELECT object_id, group_concat(tag_id) AS tag_list
		FROM object_tags
		GROUP BY object_id
	) tags ON tags.object_id = objects.id
	LEFT JOIN (
		SELECT object_id, group_concat(id) AS quantity_change_list
		FROM quantity_changes
		GROUP BY object_id
	) quantity_changes ON quantity_changes.object_id = objects.id`
	var params []interface{}

	// TODO: improve with query builder
	if id != 0 {
		query = query + "\nWHERE id = ?"
		params = append(params, id)
	}

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
		var object = Object{
			Tags:            []TagID{},
			QuantityChanges: []*QuantityChange{},
		}

		var tagIdsString string
		var quantityChangeIdsString string

		err := rows.Scan(
			&object.ID,
			&object.Name,
			&object.Description,
			&tagIdsString,
			&quantityChangeIdsString,
		)
		if err != nil {
			return nil, err
		}

		if tagIdsString != "" {
			var tags []TagID
			for _, tagID := range strings.Split(tagIdsString, ",") {
				v, _ := strconv.Atoi(tagID)
				tags = append(tags, TagID(v))
			}
			object.Tags = tags
		}

		if quantityChangeIdsString != "" {
			quantityChangeList, err := r.getQuantityChanges(quantityChangeIdsString)
			if err != nil {
				return nil, err
			}
			object.QuantityChanges = quantityChangeList
		}

		objects = append(objects, &object)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return objects, nil
}

func (r ObjectRepository) getQuantityChanges(ids string) ([]*QuantityChange, error) {
	query := fmt.Sprintf(`
	SELECT id, at, quantity, description
	FROM quantity_changes
	WHERE id IN (%s)
	`, ids)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	quantityChanges := make([]*QuantityChange, 0)

	for rows.Next() {
		var quantityChange QuantityChange

		err := rows.Scan(
			&quantityChange.ID,
			&quantityChange.At,
			&quantityChange.Quantity,
			&quantityChange.Description,
		)
		if err != nil {
			return nil, err
		}

		quantityChanges = append(quantityChanges, &quantityChange)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return quantityChanges, nil
}