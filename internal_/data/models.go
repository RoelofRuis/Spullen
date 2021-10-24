package data

type Models struct {
	DB      *DBProxy
	Objects *ObjectModel
}

func NewModels(db *DBProxy) Models {
	return Models{
		DB:      db,
		Objects: &ObjectModel{DB: db},
	}
}
