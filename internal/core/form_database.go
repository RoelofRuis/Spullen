package core

import (
	"strconv"
)

func NewDatabaseForm(finder *Finder) (*DatabaseForm, error) {
	databases, err := finder.FindDatabases()
	if err != nil {
		return nil, err
	}
	return &DatabaseForm{AvailableDatabases: databases}, nil
}

type DatabaseForm struct {
	Database        string
	Password        string
	ShowHiddenItems string

	AvailableDatabases []string

	Errors map[string]string

	ParsedShowHiddenItems bool
}

func (f *DatabaseForm) Validate(isExisting bool) bool {
	f.Errors = make(map[string]string)

	var found = false
	for _, s := range f.AvailableDatabases {
		if f.Database == s {
			found = true
			break
		}
	}

	if isExisting {
		if !found {
			f.Errors["Database"] = "Geef een bestaande database op"
		}
	} else {
		if f.Database == "" {
			f.Errors["Database"] = "Geef een databasenaam op"
		}
		if found {
			f.Errors["Database"] = "Er bestaat al een database met deze naam"
		}
	}

	if len(f.Password) == 0 {
		f.Errors["Password"] = "Wachtwoord mag niet leeg zijn"
	}

	showHiddenItems, err := strconv.ParseBool(f.ShowHiddenItems)
	if err != nil {
		f.Errors["ShowHiddenItems"] = "Verborgen items tonen moet een geldige booleaanse waarde zijn"
	} else {
		f.ParsedShowHiddenItems = showHiddenItems
	}

	return len(f.Errors) == 0
}
