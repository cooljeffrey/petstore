package model

type Pet struct {
	ID        int64     `json:"id" bson:"id"`
	Category  *Category `json:"category,omitempty" bson:"category,omitempty"`
	Name      string    `json:"name" bson:"name"`
	PhotoUrls []string  `json:"photoUrls" bson:"photoUrls,omitempty"`
	Tags      []*Tag    `json:"tags" bson:"tags,omitempty"`
	Status    string    `json:"status" bson:"status"`
}

const (
	PetStatusAvailable string = "available"
	PetStatusPending   string = "pending"
	PetStatusSold      string = "sold"
)

func NewPet(id int64, category *Category, name string, photoUrls []string, tags []*Tag, status string) *Pet {
	return &Pet{
		ID:        id,
		Category:  category,
		Name:      name,
		PhotoUrls: photoUrls,
		Tags:      tags,
		Status:    status,
	}
}

func (p *Pet) AddPhotoUrl(url string) {
	if p.PhotoUrls == nil {
		p.PhotoUrls = []string{}
	}
	p.PhotoUrls = append(p.PhotoUrls, url)
}

func (p *Pet) RemovePhotoUrl(url string) {
	if p.PhotoUrls == nil || len(p.PhotoUrls) == 0 {
		return
	}

	for i, u := range p.PhotoUrls {
		if u == url {
			p.PhotoUrls = append(p.PhotoUrls[:i], p.PhotoUrls[i+1:]...)
		}
	}
}

func (p *Pet) AddTag(tag *Tag) {
	if p.Tags == nil {
		p.Tags = []*Tag{}
	}

	p.Tags = append(p.Tags, tag)
}

func (p *Pet) RemoveTag(tag *Tag) {
	if p.Tags == nil || len(p.Tags) == 0 {
		return
	}

	for i, t := range p.Tags {
		if t.ID == tag.ID {
			p.Tags = append(p.Tags[:i], p.Tags[i+1:]...)
		}
	}
}

type UploadImageResult struct {
	Code    int32  `json:"code"`
	Type    string `json:"type"`
	Message string `json:"message"`
}
