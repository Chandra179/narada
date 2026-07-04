package main

import (
	"math"
	"strconv"
	"strings"
)

// =============================================================================
// Data tables — canonical Bazi / Wu Xing mappings
// Sources:
//   - Wikipedia "Heavenly Stems"  https://en.wikipedia.org/wiki/Heavenly_Stems
//   - Wikipedia "Four Pillars of Destiny"  https://en.wikipedia.org/wiki/Four_Pillars_of_Destiny
//   - Wikipedia "Sexagenary cycle"  https://en.wikipedia.org/wiki/Sexagenary_cycle
//   - Hong Kong Observatory "Heavenly Stems and Earthly Branches"
//     https://www.hko.gov.hk/en/gts/time/stemsandbranches.htm
//   - Ho, Peng Yoke (2003). Chinese Mathematical Astrology. Routledge.
//   - Yuan Hai Zi Ping (渊海子平), Song dynasty classical Bazi text by Xu Ziping
// =============================================================================

var heavenlyStems = [10]string{"Jia", "Yi", "Bing", "Ding", "Wu", "Ji", "Geng", "Xin", "Ren", "Gui"}
var stemElement = [10]string{"wood", "wood", "fire", "fire", "earth", "earth", "metal", "metal", "water", "water"}
var earthlyBranches = [12]string{"Zi", "Chou", "Yin", "Mao", "Chen", "Si", "Wu", "Wei", "Shen", "You", "Xu", "Hai"}
var branchElement = [12]string{"water", "earth", "wood", "wood", "earth", "fire", "fire", "earth", "metal", "metal", "earth", "water"}

var hiddenStems = [12][]int{
	{9},          // Zi  子 → Gui    癸
	{5, 9, 7},    // Chou 丑 → Ji, Gui, Xin      己 癸 辛
	{0, 2, 4},    // Yin  寅 → Jia, Bing, Wu     甲 丙 戊
	{1},          // Mao  卯 → Yi                 乙
	{4, 1, 9},    // Chen 辰 → Wu, Yi, Gui       戊 乙 癸
	{2, 4, 6},    // Si   巳 → Bing, Wu, Geng    丙 戊 庚
	{3, 5},       // Wu   午 → Ding, Ji          丁 己
	{5, 3, 1},    // Wei  未 → Ji, Ding, Yi      己 丁 乙
	{6, 4, 8},    // Shen 申 → Geng, Wu, Ren     庚 戊 壬
	{7},          // You  酉 → Xin                辛
	{4, 7, 3},    // Xu   戌 → Wu, Xin, Ding     戊 辛 丁
	{8, 0},       // Hai  亥 → Ren, Jia          壬 甲
}

// =============================================================================
// Gregorian ↔ Julian Day Number
// =============================================================================

// gregorianToJDN converts a Gregorian date to Julian Day Number at noon UTC.
// Algorithm from https://en.wikipedia.org/wiki/Julian_day
//
// Returns the JDN at 12:00 UTC.  To get midnight-based JDN, subtract 0.5.
func gregorianToJDN(year, month, day int) int {
	a := (14 - month) / 12
	y := year + 4800 - a
	m := month + 12*a - 3
	return day + (153*m+2)/5 + 365*y + y/4 - y/100 + y/400 - 32045
}

// midnightJDN returns the Julian Day Number at 00:00 UTC of the given date.
func midnightJDN(year, month, day int) float64 {
	return float64(gregorianToJDN(year, month, day)) - 0.5
}

// =============================================================================
// Solar Term Computation (节气)
// =============================================================================

