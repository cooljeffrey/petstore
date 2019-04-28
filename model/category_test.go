package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewCategory(t *testing.T) {
	c := NewCategory(1, "cat")
	assert.NotNil(t, c)
	assert.Equal(t, int64(1), c.ID)
	assert.Equal(t, "cat", c.Name)
}
