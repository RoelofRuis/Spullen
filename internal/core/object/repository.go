package object

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/roelofruis/spullen"
	"io"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type StorableObjectRepository struct {
	lock sync.RWMutex

	objects map[spullen.ObjectId]*spullen.Object
	dirty   bool
}

type Marshaller interface {
	Unmarshall(record []string) (*spullen.Object, error)
	Marshall(obj *spullen.Object) []string
}

func NewRepository() *StorableObjectRepository {
	return &StorableObjectRepository{
		lock: sync.RWMutex{},

		objects: map[spullen.ObjectId]*spullen.Object{},
		dirty:   false,
	}
}

type objectName struct {
	id   spullen.ObjectId
	name string
}

type objectNames []objectName

func (o objectNames) Len() int           { return len(o) }
func (o objectNames) Less(i, j int) bool { return o[i].name < o[j].name }
func (o objectNames) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }

func (s *StorableObjectRepository) Count() int {
	s.lock.RLock()
	defer s.lock.RUnlock()

	var total = 0
	for _, o := range s.objects {
		total += o.Quantity
	}

	return total
}

func (s *StorableObjectRepository) Get(id spullen.ObjectId) *spullen.Object {
	s.lock.RLock()
	defer s.lock.RUnlock()

	if object, found := s.objects[id]; found {
		return object
	}

	return nil
}

func (s *StorableObjectRepository) GetAll() []*spullen.Object {
	s.lock.RLock()
	defer s.lock.RUnlock()

	var identifiers objectNames = nil
	for _, o := range s.objects {
		identifiers = append(identifiers, objectName{id: o.Id, name: o.Name})
	}
	sort.Sort(identifiers)

	var objectList []*spullen.Object = nil
	for _, i := range identifiers {
		objectList = append(objectList, s.objects[i.id])
	}

	return objectList
}

func (s *StorableObjectRepository) Put(o *spullen.Object) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.objects[o.Id] = o
	s.dirty = true
}

func (s *StorableObjectRepository) GetDistinctCategories(includeHidden bool) []string {
	s.lock.RLock()
	defer s.lock.RUnlock()

	seen := map[string]struct{}{}
	var categories []string
	for _, o := range s.objects {
		if !includeHidden && o.Hidden {
			continue
		}
		for _, c := range o.Categories {
			_, found := seen[c]
			if !found {
				seen[c] = struct{}{}
				categories = append(categories, c)
			}
		}
	}

	return categories
}

func (s *StorableObjectRepository) GetDistinctTags(includeHidden bool) []string {
	s.lock.RLock()
	defer s.lock.RUnlock()

	seen := map[string]struct{}{}
	var tags []string
	for _, o := range s.objects {
		if !includeHidden && o.Hidden {
			continue
		}
		for _, t := range o.Tags {
			_, found := seen[t]
			if !found {
				seen[t] = struct{}{}
				tags = append(tags, t)
			}
		}
	}

	return tags
}

func (s *StorableObjectRepository) GetDistinctPropertyKeys(includeHidden bool) []string {
	s.lock.RLock()
	defer s.lock.RUnlock()

	seen := map[string]struct{}{}
	var propKeys []string
	for _, o := range s.objects {
		if !includeHidden && o.Hidden {
			continue
		}
		for _, p := range o.Properties {
			_, found := seen[p.Key]
			if !found {
				seen[p.Key] = struct{}{}
				propKeys = append(propKeys, p.Key)
			}
		}
	}

	return propKeys
}

func (s *StorableObjectRepository) Has(id spullen.ObjectId) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()

	_, hasKey := s.objects[id]
	return hasKey
}

func (s *StorableObjectRepository) Remove(id spullen.ObjectId) {
	s.lock.Lock()
	defer s.lock.Unlock()

	delete(s.objects, id)
	s.dirty = true
}

// --- LOADING AND SAVING
// Ensuring it is a Storable
func (s *StorableObjectRepository) IsDirty() bool {
	return s.dirty
}

// FIXME: proper multithreaded usage requires checking whether the state was changed between `ToRaw` and this call.
func (s *StorableObjectRepository) AfterPersist() {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.dirty = false
}
func (s *StorableObjectRepository) Instantiate(data []byte) error {
	r := csv.NewReader(strings.NewReader(string(data)))
	r.Comma = ';'

	var objects = make(map[spullen.ObjectId]*spullen.Object)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		form := &Form{
			Id:         spullen.ObjectId(record[0]),
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

		obj, err := form.GetObject()
		if err != nil {
			return err
		}

		objects[obj.Id] = obj
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	s.objects = objects
	s.dirty = false

	return nil
}

func (s *StorableObjectRepository) ToRaw() ([]byte, error) {
	b := &bytes.Buffer{}
	w := csv.NewWriter(b)
	w.Comma = ';'

	for _, obj := range s.GetAll() {
		var properties []string
		for _, p := range obj.Properties {
			properties = append(properties, p.String())
		}

		record := []string{
			string(obj.Id),
			strconv.FormatInt(obj.Added.Unix(), 10),
			obj.Name,
			fmt.Sprintf("%d", obj.Quantity),
			strings.Join(obj.Categories, ","),
			strings.Join(obj.Tags, ","),
			strings.Join(properties, ","),
			strconv.FormatBool(obj.Hidden),
			obj.Notes,
		}

		err := w.Write(record)
		if err != nil {
			return nil, err
		}
	}

	w.Flush()

	return b.Bytes(), nil
}
