package database

import (
	"fmt"
	"github.com/roelofruis/spullen"
	"github.com/roelofruis/spullen/internal/util"
	"log"
	"net/http"
	"strconv"
)

func NewDatabaseForm(finder *util.Finder) *Form {
	databases, err := finder.FindDatabases()
	if err != nil {
		log.Print(fmt.Sprintf("Error finding database: %s", err.Error()))
		databases = []string{}
	}
	return &Form{AvailableDatabases: databases}
}

type Form struct {
	IsExistingDatabase bool

	Database         string
	Password         string
	ShowHiddenItems  string
	ShowDeletedItems string

	AvailableDatabases []string

	Errors map[string]string

	dataFlags *spullen.DataFlags
}

func (f *Form) FillFromRequest(r *http.Request) {
	f.Database = r.PostFormValue("database")
	f.Password = r.PostFormValue("password")
	f.ShowHiddenItems = r.PostFormValue("show-hidden-items")
	f.ShowDeletedItems = r.PostFormValue("show-deleted-items")
}

func (f *Form) GetDataFlags() *spullen.DataFlags {
	return f.dataFlags
}

func (f *Form) Validate() bool {
	f.Errors = make(map[string]string)

	var found = false
	for _, s := range f.AvailableDatabases {
		if f.Database == s {
			found = true
			break
		}
	}

	if f.IsExistingDatabase {
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
	}

	showDeletedItems, err := strconv.ParseBool(f.ShowDeletedItems)
	if err != nil {
		f.Errors["ShowDeletedItems"] = "Verwijderde items tonen moet een geldige booleaanse waarde zijn"
	}

	isValid := len(f.Errors) == 0
	if isValid {
		f.dataFlags = &spullen.DataFlags{
			ShowHiddenItems:  showHiddenItems,
			ShowDeletedItems: showDeletedItems,
		}
	}

	return isValid
}
