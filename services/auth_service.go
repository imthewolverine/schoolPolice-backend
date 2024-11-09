package services

import (
	"errors"
	"time"

	"github.com/imthewolverine/schoolPolice-backend/config"
	"github.com/imthewolverine/schoolPolice-backend/models"

	"github.com/dgrijalva/jwt-go"
)

// AuthService provides authentication methods
type AuthService struct{}

func (s *AuthService) AuthenticateUser(username, password string) (string, error) {
    if username == models.SampleUser.Username && password == models.SampleUser.Password {
        return s.GenerateJWT(username)
    }
    return "", errors.New("invalid username or password")
}

// GenerateJWT creates a JWT token
func (s *AuthService) GenerateJWT(username string) (string, error) {
    secret := config.GetEnv("JWT_SECRET")
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "username": username,
        "exp":      time.Now().Add(time.Hour * 72).Unix(),
    })
    return token.SignedString([]byte(secret))
}
