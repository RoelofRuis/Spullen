package core

import (
	"github.com/roelofruis/spullen"
)

type DeleteForm struct {
	Id        spullen.ObjectId
	Reason    string
	RemovedAt string

	Errors map[string]string
}

func (f *DeleteForm) Validate() bool {
	f.Errors = make(map[string]string)

	if len(f.Id) != 16 {
		f.Errors["Id"] = "Id moet bestaan uit 16 tekens"
	}

	return len(f.Errors) == 0
}
