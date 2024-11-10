package models

type User struct {
	Name           string  `json:"name" firestore:"name"`
	Email          string  `json:"email" firestore:"email"`
	Password       string  `json:"password" firestore:"password"`
	Address        string  `json:"address" firestore:"address"`
	PhoneNumber    string  `json:"phoneNumber" firestore:"phoneNumber"`
	Rating         float64 `json:"rating" firestore:"rating"`
	TotalWorkCount int     `json:"totalWorkCount" firestore:"totalWorkCount"`
	UserID         int     `json:"userid" firestore:"userid"`
}
