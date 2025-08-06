package controller

import (
	"net/http"
	"setup-preoject/app/middleware"
	"setup-preoject/app/model/dto"
	"setup-preoject/app/model/entity"
	"setup-preoject/app/util"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type AuthController struct {
	db    *gorm.DB
	cache *redis.Client
}

func NewAuthController(db *gorm.DB, cache *redis.Client) *AuthController {
	return &AuthController{
		db:    db,
		cache: cache,
	}
}

func (ac *AuthController) Login(c *gin.Context) {
	var request dto.LoginRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var findUser entity.User
	result := ac.db.Where("username = ?", request.Username).Find(&findUser)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "username error"})
		return
	}

	if findUser.Password != util.GeneratePassword(request.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "password error"})
		return
	}

	token, err := middleware.GenerateToken(findUser, ac.cache)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"user": gin.H{
			"id":       findUser.Id,
			"username": findUser.Username,
		},
		"token": gin.H{
			"key": token,
			"exp": time.Now().Add(time.Hour * 24 * 20).Unix(),
		},
	})
	return
}

func (ac *AuthController) Register(c *gin.Context) {
	var request dto.RegisterRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if request.Password != request.ConfirmPassword {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "password error"})
		return
	}

	entityUser := entity.User{
		Username: request.Username,
		Password: util.GeneratePassword(request.Password),
	}
	result := ac.db.Create(&entityUser)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "Duplicate") {
			c.JSON(http.StatusConflict, gin.H{"error": "username already exists"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	token, err := middleware.GenerateToken(entityUser, ac.cache)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":       entityUser.Id,
			"username": entityUser.Username,
		},
		"token": gin.H{
			"key": token,
			"exp": time.Now().Add(time.Hour * 24 * 20).Unix(),
		},
	})
}

func (ac *AuthController) Me(c *gin.Context) {
	userID, exists := c.Get("ssi")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user_id not found"})
		return
	}

	var findUser entity.User
	result := ac.db.Find(&findUser, userID)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user_id not found"})
		return
	}

	c.JSON(200, gin.H{
		"user": gin.H{
			"id":       findUser.Id,
			"username": findUser.Username,
		},
	})
	return
}
