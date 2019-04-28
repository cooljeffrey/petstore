package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewOrder(t *testing.T) {
	ts := time.Now().UTC()
	order := NewOrder(1, 1, 2, ts, OrderStatusPlaced, false)
	assert.NotNil(t, order)
	assert.Equal(t, int64(1), order.ID)
	assert.Equal(t, int64(1), order.PetID)
	assert.Equal(t, int32(2), order.Quantity)
	assert.Equal(t, ts, order.ShipDate)
	assert.Equal(t, OrderStatusPlaced, order.Status)
	assert.Equal(t, false, order.Complete)
}
