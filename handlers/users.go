package handlers

import (
	"net/http"
	"time"

	"github.com/abiiranathan/acada/auth"
	"github.com/abiiranathan/acada/models"
	"github.com/abiiranathan/acada/services"
	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	GetAllUsers() gin.HandlerFunc
	CreateUser() gin.HandlerFunc
	GetUser() gin.HandlerFunc
	UpdateUser() gin.HandlerFunc
	DeleteUser() gin.HandlerFunc
	Login() gin.HandlerFunc
	VerifyToken() gin.HandlerFunc
}

type userhandler struct {
	UserService services.UserService
}

func NewUserHandler(svc services.UserService) UserHandler {
	return &userhandler{svc}
}

func (h *userhandler) GetAllUsers() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		users, err := h.UserService.GetAllUsers()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
			ctx.Abort()
			return
		}

		ctx.JSON(http.StatusOK, users)
	}
}

func (h *userhandler) GetUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := GetParam(ctx, "id")
		if err != nil {
			return
		}

		user, err := h.UserService.GetUser(id)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err})
			ctx.Abort()
			return
		}
		ctx.JSON(http.StatusOK, user)
	}
}

func (h *userhandler) CreateUser() gin.HandlerFunc {
	type NewUser struct {
		Name     string `json:"name" gorm:"not null" binding:"required"`
		Email    string `json:"email" gorm:"not null;unique" binding:"email,lowercase,required"`
		Password string `json:"password" gorm:"not null" binding:"required,min=8"`
	}

	return func(ctx *gin.Context) {
		var new_user NewUser
		if err := ctx.ShouldBindJSON(&new_user); err != nil {
			RespondWithError(ctx, err)
			return
		}

		user := &models.User{
			Name:      new_user.Name,
			Email:     new_user.Email,
			Password:  new_user.Password,
			CreatedAt: time.Now(),
		}

		err := h.UserService.CreateUser(user)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		ctx.JSON(http.StatusCreated, user)
	}
}

func (h *userhandler) UpdateUser() gin.HandlerFunc {
	type UpdateUser struct {
		Name     string `json:"name" gorm:"not null" binding:"omitempty"`
		Email    string `json:"email" gorm:"not null;unique" binding:"omitempty,email,lowercase"`
		Password string `json:"password" gorm:"not null" binding:"omitempty,min=8"`
	}

	return func(ctx *gin.Context) {
		id, err := GetParam(ctx, "id")
		if err != nil {
			return
		}

		var updates UpdateUser
		if err := ctx.ShouldBindJSON(&updates); err != nil {
			RespondWithError(ctx, err)
			return
		}

		user := models.User{
			Name:     updates.Name,
			Email:    updates.Email,
			Password: updates.Password,
		}

		updatedUser, err := h.UserService.UpdateUser(id, user)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		ctx.JSON(http.StatusOK, updatedUser)
	}
}

func (h *userhandler) DeleteUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := GetParam(ctx, "id")
		if err != nil {
			return
		}

		err = h.UserService.DeleteUser(id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		ctx.JSON(http.StatusNoContent, nil)
	}
}

func (h *userhandler) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var login models.Login

		if err := c.ShouldBindJSON(&login); err != nil {
			RespondWithError(c, err)
			return
		}

		user, err := h.UserService.GetByEmail(login.Email)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid login credentials"})
			return
		}

		// Check if user is active
		if !user.IsActive {
			c.JSON(http.StatusForbidden, gin.H{"error": "Your account is not active"})
			c.Abort()
			return
		}

		if !auth.CheckPasswordHash(login.Password, user.Password) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid login credentials"})
			return
		}

		token, err := auth.CreateToken(user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
			return
		}

		c.JSON(http.StatusOK, models.LoginResponse{
			User:  user,
			Token: token,
		})
	}
}

func (r *userhandler) VerifyToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, ok := c.Get("user")

		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "Login is required"})
			c.Abort()
			return
		}

		if sessionUser, ok := user.(models.User); ok {
			c.JSON(http.StatusOK, sessionUser)
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user type"})
	}
}
