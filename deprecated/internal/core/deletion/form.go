package deletion

import (
	"errors"
	"github.com/roelofruis/spullen"
	"strconv"
	"time"
)

type Form struct {
	Id        spullen.ObjectId
	Reason    string
	RemovedAt string

	Errors map[string]string

	deletion *spullen.Deletion
}

func (f *Form) GetDeletion() (*spullen.Deletion, error) {
	if f.deletion == nil {
		return nil, errors.New("form has not been validated")
	}

	return f.deletion, nil
}

func (f *Form) Validate() bool {
	f.Errors = make(map[string]string)

	if len(f.Id) != 16 {
		f.Errors["Id"] = "Id moet bestaan uit 16 tekens"
	}

	t, err := strconv.ParseInt(f.RemovedAt, 10, 64)
	if err != nil {
		f.Errors["TimeAdded"] = "Geen geldige Unix tijdwaarde"
	}

	isValid := len(f.Errors) == 0
	if isValid {
		f.deletion = &spullen.Deletion{
			Id:        f.Id,
			DeletedAt: time.Unix(t, 0),
			Reason:    f.Reason,
		}
	}

	return isValid
}
