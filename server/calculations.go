package main

import (
	"math"
	"strconv"
	"strings"
)

var baziCfg *BaziConfig

func computeProfile(birthdate, birthtime, gender string) Profile {
	cfg := baziCfg
	if cfg == nil {
		return fallbackProfile()
	}

	parts := strings.Split(birthdate, "-")
	if len(parts) < 3 {
		return fallbackProfile()
	}

	year, _ := strconv.Atoi(parts[0])
	month, _ := strconv.Atoi(parts[1])
	day, _ := strconv.Atoi(parts[2])

	hour := 12
	min := 0
	if birthtime != "" {
		tp := strings.Split(birthtime, ":")
		if len(tp) >= 1 {
			hour, _ = strconv.Atoi(tp[0])
		}
		if len(tp) >= 2 {
			min, _ = strconv.Atoi(tp[1])
		}
	}

	if gender == "" {
		gender = "male"
	}

	jdn := midnightJDN(year, month, day)
	dayFrac := (float64(hour) + float64(min)/60.0) / 24.0
	birthJDN := jdn + dayFrac

	solarYear := year
	if birthJDN < solarTermJDN(year, 0) {
		solarYear = year - 1
	}
	yStem, yBranch := yearStemBranch(solarYear)

	mStem, mBranch := monthStemBranch(yStem, year, month, day, hour, min)

	dy, dm, dd := baziDayPillarDate(year, month, day, hour)
	dStem, dBranch := dayStemBranch(dy, dm, dd)

	hStem, hBranch := hourStemBranch(dStem, hour)

	dominant := cfg.StemElement[dStem]
	dayMasterName := cfg.StemName[dStem]

	lp := lifePathNumber(year, month, day)

	stems := []int{yStem, mStem, dStem, hStem}
	branches := []int{yBranch, mBranch, dBranch, hBranch}
	raw := countElements(cfg, stems, branches)

	total := 0
	for _, v := range raw {
		total += v
	}

	balance := make(map[string]int, 5)
	if total > 0 {
		allocated := 0
		for _, el := range cfg.ElementOrder {
			pct := int(math.Round(float64(raw[el]) * 100.0 / float64(total)))
			balance[el] = pct
			allocated += pct
		}
		balance[dominant] += 100 - allocated
	} else {
		balance = map[string]int{"wood": 20, "fire": 20, "earth": 20, "metal": 20, "water": 20}
	}

	dominantEl := dominant
	highest := 0
	for _, el := range cfg.ElementOrder {
		if balance[el] > highest {
			highest = balance[el]
			dominantEl = el
		}
	}

	pillars := map[string]Pillar{
		"year":  buildPillar(cfg, yStem, yBranch),
		"month": buildPillar(cfg, mStem, mBranch),
		"day":   buildPillar(cfg, dStem, dBranch),
		"hour":  buildPillar(cfg, hStem, hBranch),
	}

	dayMaster := DayMaster{
		Stem:    dayMasterName,
		Element: dominant,
		YinYang: stemYinYang(dStem),
	}

	seasonPeakEl := cfg.ElementOrder[cfg.SeasonPeak[mBranch]]
	season := SeasonInfo{
		Name:    cfg.SeasonName[mBranch],
		Branch:  cfg.BranchName[mBranch],
		Element: seasonPeakEl,
	}

	dmStrength, _ := determineDMStrength(cfg, dominant, seasonPeakEl, stems, branches)

	yongShen, favorable, unfavorable := determineYongShen(cfg, dominant, dmStrength)

	tenGods := buildTenGods(cfg, dStem, yStem, mStem, dStem, hStem)

	forward := luckDirection(yStem, gender)
	startAge := luckStartAge(birthJDN, mBranch, forward, solarYear)
	luckCycles := computeLuckCycles(cfg, mStem, mBranch, forward, startAge)

	dominantTG := dominantTenGod(tenGods)
	tabText := assembleTabText(cfg, dayMasterName)
	dmProfile := assembleDayMasterProfile(cfg, dayMasterName)
	tenGodInsights := assembleTenGodInsights(cfg, tenGods)
	textModifiers := findModifiers(cfg, dayMasterName, dmStrength, dominantEl, dominantTG, season.Name)

	lpText := ""
	if t, ok := cfg.LifePath[lp]; ok {
		lpText = t
	}

	return Profile{
		Dominant:            dominant,
		Lp:                  lp,
		Balance:             balance,
		Descriptions:        map[string]map[string]string{},
		Pillars:             pillars,
		DayMaster:           dayMaster,
		DmStrength:          dmStrength,
		Season:              season,
		YongShen:            yongShen,
		FavorableElements:   favorable,
		UnfavorableElements: unfavorable,
		TenGods:             tenGods,
		LuckCycles:          luckCycles,
		TabText:             tabText,
		DayMasterProfile:    dmProfile,
		TenGodInsights:      tenGodInsights,
		TextModifiers:       textModifiers,
		LifePathText:        lpText,
	}
}

func fallbackProfile() Profile {
	return Profile{
		Dominant:            "earth",
		Lp:                  0,
		Balance:             map[string]int{"wood": 20, "fire": 20, "earth": 20, "metal": 20, "water": 20},
		Descriptions:        map[string]map[string]string{},
		Pillars:             map[string]Pillar{},
		DayMaster:           DayMaster{},
		DmStrength:          "unknown",
		Season:              SeasonInfo{},
		YongShen:            YongShenInfo{},
		FavorableElements:   []string{},
		UnfavorableElements: []string{},
		TenGods:             map[string]TenGodEntry{},
		LuckCycles:          []LuckCycle{},
		TabText:             map[string]string{},
		TextModifiers:       []string{},
	}
}