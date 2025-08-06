package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"setup-preoject/app/model/entity"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

func GenerateToken(user entity.User, cache *redis.Client) (string, error) {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	key := hex.EncodeToString(bytes)

	claims := jwt.MapClaims{
		"username": user.Username,
		"ssi":      user.Id,
		"exp":      time.Now().Add(time.Hour * 24 * 20).Unix(),
		//		"role": 	user.Role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}

	ctx := context.Background()
	expiration := time.Hour * 24 * 20
	err = cache.Set(ctx, signedToken, key, expiration).Err()
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func AuthMiddleware(cache *redis.Client, roles []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			c.Abort()
			return
		}

		signedToken := parts[1]

		token, err := validateToken(signedToken, cache)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		claims, claimsOk := token.Claims.(jwt.MapClaims)
		if claimsOk {
			c.Set("username", claims["username"])
			c.Set("ssi", claims["ssi"])
		}

		if role, ok := claims["role"]; ok && len(roles) != 0 {
			if containsRole(roles, role.(string)) {
				c.Set("role", role.(string))
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid role for token"})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

func containsRole(roles []string, role string) bool {
	for _, item := range roles {
		if item == role {
			return true
		}
	}
	return false
}

func validateToken(signedToken string, cache *redis.Client) (*jwt.Token, error) {
	ctx := context.Background()

	key, err := cache.Get(ctx, signedToken).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("token not found or expired")
		}
		return nil, err
	}

	token, err := jwt.Parse(signedToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(key), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return token, nil
}
