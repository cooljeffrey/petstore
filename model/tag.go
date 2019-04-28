package model

type Tag struct {
	ID   int64  `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
}

func NewTag(id int64, name string) *Tag {
	return &Tag{ID: id, Name: name}
}
