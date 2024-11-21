package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"plainrandom/config"
	"plainrandom/models"
	"plainrandom/sqlite"
	"strings"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"gorm.io/gorm"
)

type Server struct {
	Config *config.Config

	DB     *gorm.DB
	Server *http.Server

	ItemService *models.ItemService
}

// Stronger object to better organize functionality
func NewServer(c *config.Config) *Server {
	// Create Config and Database
	m := &Server{
		Config: config.NewConfig(),
		DB:     sqlite.NewDB(),
	}

	// Create the http.Server to serve requests
	m.Server = &http.Server{
		Addr: fmt.Sprintf("%s:%d", c.Host, c.Port),
	}

	// Load API services
	m.ItemService = models.NewItemService(m.DB)

	// Establish the routes
	m.NewRouter()

	return m
}

// Start function for server
func (m *Server) Start() {
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

// Stop function to allow graceful shutdown.
func (m *Server) Stop(ctx context.Context) error {
	// API Server to host endpoints.
	log.Printf("Stopping API Server")

	return m.Server.Shutdown(ctx)

}

// Loads Router into http.Server
func (m *Server) NewRouter() {
	// Creates a chi.Router for the server
	r := chi.NewRouter()

	// Include CORS handler to fix CORS error in API calls.
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Use(middleware.Logger)

	// Swagger API on root link
	r.Route("/", RegisterFileServer("docs/static"))

	// OpenAPI file that holds API documentation
	r.Route("/api/v1", RegisterFileServer("docs/api/v1"))

	// Items endpoint
	r.Route("/items", m.ItemService.RegisterRoutes)

	m.Server.Handler = r
}

// Register FileServer with path as argument
func RegisterFileServer(path string) func(r chi.Router) {
	// Returns a function that works with r.Route()
	return func(r chi.Router) {
		r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
			workDir, _ := os.Getwd()
			filesDir := http.Dir(filepath.Join(workDir, path))
			rctx := chi.RouteContext(r.Context())
			pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
			fs := http.StripPrefix(pathPrefix, http.FileServer(filesDir))
			fs.ServeHTTP(w, r)
		})
	}
}
