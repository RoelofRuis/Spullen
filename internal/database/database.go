package database

import (
	"errors"
	"fmt"
	"github.com/roelofruis/spullen/internal/spullen"
	"sync"
)

func NewDatabase(repoFactory spullen.ObjectRepositoryFactory) spullen.Database {
	return &FileDatabase{
		repoFactory: repoFactory,

		lock: &sync.Mutex{},
		isOpened: false,
		storage:  nil,
		objects:  nil,
	}
}

type FileDatabase struct {
	repoFactory spullen.ObjectRepositoryFactory

	lock sync.Locker
	isOpened bool
	storage  storage
	objects spullen.ObjectRepository
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

func (db *FileDatabase) Open(name string, pass []byte, isExisting bool) (spullen.ObjectRepository, error) {
	storage := &encryptedStorage{
		dbName: name,
		path: fmt.Sprintf("%s.db", name),
		pass: pass,
	}

	var repo spullen.ObjectRepository
	if isExisting {
		data, err := storage.read()
		if err != nil {
			return nil, err
		}

		repo, err = db.repoFactory.CreateFromData(data)
		if err != nil {
			return nil, err
		}
	} else {
		repo = db.repoFactory.CreateNew()
	}

	db.lock.Lock()
	db.storage = storage
	db.objects = repo
	db.isOpened = true
	db.lock.Unlock()

	return db.objects, nil
}

func (db *FileDatabase) Persist() error {
	if !db.IsOpened() {
		return errors.New("database should be opened before it can be persisted")
	}

	data, err := db.objects.ToRawData()
	if err != nil {
		return err
	}

	err = db.storage.write(data)
	if err != nil {
		return err
	}

	return nil
}

func (db *FileDatabase) Close() error {
	db.lock.Lock()
	db.storage = nil
	db.objects = nil
	db.isOpened = false
	db.lock.Unlock()

	return nil
}