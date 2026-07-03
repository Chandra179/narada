package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Profile struct {
	Dominant string         `json:"dominant"`
	Lp       int            `json:"lp"`
	Balance  map[string]int `json:"balance"`
}

type profileRequest struct {
	Birthdate string `json:"birthdate"`
	Birthtime string `json:"birthtime"` // optional, format "HH:MM", defaults to "12:00"
}

func handleProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req profileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if req.Birthdate == "" {
		http.Error(w, "birthdate required", http.StatusBadRequest)
		return
	}

	profile := computeProfile(req.Birthdate, req.Birthtime)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(profile); err != nil {
		log.Printf("error encoding response: %v", err)
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/profile", handleProfile)
	log.Println("listening on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
