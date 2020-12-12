package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

const DBPATH string = "./db/objects.csv"

type FileStorage struct {
	Objects map[string]*Object
}

func NewFileStorage() (*FileStorage, error) {
	f, err := os.Open(DBPATH)
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
			Hidden:     record[7],
			Notes:      record[8],
		})
		if err != nil {
			return nil, err
		}

		objects[object.Id] = object
	}
	return &FileStorage{Objects: objects}, nil
}

func (s *FileStorage) Get(id string) *Object {
	if object, found := s.Objects[id]; found {
		return object
	}

	return nil
}

type objectName struct {
	id string
	name string
}

type objectNames []objectName

func (o objectNames) Len() int { return len(o) }
func (o objectNames) Less(i, j int) bool { return o[i].name < o[j].name }
func (o objectNames) Swap(i, j int) { o[i], o[j] = o[j], o[i] }

func (s *FileStorage) GetAll() []*Object {
	var identifiers objectNames = nil
	for _, o := range s.Objects {
		identifiers = append(identifiers, objectName{id: o.Id, name: o.Name})
	}
	sort.Sort(identifiers)

	var objectList []*Object = nil
	for _, i := range identifiers {
		objectList = append(objectList, s.Objects[i.id])
	}

	return objectList
}

func (s *FileStorage) PutObject(o *Object) error {
	s.Objects[o.Id] = o

	return s.writeToFile()
}

func (s *FileStorage) RemoveObject(id string) error {
	delete(s.Objects, id)

	return s.writeToFile()
}

func (s *FileStorage) writeToFile() error {
	f, err := os.OpenFile(DBPATH, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	w.Comma = ';'
	defer w.Flush()

	var data []string
	for _, o := range s.Objects {
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
			strconv.FormatBool(o.Hidden),
			o.Notes,
		}
		err := w.Write(data)
		if err != nil {
			return err
		}
	}

	return nil
}
