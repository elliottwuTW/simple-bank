package api

import (
	"github.com/gin-gonic/gin"
	"github.com/simple_bank/database"
)

// Server serves HTTP requests for our service.
type Server struct {
	db     *database.Database
	router *gin.Engine
}

func NewServer(db *database.Database) *Server {
	server := &Server{db: db}
	router := gin.Default()

	router.POST("/accounts", server.createAccount)

	server.router = router
	return server
}

// Start runs the HTTP server on specific address.
func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
