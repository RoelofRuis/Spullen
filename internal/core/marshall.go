package core

import (
	"fmt"
	"github.com/roelofruis/spullen"
	"strconv"
	"strings"
)

type ObjectMarshallerImpl struct{}

func (o *ObjectMarshallerImpl) Unmarshall(record []string) (*spullen.Object, error) {
	form := &ObjectForm{
		Id:         record[0],
		TimeAdded:  record[1],
		Name:       record[2],
		Quantity:   record[3],
		Categories: record[4],
		Tags:       record[5],
		Properties: record[6],
		Hidden:     record[7],
		Notes:      record[8],
	}

	if !form.Validate() {
		return nil, fmt.Errorf("invalid object [%s]", record[0])
	}

	object, err := form.GetObject()
	if err != nil {
		return nil, err
	}

	return object, nil
}

func (o *ObjectMarshallerImpl) Marshall(obj *spullen.Object) []string {
	var properties []string
	for _, p := range obj.Properties {
		properties = append(properties, fmt.Sprintf("%s=%s", p.Key, p.Value))
	}

	record := []string{
		obj.Id,
		strconv.FormatInt(obj.Added.Unix(), 10),
		obj.Name,
		fmt.Sprintf("%d", obj.Quantity),
		strings.Join(obj.Categories, ","),
		strings.Join(obj.Tags, ","),
		strings.Join(properties, ","),
		strconv.FormatBool(obj.Hidden),
		obj.Notes,
	}

	return record
}