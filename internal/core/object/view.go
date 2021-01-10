package object

import "github.com/roelofruis/spullen"

type EditObject struct {
	ExistingTags         []string
	ExistingCategories   []string
	ExistingPropertyKeys []string

	Form *Form
}

type View struct {
	EditObject

	DatabaseIsDirty bool

	TotalCount  int
	DbName      string
	Objects     []*spullen.Object
	PrivateMode bool
}

type Split struct {
	EditObject

	Alert string

	Original *Form
}

type Edit struct {
	EditObject

	Alert string
}
