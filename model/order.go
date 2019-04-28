package model

import (
	"reflect"
	"time"
)

type Order struct {
	ID       int64     `json:"id" bson:"id"`
	PetID    int64     `json:"petId" bson:"petId"`
	Quantity int32     `json:"quantity" bson:"quantity"`
	ShipDate time.Time `json:"shipDate" bson:"shipDate"`
	Status   string    `json:"status" bson:"status"`
	Complete bool      `json:"complete" bson:"complete"`
}

const (
	OrderStatusPlaced    string = "placed"
	OrderStatusApproved  string = "approved"
	OrderStatusDelivered string = "delivered"
)

func NewOrder(id, petId int64, quantity int32, shipDate time.Time, status string, complete bool) *Order {
	return &Order{
		ID:       id,
		PetID:    petId,
		Quantity: quantity,
		ShipDate: shipDate,
		Status:   status,
		Complete: complete,
	}
}

type Inventory map[string]int64

func (inv *Inventory) Count() int64 {
	return int64(len(reflect.ValueOf(*inv).MapKeys()))
}
