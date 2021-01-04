package database

import (
	"bytes"
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"golang.org/x/crypto/scrypt"
	"io/ioutil"
	"os"
	"strconv"
)

type storage interface {
	name() string
	read() (map[string][]byte, error)
	write(data map[string][]byte) error
}

type storageImpl struct {
	useGzip       bool
	useEncryption bool

	dbName string
	path   string
	pass   []byte
}

func (s *storageImpl) name() string {
	return s.dbName
}

 func (s *storageImpl) read() (map[string][]byte, error) {
 	data, err := s.readRaw()
 	if err != nil {
 		return nil, err
	}

	dataMap := map[string][]byte{}
 	for {
 		if len(data) == 0 {
 			break
		}

		lengthParts := bytes.SplitN(data, []byte(":"), 2)
		if len (lengthParts) != 2 {
			return nil, errors.New("invalid data: unable to determine block length")
		}

		blockLength, err := strconv.ParseInt(string(lengthParts[0]), 10, 64)
		if err != nil {
			return nil, err
		}

		nameParts := bytes.SplitN(lengthParts[1], []byte("\n"), 2)
		if len(nameParts) != 2 {
			return nil, errors.New("invalid data: unable to determine block name")
		}

		name := string(nameParts[0])
		if len(nameParts[1]) < int(blockLength) {
			return nil, errors.New("invalid data: block too short")
		}

		block := nameParts[1][0:blockLength]
		dataMap[name] = block

		data = nameParts[1][blockLength:]
	}

	return dataMap, nil
 }

func (s *storageImpl) readRaw() ([]byte, error) {
	data, err := ioutil.ReadFile(s.path)
	if err != nil {
		return nil, err
	}

	if s.useEncryption {
		data, err = decrypt(s.pass, data)
		if err != nil {
			return nil, err
		}
	}

	if s.useGzip {
		r, err := gzip.NewReader(bytes.NewReader(data))
		if err != nil {
			return nil, err
		}
		defer r.Close()
		data, err = ioutil.ReadAll(r)

		if err != nil {
			return nil, err
		}
	}

	return data, nil
}

func (s *storageImpl) write(dataMap map[string][]byte) error {
	var fullData []byte
	for name, data := range dataMap {
		blockId := fmt.Sprintf("%d:%s\n", len(data), name)
		block := append([]byte(blockId), data...)
		fullData = append(fullData, block...)
	}
	return s.writeRaw(fullData)
}

func (s *storageImpl) writeRaw(data []byte) error {
	if s.useGzip {
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

		data = buf.Bytes()
	}

	if s.useEncryption {
		var err error
		data, err = encrypt(s.pass, data)
		if err != nil {
			return err
		}
	}

	f, err := os.OpenFile(s.path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(data)
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
