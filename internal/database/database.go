package database

import (
	"bytes"
	"errors"
	"fmt"
	"sync"
)

type Persistable interface {
	ToRawData() ([]byte, error)
}

type Database interface {
	IsOpened() bool
	Name() string
	Open(name string, pass []byte, mode Mode) ([]byte, error)
	Register(p Persistable)
	Persist() error
	Close()
}

type Mode int

const ModeOpenExisting Mode = 0x1
const ModeUseGzip Mode = 0x2
const ModeUseEncryption Mode = 0x4

func NewDatabase() Database {
	return &FileDatabase{
		lock:     &sync.Mutex{},
		isOpened: false,
		storage:  nil,
		persistable: nil,
	}
}

type FileDatabase struct {
	lock     sync.Locker
	isOpened bool
	storage  storage
	persistable Persistable
}

func (db *FileDatabase) IsOpened() bool {
	return db.isOpened
}

func (db *FileDatabase) Name() string {
	if db.isOpened {
		return db.storage.name()
	}

	return ""
}

func (db *FileDatabase) Open(name string, pass []byte, mode Mode) ([]byte, error) {
	if db.isOpened {
		return nil, errors.New("database is already opened")
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
			return nil, err
		}
		buffer = bytes.NewBuffer(data)
	} else {
		buffer = &bytes.Buffer{}
	}

	db.lock.Lock()
	db.storage = storage
	db.isOpened = true
	db.lock.Unlock()

	return buffer.Bytes(), nil
}

func (db *FileDatabase) Register(p Persistable) {
	db.lock.Lock()
	db.persistable = p
	db.lock.Unlock()
}

func (db *FileDatabase) Persist() error {
	if !db.IsOpened() {
		return errors.New("database should be opened before it can be persisted")
	}

	if db.persistable == nil {
		return errors.New("no persistable was registered")
	}

	data, err := db.persistable.ToRawData()
	if err != nil {
		return err
	}

	err = db.storage.write(data)
	if err != nil {
		return err
	}

	return nil
}

func (db *FileDatabase) Close() {
	db.lock.Lock()
	db.storage = nil
	db.isOpened = false
	db.lock.Unlock()
}
