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

	mux.HandleFunc("POST /api/v1/onboarding", helper.MakeHTTPHandlerFunc(s.OnBoardingHandler))
	mux.HandleFunc("POST /api/v1/create-account", helper.MakeHTTPHandlerFunc(s.CreateAccountHandler))
	mux.HandleFunc("POST /api/v1/login-account", helper.MakeHTTPHandlerFunc(s.LogInAccountHandler))

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

	if err := s.db.CreateAccountFlow(account); err != nil {
		return err
	}

	return helper.WriteJSON(w, http.StatusOK, helper.SuccessResponse{
		Data: account,
	})
}

func (s *Server) LogInAccountHandler(w http.ResponseWriter, r *http.Request) error {
	return helper.WriteJSON(w, http.StatusOK, helper.SuccessResponse{
		Data: map[string]string{"message": "Account Logged In"},
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
