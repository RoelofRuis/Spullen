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

type Storage interface {
	GetAll() *ObjectSet
	AddObject(*Object) error
	RemoveObject(id string) error
}

type FileStorage struct {
	Objects map[string]*Object
}

func NewFileStorage() (*FileStorage, error) {
	f, err := os.Open("./db/objects.csv")
	if os.IsNotExist(err) {
		return &FileStorage{Objects: make(map[string]*Object)}, nil
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

		id := record[0]
		name := record[1]
		i, err := strconv.ParseInt(record[2], 10, 64)
		if err != nil {
			return nil, err
		}
		added := time.Unix(i, 0)
		quantity, err := strconv.ParseInt(record[3], 10, 32)
		if err != nil {
			return nil, err
		}
		categories := strings.Split(record[4], ",")
		tags := strings.Split(record[5], ",")
		var properties []*Property
		for _, p := range strings.Split(record[6], ",") {
			if len(p) == 0 {
				continue
			}
			keyValue := strings.Split(p, "=")
			if len(keyValue) != 2 {
				return nil, fmt.Errorf("encountered invalid property value [%s]", record[6])
			}
			properties = append(properties, &Property{keyValue[0], keyValue[1]})
		}
		private, err := strconv.ParseBool(record[7])
		if err != nil {
			return nil, err
		}

		object := &Object{
			Id:         id,
			Name:       name,
			Added:      added,
			Quantity:   int(quantity),
			Categories: categories,
			Tags:       tags,
			Properties: properties,
			Private:    private,
		}

		objects[id] = object
	}
	return &FileStorage{Objects: objects}, nil
}

func (ol *FileStorage) GetAll() *ObjectSet {
	return &ObjectSet{ol.Objects}
}

func (ol *FileStorage) AddObject(o *Object) error {
	ol.Objects[o.Id] = o

	return ol.writeToFile()
}

func (ol *FileStorage) RemoveObject(id string) error {
	delete(ol.Objects, id)

	return ol.writeToFile()
}

func (ol *FileStorage) writeToFile() error {
	f, err := os.OpenFile("./db/objects.csv", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
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
			properties = append(properties, fmt.Sprintf("%s=%s", p.Key, p.Value))
		}
		data = []string{
			o.Id,
			o.Name,
			fmt.Sprintf("%d", o.Quantity),
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
