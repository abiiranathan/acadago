package auth

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidToken = errors.New("jwt token is invalid")

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Creates a jwt token that expires after the specified duration
// as env variable JWT_EXPIRATION_HOURS or 48 hours
func CreateToken(userId uint) (string, error) {
	var expiresHours = 48
	expiresAfter := os.Getenv("JWT_EXPIRATION_HOURS")

	if expiresAfter != "" {
		expiresAfterInt, err := strconv.Atoi(expiresAfter)
		if err != nil {
			expiresHours = expiresAfterInt
		}
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = userId
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(expiresHours)).Unix()
	return token.SignedString([]byte(os.Getenv("SECRET")))
}

// Verifies the token string
// Returns the user id from the payload or an error
func VerifyToken(tokenString string) (uint, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(os.Getenv("SECRET")), nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		id := uint(claims["id"].(float64))
		return id, nil
	}

	return 0, ErrInvalidToken
}
