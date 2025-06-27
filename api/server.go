package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/simple_bank/config"
	"github.com/simple_bank/database"
	"github.com/simple_bank/token"
	"github.com/simple_bank/worker"
)

// Server serves HTTP requests for our service.
type Server struct {
	db          database.Database
	config      config.Config
	tokenMaker  token.Maker
	distributor worker.TaskDistributor
	router      *gin.Engine
}

func NewServer(db database.Database, distbtr worker.TaskDistributor, config config.Config) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.Token.SymmetricKey)
	if err != nil {
		return nil, err
	}

	server := &Server{db: db, distributor: distbtr, tokenMaker: tokenMaker}

	// 取得 Gin 背後的 validator engine
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.setupRouter()
	return server, nil
}

func (s *Server) setupRouter() {
	router := gin.Default()

	router.POST("/users", s.createUser)
	router.POST("/users/login", s.loginUser)

	// 以下 api 需要 auth
	authRoutes := router.Group("/").Use(authMiddleware(s.tokenMaker))

	authRoutes.POST("/accounts", s.createAccount)

	s.router = router
}

// Start runs the HTTP server on specific address.
func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
