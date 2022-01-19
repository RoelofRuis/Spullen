package data

type Models struct {
	DB      *DBProxy
	Token   *TokenModel
	Objects *ObjectModel
}

func NewModels(db *DBProxy) Models {
	return Models{
		DB:      db,
		Token:   &TokenModel{},
		Objects: &ObjectModel{DB: db},
	}
}
