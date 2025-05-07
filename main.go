package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
)

type KodaktorResponse struct {
	Data struct {
		Login string `json:"login"`
	} `json:"data"`
}

type Response struct {
	Data interface{} `json:"data"`
}

func main() {
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/id/", idHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	log.Printf("Server started on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := Response{
		Data: struct {
			Login string `json:"login"`
		}{
			Login: "fedorovvlad", // Ваш логин
		},
	}

	sendJSON(w, response)
}

func idHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	re := regexp.MustCompile(`^/id/(\d+)/?$`)
	matches := re.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		sendError(w, "Invalid ID format", http.StatusBadRequest)
		return
	}
	n := matches[1]

	req, err := http.NewRequest("GET", "https://nd.kodaktor.ru/users/"+n, nil)
	if err != nil {
		sendError(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	req.Header.Del("Content-Type")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		sendError(w, "Bad Gateway", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		sendError(w, "User not found", http.StatusNotFound)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		sendError(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var kodaktorResp KodaktorResponse
	if err := json.Unmarshal(body, &kodaktorResp); err != nil {
		sendError(w, "Invalid response format", http.StatusInternalServerError)
		return
	}

	response := Response{
		Data: struct {
			Login string `json:"login"`
		}{
			Login: kodaktorResp.Data.Login,
		},
	}

	sendJSON(w, response)
}

func sendJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func sendError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": message,
	})
}
