package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		DB: db,
	}
}

func (r *Repository) List() (Items, error) {
	items := make([]*Item, 0)
	if err := r.DB.Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

func (r *Repository) Create(item *Item) (*Item, error) {
	if err := r.DB.Create(item).Error; err != nil {
		return nil, err
	}

	return item, nil
}

func (r *Repository) Read(id uuid.UUID) (*Item, error) {
	item := &Item{}
	if err := r.DB.Where("id = ?", id).First(&item).Error; err != nil {
		return nil, err
	}

	return item, nil
}

func (r *Repository) Delete(id uuid.UUID) (int64, error) {
	result := r.DB.Where("id = ?", id).Delete(&Item{})

	return result.RowsAffected, result.Error
}
