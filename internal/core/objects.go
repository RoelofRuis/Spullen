package core

import (
	"github.com/roelofruis/spullen"
	"github.com/roelofruis/spullen/internal/core/object"
)

func NewObjectViewer(objects spullen.ObjectRepository, deletions spullen.DeletionRepository) *ObjectViewer {
	return &ObjectViewer{
		objects:   objects,
		deletions: deletions,
	}
}

type ObjectViewer struct {
	objects   spullen.ObjectRepository
	deletions spullen.DeletionRepository
}

func (v *ObjectViewer) GetAll(flags *spullen.DataFlags) []object.ObjectView {
	// TODO: implement
}
