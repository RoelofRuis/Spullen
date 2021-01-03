package spullen

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
)

func LoadRepository(data []byte) (*ObjectRepositoryImpl, error) {
	r := csv.NewReader(strings.NewReader(string(data)))
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

		objects[object.Id] = object
	}

	return &ObjectRepositoryImpl{
		objects: objects,
	}, nil
}

type ObjectRepositoryImpl struct {
	objects map[string]*Object
}

func (s *ObjectRepositoryImpl) Get(id string) *Object {
	if object, found := s.objects[id]; found {
		return object
	}

	return nil
}

type objectName struct {
	id   string
	name string
}

type objectNames []objectName

func (o objectNames) Len() int           { return len(o) }
func (o objectNames) Less(i, j int) bool { return o[i].name < o[j].name }
func (o objectNames) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }

func (s *ObjectRepositoryImpl) GetAll() []*Object {
	var identifiers objectNames = nil
	for _, o := range s.objects {
		identifiers = append(identifiers, objectName{id: o.Id, name: o.Name})
	}
	sort.Sort(identifiers)

	var objectList []*Object = nil
	for _, i := range identifiers {
		objectList = append(objectList, s.objects[i.id])
	}

	return objectList
}

func (s *ObjectRepositoryImpl) Put(o *Object) {
	s.objects[o.Id] = o
}

func (s *ObjectRepositoryImpl) Has(id string) bool {
	_, hasKey := s.objects[id]
	return hasKey
}

func (s *ObjectRepositoryImpl) Remove(id string) {
	delete(s.objects, id)
}

func (s *ObjectRepositoryImpl) ToRawData() ([]byte, error) {
	b := &bytes.Buffer{}
	w := csv.NewWriter(b)
	w.Comma = ';'

	var data []string
	for _, o := range s.GetAll() {
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
			return nil, err
		}
	}

	w.Flush()

	return b.Bytes(), nil
}
