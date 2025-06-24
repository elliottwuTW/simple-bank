package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/simple_bank/database"
)

// Server serves HTTP requests for our service.
type Server struct {
	db     database.Database
	router *gin.Engine
}

func NewServer(db database.Database) *Server {
	server := &Server{db: db}
	router := gin.Default()

	// 取得 Gin 背後的 validator engine
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	router.POST("/accounts", server.createAccount)
	router.POST("/users", server.createUser)

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
