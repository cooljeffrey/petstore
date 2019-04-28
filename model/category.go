package model

type Category struct {
	ID   int64  `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
}

func NewCategory(id int64, name string) *Category {
	return &Category{ID: id, Name: name}
}
