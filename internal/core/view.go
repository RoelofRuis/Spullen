package core

import (
	"github.com/roelofruis/spullen"
	"github.com/roelofruis/spullen/internal/core/object"
)

type AppInfo struct {
	DevMode bool
	DbOpen  bool
	Version Version
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

	Form *object.ObjectForm
}

type View struct {
	AppInfo
	EditObject

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

	Original *object.ObjectForm
}

type Edit struct {
	AppInfo
	EditObject

	Alert string
}

type Delete struct {
	AppInfo

	Alert string

	Original *object.ObjectForm
	Form     *DeleteForm
}
