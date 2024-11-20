package server

import (
	"context"
	"log"
	"net/http"
	"plainrandom/models"
	"plainrandom/sqlite"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type Server struct {
	DB     *gorm.DB
	Server *http.Server

	ItemService *models.ItemService
}

// Stronger object to better organize functionality
func NewServer() *Server {
	m := &Server{

		DB: sqlite.NewDB(),

		Server: &http.Server{
			Addr: "localhost:5050",
		},
	}

	m.ItemService = models.NewItemService(m.DB)

	m.NewRouter()

	return m
}

func (m *Server) Start() {
	// API Server to host endpoints.
	log.Printf("Starting API Server")

	// Start with connecting to the Database
	go func() {

		log.Printf("Listening at: http://%v", m.Server.Addr)
		// always returns error. ErrServerClosed on graceful close
		if err := m.Server.ListenAndServe(); err != http.ErrServerClosed {
			// unexpected error. port in use?
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()
}

func (m *Server) Stop(ctx context.Context) error {
	// API Server to host endpoints.
	log.Printf("Stopping API Server")

	return m.Server.Shutdown(ctx)

}

func (m *Server) NewRouter() {
	// Creates a chi.Router for the server
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Route("/", m.ItemService.RegisterRoutes)

	m.Server.Handler = r
}
