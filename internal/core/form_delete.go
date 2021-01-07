package core

import (
	"github.com/roelofruis/spullen"
)

type DeleteForm struct {
	Id        spullen.ObjectId
	Reason    string
	RemovedAt string
}
