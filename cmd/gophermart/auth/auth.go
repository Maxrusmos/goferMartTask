package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

var users = make(map[string]string)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Парсинг JSON из запроса
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Проверка наличия обязательных полей
	if user.Login == "" || user.Password == "" {
		http.Error(w, "Login and password are required", http.StatusBadRequest)
		return
	}

	// Проверка уникальности логина
	if _, exists := users[user.Login]; exists {
		http.Error(w, "Login already taken", http.StatusConflict)
		return
	}

	// Сохранение пользователя
	users[user.Login] = user.Password

	// Создание токена аутентификации (здесь просто пример, на практике используйте более безопасные методы)
	token := generateToken(user.Login)

	// Установка куки с токеном аутентификации
	expiration := time.Now().Add(24 * time.Hour) // Например, токен действителен 24 часа
	cookie := http.Cookie{Name: "token", Value: token, Expires: expiration, HttpOnly: true}
	http.SetCookie(w, &cookie)

	// Отправка ответа
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "User %s successfully registered and authenticated", user.Login)
}

// Генерация токена аутентификации (здесь просто пример, на практике используйте более безопасные методы)
func generateToken(login string) string {
	// Здесь может быть ваша логика генерации токена, например, JWT или что-то подобное
	return login + "-token"
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Проверка метода запроса
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Парсинг JSON из запроса
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Проверка наличия обязательных полей
	if user.Login == "" || user.Password == "" {
		http.Error(w, "Login and password are required", http.StatusBadRequest)
		return
	}

	// Проверка соответствия переданных учетных данных данным пользователя в хранилище
	storedPassword, ok := users[user.Login]
	if !ok || storedPassword != user.Password {
		http.Error(w, "Invalid login or password", http.StatusUnauthorized)
		return
	}

	// Создание токена аутентификации (здесь просто пример, на практике используйте более безопасные методы)
	token := generateToken(user.Login)

	// Установка куки с токеном аутентификации
	expiration := time.Now().Add(24 * time.Hour) // Например, токен действителен 24 часа
	cookie := http.Cookie{Name: "token", Value: token, Expires: expiration, HttpOnly: true}
	http.SetCookie(w, &cookie)

	// Отправка ответа
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "User %s successfully authenticated", user.Login)
}
