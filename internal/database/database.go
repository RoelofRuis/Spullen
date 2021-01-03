package database

import (
	"bytes"
	"errors"
	"fmt"
	"sync"
)

type Storable interface {
	Instantiate([]byte) error
	ToRaw() ([]byte, error)
	IsDirty() bool
}

type Database interface {
	IsOpened() bool
	Name() string
	Open(name string, pass []byte, mode Mode) error
	IsDirty() bool
	Register(h Storable)
	Persist() error
	Close()
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
		storable: nil,
	}
}

type fileDatabase struct {
	lock     sync.Locker

	isOpened bool
	storage  storage

	storable Storable
}

func (db *fileDatabase) IsOpened() bool {
	return db.isOpened
}

func (db *fileDatabase) IsDirty() bool {
	if !db.isOpened || db.storable == nil{
		return false
	}

	return db.storable.IsDirty()
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

	if db.storable == nil {
		return errors.New("no storable is registered")
	}

	openExisting := mode & ModeOpenExisting == ModeOpenExisting
	useGzip := mode & ModeUseGzip == ModeUseGzip
	useEncryption := mode & ModeUseEncryption == ModeUseEncryption

	storage := &storageImpl{
		useGzip:       useGzip,
		useEncryption: useEncryption,
		dbName:        name,
		path:          fmt.Sprintf("%s.db", name),
		pass:          pass,
	}

	var buffer *bytes.Buffer
	if openExisting {
		data, err := storage.read()
		if err != nil {
			return err
		}
		buffer = bytes.NewBuffer(data)
	} else {
		buffer = &bytes.Buffer{}
	}

	err := db.storable.Instantiate(buffer.Bytes())
	if err != nil {
		return err
	}

	db.lock.Lock()
	db.storage = storage
	db.isOpened = true
	db.lock.Unlock()

	return nil
}

func (db *fileDatabase) Register(p Storable) {
	db.lock.Lock()
	db.storable = p
	db.lock.Unlock()
}

func (db *fileDatabase) Persist() error {
	if !db.IsOpened() {
		return errors.New("database should be opened before it can be persisted")
	}

	if db.storable == nil {
		return nil
	}

	if !db.IsDirty() {
		return nil
	}

	data, err := db.storable.ToRaw()
	if err != nil {
		return err
	}

	err = db.storage.write(data)
	if err != nil {
		return err
	}

	return nil
}

func (db *fileDatabase) Close() {
	db.lock.Lock()
	db.storage = nil
	db.isOpened = false
	db.lock.Unlock()
}
