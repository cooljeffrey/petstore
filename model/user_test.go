package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewUser(t *testing.T) {
	user := NewUser(
		10,
		"username",
		"firstname",
		"lastname",
		"email@non.email",
		"password",
		"1234567",
		0)

	assert.NotNil(t, user)
	assert.Equal(t, int64(10), user.ID)
	assert.Equal(t, "username", user.Username)
	assert.Equal(t, "firstname", user.Firstname)
	assert.Equal(t, "lastname", user.Lastname)
	assert.Equal(t, "email@non.email", user.Email)
	assert.Equal(t, "password", user.Password)
	assert.Equal(t, "1234567", user.Phone)
	assert.Equal(t, int32(0), user.UserStatus)
}
