package services

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/imthewolverine/schoolPolice-backend/models"
)

var client *firestore.Client

func init() {
    ctx := context.Background()
    app, err := firebase.NewApp(ctx, nil)
    if err != nil {
        log.Fatalf("Error initializing Firebase app: %v", err)
    }

    client, err = app.Firestore(ctx)
    if err != nil {
        log.Fatalf("Error initializing Firestore client: %v", err)
    }
}

func CreateUser(user models.User) (models.User, error) {
    _, _, err := client.Collection("users").Add(context.Background(), user)
    if err != nil {
        return user, err
    }
    return user, nil
}
