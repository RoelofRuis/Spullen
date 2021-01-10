package object

type ObjectView struct {
	Id         string
	Name       string
	AddedAt    string
	Quantity   int
	Categories []string
	Tags       []string
	Properties []string
	Hidden     bool
	Deleted    bool
}

type EditableObjectForm struct {
	ExistingTags         []string
	ExistingCategories   []string
	ExistingPropertyKeys []string

	Form *Form
}

type View struct {
	EditableObjectForm

	DatabaseIsDirty bool

	TotalCount          int
	DbName              string
	Objects             []ObjectView
	ShowingHiddenItems  bool
	ShowingDeletedItems bool
}

type Split struct {
	EditableObjectForm

	Alert string

	Original *Form
}

type Edit struct {
	EditableObjectForm

	Alert string
}
