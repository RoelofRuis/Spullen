package main

type IndexForm struct {
	DatabaseName string
	Password     string
	IsExisting   bool

	Errors map[string]string
}

func (f *IndexForm) Validate() bool {
	f.Errors = make(map[string]string)

	if len(f.DatabaseName) == 0 {
		if f.IsExisting {
			f.Errors["LoadDatabaseName"] = "Selecteer een bestaande database"
		} else {
			f.Errors["NewDatabaseName"] = "Geef een database op"
		}
	}

	if len(f.Password) == 0 {
		if f.IsExisting {
			f.Errors["LoadPassword"] = "Wachtwoord mag niet leeg zijn"
		} else {
			f.Errors["NewPassword"] = "Wachtwoord mag niet leeg zijn"
		}
	}

	return len(f.Errors) == 0
}