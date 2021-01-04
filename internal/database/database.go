package database

import (
	"errors"
	"fmt"
	"sync"
)

// Implement this interface for any class that can be saved as and loaded from a raw binary format.
//
// Implementations should then be registered with a database.
type Storable interface {
	// Should return a name by which this storable can be uniquely identified within all the
	// storables registered with a database.
	Identifier() string

	// Called with the raw data when the storable is instantiated by the database.
	Instantiate([]byte) error

	// Called by the database to request the raw data representation for storage.
	ToRaw() ([]byte, error)

	// Called by the database to check whether the storable contains dirty data, to allow for storage optimizations.
	IsDirty() bool

	// Called by the database right after data was successfully persisted.
	AfterPersist()
}

type Database interface {
	// Whether the database is opened.
	IsOpened() bool

	// The name with which the database was opened. Can be empty if the database is not open.
	Name() string

	// Open the database by passing in the required information.
	// A database cannot be opened twice, and should be closed first before reopening.
	Open(name string, pass []byte, mode Mode) error

	// Whether the database is dirty. A closed database is never dirty.
	IsDirty() bool

	// Register a storable to this database.
	// The Storable will be instantiated when the database is opened and persisted when the database is persisted.
	Register(h Storable)

	// Persist the database.
	Persist() error

	// Close the database. A database should be opened before it can be closed.
	Close() error
}

type Mode int

const ModeOpenExisting Mode = 0x1
const ModeUseGzip Mode = 0x2
const ModeUseEncryption Mode = 0x4

func NewDatabase() Database {
	return &fileDatabase{
		lock:     &sync.Mutex{},
		isOpened: false,
		storage:  nil,
		storables: nil,
	}
}

type fileDatabase struct {
	lock sync.Locker

	isOpened bool
	storage  storage

	storables []Storable
}

func (db *fileDatabase) IsOpened() bool {
	return db.isOpened
}

func (db *fileDatabase) IsDirty() bool {
	if !db.isOpened {
		return false
	}

	for _, s := range db.storables {
		if s.IsDirty() {
			return true
		}
	}

	return false
}

func (db *fileDatabase) Name() string {
	if db.isOpened {
		return db.storage.name()
	}

	return ""
}

func (db *fileDatabase) Open(name string, pass []byte, mode Mode) error {
	if db.isOpened {
		return errors.New("database is already opened")
	}

	if len(db.storables) == 0 {
		return errors.New("no storable is registered")
	}

	openExisting := mode&ModeOpenExisting == ModeOpenExisting
	useGzip := mode&ModeUseGzip == ModeUseGzip
	useEncryption := mode&ModeUseEncryption == ModeUseEncryption

	storage := &storageImpl{
		useGzip:       useGzip,
		useEncryption: useEncryption,
		dbName:        name,
		path:          fmt.Sprintf("%s.db", name),
		pass:          pass,
	}

	if openExisting {
		dataMap, err := storage.read()
		if err != nil {
			return err
		}
		for _, s := range db.storables {
			data, hasKey := dataMap[s.Identifier()]
			if !hasKey {
				return fmt.Errorf("data missing for storable [%s]", s.Identifier())
			}
			err := s.Instantiate(data)
			if err != nil {
				return err
			}
		}
	} else {
		for _, s := range db.storables {
			err := s.Instantiate([]byte{})
			if err != nil {
				return err
			}
		}
	}

	db.lock.Lock()
	db.storage = storage
	db.isOpened = true
	db.lock.Unlock()

	return nil
}

func (db *fileDatabase) Register(p Storable) {
	db.lock.Lock()
	db.storables = append(db.storables, p)
	db.lock.Unlock()
}

func (db *fileDatabase) Persist() error {
	if !db.IsOpened() {
		return errors.New("database should be opened before it can be persisted")
	}

	if !db.IsDirty() {
		return nil
	}

	var dataMap = map[string][]byte{}
	for _, s := range db.storables {
		data, err := s.ToRaw()
		if err != nil {
			return err
		}

		dataMap[s.Identifier()] = data
	}

	err := db.storage.write(dataMap)
	if err != nil {
		return err
	}

	for _, s := range db.storables {
		s.AfterPersist()
	}

	return nil
}

func (db *fileDatabase) Close() error {
	if !db.IsOpened() {
		return errors.New("database should be opened before it can be closed")
	}

	db.lock.Lock()
	db.storage = nil
	db.isOpened = false
	db.lock.Unlock()

	return nil
}
