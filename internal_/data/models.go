package data

type Models struct {
	Objects *ObjectModel
}

func NewModels(db DBProxy) Models {
	return Models{
		Objects: &ObjectModel{DB: db},
	}
}