// solarTermJDN returns the approximate Julian Day Number of the given jie
// solar term for a given solar year.
//
// termIndex: 0=Lichun(立春), 1=Jingzhe(惊蛰), …, 11=Xiaohan(小寒).
//
// The 12 jie terms partition the solar year starting at Lichun. Terms 0–10
// fall within the same calendar year; term 11 (Xiaohan) falls in January of
// the following calendar year.
//
// Uses exact Gregorian calendar dates (handling leap years) plus a small
// sub-day correction for the tropical-year drift.  Accuracy is ~6 hours for
// years 1800–2200.
//
// Source: https://en.wikipedia.org/wiki/Solar_term
func solarTermJDN(solarYear, termIndex int) float64 {
	// Midnight JDN of Feb 4 in the given calendar year.
	lichunMidnight := midnightJDN(solarYear, 2, 4)

	// Lichun 2000 occurred at ~20:32 UTC (0.355 past the midnight of Feb 4).
	// The tropical year (365.242189 d) is shorter than the calendar year
	// (365.25 d mean), so Lichun shifts ~11 min earlier each year.
	const referenceYear = 2000
	const refLichunFrac = 0.355       // fraction of day past midnight
	driftPerYear := 365.242189 - 365.25 // ≈ -0.007811 d/yr

	lichunJDN := lichunMidnight + refLichunFrac + float64(solarYear-referenceYear)*driftPerYear

	// Offsets from Lichun to each subsequent jie term in year 2000.
	offsets := [12]float64{0, 30, 60, 91, 122, 154, 185, 216, 247, 277, 307, 336}
	return lichunJDN + offsets[termIndex]
}

