package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Pet represents a pet in the store
type Pet struct {
	Id     int64  `json:"id" yaml:"id" validate:"required"`
	Name   string `json:"name" yaml:"name" validate:"required"`
	Tag    string `json:"tag,omitempty" yaml:"tag,omitempty"`
	Status string `json:"status,omitempty" yaml:"status,omitempty"`
}

// Server is the API server
type Server struct {
	router *gin.Engine
	pets   []Pet // In-memory storage for pets
}

// NewServer creates a new API server
func NewServer() *Server {
	s := &Server{
		router: gin.Default(),
		pets: []Pet{
			{Id: 1, Name: "Max", Tag: "dog", Status: "available"},
			{Id: 2, Name: "Whiskers", Tag: "cat", Status: "available"},
		},
	}
	s.setupRoutes()
	return s
}

// setupRoutes configures the routes for the API
func (s *Server) setupRoutes() {
	// Health check endpoint
	s.router.GET("/healthCheck", s.HandleHealthCheck)

	// Routes for /pets
	s.router.GET("/pets", s.HandleListPets)
	s.router.POST("/pets", s.HandleCreatePet)

	// Routes for /pets/:petId
	s.router.GET("/pets/:petId", s.HandleGetPet)
}

// Start starts the server on the specified address
func (s *Server) Start(addr string) error {
	log.Printf("Server starting on %s", addr)
	return s.router.Run(addr)
}

// HandleHealthCheck Health check endpoint
func (s *Server) HandleHealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

// HandleListPets List all pets
func (s *Server) HandleListPets(c *gin.Context) {
	log.Println("Handling GET /pets request")
	c.JSON(http.StatusOK, s.pets)
}

// HandleCreatePet Create a pet
func (s *Server) HandleCreatePet(c *gin.Context) {
	log.Println("Handling POST /pets request")
	
	// Parse request body
	var pet Pet
	if err := c.ShouldBindJSON(&pet); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	
	// Assign a new ID (simple implementation)
	pet.Id = int64(len(s.pets) + 1)
	
	// Add to pets
	s.pets = append(s.pets, pet)
	
	// Return success status
	c.Status(http.StatusCreated)
}

// HandleGetPet Get a pet by ID
func (s *Server) HandleGetPet(c *gin.Context) {
	log.Println("Handling GET /pets/:petId request")
	
	// Get path parameters
	petIdStr := c.Param("petId")
	petId, err := strconv.Atoi(petIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid petId"})
		return
	}
	
	// Find the pet
	for _, pet := range s.pets {
		if pet.Id == int64(petId) {
			c.JSON(http.StatusOK, pet)
			return
		}
	}
	
	// Pet not found
	c.JSON(http.StatusNotFound, gin.H{"error": "Pet not found"})
}

func main() {
	// Start the server on port 8080
	log.Println("Starting Petstore API server on :8080")
	server := NewServer()
	if err := server.Start(":8080"); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
