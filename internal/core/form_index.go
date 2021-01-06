package core

import "strconv"

type IndexForm struct {
	ExistingDatabaseName string
	NewDatabaseName      string
	Password             string
	PrivateMode          string

	Errors map[string]string

	database      string
	isNew         bool
	isPrivateMode bool
}

func (f *IndexForm) Validate() bool {
	f.Errors = make(map[string]string)

	existingSelected := len(f.ExistingDatabaseName) > 0
	newSelected := len(f.NewDatabaseName) > 0

	if !existingSelected && !newSelected {
		f.Errors["Database"] = "Geef een database op"
	}

	f.isNew = newSelected
	if f.isNew {
		f.database = f.NewDatabaseName
	} else {
		f.database = f.ExistingDatabaseName
	}

	if len(f.Password) == 0 {
		f.Errors["Password"] = "Wachtwoord mag niet leeg zijn"
	}

	isPrivate, err := strconv.ParseBool(f.PrivateMode)
	if err != nil {
		f.Errors["PrivateMode"] = "Privémodus moet een geldige booleaanse waarde zijn"
	} else {
		f.isPrivateMode = isPrivate
	}

	return len(f.Errors) == 0
}
