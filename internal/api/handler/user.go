package handler

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"idempotency/user"
	"net/http"
)

type (
	UserUseCase interface {
		Save(ctx context.Context, user user.User) (user.User, error)
		Get(ctx context.Context, id uint) (user.User, error)
	}

	UserHandler struct {
		uc UserUseCase
	}

	createUserRequest struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	getUserRequest struct {
		ID uint `uri:"id"`
	}

	response struct {
		ID    uint   `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}
)

func NewUserHandler(uc UserUseCase) *UserHandler {
	return &UserHandler{
		uc: uc,
	}
}

func (h UserHandler) CreateUser(ctx *gin.Context) {
	var request createUserRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.String(http.StatusBadRequest, "invalid request")
		return
	}

	u, err := h.uc.Save(ctx.Request.Context(), user.User{
		Name:  request.Name,
		Email: request.Email,
	})

	if err != nil {
		ctx.String(http.StatusInternalServerError, "internal server error")
		return
	}

	ctx.JSON(http.StatusCreated, response{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
	})
}

func (h UserHandler) GetUser(ctx *gin.Context) {
	var request getUserRequest
	if err := ctx.ShouldBindUri(&request); err != nil {
		ctx.String(http.StatusBadRequest, "id is required")
		return
	}

	u, err := h.uc.Get(ctx.Request.Context(), request.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.String(http.StatusNotFound, "user not found")
			return
		}

		ctx.String(http.StatusInternalServerError, "internal server error")
		return
	}

	ctx.JSON(http.StatusCreated, response{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
	})
}
