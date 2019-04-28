package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/cooljeffrey/petstore/model"
	"github.com/go-kit/kit/log"
	"io/ioutil"
	"strings"
	"time"
)

type PetService interface {
	AddPet(ctx context.Context, pet *model.Pet) error
	UpdatePet(ctx context.Context, pet *model.Pet) error
	FindPetsByStatus(ctx context.Context, statuses []string) ([]*model.Pet, error)
	FindPetByID(ctx context.Context, id int64) (*model.Pet, error)
	UpdatePetByID(ctx context.Context, id int64, name string, status string) error
	AddImageUrlForPetByID(ctx context.Context, id int64, filename string, file []byte) error
	DeletePetByID(ctx context.Context, id int64) error
}

type petService struct {
	logger   log.Logger
	storage  model.Storage
	baseUri  string
	filePath string
}

func NewPetService(logger log.Logger, storage model.Storage, baseUri, filePath string) PetService {
	return &petService{
		logger:   logger,
		storage:  storage,
		baseUri:  baseUri,
		filePath: filePath,
	}
}

func (s petService) AddPet(ctx context.Context, pet *model.Pet) error {
	return s.storage.CreatePet(pet)
}

func (s petService) UpdatePet(ctx context.Context, pet *model.Pet) error {
	return s.storage.UpdatePetByID(pet)
}

func (s petService) FindPetsByStatus(ctx context.Context, statuses []string) ([]*model.Pet, error) {
	return s.storage.FindPetsByStatus(statuses)
}

func (s petService) FindPetByID(ctx context.Context, id int64) (*model.Pet, error) {
	return s.storage.RetrievePetByID(id)
}

func (s petService) UpdatePetByID(ctx context.Context, id int64, name, status string) error {
	if name == "" && status != "" {
		return s.storage.UpdatePetStatusByID(id, status)
	}
	if name != "" && status == "" {
		return s.storage.UpdatePetNameByID(id, name)
	}
	if name != "" && status != "" {
		return s.storage.UpdatePetNameAndStatusByID(id, name, status)
	}
	return errors.New("both name and status are empty")
}

func (s petService) AddImageUrlForPetByID(ctx context.Context, id int64, filename string, file []byte) error {
	newfilename := fmt.Sprintf(
		"%d.%s",
		time.Now().UTC().UnixNano(),
		strings.Split(filename, ".")[1])
	err := ioutil.WriteFile(fmt.Sprintf("%s/%s", s.filePath, newfilename), file, 0644)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s/%s", s.baseUri, newfilename)
	_, err = s.storage.AddImageUrlByPetID(id, url)
	return err
}

func (s petService) DeletePetByID(ctx context.Context, id int64) error {
	return s.storage.DeletePetByID(id)
}
