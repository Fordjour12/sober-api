package helper

import (
	"time"
)

type CreateAccountRequest struct {
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"createdAt"`
}

type Sobriety struct {
	ReasonForJoining string `json:"reason"`
	SoberDate        string `json:"soberDate"`
}

type OnBoardingRequest struct {
	UserId    int       `json:"userId"`
	Sobriety  Sobriety  `json:"sobriety"`
	CreatedAt time.Time `json:"createdAt"`
}

func CreateUserAccount(username, email, password string) (*CreateAccountRequest, error) {
	return &CreateAccountRequest{
		Username: username,
		Email:    email,
		Password: password,
	}, nil

}

func AddOnBoardingFlow(userId int, reason, date string) (*OnBoardingRequest, error) {

	return &OnBoardingRequest{
		UserId: userId,
		Sobriety: Sobriety{
			ReasonForJoining: reason,
			SoberDate:        date,
		},
	}, nil

}
