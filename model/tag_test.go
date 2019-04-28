package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewTag(t *testing.T) {
	tag := NewTag(1, "best")
	assert.NotNil(t, tag)
	assert.Equal(t, int64(1), tag.ID)
	assert.Equal(t, "best", tag.Name)
}
