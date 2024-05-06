package helper

import (
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type CreateAccountRequest struct {
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
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

type CreateNotesRequest struct {
	UserId    int       `json:"userId"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
	UpdateAt  time.Time `json:"updateAt"`
}

func (c *CreateAccountRequest) ValidatePassword(pw string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(c.Password), []byte(pw))
	if err != nil {
		return false, err
	}

	return true, nil
}

func CreateUserAccount(username, email, password string) (*CreateAccountRequest, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &CreateAccountRequest{
		Username:  username,
		Email:     email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now().UTC(),
	}, nil

}

func AddOnBoardingFlow(userId int, reason, date string) (*OnBoardingRequest, error) {
	return &OnBoardingRequest{
		UserId: userId,
		Sobriety: Sobriety{
			ReasonForJoining: reason,
			SoberDate:        date,
		},
		CreatedAt: time.Now().UTC(),
	}, nil
}

func CreateNewNotes(userId int, content string) (*CreateNotesRequest, error) {
	return &CreateNotesRequest{
		UserId:    userId,
		Content:   content,
		CreatedAt: time.Now().UTC(),
		UpdateAt:  time.Now().UTC(),
	}, nil
}

func CreateJWTToken(account *CreateAccountRequest) (string, error) {

	claims := &jwt.MapClaims{
		"expiresAt": time.Now().Add(time.Minute * 15).Unix(),
		"username":  account.Username,
		"email":     account.Email,
	}

	secret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}
