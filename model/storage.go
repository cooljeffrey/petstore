package model

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-kit/kit/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

const (
	CollectionPets       string = "pets"
	CollectionUsers      string = "users"
	CollectionOrders     string = "orders"
	CollectionCategories string = "categories"
	CollectionTags       string = "tags"
)

type Storage interface {
	// Create user
	CreateUser(user *User) error
	// Create users from slice
	CreateManyUsers(users []*User) error
	// Fetch user by username
	RetrieveUserByUsername(username string) (*User, error)
	// Fetch user by user id
	RetrieveUserByID(id int64) (*User, error)
	// Update user by username
	UpdateUserByUsername(username string, user *User) (*User, error)
	// Delete user by username
	DeleteUserByUsername(username string) error

	// Create pet
	CreatePet(pet *Pet) error
	// Create pets from slice
	CreateManyPets(pets []*Pet) error
	// Update pet by pet id
	UpdatePetByID(pet *Pet) error
	// Fetch pet by id
	RetrievePetByID(id int64) (*Pet, error)
	// Find pets by given statuses slice
	FindPetsByStatus(statuses []string) ([]*Pet, error)
	// Update pet naem and status by given pet id
	UpdatePetNameAndStatusByID(id int64, name string, status string) error
	// Update pet name by given id
	UpdatePetNameByID(id int64, name string) error
	// Update pet status by given id
	UpdatePetStatusByID(id int64, status string) error
	// Add image url to pet by give pet id
	AddImageUrlByPetID(id int64, url string) (*Pet, error)
	// Fetch store inventory of all statuses
	RetrieveStoreInventoriesByStatus() (map[string]int64, error)
	// Delete pet by given ID
	DeletePetByID(id int64) error

	// Create order
	CreateOrder(order *Order) (*Order, error)
	// Fetch order by given order id
	RetrieveOrderByID(id int64) (*Order, error)
	// Delete order by given order id
	DeleteOrderByID(id int64) error

	// Drop whole specified collection
	EmptyCollection(collection string) error
}

type MongoStorage struct {
	client   *mongo.Client
	URI      string
	Timeout  int64
	Database string
	Logger   log.Logger
}

func NewMongoStorage(uri, database string, timeout int64, logger log.Logger) (Storage, error) {
	storage := MongoStorage{
		URI:      uri,
		Database: database,
		Timeout:  timeout,
		Logger:   logger,
	}
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}
	storage.client = client
	return storage, nil
}

func toBsonD(val interface{}) (*bson.D, error) {
	b, err := bson.Marshal(val)
	if err != nil {
		return nil, err
	}
	d := bson.D{}
	err = bson.Unmarshal(b, &d)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (m MongoStorage) CreateUser(user *User) error {
	collection := m.client.Database(m.Database).Collection(CollectionUsers)
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(m.Timeout)*time.Second)

	u, err := m.RetrieveUserByID(user.ID)
	if err == nil && u != nil {
		return errors.New("duplicate user id exists")
	}

	u, err = m.RetrieveUserByUsername(user.Username)
	if err == nil && u != nil {
		return errors.New("duplicate username exists")
	}

	d, err := toBsonD(user)
	if err != nil {
		return err
	}
	_, err = collection.InsertOne(ctx, d)
	if err != nil {
		return err
	}
	return nil
}

func (m MongoStorage) CreateManyUsers(users []*User) error {
	collection := m.client.Database(m.Database).Collection(CollectionUsers)
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(m.Timeout)*time.Second)

	var docs []interface{}
	for _, u := range users {
		u2, err := m.RetrieveUserByID(u.ID)
		if err == nil && u2 != nil {
			return fmt.Errorf("duplicate user id exists for %d", u.ID)
		}
		u2, err = m.RetrieveUserByUsername(u.Username)
		if err == nil && u2 != nil {
			return fmt.Errorf("duplicate username exists for %s", u.Username)
		}
		d, err := toBsonD(u)
		if err != nil {
			return err
		}
		docs = append(docs, d)
	}

	_, err := collection.InsertMany(ctx, docs)
	if err != nil {
		return err
	}
	return nil
}

func (m MongoStorage) RetrieveUserByUsername(username string) (*User, error) {
	collection := m.client.Database(m.Database).Collection(CollectionUsers)
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(m.Timeout)*time.Second)

	var user User
	err := collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (m MongoStorage) RetrieveUserByID(id int64) (*User, error) {
	collection := m.client.Database(m.Database).Collection(CollectionUsers)
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(m.Timeout)*time.Second)

	var user User
	err := collection.FindOne(ctx, bson.M{"id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (m MongoStorage) UpdateUserByUsername(username string, user *User) (*User, error) {
	collection := m.client.Database(m.Database).Collection(CollectionUsers)
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(m.Timeout)*time.Second)

	d, err := toBsonD(user)
	if err != nil {
		return nil, err
	}
	err = collection.FindOneAndReplace(ctx, bson.M{"username": username}, d).Err()
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (m MongoStorage) DeleteUserByUsername(username string) error {
	collection := m.client.Database(m.Database).Collection(CollectionUsers)
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(m.Timeout)*time.Second)

	return collection.FindOneAndDelete(ctx, bson.M{"username": username}).Err()
}

func (m MongoStorage) RetrieveStoreInventoriesByStatus() (map[string]int64, error) {
	collection := m.client.Database(m.Database).Collection(CollectionPets)
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(m.Timeout)*time.Second)
	cur, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	defer cur.Close(context.Background())

	inv := map[string]int64{}
	for cur.Next(context.Background()) {
		pet := Pet{}
		err = cur.Decode(&pet)
		if err != nil {
			return nil, err
		}

		i, ok := inv[pet.Status]
		if ok {
			inv[pet.Status] = int64(i + 1)
		} else {
			inv[pet.Status] = int64(1)
		}
	}
	return inv, nil
}

