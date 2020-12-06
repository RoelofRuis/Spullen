package main

import (
	"encoding/csv"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

func LoadObjectList() (*ObjectList, error) {
	f, err := os.Open("./data/objects.csv")
	if os.IsNotExist(err) {
		return &ObjectList{Objects: nil}, nil
	}
	if err != nil {
		return nil, err
	}

	defer f.Close()

	r := csv.NewReader(f)
	r.Comma = ';'

	var objects []*Object = nil
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
		added := time.Unix(i, 0)
		object := &Object{
			Id:    record[0],
			Name:  record[1],
			Added: added,
			Tags:  strings.Split(record[3], ","),
		}

		objects = append(objects, object)
	}
	return &ObjectList{Objects: objects}, nil
}

func (ol *ObjectList) AddObject(o *Object) {
	ol.Objects = append(ol.Objects, o)
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
		data = []string{
			o.Id,
			o.Name,
			strconv.FormatInt(o.Added.Unix(), 10),
			strings.Join(o.Tags, ","),
		}
		err := w.Write(data)
		if err != nil {
			return err
		}
	}

	return nil
}
