package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type Storage interface {
	GetAll() *map[string]*Object
	Get(id string) *Object
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

		object, err := ParseObjectForm(&ObjectForm{
			Id:         record[0],
			TimeAdded:  record[1],
			Name:       record[2],
			Quantity:   record[3],
			Categories: record[4],
			Tags:       record[5],
			Properties: record[6],
			Private:    record[7],
		})
		if err != nil {
			return nil, err
		}

		objects[object.Id] = object
	}
	return &FileStorage{Objects: objects}, nil
}

func (ol *FileStorage) Get(id string) *Object {
	if object, found := ol.Objects[id]; found {
		return object
	}

	return nil
}

func (ol *FileStorage) GetAll() *map[string]*Object {
	return &ol.Objects
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
			strconv.FormatInt(o.Added.Unix(), 10),
			o.Name,
			fmt.Sprintf("%d", o.Quantity),
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
