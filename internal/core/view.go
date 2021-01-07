package core

import "github.com/roelofruis/spullen"

type AppInfo struct {
	DevMode bool
	DbOpen bool
	StoredVersion int
	AppVersion int
}

type Database struct {
	AppInfo

	Alert string

	Form *DatabaseForm
}

type EditObject struct {
	ExistingTags         []string
	ExistingCategories   []string
	ExistingPropertyKeys []string

	Form *ObjectForm
}

type View struct {
	AppInfo
	EditObject

	Alert string

	DatabaseIsDirty bool

	TotalCount  int
	DbName      string
	Objects     []*spullen.Object
	PrivateMode bool
}

type Split struct {
	AppInfo
	EditObject

	Alert string

	Original *ObjectForm
}

type Edit struct {
	AppInfo
	EditObject

	Alert string
}

type Delete struct {
	AppInfo

	Alert string

	Original *ObjectForm
	Form *DeleteForm
}