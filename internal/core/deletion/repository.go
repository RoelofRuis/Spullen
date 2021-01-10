package deletion

import (
	"github.com/roelofruis/spullen"
	"sync"
)

type StorableDeletionRepository struct {
	lock sync.RWMutex

	deletions  map[spullen.ObjectId]*spullen.Deletion
	dirty      bool
}

func NewRepository() *StorableDeletionRepository {
	return &StorableDeletionRepository{
		lock: sync.RWMutex{},

		dirty:   false,
	}
}

func (s *StorableDeletionRepository) Get(id spullen.ObjectId) *spullen.Deletion {
	// TODO: implement
}

func (s *StorableDeletionRepository) Put(*spullen.Deletion) {
	// TODO: implement
}

func (s *StorableDeletionRepository) Has(id spullen.ObjectId) bool {
	// TODO: implement
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
	// TODO: implement
}

func (s *StorableDeletionRepository) ToRaw() ([]byte, error) {
	// TODO: implement
}
