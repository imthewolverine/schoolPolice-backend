package services

import (
	"context"
	"errors"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/dgrijalva/jwt-go"
	"github.com/imthewolverine/schoolPolice-backend/config"
)

type AuthService struct {
    FirestoreClient *firestore.Client
}

func NewAuthService(client *firestore.Client) *AuthService {
    return &AuthService{FirestoreClient: client}
}

// AuthenticateUser checks the provided credentials against Firestore
func (s *AuthService) AuthenticateUser(ctx context.Context, usernameOrEmail, password string) (string, error) {
    userDoc, err := s.findUserByUsernameOrEmail(ctx, usernameOrEmail)
    if err != nil {
        return "", err
    }

    // Validate password
    if userDoc["password"] != password {
        return "", errors.New("invalid username or password")
    }

    // Generate JWT if authentication succeeds
    return s.GenerateJWT(userDoc["name"].(string))
}

// findUserByUsernameOrEmail searches Firestore for a user by name or email
func (s *AuthService) findUserByUsernameOrEmail(ctx context.Context, usernameOrEmail string) (map[string]interface{}, error) {
    users := s.FirestoreClient.Collection("user")

    // Check by email
    emailQuery := users.Where("email", "==", usernameOrEmail).Documents(ctx)
    emailDocs, err := emailQuery.GetAll()
    if err != nil {
        return nil, err
    }
    if len(emailDocs) > 0 {
        return emailDocs[0].Data(), nil
    }

    // Check by name
    nameQuery := users.Where("name", "==", usernameOrEmail).Documents(ctx)
    nameDocs, err := nameQuery.GetAll()
    if err != nil {
        return nil, err
    }
    if len(nameDocs) > 0 {
        return nameDocs[0].Data(), nil
    }

    return nil, errors.New("user not found")
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


