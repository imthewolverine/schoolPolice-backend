package models

// User - Example user model
type User struct {
    ID       string `json:"id" firestore:"id,omitempty"`
    Name     string `json:"name" firestore:"name"`
    Email    string `json:"email" firestore:"email"`
    Password string `json:"-" firestore:"password,omitempty"` // Hide password in JSON
}
