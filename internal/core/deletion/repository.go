package deletion

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/roelofruis/spullen"
	"io"
	"strconv"
	"strings"
	"sync"
)

type StorableDeletionRepository struct {
	lock sync.RWMutex

	deletions map[spullen.ObjectId]*spullen.Deletion
	dirty     bool
}

func NewRepository() *StorableDeletionRepository {
	return &StorableDeletionRepository{
		lock: sync.RWMutex{},

		deletions: map[spullen.ObjectId]*spullen.Deletion{},
		dirty:     false,
	}
}

func (s *StorableDeletionRepository) Get(id spullen.ObjectId) *spullen.Deletion {
	s.lock.RLock()
	defer s.lock.RUnlock()

	if deletion, found := s.deletions[id]; found {
		return deletion
	}

	return nil
}

func (s *StorableDeletionRepository) Put(d *spullen.Deletion) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.deletions[d.Id] = d
	s.dirty = true
}

func (s *StorableDeletionRepository) Has(id spullen.ObjectId) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()

	_, hasKey := s.deletions[id]
	return hasKey
}

// --- LOADING AND SAVING
// Ensuring it is a Storable
func (s *StorableDeletionRepository) IsDirty() bool {
	return s.dirty
}

// FIXME: proper multithreaded usage requires checking whether the state was changed between `ToRaw` and this call.
func (s *StorableDeletionRepository) AfterPersist() {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.dirty = false
}
func (s *StorableDeletionRepository) Instantiate(data []byte) error {
	r := csv.NewReader(strings.NewReader(string(data)))
	r.Comma = ';'

	var deletions = make(map[spullen.ObjectId]*spullen.Deletion)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		form := &Form{
			Id:        spullen.ObjectId(record[0]),
			RemovedAt: record[1],
			Reason:    record[2],
		}

		if !form.Validate() {
			return fmt.Errorf("invalid object[%s]", record[0])
		}

		deletion, err := form.GetDeletion()
		if err != nil {
			return err
		}

		deletions[deletion.Id] = deletion
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	s.deletions = deletions
	s.dirty = false

	return nil
}

func (s *StorableDeletionRepository) ToRaw() ([]byte, error) {
	b := &bytes.Buffer{}
	w := csv.NewWriter(b)
	w.Comma = ';'

	for _, d := range s.deletions {
		record := []string{
			string(d.Id),
			strconv.FormatInt(d.DeletedAt.Unix(), 10),
			d.Reason,
		}

		err := w.Write(record)
		if err != nil {
			return nil, err
		}
	}

	w.Flush()

	return b.Bytes(), nil
}
