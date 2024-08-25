package model

type User struct {
	UserId       int64  `json:"userid"`
	Email        string `json:"email"`
	RefreshToken string `json:"refreshtoken"`
}
