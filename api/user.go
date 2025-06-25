package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/simple_bank/database"
	"github.com/simple_bank/model"
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

type UserResponse struct {
	Username string `json:"username" binding:"required,alphanum"`
	Email    string `json:"email"    binding:"required,email"`
}

func newUserResponse(user model.User) UserResponse {
	return UserResponse{
		Username: user.Username,
		Email:    user.Email,
	}
}

type LoginUserReq struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginUserRsp struct {
	AccessToken string `json:"access_token"`
	User        UserResponse
}

func (s *Server) loginUser(ctx *gin.Context) {
	var req LoginUserReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// 去 DB 找 user，使用 util.CheckPassword
	// 不滿足 => ctx.JSON(http.StatusUnauthorized, errorResponse(err))

	// password 通過後
	token, err := s.tokenMaker.SignToken("username", s.config.Token.Duration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(
		http.StatusOK,
		LoginUserRsp{AccessToken: token, User: newUserResponse(model.User{})},
	)
}
