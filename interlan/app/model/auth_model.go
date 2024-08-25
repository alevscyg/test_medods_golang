package model

type Auth struct {
	RefreshToken string `json:"refreshtoken"`
	AccessToken  string `json:"accesstoken"`
}
