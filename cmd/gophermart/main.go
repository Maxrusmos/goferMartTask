package main

import (
	"goferMartTask/cmd/gophermart/auth"
	"goferMartTask/cmd/gophermart/loyalty"
	"goferMartTask/cmd/gophermart/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/api/user/register", auth.RegisterHandler).Methods("POST")

	r.HandleFunc("/api/user/login", auth.LoginHandler).Methods("POST")

	s := r.PathPrefix("/api/user").Subrouter()
	s.Use(middleware.AuthMiddleware)
	s.HandleFunc("/orders", loyalty.UploadOrderHandler).Methods("POST")

	s.HandleFunc("/orders", loyalty.GetOrdersHandler).Methods("GET")

	s.HandleFunc("/balance", loyalty.GetBalanceHandler).Methods("GET")

	s.HandleFunc("/balance/withdraw", loyalty.WithdrawHandler).Methods("POST")

	s.HandleFunc("/withdrawals", loyalty.GetWithdrawalsHandler).Methods("GET")

	http.ListenAndServe(":8080", r)
}
