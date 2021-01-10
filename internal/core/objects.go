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
	var objects []object.ObjectView
	for _, o := range v.objects.GetAll() {
		if !flags.ShowHiddenItems && o.Hidden {
			continue
		}

		isDeleted := v.deletions.Has(o.Id)

		if !flags.ShowDeletedItems && isDeleted {
			continue
		}

		var properties []string
		for _, p := range o.Properties {
			properties = append(properties, p.String())
		}

		objects = append(objects, object.ObjectView{
			Id:         string(o.Id),
			Name:       o.Name,
			AddedAt:    o.Added.Format("02-01-2006"),
			Quantity:   o.Quantity,
			Categories: o.Categories,
			Tags:       o.Tags,
			Properties: properties,
			Hidden:     o.Hidden,
			Deleted:    isDeleted,
		})
	}
	return objects
}
