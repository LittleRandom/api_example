package models_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"plainrandom/models"
	"plainrandom/server"
	"plainrandom/sqlite"
	"testing"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"

	"github.com/google/uuid"
)

func MustOpenDatabase(t *testing.T) *gorm.DB {
	err := os.MkdirAll("./testing_data", os.ModePerm)
	if err != nil {
		log.Printf("error when creating a directory: %v", err)
		panic(err)
	}
	// Use stdlib to open a connection to postgres db.
	path := filepath.Join("./testing_data", "api.db")
	DB, err := sqlite.OpenDatabase(path)
	if err != nil {
		log.Printf("error when connecting to database %v", err)
		panic(err)
	}

	t.Cleanup(func() {
		os.RemoveAll("./testing_data")

		if DB != nil {
			db, err := DB.DB()
			if err != nil {
				db.Close()
			}
		}
	})

	return DB

}

func MustOpenServer(t *testing.T) *server.Server {
	DB := MustOpenDatabase(t)
	s := &server.Server{

		DB: DB,

		Server: &http.Server{
			Addr: "https://api.plainrandom.com",
			// Addr: "localhost:5060",
		},
	}

	s.ItemService = models.NewItemService(DB)

	s.NewRouter()

	t.Cleanup(func() {
		s.Stop(context.Background())
	})

	return s

}

func TestItemService_HandleGetRoot(t *testing.T) {
	s := MustOpenServer(t)
	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	s.ItemService.HandleGetRoot(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %v, got %v", http.StatusOK, w.Code)
	}
}

func TestItemService_HandleGetItem(t *testing.T) {
	s := MustOpenServer(t)

	// Create a new item in database
	newUUID, _ := uuid.NewUUID()
	item, err := s.ItemService.Repository.Create(&models.Item{
		ID:          newUUID,
		Title:       "HandleGetItemTest",
		Description: "Testing Handle Get Item",
	})
	if err != nil {
		t.Errorf("Error when creating an item in database: %v", err)
	}

	// Create new Delete request with context
	r := httptest.NewRequest("GET", "/{id}", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", newUUID.String())
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	// Run handler
	s.ItemService.HandleGetItem(w, r)

	// Check response status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %v, got %v", http.StatusOK, w.Code)
	}

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["id"] != newUUID.String() {
		t.Errorf("Expected %v, got %v", newUUID.String(), response["id"])
	}
	if response["title"] != item.Title {
		t.Errorf("Expected %v, got %v", item.Title, response["title"])
	}
	log.Printf("🌮 ~ %v", response["title"])
	if response["description"] != item.Description {
		t.Errorf("Expected %v, got %v", item.Description, response["description"])
	}

}

func makePostRequest(item models.Item) (*http.Request, error) {
	body, err := json.Marshal(item)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, "localhost:5060", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	return req, nil
}

func TestItemService_HandleImportItem(t *testing.T) {
	s := MustOpenServer(t)

	// Create POST request body
	item := models.Item{
		Title:       "Test Title",
		Description: "Some Lorem Ipsum",
	}
	r, err := makePostRequest(item)
	if err != nil {
		t.Error(err)
	}
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Accept", "application/json")
	w := httptest.NewRecorder()

	// Run Handler
	s.ItemService.HandleImportItem(w, r)

	// Check response status code
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %v, got %v", http.StatusCreated, w.Code)
	}

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["id"] == uuid.Nil.String() {
		t.Errorf("Id got %v", response["id"])
	}
	if response["title"] != item.Title {
		t.Errorf("Expected %v, got %v", item.Title, response["title"])
	}
	if response["description"] != item.Description {
		t.Errorf("Expected %v, got %v", item.Description, response["description"])
	}
}

func TestItemService_HandleDeleteItem(t *testing.T) {
	s := MustOpenServer(t)

	// Create a new item in database
	newUUID, _ := uuid.NewUUID()
	_, err := s.ItemService.Repository.Create(&models.Item{
		ID:          newUUID,
		Title:       "Test Title",
		Description: "Some Lorem Ipsum",
	})
	if err != nil {
		t.Errorf("Error when creating an item in database: %v", err)
	}

	// Create new Delete request with context
	r := httptest.NewRequest("DELETE", "/{id}", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", newUUID.String())
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	// Run handler
	s.ItemService.HandleDeleteItem(w, r)

	// Check response status code
	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status code %v, got %v", http.StatusNoContent, w.Code)
	}

}
