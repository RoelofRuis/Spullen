package data

type Models struct {
	DB      *DBProxy
	Token   *TokenModel
	Objects *ObjectModel
	Tags    *TagModel
}

func NewModels(db *DBProxy) Models {
	return Models{
		DB:      db,
		Token:   &TokenModel{},
		Objects: &ObjectModel{DB: db},
		Tags:    &TagModel{DB: db},
	}
}
