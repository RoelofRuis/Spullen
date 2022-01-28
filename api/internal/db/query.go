package db

import (
	"fmt"
	"strings"
)

type InsertStatement struct {
	table    string
	columns  []string
	values   []interface{}
	idColumn string
	idValue  interface{}
}

func Insert(table string, args map[string]interface{}) *InsertStatement {
	var columns []string
	var values []interface{}

	for col, val := range args {
		columns = append(columns, col)
		values = append(values, val)
	}

	return &InsertStatement{
		table:  table,
		columns: columns,
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
		var params []string

		for x := 0; x < len(i.columns); x++ {
			params = append(params, "?")
		}
		return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
			i.table,
			strings.Join(i.columns, ","),
			strings.Join(params, ","),
		)
	}

	var params []string
	for _, column := range i.columns {
		params = append(params, fmt.Sprintf("%s = ?", column))
	}
	return fmt.Sprintf("UPDATE %s SET %s WHERE %s = ?",
		i.table,
		strings.Join(params, ","),
		i.idColumn,
	)
}

func (i *InsertStatement) args() []interface{} {
	var values = i.values

	if i.idColumn != "" {
		values = append(values, i.idValue)
	}

	return values
}
