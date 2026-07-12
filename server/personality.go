package main

// =============================================================================
// Personality Assembly
// =============================================================================

func dominantTenGod(tenGods map[string]TenGodEntry) string {
	for _, key := range []string{"monthStem", "yearStem", "hourStem"} {
		if tg, ok := tenGods[key]; ok && tg.God != "Day Master" {
			return tg.God
		}
	}
	return ""
}

func findModifiers(cfg *BaziConfig, dayMasterName, dmStrength, dominantEl, dominantTG, season string) []string {
	result := []string{}
	for _, m := range cfg.ClashModifiers {
		c := m.Condition

		if c.DayMaster != "" && c.DayMaster != dayMasterName {
			continue
		}
		if c.DominantElement != "" && c.DominantElement != dominantEl {
			continue
		}
		if c.DominantTenGod != "" && c.DominantTenGod != dominantTG {
			continue
		}
		if c.DmStrength != "" && c.DmStrength != dmStrength {
			continue
		}
		if c.Season != "" && c.Season != season {
			continue
		}

		result = append(result, m.Output.TextModifier)
	}
	return result
}

func assembleTabText(cfg *BaziConfig, dayMasterName string) map[string]string {
	profile, ok := cfg.DayMasterProfiles[dayMasterName]
	if !ok {
		return nil
	}
	text := make(map[string]string, len(profile.TabText))
	for k, v := range profile.TabText {
		text[k] = v
	}
	return text
}

func assembleDayMasterProfile(cfg *BaziConfig, dayMasterName string) *DayMasterProfileInfo {
	profile, ok := cfg.DayMasterProfiles[dayMasterName]
	if !ok {
		return nil
	}
	return &DayMasterProfileInfo{
		Archetype:  profile.Archetype,
		CoreTraits: profile.CoreTraits,
		Summary:    profile.Summary,
	}
}

func assembleTenGodInsights(cfg *BaziConfig, tenGods map[string]TenGodEntry) map[string]TenGodInsight {
	insights := make(map[string]TenGodInsight)
	for pillar, entry := range tenGods {
		tg, ok := cfg.TenGods[entry.God]
		if !ok {
			insights[pillar] = TenGodInsight{
				God:         entry.God,
				Description: "",
			}
			continue
		}
		insights[pillar] = TenGodInsight{
			God:         entry.God,
			Category:    tg.Category,
			Description: tg.Description,
		}
	}
	return insights
}