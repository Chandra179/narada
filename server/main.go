package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Descriptions map[string]map[string]string

type Pillar struct {
	Stem          string   `json:"stem"`
	StemElement   string   `json:"stemElement"`
	StemYinYang   string   `json:"stemYinYang"`
	Branch        string   `json:"branch"`
	BranchElement string   `json:"branchElement"`
	HiddenStems   []string `json:"hiddenStems"`
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

// --- New types for personality assembly ---

type DayMasterProfileInfo struct {
	Archetype  string   `json:"archetype"`
	CoreTraits []string `json:"coreTraits"`
	Summary    string   `json:"summary"`
}

type TenGodInsight struct {
	God         string `json:"god"`
	Category    string `json:"category,omitempty"`
	Description string `json:"description"`
}

type Profile struct {
	Dominant            string                    `json:"dominant"`
	Lp                  int                       `json:"lp"`
	Balance             map[string]int            `json:"balance"`
	Descriptions        Descriptions              `json:"descriptions"`
	Pillars             map[string]Pillar         `json:"pillars"`
	DayMaster           DayMaster                 `json:"dayMaster"`
	DmStrength          string                    `json:"dmStrength"`
	Season              SeasonInfo                `json:"season"`
	YongShen            YongShenInfo              `json:"yongShen"`
	FavorableElements   []string                  `json:"favorableElements"`
	UnfavorableElements []string                  `json:"unfavorableElements"`
	TenGods             map[string]TenGodEntry    `json:"tenGods"`
	LuckCycles          []LuckCycle               `json:"luckCycles"`
	TabText             map[string]string         `json:"tabText"`
	DayMasterProfile    *DayMasterProfileInfo     `json:"dayMasterProfile"`
	TenGodInsights      map[string]TenGodInsight  `json:"tenGodInsights"`
	TextModifiers       []string                  `json:"textModifiers"`
	LifePathText        string                    `json:"lifePathText"`
}

type profileRequest struct {
	Birthdate string `json:"birthdate"`
	Birthtime string `json:"birthtime"`
	Gender    string `json:"gender"`
}

// --- Config response types ---

type TabResponse struct {
	Key   string `json:"key"`
	Label string `json:"label"`
	Icon  string `json:"icon"`
}

type ElementMetaResponse struct {
	Label string `json:"label"`
	Short string `json:"short"`
	Color string `json:"color"`
}

type CityResponse struct {
	Name      string  `json:"name"`
	Country   string  `json:"country"`
	Timezone  float64 `json:"timezone"`
}

type ConfigResponse struct {
	Tabs            []TabResponse                   `json:"tabs"`
	Elements        map[string]ElementMetaResponse  `json:"elements"`
	Cities          []CityResponse                  `json:"cities"`
	LifePathNumbers map[int]string                  `json:"lifePathNumbers"`
}

func main() {
	var err error
	baziCfg, err = loadBaziConfig("bazi.yaml")
	if err != nil {
		log.Fatalf("failed to load bazi.yaml: %v", err)
	}
	log.Printf("loaded bazi config: %d stems, %d branches, %d day masters, %d clash modifiers",
		len(baziCfg.HeavenlyStems), len(baziCfg.EarthlyBranches),
		len(baziCfg.DayMasterProfiles), len(baziCfg.ClashModifiers))

	mux := http.NewServeMux()
	mux.HandleFunc("/api/profile", handleProfile)
	mux.HandleFunc("/api/config", handleConfig)
	mux.HandleFunc("/api/cities", handleCities)
	log.Println("listening on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
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

func handleConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cfg := baziCfg
	if cfg == nil {
		http.Error(w, "config not loaded", http.StatusInternalServerError)
		return
	}

	tabs := make([]TabResponse, len(cfg.Tabs))
	for i, t := range cfg.Tabs {
		tabs[i] = TabResponse{Key: t.Key, Label: t.Label, Icon: t.Icon}
	}

	elements := make(map[string]ElementMetaResponse, len(cfg.Elements))
	for k, v := range cfg.Elements {
		elements[k] = ElementMetaResponse{Label: v.Label, Short: v.Short, Color: v.Color}
	}

	cities := make([]CityResponse, len(cfg.Cities))
	for i, c := range cfg.Cities {
		cities[i] = CityResponse{Name: c.Name, Country: c.Country, Timezone: c.Timezone}
	}

	resp := ConfigResponse{
		Tabs:            tabs,
		Elements:        elements,
		Cities:          cities,
		LifePathNumbers: cfg.LifePath,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("error encoding config response: %v", err)
	}
}

func handleCities(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cfg := baziCfg
	if cfg == nil {
		http.Error(w, "config not loaded", http.StatusInternalServerError)
		return
	}

	cities := make([]CityResponse, len(cfg.Cities))
	for i, c := range cfg.Cities {
		cities[i] = CityResponse{Name: c.Name, Country: c.Country, Timezone: c.Timezone}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(cities); err != nil {
		log.Printf("error encoding cities response: %v", err)
	}
}
