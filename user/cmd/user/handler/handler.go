package handler

import (
	"net/http"
	"user/cmd/user/usecase"
	"user/infrastructure/log"
	"user/models"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	UserUseCase usecase.UserUseCase
}

func NewUserHandler(userUseCase usecase.UserUseCase) *UserHandler {
	return &UserHandler{
		UserUseCase: userUseCase,
	}

}

func (h *UserHandler) RegisterRoutes(c *gin.Context) {
	var param models.RegisterParameter
	if err := c.ShouldBindJSON(&param); err != nil {
		log.Logger.Info("invalid parameter")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(param.Password) < 8 || len(param.ConfirmPassword) < 8 {
		log.Logger.Info("Password must be at least 8 characters long")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 8 characters long"})
		return
	}

	if param.Password != param.ConfirmPassword {
		log.Logger.Info("Password and confirm Password do not match")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password and confirm Password do not match"})
		return
	}

	//call the usecase
	user, err := h.UserUseCase.GetUserByEmail(c.Request.Context(), param.Email)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if user != nil && user.ID != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already registered"})
		return
	}
	user = &models.User{
		Name:     param.Name,
		Email:    param.Email,
		Password: param.Password,
	}

	if err := h.UserUseCase.RegisterUser(c.Request.Context(), user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Register success"})
}

func (h *UserHandler) LoginRoutes(c *gin.Context) {
	var param models.LoginParameter
	if err := c.ShouldBindJSON(&param); err != nil {
		log.Logger.Info("invalid parameter")
		c.JSON(http.StatusBadRequest, gin.H{"error_message": "Invalid Input Parameter"})
		return
	}

	if len(param.Password) < 8 {
		log.Logger.Info("Password must be at least 8 characters long")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 8 characters long"})
		return
	}

	token, err := h.UserUseCase.Login(c.Request.Context(), &param)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Wrong Email or Password ",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})

}

func (h *UserHandler) GetUserInfo(c *gin.Context) {
	userIDStr, isExists := c.Get("user_id")
	if !isExists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}
	userId, ok := userIDStr.(float64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid User id",
		})
		return
	}

	user, err := h.UserUseCase.GetUserById(c.Request.Context(), int64(userId))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"User": user,
	})
}

func (h *UserHandler) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "pong",
	})
}
