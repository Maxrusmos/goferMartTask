package loyalty

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sort"
	"time"
)

type Order struct {
	Number     string    `json:"number"`
	Status     string    `json:"status"`
	Accrual    int       `json:"accrual,omitempty"`
	UploadedAt time.Time `json:"uploaded_at"`
}

func UploadOrderHandler(w http.ResponseWriter, r *http.Request) {
	// Проверка аутентификации пользователя
	if !isAuthenticated(r) {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	// Проверка метода запроса
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	orderNumber := string(body)

	// Проверка наличия номера заказа в теле запроса
	if orderNumber == "" {
		http.Error(w, "Order number is required", http.StatusBadRequest)
		return
	}

	// Проверка, что номер заказа не был загружен ранее этим пользователем
	if orderExists(orderNumber) {
		http.Error(w, "Order number already uploaded by this user", http.StatusConflict)
		return
	}

	// Проверка корректности номера заказа с помощью алгоритма Луна
	if !isValidLuhn(orderNumber) {
		http.Error(w, "Invalid order number format", http.StatusUnprocessableEntity)
		return
	}

	// Обработка номера заказа (здесь может быть ваша логика обработки заказа)

	w.WriteHeader(http.StatusAccepted)
}

///////////////////////////////////////////////////////////////////////////////////////////////////////

// GetOrdersHandler обработчик для получения списка заказов пользователя
func GetOrdersHandler(w http.ResponseWriter, r *http.Request) {
	// Проверка авторизации пользователя
	if !isAuthenticated(r) {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	// Получение заказов для пользователя (в данном примере используются моковые данные)
	orders := getOrdersForUser(r)

	// Сортировка заказов по времени загрузки от самых старых к самым новым
	sortOrdersByUploadTime(orders)

	// Формирование JSON-ответа
	w.Header().Set("Content-Type", "application/json")
	if len(orders) > 0 {
		json.NewEncoder(w).Encode(orders)
	} else {
		// Если заказов нет, отправляем статус No Content (204)
		w.WriteHeader(http.StatusNoContent)
	}
}

// getOrdersForUser функция возвращает список заказов для пользователя
func getOrdersForUser(r *http.Request) []Order {
	// Здесь должна быть ваша логика получения заказов для пользователя
	// В данном примере используются моковые данные
	return []Order{
		{Number: "9278923470", Status: "PROCESSED", Accrual: 500, UploadedAt: time.Now()},
		{Number: "12345678903", Status: "PROCESSING", UploadedAt: time.Now()},
		{Number: "346436439", Status: "INVALID", UploadedAt: time.Now()},
	}
}

// sortOrdersByUploadTime функция сортирует заказы по времени загрузки
func sortOrdersByUploadTime(orders []Order) {
	// Здесь можно использовать любой подход к сортировке
	// В данном примере используем встроенную сортировку сортировку по возрастанию времени загрузки
	// Поскольку мы хотим отсортировать заказы от самых старых к самым новым, мы инвертируем порядок сортировки
	sort.Slice(orders, func(i, j int) bool {
		return orders[i].UploadedAt.Before(orders[j].UploadedAt)
	})
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type Balance struct {
	Current   float64 `json:"current"`   // Текущая сумма баллов лояльности
	Withdrawn int     `json:"withdrawn"` // Сумма использованных за весь период баллов
}

func GetBalanceHandler(w http.ResponseWriter, r *http.Request) {
	// Проверка авторизации пользователя
	if !isAuthenticated(r) {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	// Получение данных о балансе пользователя (в данном примере используются моковые данные)
	balance := getBalanceForUser(r)

	// Формирование JSON-ответа
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(balance)
}

// getBalanceForUser функция возвращает данные о балансе пользователя
func getBalanceForUser(r *http.Request) Balance {
	// Здесь должна быть ваша логика получения данных о балансе пользователя
	// В данном примере используются моковые данные
	return Balance{
		Current:   500.5,
		Withdrawn: 42,
	}
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type WithdrawRequest struct {
	Order string `json:"order"` // Номер заказа
	Sum   int    `json:"sum"`   // Сумма баллов к списанию
}

// WithdrawResponse структура представляет собой данные ответа на запрос на списание средств
type WithdrawResponse struct {
	Message string `json:"message"` // Сообщение об успешном списании средств
}

// WithdrawHandler обработчик для запроса на списание средств
func WithdrawHandler(w http.ResponseWriter, r *http.Request) {
	// Проверка авторизации пользователя
	if !isAuthenticated(r) {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	// Распаковка данных запроса
	var withdrawReq WithdrawRequest
	if err := json.NewDecoder(r.Body).Decode(&withdrawReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Проверка наличия достаточного количества средств на балансе пользователя
	balance := getBalanceForUser(r)
	if balance.Current < float64(withdrawReq.Sum) {
		http.Error(w, "Insufficient funds", http.StatusPaymentRequired)
		return
	}

	// Обработка запроса на списание средств (в данном примере не выполняется реальное списание)
	// В случае успешного списания генерируем ответ
	response := WithdrawResponse{
		Message: "Funds successfully withdrawn",
	}

	// Формирование JSON-ответа
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type Withdrawal struct {
	Order       string    `json:"order"`        // Номер заказа
	Sum         int       `json:"sum"`          // Сумма баллов
	ProcessedAt time.Time `json:"processed_at"` // Время обработки вывода средств
}

// GetWithdrawalsHandler обработчик для получения информации о выводах средств
func GetWithdrawalsHandler(w http.ResponseWriter, r *http.Request) {
	// Проверка авторизации пользователя
	if !isAuthenticated(r) {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	// Получение информации о выводах средств для пользователя (в данном примере используются моковые данные)
	withdrawals := getWithdrawalsForUser(r)

	// Сортировка выводов средств по времени вывода от самых старых к самым новым
	sortWithdrawalsByProcessedTime(withdrawals)

	// Формирование JSON-ответа
	w.Header().Set("Content-Type", "application/json")
	if len(withdrawals) > 0 {
		json.NewEncoder(w).Encode(withdrawals)
	} else {
		// Если нет ни одного вывода средств, отправляем статус No Content (204)
		w.WriteHeader(http.StatusNoContent)
	}
}

// getWithdrawalsForUser функция возвращает информацию о выводах средств для пользователя
func getWithdrawalsForUser(r *http.Request) []Withdrawal {
	// Здесь должна быть ваша логика получения информации о выводах средств для пользователя
	// В данном примере используются моковые данные
	return []Withdrawal{
		{Order: "2377225624", Sum: 500, ProcessedAt: time.Now()},
	}
}

// sortWithdrawalsByProcessedTime функция сортирует выводы средств по времени обработки
func sortWithdrawalsByProcessedTime(withdrawals []Withdrawal) {
	// Здесь можно использовать любой подход к сортировке
	// В данном примере используем встроенную сортировку сортировку по возрастанию времени обработки
	// Поскольку мы хотим отсортировать выводы средств от самых старых к самым новым, мы инвертируем порядок сортировки
	sort.Slice(withdrawals, func(i, j int) bool {
		return withdrawals[i].ProcessedAt.Before(withdrawals[j].ProcessedAt)
	})
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// Функция для проверки аутентификации пользователя
func isAuthenticated(r *http.Request) bool {
	// Ваша логика проверки аутентификации, например, проверка наличия токена аутентификации в запросе
	return true // Возвращаем true для примера, замените на вашу логику
}

// Функция для проверки существования номера заказа
func orderExists(orderNumber string) bool {
	// Ваша логика проверки существования заказа, например, проверка в базе данных
	return false // Возвращаем false для примера, замените на вашу логику
}

// Функция для проверки корректности номера заказа с помощью алгоритма Луна
func isValidLuhn(orderNumber string) bool {
	// Ваша логика проверки корректности номера заказа с помощью алгоритма Луна
	return true // Возвращаем true для примера, замените на вашу логику
}
