package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Descriptions map[string]map[string]string

type Pillar struct {
	Stem         string   `json:"stem"`
	StemElement  string   `json:"stemElement"`
	StemYinYang  string   `json:"stemYinYang"`
	Branch       string   `json:"branch"`
	BranchElement string  `json:"branchElement"`
	HiddenStems  []string `json:"hiddenStems"`
}

type DayMaster struct {
	Stem    string `json:"stem"`
	Element string `json:"element"`
	YinYang string `json:"yinYang"`
}

type SeasonInfo struct {
	Name    string `json:"name"`
	Branch  string `json:"branch"`
	Element string `json:"element"`
}

type YongShenInfo struct {
	Element string `json:"element"`
	Reason  string `json:"reason"`
}

type TenGodEntry struct {
	Element string `json:"element"`
	God     string `json:"god"`
}

type LuckCycle struct {
	StartAge int    `json:"startAge"`
	EndAge   int    `json:"endAge"`
	Stem     string `json:"stem"`
	Branch   string `json:"branch"`
	Element  string `json:"element"`
}

type Profile struct {
	Dominant           string                   `json:"dominant"`
	Lp                 int                      `json:"lp"`
	Balance            map[string]int           `json:"balance"`
	Descriptions       Descriptions             `json:"descriptions"`
	Pillars            map[string]Pillar        `json:"pillars"`
	DayMaster          DayMaster                `json:"dayMaster"`
	DmStrength         string                   `json:"dmStrength"`
	Season             SeasonInfo               `json:"season"`
	YongShen           YongShenInfo             `json:"yongShen"`
	FavorableElements  []string                 `json:"favorableElements"`
	UnfavorableElements []string                `json:"unfavorableElements"`
	TenGods            map[string]TenGodEntry   `json:"tenGods"`
	LuckCycles         []LuckCycle              `json:"luckCycles"`
}

type profileRequest struct {
	Birthdate string `json:"birthdate"`
	Birthtime string `json:"birthtime"`
	Gender    string `json:"gender"`
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

	log.Printf("profile request: birthdate=%q birthtime=%q gender=%q", req.Birthdate, req.Birthtime, req.Gender)

	profile := computeProfile(req.Birthdate, req.Birthtime, req.Gender)

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
