package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/simple_bank/database"
)

type CreateAccountReq struct {
	// https://pkg.go.dev/github.com/go-playground/validator/v10#section-readme
	Owner    string `json:"owner"    binding:"required"`
	Currency string `json:"currency" binding:"required,currency"`
}

func (s *Server) createAccount(ctx *gin.Context) {
	var req CreateAccountReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// gin will serialize the key-value pair object to JSON data.
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := database.CreateAccountParams{
		Owner:    req.Owner,
		Currency: req.Currency,
		Balance:  0,
	}
	account, err := s.db.CreateAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, account)
}
