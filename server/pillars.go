package main

// =============================================================================
// Gregorian ↔ Julian Day Number
// =============================================================================

func gregorianToJDN(year, month, day int) int {
	a := (14 - month) / 12
	y := year + 4800 - a
	m := month + 12*a - 3
	return day + (153*m+2)/5 + 365*y + y/4 - y/100 + y/400 - 32045
}

func midnightJDN(year, month, day int) float64 {
	return float64(gregorianToJDN(year, month, day)) - 0.5
}

// =============================================================================
// Solar Term Computation (节气)
// =============================================================================

func solarTermJDN(solarYear, termIndex int) float64 {
	lichunMidnight := midnightJDN(solarYear, 2, 4)

	const referenceYear = 2000
	const refLichunFrac = 0.355
	driftPerYear := 365.242189 - 365.25

	lichunJDN := lichunMidnight + refLichunFrac + float64(solarYear-referenceYear)*driftPerYear

	offsets := [12]float64{0, 30, 60, 91, 122, 154, 185, 216, 247, 277, 307, 336}
	return lichunJDN + offsets[termIndex]
}

func solarMonthBranch(year, month, day, hour, min int) int {
	jdn := midnightJDN(year, month, day)
	dayFrac := (float64(hour) + float64(min)/60.0) / 24.0
	birthJDN := jdn + dayFrac

	solarYear := year
	if birthJDN < solarTermJDN(year, 0) {
		solarYear = year - 1
	}

	for i := 0; i < 12; i++ {
		branch := (i + 2) % 12

		start := solarTermJDN(solarYear, i)

		var end float64
		if i < 11 {
			end = solarTermJDN(solarYear, i+1)
		} else {
			end = solarTermJDN(solarYear+1, 0)
		}

		if birthJDN >= start && birthJDN < end {
			return branch
		}
	}

	return 2
}

func jieTermForBranch(branch int) int {
	return (branch - 2 + 12) % 12
}

// =============================================================================
// Year Pillar (年柱)
// =============================================================================

func yearStemBranch(adjustedYear int) (stem, branch int) {
	y := adjustedYear
	stem = ((y-4)%10 + 10) % 10
	branch = ((y-4)%12 + 12) % 12
	return
}

// =============================================================================
// Month Pillar (月柱)
// =============================================================================

func fiveTigersEscape(yearStem int) int {
	return (yearStem%5)*2 + 2
}

func monthStemBranch(yearStem, year, month, day, hour, min int) (stem, branch int) {
	branch = solarMonthBranch(year, month, day, hour, min)
	monthIndex := (branch - 2 + 12) % 12
	firstMonthStem := fiveTigersEscape(yearStem)
	stem = (firstMonthStem + monthIndex) % 10
	return
}

// =============================================================================
// Day Pillar (日柱)
// =============================================================================

func isLeapYear(y int) bool {
	return y%4 == 0 && (y%100 != 0 || y%400 == 0)
}

func daysFrom1900Jan1(year, month, day int) int {
	total := 0
	if year >= 1900 {
		for y := 1900; y < year; y++ {
			total += 365
			if isLeapYear(y) {
				total++
			}
		}
	} else {
		for y := year; y < 1900; y++ {
			total -= 365
			if isLeapYear(y) {
				total--
			}
		}
	}

	monthDays := [12]int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
	for m := 1; m < month; m++ {
		total += monthDays[m-1]
	}
	if month > 2 && isLeapYear(year) {
		total++
	}
	total += day - 1
	return total
}

func dayStemBranch(year, month, day int) (stem, branch int) {
	days := daysFrom1900Jan1(year, month, day)
	idx := (10 + days) % 60
	if idx < 0 {
		idx += 60
	}
	stem = idx % 10
	branch = idx % 12
	return
}

func baziDayPillarDate(year, month, day, hour int) (int, int, int) {
	if hour < 23 {
		return year, month, day
	}
	return addOneDay(year, month, day)
}

func addOneDay(year, month, day int) (int, int, int) {
	monthDays := [12]int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
	if isLeapYear(year) {
		monthDays[1] = 29
	}
	day++
	if day > monthDays[month-1] {
		day = 1
		month++
		if month > 12 {
			month = 1
			year++
		}
	}
	return year, month, day
}

// =============================================================================
// Hour Pillar (时柱)
// =============================================================================

func hourBranch(hour24 int) int {
	return (hour24 + 1) % 24 / 2
}

func fiveRatsEscape(dayStem int) int {
	return (dayStem % 5) * 2
}

func hourStemBranch(dayStem, hour24 int) (stem, branch int) {
	branch = hourBranch(hour24)
	ziStem := fiveRatsEscape(dayStem)
	stem = (ziStem + branch) % 10
	return
}

// =============================================================================
// Build enriched pillar data
// =============================================================================

func buildPillar(cfg *BaziConfig, stem, branch int) Pillar {
	hs := make([]string, len(cfg.BranchHidden[branch]))
	for i, s := range cfg.BranchHidden[branch] {
		hs[i] = cfg.StemName[s]
	}
	return Pillar{
		Stem:          cfg.StemName[stem],
		StemElement:   cfg.StemElement[stem],
		StemYinYang:   stemYinYang(stem),
		Branch:        cfg.BranchName[branch],
		BranchElement: cfg.BranchElement[branch],
		HiddenStems:   hs,
	}
}

func buildTenGods(cfg *BaziConfig, dayStem int, yStem, mStem, dStem, hStem int) map[string]TenGodEntry {
	return map[string]TenGodEntry{
		"yearStem":  {Element: cfg.StemElement[yStem], God: tenGod(cfg, dayStem, yStem)},
		"monthStem": {Element: cfg.StemElement[mStem], God: tenGod(cfg, dayStem, mStem)},
		"dayStem":   {Element: cfg.StemElement[dStem], God: "Day Master"},
		"hourStem":  {Element: cfg.StemElement[hStem], God: tenGod(cfg, dayStem, hStem)},
	}
}