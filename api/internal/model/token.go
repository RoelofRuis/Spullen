package model

import (
	"crypto/rand"
	"encoding/base32"
	"github.com/roelofruis/spullen/internal/validator"
)

func ValidateTokenPlaintext(v *validator.Validator, tokenPlaintext string) {
	v.Check(tokenPlaintext != "", "token", "must be provided")
	v.Check(len(tokenPlaintext) == 26, "token", "must be 26 bytes long")
}

type Token struct {
	plaintext string
}

func (r *Token) Refresh() error {
	randomBytes := make([]byte, 16)

	_, err := rand.Read(randomBytes)
	if err != nil {
		return err
	}

	r.plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)

	return nil
}

func (r *Token) Get() string {
	return r.plaintext
}