func (m MongoStorage) CreateOrder(order *Order) (*Order, error) {
	collection := m.client.Database(m.Database).Collection(CollectionOrders)
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(m.Timeout)*time.Second)

	o, err := m.RetrieveOrderByID(order.ID)
	if err == nil && o != nil {
		return nil, errors.New("duplicate order id exists")
	}

	d, err := toBsonD(order)
	if err != nil {
		return nil, err
	}
	_, err = collection.InsertOne(ctx, d)

	return order, nil
}

func (m MongoStorage) RetrieveOrderByID(id int64) (*Order, error) {
	collection := m.client.Database(m.Database).Collection(CollectionOrders)
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(m.Timeout)*time.Second)

	var order Order
	err := collection.FindOne(ctx, bson.M{"id": id}).Decode(&order)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (m MongoStorage) DeleteOrderByID(id int64) error {
	collection := m.client.Database(m.Database).Collection(CollectionOrders)
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(m.Timeout)*time.Second)

	return collection.FindOneAndDelete(ctx, bson.M{"id": id}).Err()
}

func (m MongoStorage) CreatePet(pet *Pet) error {
	collection := m.client.Database(m.Database).Collection(CollectionPets)
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(m.Timeout)*time.Second)

	u, err := m.RetrievePetByID(pet.ID)
	if err == nil && u != nil {
		return errors.New("duplicate pet id exists")
	}

	d, err := toBsonD(pet)
	if err != nil {
		return err
	}
	_, err = collection.InsertOne(ctx, d)

	return nil
}

func (m MongoStorage) CreateManyPets(pets []*Pet) error {
	collection := m.client.Database(m.Database).Collection(CollectionPets)
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(m.Timeout)*time.Second)

	var docs []interface{}
	for _, p := range pets {
		u, err := m.RetrievePetByID(p.ID)
		if err == nil && u != nil {
			return fmt.Errorf("duplicate pet id exists for %d", p.ID)
		}

		d, err := toBsonD(p)
		if err != nil {
			return err
		}
		docs = append(docs, d)
	}

	_, err := collection.InsertMany(ctx, docs)
	if err != nil {
		return err
	}
	return nil
}

func (m MongoStorage) UpdatePetByID(pet *Pet) error {
	collection := m.client.Database(m.Database).Collection(CollectionPets)
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(m.Timeout)*time.Second)

	d, err := toBsonD(pet)
	if err != nil {
		return err
	}
	return collection.FindOneAndReplace(ctx, bson.M{"id": pet.ID}, d).Err()
}

func (m MongoStorage) RetrievePetByID(id int64) (*Pet, error) {
	collection := m.client.Database(m.Database).Collection(CollectionPets)
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(m.Timeout)*time.Second)

	var pet Pet
	err := collection.FindOne(ctx, bson.M{"id": id}).Decode(&pet)
	if err != nil {
		return nil, err
	}
	return &pet, nil
}

func (m MongoStorage) FindPetsByStatus(statuses []string) ([]*Pet, error) {
	collection := m.client.Database(m.Database).Collection(CollectionPets)
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(m.Timeout)*time.Second)

	var A bson.A
	for _, s := range statuses {
		A = append(A, s)
	}

	var pets []*Pet
	cur, err := collection.Find(ctx, bson.M{"status": bson.M{"$in": A}})
	if err != nil {
		return nil, err
	}

	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		var pet Pet
		err = cur.Decode(&pet)
		if err != nil {
			return nil, err
		}
		pets = append(pets, &pet)
	}
	return pets, nil
}

func (m MongoStorage) UpdatePetNameAndStatusByID(id int64, name string, status string) error {
	collection := m.client.Database(m.Database).Collection(CollectionPets)
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(m.Timeout)*time.Second)

	return collection.FindOneAndUpdate(ctx, bson.M{"id": id}, bson.D{
		{"$set", bson.D{
			{"name", name},
			{"status", status},
		}},
	}).Err()
}

func (m MongoStorage) UpdatePetNameByID(id int64, name string) error {
	collection := m.client.Database(m.Database).Collection(CollectionPets)
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(m.Timeout)*time.Second)

	return collection.FindOneAndUpdate(ctx, bson.M{"id": id}, bson.D{
		{"$set", bson.D{
			{"name", name},
		}},
	}).Err()
}

func (m MongoStorage) UpdatePetStatusByID(id int64, status string) error {
	collection := m.client.Database(m.Database).Collection(CollectionPets)
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(m.Timeout)*time.Second)

	return collection.FindOneAndUpdate(ctx, bson.M{"id": id}, bson.D{
		{"$set", bson.D{
			{"status", status},
		}},
	}).Err()
}

func (m MongoStorage) AddImageUrlByPetID(id int64, url string) (*Pet, error) {
	collection := m.client.Database(m.Database).Collection(CollectionPets)
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(m.Timeout)*time.Second)
	err := collection.FindOneAndUpdate(ctx, bson.M{"id": id}, bson.D{
		{"$addToSet", bson.M{"photoUrls": url}}}).Err()
	if err != nil {
		return nil, err
	}
	return m.RetrievePetByID(id)
}

func (m MongoStorage) DeletePetByID(id int64) error {
	collection := m.client.Database(m.Database).Collection(CollectionPets)
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(m.Timeout)*time.Second)

	return collection.FindOneAndDelete(ctx, bson.M{"id": id}).Err()
}

func (m MongoStorage) EmptyCollection(collection string) error {
	coll := m.client.Database(m.Database).Collection(collection)
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(m.Timeout)*time.Second)
	return coll.Drop(ctx)
}
