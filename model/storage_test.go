package model

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

var logger = log.NewJSONLogger(os.Stderr)
var storage Storage

func TestNewMongoStorage(t *testing.T) {
	db, err := NewMongoStorage("mongodb://127.0.0.1:27017", "petstore", 10, logger)
	assert.NoError(t, err)
	assert.NotNil(t, db)
	storage = db

	assert.NoError(t, storage.EmptyCollection(CollectionUsers))
	assert.NoError(t, storage.EmptyCollection(CollectionPets))
	assert.NoError(t, storage.EmptyCollection(CollectionOrders))
}

func TestMongoStorageUserActions(t *testing.T) {
	err := storage.CreateUser(&User{
		ID:         1,
		Username:   "username",
		Firstname:  "firstname",
		Lastname:   "lastname",
		Email:      "email@email.com",
		Password:   "password",
		Phone:      "123",
		UserStatus: 1,
	})
	assert.NoError(t, err)

	u, err := storage.RetrieveUserByUsername("username")
	assert.NoError(t, err)
	assert.NotNil(t, u)
	assert.Equal(t, int64(1), u.ID)

	u, err = storage.RetrieveUserByID(1)
	assert.NoError(t, err)
	assert.NotNil(t, u)
	assert.Equal(t, int64(1), u.ID)
	assert.Equal(t, "username", u.Username)

	u.UserStatus = 0
	u.Phone = "456"
	u.Firstname = "a"
	u.Lastname = "b"
	u.Email = "test@test.com"
	u.ID = 10
	u.Username = "username-1"

	u2, err := storage.UpdateUserByUsername("username", u)
	assert.NoError(t, err)
	assert.NotNil(t, u2)
	assert.Equal(t, "username-1", u2.Username)
	assert.Equal(t, int64(10), u2.ID)
	assert.Equal(t, "test@test.com", u2.Email)
	assert.Equal(t, "b", u2.Lastname)
	assert.Equal(t, "a", u2.Firstname)
	assert.Equal(t, "456", u2.Phone)
	assert.Equal(t, int32(0), u2.UserStatus)

	err = storage.DeleteUserByUsername("username-1")
	assert.NoError(t, err)
}

func TestMongoStoragePetActions(t *testing.T) {
	pet := Pet{
		ID:        1,
		Name:      "cat1",
		Status:    PetStatusAvailable,
		Category:  NewCategory(1, "cat"),
		PhotoUrls: []string{},
	}

	err := storage.CreatePet(&pet)
	assert.NoError(t, err)

	pet.Name = "cat2"
	pet.Status = PetStatusPending
	err = storage.UpdatePetByID(&pet)
	assert.NoError(t, err)

	assert.NoError(t, storage.UpdatePetNameByID(1, "cat3"))
	assert.NoError(t, storage.UpdatePetStatusByID(1, PetStatusSold))
	assert.NoError(t, storage.UpdatePetNameAndStatusByID(1, "cat4", PetStatusAvailable))

	url := "http://localhost:8080/images/1.jpg"
	p, err := storage.AddImageUrlByPetID(1, url)
	assert.NoError(t, err)
	assert.NotNil(t, p)
	assert.Equal(t, []string{url}, p.PhotoUrls)

	var ps []*Pet
	for _, id := range []int64{2, 3, 4, 5} {
		pet := Pet{
			ID:        id,
			Name:      fmt.Sprintf("cat%d", id),
			Status:    PetStatusAvailable,
			Category:  NewCategory(1, "cat"),
			PhotoUrls: []string{},
		}
		ps = append(ps, &pet)
	}
	assert.NoError(t, storage.CreateManyPets(ps))

	pets, err := storage.FindPetsByStatus([]string{PetStatusAvailable, PetStatusPending})
	assert.NoError(t, err)
	assert.NotNil(t, pets)
	assert.True(t, len(pets) > 1)

	inv, err := storage.RetrieveStoreInventoriesByStatus()
	assert.NoError(t, err)
	assert.NotNil(t, inv)
	assert.Equal(t, 1, len(inv))
	assert.Equal(t, int64(5), inv[PetStatusAvailable])
}

func TestMongoStorageStoreActions(t *testing.T) {
	order := Order{
		ID:       1,
		PetID:    1,
		Quantity: int32(2),
		ShipDate: time.Now().UTC(),
		Status:   OrderStatusPlaced,
		Complete: false,
	}
	o, err := storage.CreateOrder(&order)
	assert.NoError(t, err)
	assert.NotNil(t, o)

	o, err = storage.RetrieveOrderByID(1)
	assert.NoError(t, err)
	assert.NotNil(t, o)
	assert.ObjectsAreEqualValues(o, order)

	assert.NoError(t, storage.DeleteOrderByID(1))
	o, err = storage.RetrieveOrderByID(1)
	assert.Error(t, err)
	assert.Nil(t, o)
}

func TestCleanUp(t *testing.T) {
	assert.NoError(t, storage.EmptyCollection(CollectionUsers))
	assert.NoError(t, storage.EmptyCollection(CollectionPets))
	assert.NoError(t, storage.EmptyCollection(CollectionOrders))
}
