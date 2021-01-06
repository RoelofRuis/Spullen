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
	Open(name string, pass []byte, openExisting bool) error

	// Whether the database is dirty. A closed database is never dirty.
	IsDirty() bool

	// Register a storable to this database.
	// Upon opening the database, all Storables will receive a call to Initialize() with their appropriate data.
	// After a successful persist, all Storables will receive a call to their AfterPersist callback.
	//
	// Calling register twice with the same id will result in an error.
	Register(id string, h Storable) error

	// Persist the database.
	Persist() error

	// Close the database. A database should be opened before it can be closed.
	Close() error
}

func NewDatabase(useGzip bool, useEncryption bool) Database {
	return &fileDatabase{
		lock:      &sync.Mutex{},
		useGzip: useGzip,
		useEncryption: useEncryption,
		isOpened:  false,
		storage:   nil,
		storables: map[string]Storable{},
	}
}

type fileDatabase struct {
	lock sync.Locker

	useGzip bool
	useEncryption bool
	isOpened bool
	storage  storage

	storables map[string]Storable
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

func (db *fileDatabase) Open(name string, pass []byte, openExisting bool) error {
	if db.isOpened {
		return errors.New("database is already opened")
	}

	storage := &storageImpl{
		useGzip:       db.useGzip,
		useEncryption: db.useEncryption,
		dbName:        name,
		path:          fmt.Sprintf("%s.db", name),
		pass:          pass,
	}

	var dataMap = map[string][]byte{}
	if openExisting {
		data, err := storage.read()
		if err != nil {
			return err
		}
		dataMap = data
	}

	for name, s := range db.storables {
		data, hasKey := dataMap[name]
		if openExisting && !hasKey {
			return fmt.Errorf("data missing for storable [%s]", name)
		}
		err := s.Instantiate(data)
		if err != nil {
			return err
		}
	}

	db.lock.Lock()
	defer db.lock.Unlock()

	db.storage = storage
	db.isOpened = true

	return nil
}

func (db *fileDatabase) Register(id string, p Storable) error {
	_, exists := db.storables[id]
	if exists {
		return fmt.Errorf("storable with id [%s] was already registered", id)
	}

	db.lock.Lock()
	defer db.lock.Unlock()

	db.storables[id] = p

	return nil
}

func (db *fileDatabase) Persist() error {
	if !db.IsOpened() {
		return errors.New("database should be opened before it can be persisted")
	}

	if !db.IsDirty() {
		return nil
	}

	var dataMap = map[string][]byte{}
	for name, s := range db.storables {
		data, err := s.ToRaw()
		if err != nil {
			return err
		}

		dataMap[name] = data
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
	defer db.lock.Unlock()

	db.storage = nil
	db.isOpened = false

	return nil
}
