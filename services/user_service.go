package services

import (
	"context"
	"errors"

	"cloud.google.com/go/firestore"
	"github.com/imthewolverine/schoolPolice-backend/models"
)

type UserService struct {
    FirestoreClient *firestore.Client
}

func NewUserService(client *firestore.Client) *UserService {
    return &UserService{FirestoreClient: client}
}

// CreateUser adds a new user to Firestore, checking for duplicates by name or email
func (s *UserService) CreateUser(ctx context.Context, user models.User) error {
    users := s.FirestoreClient.Collection("user")

    // Check if a user with the same email or name already exists
    existingUser, err := s.findUserByUsernameOrEmail(ctx, user.Email)
    if err == nil && existingUser != nil {
        return errors.New("user with this email already exists")
    }
    
    existingUser, err = s.findUserByUsernameOrEmail(ctx, user.Name)
    if err == nil && existingUser != nil {
        return errors.New("user with this name already exists")
    }

    // Step 1: Create a new account document with an initial balance of 0
    accountRef := s.FirestoreClient.Collection("account").NewDoc()
    _, err = accountRef.Set(ctx, map[string]interface{}{
        "balance": 0,  // Initial balance
    })
    if err != nil {
        return errors.New("failed to create account for user")
    }

    // Step 2: Add the new user to Firestore, including the account reference
    _, _, err = users.Add(ctx, map[string]interface{}{
        "name":           user.Name,
        "email":          user.Email,
        "password":       user.Password,  // In production, hash this password
        "address":        user.Address,
        "phoneNumber":    user.PhoneNumber,
        "rating":         user.Rating,
        "totalWorkCount": user.TotalWorkCount,
        "userid":         user.UserID,
        "account":        accountRef,     // Reference to the new account document
    })

    return err
}

// findUserByUsernameOrEmail searches Firestore for a user by name or email
func (s *UserService) findUserByUsernameOrEmail(ctx context.Context, usernameOrEmail string) (map[string]interface{}, error) {
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
