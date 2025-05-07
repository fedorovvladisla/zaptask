package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type Response struct {
	Login string `json:"login"`
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/login", loginHandler)

	router.HandleFunc("/id/{N}", nHandler)

	port := "5000"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	log.Printf("Сервер запущен на порту %s", port)

	err := http.ListenAndServe(":"+port, router)
	if err != nil {
		log.Fatal(err)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("vlados219"))
}

func nHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	n := vars["N"]
	res, err := http.Get("https://nd.kodaktor.ru/users/" + n)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var resp Response
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Write([]byte(resp.Login))
}
