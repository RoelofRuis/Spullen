package main

import (
	"bytes"
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"golang.org/x/crypto/scrypt"
	"io/ioutil"
	"log"
	"os"
)

type EncryptedStorage struct {
	dbName string
	path   string
	pass   []byte
}

func (s *EncryptedStorage) Name() string {
	return s.dbName
}

func (s *EncryptedStorage) Read() ([]byte, error) {
	// READ
	data, err := ioutil.ReadFile(s.path)
	if err != nil {
		return nil, err
	}

	// DECRYPT
	plain, err := decrypt(s.pass, data)
	if err != nil {
		return nil, err
	}

	// GUNZIP
	r, err := gzip.NewReader(bytes.NewReader(plain))
	if err != nil {
		return nil, err
	}
	defer r.Close()

	unzipped, err := ioutil.ReadAll(r)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return unzipped, nil
}

func (s *EncryptedStorage) Write(data []byte) error {
	// GZIP
	buf := &bytes.Buffer{}

	w := gzip.NewWriter(buf)

	_, err := w.Write(data)
	if err != nil {
		return err
	}
	err = w.Flush()
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}

	// ENCRYPT
	encrypted, err := encrypt(s.pass, buf.Bytes())
	if err != nil {
		return err
	}

	// WRITE
	f, err := os.OpenFile(s.path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(encrypted)
	if err != nil {
		return err
	}

	return nil
}

func encrypt(pass, data []byte) ([]byte, error) {
	key, salt, err := deriveKey(pass, nil)
	if err != nil {
		return nil, err
	}

	blockCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	ciphertext = append(ciphertext, salt...)

	return ciphertext, nil
}

func decrypt(pass, data []byte) ([]byte, error) {
	salt, data := data[len(data)-32:], data[:len(data)-32]

	key, _, err := deriveKey(pass, salt)
	if err != nil {
		return nil, err
	}

	blockCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return nil, err
	}

	nonce, ciphertext := data[:gcm.NonceSize()], data[gcm.NonceSize():]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func deriveKey(password, salt []byte) ([]byte, []byte, error) {
	if salt == nil {
		salt = make([]byte, 32)
		if _, err := rand.Read(salt); err != nil {
			return nil, nil, err
		}
	}

	key, err := scrypt.Key(password, salt, 1048576, 8, 1, 32)
	if err != nil {
		return nil, nil, err
	}

	return key, salt, nil
}
