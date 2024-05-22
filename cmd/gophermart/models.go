package main

// type User struct {
// 	ID       int    `json:"id"`
// 	Username string `json:"username"`
// 	Password string `json:"password"`
// 	Balance  int    `json:"balance"`
// }

type Order struct {
	ID       int    `json:"id"`
	UserID   int    `json:"user_id"`
	OrderNum string `json:"order_num"`
	Status   string `json:"status"`
	Points   int    `json:"points"`
}

type Withdrawal struct {
	ID       int    `json:"id"`
	UserID   int    `json:"user_id"`
	Amount   int    `json:"amount"`
	DateTime string `json:"datetime"`
}
