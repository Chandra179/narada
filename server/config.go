package main

import (
	"os"

	"gopkg.in/yaml.v3"
)

// =============================================================================
// YAML Config Structs
// =============================================================================

type TabDef struct {
	Key   string `yaml:"key"`
	Label string `yaml:"label"`
	Icon  string `yaml:"icon"`
}

type ElementMeta struct {
	Label string `yaml:"label"`
	Short string `yaml:"short"`
	Color string `yaml:"color"`
}

type StemDef struct {
	Index    int    `yaml:"index"`
	Name     string `yaml:"name"`
	Element  string `yaml:"element"`
	Polarity string `yaml:"polarity"`
}

type BranchDef struct {
	Index       int    `yaml:"index"`
	Name        string `yaml:"name"`
	Animal      string `yaml:"animal"`
	Element     string `yaml:"element"`
	Polarity    string `yaml:"polarity"`
	HiddenStems []int  `yaml:"hidden_stems"`
}

type ElementCyclesDef struct {
	Productive  []int `yaml:"productive"`
	Controlling []int `yaml:"controlling"`
}

type SeasonDef struct {
	Name        string `yaml:"name"`
	PeakElement string `yaml:"peak_element"`
	Branches    []int  `yaml:"branches"`
}

type CityDef struct {
	Name     string  `yaml:"name"`
	Country  string  `yaml:"country"`
	Timezone float64 `yaml:"timezone"`
}

type TabText map[string]string

type DayMasterProfile struct {
	Archetype  string   `yaml:"archetype"`
	CoreTraits []string `yaml:"core_traits"`
	Summary    string   `yaml:"summary"`
	TabText    TabText  `yaml:"tab_text"`
}

type TenGodDef struct {
	Category      string   `yaml:"category"`
	Drive         string   `yaml:"drive"`
	PositiveTraits []string `yaml:"positive_traits"`
	NegativeTraits []string `yaml:"negative_traits"`
	Description   string   `yaml:"description"`
}

type ModifierCondition struct {
	DayMaster       string `yaml:"day_master,omitempty"`
	DominantElement string `yaml:"dominant_element,omitempty"`
	DominantTenGod  string `yaml:"dominant_ten_god,omitempty"`
	DmStrength      string `yaml:"dm_strength,omitempty"`
	Season          string `yaml:"season,omitempty"`
}

type ModifierOutput struct {
	TextModifier string `yaml:"text_modifier"`
}

type ClashModifier struct {
	ID        string           `yaml:"id"`
	Condition ModifierCondition `yaml:"condition"`
	Output    ModifierOutput   `yaml:"output"`
}

type BaziConfig struct {
	Tabs            []TabDef                   `yaml:"tabs"`
	Elements        map[string]ElementMeta     `yaml:"elements"`
	HeavenlyStems   []StemDef                  `yaml:"heavenly_stems"`
	EarthlyBranches []BranchDef                `yaml:"earthly_branches"`
	ElementCycles   ElementCyclesDef           `yaml:"element_cycles"`
	Seasons         []SeasonDef                `yaml:"seasons"`
	Cities          []CityDef                  `yaml:"cities"`
	LifePath        map[int]string             `yaml:"life_path_numbers"`
	DayMasterProfiles map[string]DayMasterProfile `yaml:"day_master_profiles"`
	TenGods         map[string]TenGodDef       `yaml:"ten_gods"`
	ClashModifiers  []ClashModifier            `yaml:"clash_modifiers"`

	// --- Precomputed lookup tables (set after Load) ---
	StemName      [10]string
	StemElement   [10]string
	BranchName    [12]string
	BranchElement [12]string
	BranchHidden  [12][]int

	// Element order (canonical): wood=0, fire=1, earth=2, metal=3, water=4
	ElementOrder []string

	// Season peak element index per branch (0-11)
	SeasonPeak [12]int // element ID (0-4)

	// Season name per branch
	SeasonName [12]string
}

// loadBaziConfig reads and parses the YAML config file.
func loadBaziConfig(path string) (*BaziConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg BaziConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	// Canonical element order
	cfg.ElementOrder = []string{"wood", "fire", "earth", "metal", "water"}

	// Index stems by position (they should already be in index order)
	for _, s := range cfg.HeavenlyStems {
		if s.Index >= 0 && s.Index < 10 {
			cfg.StemName[s.Index] = s.Name
			cfg.StemElement[s.Index] = s.Element
		}
	}

	// Index branches by position
	for _, b := range cfg.EarthlyBranches {
		if b.Index >= 0 && b.Index < 12 {
			cfg.BranchName[b.Index] = b.Name
			cfg.BranchElement[b.Index] = b.Element
			cfg.BranchHidden[b.Index] = b.HiddenStems
		}
	}

	// Build season lookup tables (per branch)
	for _, s := range cfg.Seasons {
		peakID := elementID(&cfg, s.PeakElement)
		for _, b := range s.Branches {
			if b >= 0 && b < 12 {
				cfg.SeasonPeak[b] = peakID
				cfg.SeasonName[b] = s.Name
			}
		}
	}

	return &cfg, nil
}
