package main

import "strconv"

// =============================================================================
// Element Utilities
// =============================================================================

func stemYinYang(stem int) string {
	if stem%2 == 0 {
		return "yang"
	}
	return "yin"
}

func elementID(cfg *BaziConfig, el string) int {
	for i, e := range cfg.ElementOrder {
		if e == el {
			return i
		}
	}
	return -1
}

func productiveCycle(cfg *BaziConfig, el int) int {
	return cfg.ElementCycles.Productive[el]
}

func controllingCycle(cfg *BaziConfig, el int) int {
	return cfg.ElementCycles.Controlling[el]
}

func controlledBy(cfg *BaziConfig, el int) int {
	for i, target := range cfg.ElementCycles.Productive {
		if target == el {
			return i
		}
	}
	return (el + 2) % 5
}

// =============================================================================
// Ten Gods (十神)
// =============================================================================

func tenGod(cfg *BaziConfig, dayStem, otherStem int) string {
	dmEl := elementID(cfg, cfg.StemElement[dayStem])
	otherEl := elementID(cfg, cfg.StemElement[otherStem])
	dmYang := stemYinYang(dayStem)
	otherYang := stemYinYang(otherStem)
	samePolarity := dmYang == otherYang

	if dmEl == otherEl {
		if samePolarity {
			return "Peer"
		}
		return "Rob Wealth"
	}

	if productiveCycle(cfg, dmEl) == otherEl {
		if samePolarity {
			return "Eating God"
		}
		return "Hurting Officer"
	}

	if productiveCycle(cfg, otherEl) == dmEl {
		if samePolarity {
			return "Direct Resource"
		}
		return "Indirect Resource"
	}

	if controllingCycle(cfg, dmEl) == otherEl {
		if samePolarity {
			return "Direct Wealth"
		}
		return "Indirect Wealth"
	}

	if controllingCycle(cfg, otherEl) == dmEl {
		if samePolarity {
			return "Direct Officer"
		}
		return "Seven Kill"
	}

	return "Other"
}

// =============================================================================
// Seasonal Strength
// =============================================================================

func seasonalStrength(cfg *BaziConfig, seasonPeak, element string) int {
	peakID := elementID(cfg, seasonPeak)
	elID := elementID(cfg, element)

	if peakID == elID {
		return 4
	}
	if productiveCycle(cfg, peakID) == elID {
		return 3
	}
	if productiveCycle(cfg, elID) == peakID {
		return 2
	}
	if controllingCycle(cfg, elID) == peakID {
		return 1
	}
	return 0
}

// =============================================================================
// Day Master Strength Analysis
// =============================================================================

func determineDMStrength(cfg *BaziConfig, dmElement string, seasonPeak string, stems, branches []int) (string, string) {
	dmID := elementID(cfg, dmElement)

	seasonScore := seasonalStrength(cfg, seasonPeak, dmElement)

	allyCount := 0
	enemyCount := 0

	for _, s := range stems {
		elID := elementID(cfg, cfg.StemElement[s])
		if elID == dmID || productiveCycle(cfg, elID) == dmID {
			allyCount++
		} else if controllingCycle(cfg, elID) == dmID {
			enemyCount++
		}
	}
	for _, b := range branches {
		elID := elementID(cfg, cfg.BranchElement[b])
		if elID == dmID || productiveCycle(cfg, elID) == dmID {
			allyCount++
		} else if controllingCycle(cfg, elID) == dmID {
			enemyCount++
		}
		for _, hs := range cfg.BranchHidden[b] {
			hsID := elementID(cfg, cfg.StemElement[hs])
			if hsID == dmID || productiveCycle(cfg, hsID) == dmID {
				allyCount++
			} else if controllingCycle(cfg, hsID) == dmID {
				enemyCount++
			}
		}
	}

	strength := float64(seasonScore)*1.5 + float64(allyCount)*1.0 - float64(enemyCount)*0.8

	var verdict, reasoning string
	if strength >= 6.0 {
		verdict = "strong"
		reasoning = "The Day Master has strong seasonal support and multiple allies in the chart, making it robust and self-sufficient."
	} else if strength <= 3.0 {
		verdict = "weak"
		reasoning = "The Day Master lacks seasonal support and has few allies, with significant controlling elements present."
	} else {
		verdict = "balanced"
		reasoning = "The Day Master has moderate support — neither overwhelmingly strong nor critically weak."
	}
	return verdict, reasoning
}

// =============================================================================
// Yong Shen (Useful God)
// =============================================================================

func determineYongShen(cfg *BaziConfig, dmElement, dmStrength string) (YongShenInfo, []string, []string) {
	dmID := elementID(cfg, dmElement)

	var yongShenEl string
	var reason string

	switch dmStrength {
	case "strong":
		yongShenEl = cfg.ElementOrder[controllingCycle(cfg, dmID)]
		reason = "Strong " + dmElement + " DM needs " + yongShenEl + " to control and balance the excess energy."
	case "weak":
		yongShenEl = cfg.ElementOrder[controlledBy(cfg, dmID)]
		reason = "Weak " + dmElement + " DM needs " + yongShenEl + " to produce and nourish it."
	default:
		yongShenEl = cfg.ElementOrder[controlledBy(cfg, dmID)]
		reason = "Balanced " + dmElement + " DM benefits from " + yongShenEl + " to maintain harmony."
	}

	ysID := elementID(cfg, yongShenEl)

	var favorable []string
	var unfavorable []string

	favorable = append(favorable, yongShenEl)
	producerOfYS := cfg.ElementOrder[controlledBy(cfg, ysID)]
	if producerOfYS != yongShenEl {
		favorable = append(favorable, producerOfYS)
	}

	controllerOfYS := cfg.ElementOrder[controllingCycle(cfg, ysID)]
	unfavorable = append(unfavorable, controllerOfYS)
	drainOfYS := cfg.ElementOrder[productiveCycle(cfg, ysID)]
	if drainOfYS != controllerOfYS {
		unfavorable = append(unfavorable, drainOfYS)
	}

	return YongShenInfo{Element: yongShenEl, Reason: reason}, favorable, unfavorable
}

// =============================================================================
// Elemental Analysis
// =============================================================================

func countElements(cfg *BaziConfig, stems, branches []int) map[string]int {
	counts := make(map[string]int)
	for _, s := range stems {
		counts[cfg.StemElement[s]]++
	}
	for _, b := range branches {
		counts[cfg.BranchElement[b]]++
		for _, hs := range cfg.BranchHidden[b] {
			counts[cfg.StemElement[hs]]++
		}
	}
	return counts
}

func lifePathNumber(year, month, day int) int {
	sum := 0
	for _, c := range strconv.Itoa(year) {
		sum += int(c - '0')
	}
	for _, c := range strconv.Itoa(month) {
		sum += int(c - '0')
	}
	for _, c := range strconv.Itoa(day) {
		sum += int(c - '0')
	}
	for sum > 9 && sum != 11 && sum != 22 && sum != 33 {
		sum = sum/10 + sum%10
	}
	return sum
}