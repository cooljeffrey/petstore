package model

type User struct {
	ID         int64  `json:"id" bson:"id"`
	Username   string `json:"username" bson:"username"`
	Firstname  string `json:"firstName" bson:"firstName"`
	Lastname   string `json:"lastName" bson:"lastName"`
	Email      string `json:"email" bson:"email"`
	Password   string `json:"password" bson:"password"`
	Phone      string `json:"phone" bson:"phone"`
	UserStatus int32  `json:"userStatus" bson:"userStatus"`
}

func NewUser(id int64, username, firstname, lastname, email, password, phone string, status int32) *User {
	return &User{
		ID:         id,
		Username:   username,
		Firstname:  firstname,
		Lastname:   lastname,
		Email:      email,
		Password:   password,
		Phone:      phone,
		UserStatus: status,
	}
}
