package core

import "github.com/roelofruis/spullen"

type appInfoModel struct {
	DevMode bool
}

type alertModel struct {
	Alert string
}

type IndexModel struct {
	appInfoModel
	alertModel

	Databases []string
	From *IndexForm
}

type editableObjectModel struct {
	ExistingTags []string
	ExistingCategories []string
	ExistingPropertyKeys []string

	Form *ObjectForm
}

type ViewModel struct {
	appInfoModel
	alertModel
	editableObjectModel

	DatabaseIsDirty bool

	TotalCount int
	DbName string
	Objects []*spullen.Object
	PrivateMode bool
}

type SplitModel struct {
	appInfoModel
	alertModel
	editableObjectModel

	Original *ObjectForm
}

type EditModel struct {
	appInfoModel
	alertModel
	editableObjectModel
}