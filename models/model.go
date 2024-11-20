package models

import (
	"github.com/google/uuid"
)

type Item struct {
	ID          uuid.UUID `json:"id" gorm:"<-:create"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
}

type Items []*Item