// solarMonthBranch returns the Earthly Branch index (0=子 .. 11=亥) of the
// solar month containing the given birth date and time.
//
// This is the central function that maps a Gregorian birth moment to a Bazi
// month-branch.  It first determines the correct solar year by comparing
// against the Lichun boundary, then scans the 12 jie-term intervals.
func solarMonthBranch(year, month, day, hour, min int) int {
	jdn := midnightJDN(year, month, day)
	dayFrac := (float64(hour) + float64(min)/60.0) / 24.0
	birthJDN := jdn + dayFrac

	// Determine the solar year: if before this calendar year's Lichun,
	// the birth is still part of the previous solar year.
	solarYear := year
	if birthJDN < solarTermJDN(year, 0) {
		solarYear = year - 1
	}

	// Scan the 12 jie-term intervals that partition one solar year.
	for i := 0; i < 12; i++ {
		branch := (i + 2) % 12 // Lichun (i=0) → Yin (寅 = 2)

		start := solarTermJDN(solarYear, i)

		var end float64
		if i < 11 {
			end = solarTermJDN(solarYear, i+1)
		} else {
			end = solarTermJDN(solarYear+1, 0) // next year's Lichun
		}

		if birthJDN >= start && birthJDN < end {
			return branch
		}
	}

	return 2 // fallback: Yin
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

// fiveTigersEscape (五虎遁) returns the Heavenly Stem index of the first solar
// month (寅月) based on the year's Heavenly Stem.
//
//	Year Stem   | 1st Month Stem
//	甲/己 Jia/Ji   → 丙 Bing (2)
//	乙/庚 Yi/Geng  → 戊 Wu  (4)
//	丙/辛 Bing/Xin → 庚 Geng (6)
//	丁/壬 Ding/Ren → 壬 Ren  (8)
//	戊/癸 Wu/Gui   → 甲 Jia  (0)
//
// Source: HKO Table 4, Yuan Hai Zi Ping / 渊海子平
func fiveTigersEscape(yearStem int) int {
	return (yearStem%5)*2 + 2
}

// monthStemBranch computes the Heavenly Stem and Earthly Branch for the month
// pillar.  The branch uses JDN-based solar-term lookup; the stem uses the
// Five Tigers Escape formula.
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

// dayStemBranch computes the Heavenly Stem and Earthly Branch for the day
// pillar.  Uses the unbroken sexagenary day count with reference
// Jan 1, 1900 = 甲戌 (Jia-Xu, index 10).
//
// Source: https://en.wikipedia.org/wiki/Sexagenary_cycle#Sexagenary_days
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

// baziDayPillarDate returns the effective date for the Bazi Day Pillar.
//
// In Bazi the day begins at 子时 (23:00), not midnight.  When the birth
// hour is ≥ 23, the Day Pillar rolls forward to the next calendar date while
// the Hour Pillar still uses branch 子 (0) for 23:00–00:59.
//
// Source: https://www.hko.gov.hk/en/gts/time/stemsandbranches.htm Table 3
func baziDayPillarDate(year, month, day, hour int) (int, int, int) {
	if hour < 23 {
		return year, month, day
	}
	return addOneDay(year, month, day)
}

// addOneDay returns the Gregorian date one day after (year, month, day).
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

// fiveRatsEscape (五鼠遁) returns the Heavenly Stem index of 子时 based on
// the day's Heavenly Stem.
//
//	Day Stem     → Zi hour Stem
//	甲/己 Jia/Ji   → 甲 Jia (0)
//	乙/庚 Yi/Geng  → 丙 Bing (2)
//	丙/辛 Bing/Xin → 戊 Wu (4)
//	丁/壬 Ding/Ren → 庚 Geng (6)
//	戊/癸 Wu/Gui   → 壬 Ren (8)
//
// Source: HKO Table 5, Yuan Hai Zi Ping / 渊海子平
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
// Elemental Analysis
// =============================================================================

func elementFromYear(year int) string {
	return stemElement[year%10]
}

func countElements(yearStem, yearBranch, monthStem, monthBranch, dayStem, dayBranch, hourStem, hourBranch int) map[string]int {
	counts := map[string]int{
		"wood": 0, "fire": 0, "earth": 0, "metal": 0, "water": 0,
	}

	stems := [4]int{yearStem, monthStem, dayStem, hourStem}
	branches := [4]int{yearBranch, monthBranch, dayBranch, hourBranch}

	for _, s := range stems {
		counts[stemElement[s]]++
	}
	for _, b := range branches {
		counts[branchElement[b]]++
		for _, hs := range hiddenStems[b] {
			counts[stemElement[hs]]++
		}
	}
	return counts
}

func lifePathNumber(dateStr string) int {
	sum := 0
	for _, ch := range dateStr {
		if ch >= '0' && ch <= '9' {
			sum += int(ch - '0')
		}
	}
	for sum > 9 && sum != 11 && sum != 22 && sum != 33 {
		n := 0
		for s := sum; s > 0; s /= 10 {
			n += s % 10
		}
		sum = n
	}
	return sum
}

// =============================================================================
// Element Descriptions
// =============================================================================

func elementDescriptions() Descriptions {
	return Descriptions{
		"wood": {
			"romance": "You love the way a tree loves — reaching outward, always growing the bond rather than settling into it. Partners feel your steady, expansive warmth, but you can outgrow connections that stop stretching with you.",
			"health":  "Vitality flows through your tendons and decision-making organs. Movement — real movement, not just intention — keeps your energy from turning brittle. Stagnation is your body's least favorite word.",
			"career":  "You thrive where there's room to expand: building, planting, initiating. Rigid hierarchies frustrate you; you do your best work when you're allowed to grow the role, not just fill it.",
			"wealth":  "Your money grows the way you do — slowly, then suddenly. You're better at cultivating long-term value than chasing quick wins, though impatience can tempt you to pull up the roots too early.",
		},
		"fire": {
			"romance": "You love loudly and visibly. Chemistry matters enormously to you, and you'd rather burn brightly for a short season than dim yourself for stability. Partners either match your heat or feel scorched by it.",
			"health":  "Your circulatory system and heart carry your signature. Passion fuels you, but unchecked intensity burns through your reserves fast — rest isn't optional, it's fuel.",
			"career":  "You're the spark in the room: charismatic, quick, magnetic to opportunity. You do best in visible, expressive roles, and worst in ones that ask you to stay quiet and wait your turn.",
			"wealth":  "Money moves fast around you — earned in bursts, spent with flair. Building lasting wealth means pairing your instinct for opportunity with someone (or something) that can bank the embers.",
		},
		"earth": {
			"romance": "You love like ground you can build a house on. Steady, dependable, unglamorous in the best way — you're the partner people come home to, not just the one they chase.",
			"health":  "Your digestive center is your barometer. When you're overextended caring for everyone else, it shows up first in your gut. Nourishment, including your own, is not indulgence.",
			"career":  "You're the stabilizer — the one who makes plans actually work. Reliable execution is your edge; just watch for being handed everyone else's unfinished business because you never say no.",
			"wealth":  "You save before you spend and plan before you leap. Your wealth grows through patience and diversification rather than bold bets — steady compounding is your natural mode.",
		},
		"metal": {
			"romance": "You love with precision — you notice details others miss and hold high standards for what a relationship should be. That clarity is a gift, but it can read as distance if you don't name what you feel.",
			"health":  "Your lungs and skin carry your tension first. Structure and clean air calm you; chaos and clutter — physical or emotional — drain you faster than almost anything else.",
			"career":  "You bring order to disorder. Editing, refining, systems, quality control — you make things sharper. You do poorly in loose, undefined roles with no clear standard to meet.",
			"wealth":  "You're disciplined with money to the point of austerity sometimes. You rarely overspend, and your instinct to cut what isn't working serves your portfolio as well as your closet.",
		},
		"water": {
			"romance": "You love the way water moves — adapting to whoever you're with, reading what they need before they say it. That flexibility draws people in, but you can lose your own shape if you're not careful.",
			"health":  "Your kidneys and adrenal reserves are your tell. You run on deep reserves rather than quick bursts, which means burnout, when it comes, comes from the bottom up. Protect your sleep.",
			"career":  "You're the strategist — comfortable with ambiguity, good at finding the path of least resistance to a goal. You underperform in rigid, script-driven roles that don't let you improvise.",
			"wealth":  "Your money finds unconventional channels — you're drawn to opportunities others overlook. The risk is spreading too thin; depth in a few currents beats width across all of them.",
		},
	}
}

// =============================================================================
// Main API
// =============================================================================

// computeProfile takes a birthdate (YYYY-MM-DD) and optional birthtime
// (HH:MM, defaults to "12:00") and returns a fully computed Bazi profile.
func computeProfile(birthdate, birthtime string) Profile {
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

	// --- Year Pillar (Lichun boundary via JDN comparison) ---
	jdn := midnightJDN(year, month, day)
	dayFrac := (float64(hour) + float64(min)/60.0) / 24.0
	birthJDN := jdn + dayFrac

	solarYear := year
	if birthJDN < solarTermJDN(year, 0) {
		solarYear = year - 1
	}
	yStem, yBranch := yearStemBranch(solarYear)

	// --- Month Pillar ---
	mStem, mBranch := monthStemBranch(yStem, year, month, day, hour, min)

	// --- Day Pillar (23:00 boundary) ---
	dy, dm, dd := baziDayPillarDate(year, month, day, hour)
	dStem, dBranch := dayStemBranch(dy, dm, dd)

	// --- Hour Pillar ---
	hStem, hBranch := hourStemBranch(dStem, hour)

	// --- Day Master (日主) ---
	dominant := stemElement[dStem]

	lp := lifePathNumber(birthdate)

	// --- Elemental Balance ---
	raw := countElements(yStem, yBranch, mStem, mBranch, dStem, dBranch, hStem, hBranch)

	total := 0
	for _, v := range raw {
		total += v
	}

	balance := make(map[string]int, 5)
	if total > 0 {
		allocated := 0
		elements := []string{"wood", "fire", "earth", "metal", "water"}
		for _, el := range elements {
			pct := int(math.Round(float64(raw[el]) * 100.0 / float64(total)))
			balance[el] = pct
			allocated += pct
		}
		balance[dominant] += 100 - allocated
	} else {
		balance = map[string]int{"wood": 20, "fire": 20, "earth": 20, "metal": 20, "water": 20}
	}

	return Profile{
		Dominant:     dominant,
		Lp:           lp,
		Balance:      balance,
		Descriptions: elementDescriptions(),
	}
}

func fallbackProfile() Profile {
	return Profile{
		Dominant:     "earth",
		Lp:           0,
		Balance:      map[string]int{"wood": 20, "fire": 20, "earth": 20, "metal": 20, "water": 20},
		Descriptions: elementDescriptions(),
	}
}
