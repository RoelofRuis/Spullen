package core

import "github.com/roelofruis/spullen"

type AppInfo struct {
	DevMode bool

	Alert string
}

type Index struct {
	AppInfo

	Databases []string
	Form      *OpenDatabaseForm
}

type EditObject struct {
	ExistingTags []string
	ExistingCategories []string
	ExistingPropertyKeys []string

	Form *ObjectForm
}

type View struct {
	AppInfo
	EditObject

	DatabaseIsDirty bool

	TotalCount int
	DbName string
	Objects []*spullen.Object
	PrivateMode bool
}

type Split struct {
	AppInfo
	EditObject

	Original *ObjectForm
}

type Edit struct {
	AppInfo
	EditObject
}