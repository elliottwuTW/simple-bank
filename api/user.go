package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/simple_bank/database"
	"github.com/simple_bank/util"
)

type CreateUserReq struct {
	Username string `json:"username" binding:"required,alphanum"` // only ASCII alphanumeric characters
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email"    binding:"required,email"`
}

func (s *Server) createUser(ctx *gin.Context) {
	var req CreateUserReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := database.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		Email:          req.Email,
	}
	user, err := s.db.CreateUser(ctx, arg)
	if err != nil {
		// DB error 特殊處理
		// if pgErr, ok := err.(*pq.Error);ok {
		// 	switch pqErr.Code.Name():
		// case "xxx":
		// ctx.JSON(http.StatusForbidden, errorResponse(err))
		// return
		// }
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// 亦可定義 CreateUserRsp 來把 hashedPassword 隱藏
	ctx.JSON(http.StatusCreated, user)
}
