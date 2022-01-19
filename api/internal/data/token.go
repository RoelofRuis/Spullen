package data

import (
	"crypto/rand"
	"encoding/base32"
	"github.com/roelofruis/spullen/internal/validator"
)

type Token struct {
	Plaintext string
}

func ValidateTokenPlaintext(v *validator.Validator, tokenPlaintext string) {
	v.Check(tokenPlaintext != "", "token", "must be provided")
	v.Check(len(tokenPlaintext) == 26, "token", "must be 26 bytes long")
}

type TokenModel struct {
	token *Token
}

func (r *TokenModel) Refresh() error {
	randomBytes := make([]byte, 16)

	_, err := rand.Read(randomBytes)
	if err != nil {
		return err
	}

	r.token = &Token{
		Plaintext: base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes),
	}

	return nil
}

func (r *TokenModel) Get() *Token {
	return r.token
}
