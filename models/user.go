package models

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// A sample user for demonstration purposes
var SampleUser = User{
	Username: "admin",
	Password: "123", // Store passwords securely with hashing in a real app
}
