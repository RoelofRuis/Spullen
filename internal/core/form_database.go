package core

import "strconv"

type OpenDatabaseForm struct {
	Database        string
	Password        string
	ShowHiddenItems string

	Errors map[string]string

	showHiddenItems bool
}

func (f *OpenDatabaseForm) Validate() bool {
	f.Errors = make(map[string]string)

	if f.Database == "" {
		f.Errors["Database"] = "Geef een database op"
	}

	if len(f.Password) == 0 {
		f.Errors["Password"] = "Wachtwoord mag niet leeg zijn"
	}

	showHiddenItems, err := strconv.ParseBool(f.ShowHiddenItems)
	if err != nil {
		f.Errors["ShowHiddenItems"] = "Verborgen items tonen moet een geldige booleaanse waarde zijn"
	} else {
		f.showHiddenItems = showHiddenItems
	}

	return len(f.Errors) == 0
}
