package helper

import (
	"log"
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Account struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"createdAt"`
}

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

type CreateNotesRequest struct {
	UserId    int       `json:"userId"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
	UpdateAt  time.Time `json:"updateAt"`
}

type LoginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Token    string `json:"token"`
}

func CreateNewUserAccount(username, email, password string) (*Account, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return &Account{
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

func LoginUserAccount(email, password string) (*LoginUserRequest, error) {
	return &LoginUserRequest{
		Email:    email,
		Password: password,
	}, nil
}

func (ac *Account) ValidPassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(ac.Password), []byte(password)) == nil
}

func CreateJWTToken(ac *Account) (string, error) {

	claims := &jwt.MapClaims{
		"expiresAt": time.Now().Add(time.Minute * 15).Unix(),
		"username":  ac.Username,
		"email":     ac.Email,
	}

	secret := os.Getenv("JWT_SECRET")
	log.Printf("jwt %+v:", secret)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

func AccountOnBoarding(userId int, reason, date string) (*OnBoardingRequest, error) {
	return &OnBoardingRequest{
		UserId: userId,
		Sobriety: Sobriety{
			ReasonForJoining: reason,
			SoberDate:        date,
		},
		CreatedAt: time.Now().UTC(),
	}, nil
}
