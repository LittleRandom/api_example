package models

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ItemService struct {
	Repository *Repository
}

func NewItemService(db *gorm.DB) *ItemService {
	return &ItemService{
		Repository: NewRepository(db),
	}
}

func (s *ItemService) RegisterRoutes(r chi.Router) {
	r.Get("/", s.HandleGetRoot)
	r.Get("/{id}", s.HandleGetItem)
	r.Post("/", s.HandleImportItem)
	r.Delete("/{id}", s.HandleDeleteItem)
}

func (s *ItemService) HandleGetRoot(w http.ResponseWriter, r *http.Request) {

	books, err := s.Repository.List()

	if err != nil {
		log.Printf("error reading rows %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	j, err := json.Marshal(books)
	if err != nil {
		log.Printf("error marshalling books into json %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(j)
}

func (s *ItemService) HandleGetItem(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	UUID, err := uuid.Parse(id)
	if err != nil {
		log.Printf("error parsing uuid from url: %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	item, err := s.Repository.Read(UUID)
	if err != nil {
		log.Printf("error finding book in db: %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	j, err := json.Marshal(item)
	if err != nil {
		log.Printf("error marshalling books into json %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

func (s *ItemService) HandleImportItem(w http.ResponseWriter, r *http.Request) {
	var item Item

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&item)
	if err != nil {
		log.Printf("Error with decoding json from body: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	item.ID, err = uuid.NewUUID()
	if err != nil {
		log.Printf("Error when assigning new UUID in database: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	retItem, err := s.Repository.Create(&item)
	if err != nil {
		log.Printf("Error when creating item in database: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	j, err := json.Marshal(retItem)
	if err != nil {
		log.Printf("Error response from database: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(j)
}

func (s *ItemService) HandleDeleteItem(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")
	itemUUID, err := uuid.Parse(id)
	if err != nil {
		log.Printf("Error when trying to parse UUID: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = s.Repository.Delete(itemUUID)
	if err != nil {
		log.Printf("Error when trying to delete from table: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
