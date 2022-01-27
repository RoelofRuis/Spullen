package db

import (
	"fmt"
	"strings"
)

type InsertStatement struct {
	table    string
	values   map[string]interface{}
	idColumn string
	idValue  interface{}
}

func Insert(table string, values map[string]interface{}) *InsertStatement {
	return &InsertStatement{
		table:  table,
		values: values,
	}
}

func (i *InsertStatement) Update(idColumn string, idValue interface{}) *InsertStatement {
	return &InsertStatement{
		idColumn: idColumn,
		idValue: idValue,
	}
}

func (i *InsertStatement) query() string {
	if i.idColumn == "" {
		var columns []string
		var params []string

		for column := range i.values {
			columns = append(columns, column)
			params = append(params, "?")
		}
		return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
			i.table,
			strings.Join(columns, ","),
			strings.Join(params, ","),
		)
	}

	var params []string
	for column := range i.values {
		params = append(params, fmt.Sprintf("%s = ?", column))
	}
	return fmt.Sprintf("UPDATE %s SET %s WHERE %s = ?",
		i.table,
		strings.Join(params, ","),
		i.idColumn,
	)
}

func (i *InsertStatement) args() []interface{} {
	var values []interface{}

	for _, value := range i.values {
		values = append(values, value)
	}

	if i.idColumn != "" {
		values = append(values, i.idValue)
	}

	return values
}
