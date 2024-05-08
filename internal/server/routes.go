package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sober-api/internal/helper"

	"golang.org/x/crypto/bcrypt"
)

func (s *Server) RegisterRoutes() http.Handler {

	mux := http.NewServeMux()
	mux.HandleFunc("/", s.HelloWorldHandler)

	mux.HandleFunc("/health", s.healthHandler)

	mux.HandleFunc("POST /api/v1/onboarding", helper.MakeHTTPHandlerFunc(s.OnBoardingHandler))
	mux.HandleFunc("POST /api/v1/create-account", helper.MakeHTTPHandlerFunc(s.CreateAccountHandler))
	mux.HandleFunc("POST /api/v1/login-account", helper.MakeHTTPHandlerFunc(s.LogInAccountHandler))
	mux.HandleFunc("POST /api/v1/add-notes", helper.MakeHTTPHandlerFunc(s.AddNotesHandler))

	return mux
}

func (s *Server) OnBoardingHandler(w http.ResponseWriter, r *http.Request) error {

	onBoardingReq := &helper.OnBoardingRequest{}
	if err := json.NewDecoder(r.Body).Decode(onBoardingReq); err != nil {
		return err
	}

	onBoard, err := helper.AddOnBoardingFlow(
		onBoardingReq.UserId,
		onBoardingReq.Sobriety.ReasonForJoining,
		onBoardingReq.Sobriety.SoberDate,
	)

	if err != nil {
		return err
	}

	if err := s.db.CreateOnBoardingFlow(onBoard); err != nil {
		return err
	}

	return helper.WriteJSON(w, http.StatusOK, helper.SuccessResponse{
		Data: onBoard,
	})
}

func (s *Server) CreateAccountHandler(w http.ResponseWriter, r *http.Request) error {

	createAccountReq := &helper.CreateAccountRequest{}
	if err := json.NewDecoder(r.Body).Decode(createAccountReq); err != nil {
		return err
	}
	account, err := helper.CreateUserAccount(
		createAccountReq.Username,
		createAccountReq.Email,
		createAccountReq.Password,
	)

	fmt.Printf("Account: %+v\n", account.Password)

	if err != nil {
		return err
	}

	dst, err := s.db.CreateAccountFlow(account)

	if err != nil {
		return err
	}

	return helper.WriteJSON(w, http.StatusOK, helper.SuccessResponse{
		// Data: map[string]any{"id": dst, "account": account},
		Data: struct {
			ID      int                          `json:"id"`
			Account *helper.CreateAccountRequest `json:"account"`
		}{
			ID:      dst,
			Account: account,
		},
	})
}

func (s *Server) AddNotesHandler(w http.ResponseWriter, r *http.Request) error {

	createNotesReq := &helper.CreateNotesRequest{}
	if err := json.NewDecoder(r.Body).Decode(createNotesReq); err != nil {
		return err
	}

	notes, err := helper.CreateNewNotes(
		createNotesReq.UserId,
		createNotesReq.Content,
	)

	fmt.Printf("Account: %+v\n", notes)

	if err != nil {
		return err
	}

	userId, err := s.db.CreateNotesFlow(notes)
	if err != nil {
		return err
	}

	return helper.WriteJSON(w, http.StatusCreated, helper.SuccessResponse{
		Data: struct {
			ID      int                        `json:"id"`
			Account *helper.CreateNotesRequest `json:"notes"`
		}{
			ID:      userId,
			Account: notes,
		},
	})

}

func (s *Server) LogInAccountHandler(w http.ResponseWriter, r *http.Request) error {

	//loginReq := &helper.LoginUserRequest{}
	var loginReq helper.LoginUserRequest
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		return err
	}

	account, err := s.db.GetUserByEmail(loginReq.Email)
	if err != nil {
		return err
	}

	fmt.Printf("Account: %+v\n", account.Password)

	isValid := bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(loginReq.Password)) == nil

	fmt.Printf("Account: %+v\n", isValid)

	// if !account.Password(loginReq.Password) {
	// 	return fmt.Errorf("not Auth %s", loginReq.Password)
	// }

	//	token, err := helper.CreateJWTToken(account)
	//	if err != nil {
	//		return err
	//	}
	//
	//	response := helper.LoginResponse{
	//		Username: account.Username,
	//		Email:    account.Email,
	//		Token:    token,
	//	}
	//
	return helper.WriteJSON(w, http.StatusOK, helper.SuccessResponse{
		Data: account,
	})

}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, err := json.Marshal(s.db.Health())

	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}
