package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewErrResponse(t *testing.T) {
	e := NewErrResponse(200, "type", "msg")
	assert.NotNil(t, e)
	assert.Equal(t, "{\"code\":200,\"type\":\"type\",\"message\":\"msg\"}", e.Error())
}
