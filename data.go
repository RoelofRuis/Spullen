package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

func LoadObjectList() (*ObjectList, error) {
	f, err := os.Open("./data/objects.csv")
	if os.IsNotExist(err) {
		return &ObjectList{Objects: make(map[string]*Object)}, nil
	}
	if err != nil {
		return nil, err
	}

	defer f.Close()

	r := csv.NewReader(f)
	r.Comma = ';'

	var objects = make(map[string]*Object)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		i, err := strconv.ParseInt(record[2], 10, 64)
		if err != nil {
			return nil, err
		}
		private, err := strconv.ParseBool(record[6])
		if err != nil {
			return nil, err
		}
		id := record[0]
		var properties []*Property
		for _, p := range strings.Split(record[5], ",") {
			keyValue := strings.Split(p, "=")
			if len(keyValue) != 2 {
				return nil, fmt.Errorf("encountered invalid property value [%s]", record[5])
			}
			properties = append(properties, &Property{keyValue[0], keyValue[1]})
		}
		added := time.Unix(i, 0)
		object := &Object{
			Id:    id,
			Name:  record[1],
			Added: added,
			Categories: strings.Split(record[3], ","),
			Tags: strings.Split(record[4], ","),
			Properties: properties,
			Private: private,
		}

		objects[id] = object
	}
	return &ObjectList{Objects: objects}, nil
}

func (ol *ObjectList) AddObject(o *Object) {
	ol.Objects[o.Id] = o
}

func (ol *ObjectList) RemoveObject(id string) {
	delete(ol.Objects, id)
}

func (ol *ObjectList) Save() error {
	f, err := os.OpenFile("./data/objects.csv", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	w.Comma = ';'
	defer w.Flush()

	var data []string
	for _, o := range ol.Objects {
		var properties []string
		for _, p := range o.Properties {
			properties = append(properties,  fmt.Sprintf("%s=%s", p.Key, p.Value))
		}
		data = []string{
			o.Id,
			o.Name,
			strconv.FormatInt(o.Added.Unix(), 10),
			strings.Join(o.Categories, ","),
			strings.Join(o.Tags, ","),
			strings.Join(properties, ","),
			strconv.FormatBool(o.Private),
		}
		err := w.Write(data)
		if err != nil {
			return err
		}
	}

	return nil
}
