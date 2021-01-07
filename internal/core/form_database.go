package core

import "strconv"

type DatabaseForm struct {
	Database        string
	Password        string
	ShowHiddenItems string

	Errors map[string]string

	ParsedShowHiddenItems bool
}

type OpenDatabaseForm struct {
	DatabaseForm

	AvailableDatabases []string
}

func (f *OpenDatabaseForm) Validate() bool {
	f.Errors = make(map[string]string)

	var found = false
	for _, s := range f.AvailableDatabases {
		if f.Database == s {
			found = true
			break
		}
	}
	if !found {
		f.Errors["Database"] = "Geef een bestaande database op"
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

type NewDatabaseForm struct {
	DatabaseForm
}

func (f *NewDatabaseForm) Validate() bool {
	f.Errors = make(map[string]string)

	if f.Database == "" {
		f.Errors["Database"] = "Geef een databasenaam op"
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