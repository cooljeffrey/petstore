package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewPet(t *testing.T) {
	p := NewPet(
		1,
		&Category{ID: 1, Name: "cat"},
		"cat 1",
		[]string{"http://localhost/image/1.jpg"},
		[]*Tag{{ID: 1, Name: "bset selling"}},
		PetStatusAvailable)
	assert.NotNil(t, p)
	assert.Equal(t, int64(1), p.ID)
	assert.Equal(t, "cat", p.Category.Name)
	assert.Equal(t, "http://localhost/image/1.jpg", p.PhotoUrls[0])
	assert.Equal(t, PetStatusAvailable, p.Status)
	assert.Equal(t, int64(1), p.Tags[0].ID)
}
