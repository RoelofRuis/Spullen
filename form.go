package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type ObjectForm struct {
	Id string
	TimeAdded string

	Name string
	Quantity string
	Categories string
	Tags string
	Properties string
	Private string
}

func MakeForm(o *Object) *ObjectForm {
	var propertyStrings []string = nil
	for _, p := range o.Properties {
		propertyStrings = append(propertyStrings, p.Key + "=" + p.Value)
	}

	var private = ""
	if o.Private {
		private = "on"
	}
	return &ObjectForm{
		Id:         o.Id,
		TimeAdded:  strconv.FormatInt(o.Added.Unix(), 10),
		Name:       o.Name,
		Quantity:   strconv.FormatInt(int64(o.Quantity), 10),
		Categories: strings.Join(o.Categories, ","),
		Tags:       strings.Join(o.Tags, ","),
		Properties: strings.Join(propertyStrings, ","),
		Private:    private,
	}
}

func ParseObjectForm(r *ObjectForm) (*Object, error) {
	id := r.Id
	if id == "" {
		id = randSeq(16)
	}

	var timeAdded time.Time
	if r.TimeAdded == "" {
		timeAdded = time.Now().Truncate(time.Second)
	} else {
		t, err := strconv.ParseInt(r.TimeAdded, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("unable to parse time value %s", r.TimeAdded)
		}
		timeAdded = time.Unix(t, 0)
	}

	quantity, err := strconv.ParseInt(r.Quantity, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("unable to parse quantity value %s", r.Quantity)
	}

	var categories []string
	for _, c := range strings.Split(r.Categories, ",") {
		categories = append(categories, strings.ToLower(c))
	}

	var tags []string
	for _, t := range strings.Split(r.Tags, ",") {
		tags = append(tags, strings.ToLower(t))
	}

	var properties []*Property
	for _, p := range strings.Split(r.Properties, ",") {
		if len(p) == 0 {
			continue
		}
		keyValue := strings.Split(p, "=")
		if len(keyValue) != 2 {
			return nil, fmt.Errorf("invalid property value %s", p)
		}
		properties = append(properties, &Property{
			strings.ToLower(keyValue[0]),
			strings.ToLower(keyValue[1]),
		})
	}

	var private bool
	switch r.Private {
	case "on", "true":
		private = true
	case "", "false":
		private = false
	default:
		return nil, fmt.Errorf("invalid private value %s", r.Private)
	}

	return &Object{
		Id: id,
		Name: strings.ToLower(r.Name),
		Quantity: int(quantity),
		Added: timeAdded,
		Categories: categories,
		Tags: tags,
		Properties: properties,
		Private: private,
	}, nil
}