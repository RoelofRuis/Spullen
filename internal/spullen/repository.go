package spullen

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
	"sync"
)

func NewStorableObjectRepository() *StorableObjectRepository {
	return &StorableObjectRepository{
		lock: sync.Mutex{}, // FIXME: use RWMutex!

		objects: map[string]*Object{},
		dirty:   false,
	}
}

type StorableObjectRepository struct {
	lock sync.Mutex

	objects map[string]*Object
	dirty   bool
}

func (s *StorableObjectRepository) Get(id string) *Object {
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

func (s *StorableObjectRepository) GetAll() []*Object {
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

func (s *StorableObjectRepository) Put(o *Object) {
	s.lock.Lock()
	s.objects[o.Id] = o
	s.dirty = true
	s.lock.Unlock()
}

func (s *StorableObjectRepository) Has(id string) bool {
	_, hasKey := s.objects[id]
	return hasKey
}

func (s *StorableObjectRepository) Remove(id string) {
	s.lock.Lock()
	delete(s.objects, id)
	s.dirty = true
	s.lock.Unlock()
}

func (s *StorableObjectRepository) IsDirty() bool {
	return s.dirty
}

func (s *StorableObjectRepository) WasPersisted() {
	// FIXME: proper multithreaded usage requires checking whether the state was changed between `ToRaw` and this call.
	s.lock.Lock()
	s.dirty = false
	s.lock.Unlock()
}

func (s *StorableObjectRepository) Instantiate(data []byte) error {
	r := csv.NewReader(strings.NewReader(string(data)))
	r.Comma = ';'

	var objects = make(map[string]*Object)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
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
			return fmt.Errorf("invalid object [%s]", record[0])
		}

		object, err := form.GetObject()
		if err != nil {
			return err
		}

		objects[object.Id] = object
	}

	s.lock.Lock()
	s.objects = objects
	s.dirty = false
	s.lock.Unlock()

	return nil
}

func (s *StorableObjectRepository) ToRaw() ([]byte, error) {
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
