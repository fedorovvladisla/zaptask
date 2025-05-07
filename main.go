package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

// Структура для ответа от Kodaktor.ru
type KodaktorUser struct {
	Data struct {
		Login string `json:"login"`
	} `json:"data"`
}

func main() {
	r := mux.NewRouter()

	// Маршрут /login - возвращает ваш логин в MOODLE
	r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		response := struct {
			Login string `json:"login"`
		}{
			Login: "fedorovvlad", // Замените на ваш логин
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}).Methods("GET")

	// Маршрут /id/{N} - получает логин пользователя с Kodaktor.ru
	r.HandleFunc("/id/{N:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		vars := mux.Vars(r)
		n := vars["N"]

		// Создаем запрос без Content-Type
		req, err := http.NewRequest("GET", "https://nd.kodaktor.ru/users/"+n, nil)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		req.Header.Del("Content-Type")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		var user KodaktorUser
		if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		response := struct {
			Login string `json:"login"`
		}{
			Login: user.Data.Login,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	log.Printf("Сервер запущен на порту %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
