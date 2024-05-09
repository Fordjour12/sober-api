package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sober-api/internal/helper"
)

func (s *Server) RegisterRoutes() http.Handler {

	mux := http.NewServeMux()
	mux.HandleFunc("/", s.HelloWorldHandler)
	mux.HandleFunc("/health", s.healthHandler)

	mux.HandleFunc("POST /api/v1/create-account", helper.MakeHTTPHandlerFunc(s.CreateUserAccountHandler))
	mux.HandleFunc("POST /api/v1/login-account", helper.MakeHTTPHandlerFunc(s.LogInUserAccountHandler))
	mux.HandleFunc("POST /api/v1/onboarding", helper.MakeHTTPHandlerFunc(s.CreateOnBoardingHandler))
	mux.HandleFunc("POST /api/v1/add-notes", helper.MakeHTTPHandlerFunc(s.AddNotesHandler))
	//
	return mux
}

func (s *Server) CreateUserAccountHandler(w http.ResponseWriter, r *http.Request) error {
	createReq := &helper.CreateAccountRequest{}
	if err := json.NewDecoder(r.Body).Decode(createReq); err != nil {
		return err
	}

	account, err := helper.CreateNewUserAccount(
		createReq.Username,
		createReq.Email,
		createReq.Password,
	)
	if err != nil {
		return err
	}

	userAccount, err := s.db.CreateAccount(account)
	if err != nil {
		return err
	}

	return helper.WriteJSON(w, http.StatusCreated, helper.SuccessResponse{
		Status: http.StatusCreated,
		Data:   userAccount,
	})
}

func (s *Server) LogInUserAccountHandler(w http.ResponseWriter, r *http.Request) error {
	loginReq := &helper.LoginUserRequest{}
	if err := json.NewDecoder(r.Body).Decode(loginReq); err != nil {
		return err
	}

	account, err := s.db.GetUserByEmail(loginReq.Email)
	if err != nil {
		return err
	}

	if !account.ValidPassword(loginReq.Password) {
		log.Printf("not Authenticated %s", loginReq.Password)
		return fmt.Errorf("not Authenticated %s", loginReq.Password)
	}

	token, err := helper.CreateJWTToken(account)
	if err != nil {
		return err
	}

	response := helper.LoginResponse{
		Username: account.Username,
		Email:    account.Email,
		Token:    token,
	}

	return helper.WriteJSON(w, http.StatusOK, helper.SuccessResponse{
		Status: http.StatusOK,
		Data:   response,
	})

}

func (s *Server) CreateOnBoardingHandler(w http.ResponseWriter, r *http.Request) error {
	onBoardingReq := &helper.OnBoardingRequest{}
	if err := json.NewDecoder(r.Body).Decode(onBoardingReq); err != nil {
		return err
	}

	onBoard, err := helper.AccountOnBoarding(
		onBoardingReq.UserId,
		onBoardingReq.Sobriety.ReasonForJoining,
		onBoardingReq.Sobriety.SoberDate,
	)
	if err != nil {
		return err
	}

	if err := s.db.CreateUserOnBoarding(onBoard); err != nil {
		return nil
	}

	return helper.WriteJSON(w, http.StatusOK, helper.SuccessResponse{
		Status: http.StatusOK,
		Data:   onBoard,
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
	if err != nil {
		return err
	}

	if err := s.db.CreateNotes(notes); err != nil {
		return err
	}

	return helper.WriteJSON(w, http.StatusCreated, helper.SuccessResponse{
		Status: http.StatusCreated,
		Data:   notes,
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
