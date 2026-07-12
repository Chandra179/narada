package main

import "math"

// =============================================================================
// Luck Cycles (大运 / Da Yun)
// =============================================================================

func luckDirection(yearStem int, gender string) bool {
	stemIsYang := yearStem%2 == 0
	male := gender == "male"
	return (stemIsYang && male) || (!stemIsYang && !male)
}

func luckStartAge(birthJDN float64, monthBranch int, forward bool, solarYear int) int {
	currentJie := jieTermForBranch(monthBranch)
	var targetJDN float64

	if forward {
		nextJie := (currentJie + 1) % 12
		targetJDN = solarTermJDN(solarYear, nextJie)
		if nextJie < currentJie {
			targetJDN = solarTermJDN(solarYear+1, nextJie)
		}
	} else {
		prevJie := (currentJie - 1 + 12) % 12
		targetJDN = solarTermJDN(solarYear, prevJie)
		if prevJie > currentJie {
			targetJDN = solarTermJDN(solarYear-1, prevJie)
		}
	}

	days := targetJDN - birthJDN
	if days < 0 {
		days = -days
	}

	age := int(math.Round(days / 3.0))
	if age < 0 {
		age = 0
	}
	if age > 10 {
		age = 10
	}
	return age
}

func computeLuckCycles(cfg *BaziConfig, monthStem, monthBranch int, forward bool, startAge int) []LuckCycle {
	cycles := make([]LuckCycle, 0, 8)
	step := 1
	if !forward {
		step = -1
	}

	stem := monthStem
	branch := monthBranch
	age := startAge

	for i := 0; i < 8; i++ {
		stem = (stem + step + 10) % 10
		branch = (branch + step + 12) % 12

		endAge := age + 9
		if endAge > 120 {
			endAge = 120
		}

		cycles = append(cycles, LuckCycle{
			StartAge: age,
			EndAge:   endAge,
			Stem:     cfg.StemName[stem],
			Branch:   cfg.BranchName[branch],
			Element:  cfg.StemElement[stem],
		})

		age = endAge + 1
		if age > 120 {
			break
		}
	}
	return cycles
}