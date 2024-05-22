package main

import (
	"goferMartTask/cmd/gophermart/auth"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/api/user/register", auth.RegisterHandler).Methods("POST")
	r.HandleFunc("/api/user/login", auth.LoginHandler).Methods("POST")

	// s := r.PathPrefix("/api/user").Subrouter()
	// s.Use(AuthMiddleware)
	// s.HandleFunc("/orders", UploadOrderHandler).Methods("POST")
	// s.HandleFunc("/orders", GetOrdersHandler).Methods("GET")
	// s.HandleFunc("/balance", GetBalanceHandler).Methods("GET")
	// s.HandleFunc("/balance/withdraw", WithdrawHandler).Methods("POST")
	// s.HandleFunc("/withdrawals", GetWithdrawalsHandler).Methods("GET")

	http.ListenAndServe(":8080", r)
}
