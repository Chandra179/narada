package main

import (
	"math"
	"strconv"
	"strings"
)

// =============================================================================
// Data tables — canonical Bazi / Wu Xing mappings
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

var elementOrder = []string{"wood", "fire", "earth", "metal", "water"}

var seasonTable = [12]struct {
	season string
	peak   string
}{
	{season: "winter", peak: "water"},     // Zi
	{season: "winter", peak: "water"},     // Chou
	{season: "spring", peak: "wood"},      // Yin
	{season: "spring", peak: "wood"},      // Mao
	{season: "spring", peak: "wood"},      // Chen
	{season: "summer", peak: "fire"},      // Si
	{season: "summer", peak: "fire"},      // Wu
	{season: "summer", peak: "fire"},      // Wei
	{season: "autumn", peak: "metal"},     // Shen
	{season: "autumn", peak: "metal"},     // You
	{season: "autumn", peak: "metal"},     // Xu
	{season: "winter", peak: "water"},     // Hai
}

// jieTermForBranch maps each Earthly Branch to its starting jie solar term index.
// branch 2 (Yin) starts at Lichun (index 0), etc.
func jieTermForBranch(branch int) int {
	return (branch - 2 + 12) % 12
}

// =============================================================================
// Utilities
// =============================================================================

func stemYinYang(stem int) string {
	if stem%2 == 0 {
		return "yang"
	}
	return "yin"
}

func elementID(el string) int {
	for i, e := range elementOrder {
		if e == el {
			return i
		}
	}
	return -1
}

// productiveCycle returns the element produced by the given element index.
// Wood→Fire→Earth→Metal→Water→Wood
func productiveCycle(el int) int {
	return (el + 1) % 5
}

// controllingCycle returns the element controlled by the given element index.
// Wood→Earth→Water→Fire→Metal→Wood
func controllingCycle(el int) int {
	return (el + 2) % 5
}

// controlledBy returns the element that controls the given element.
func controlledBy(el int) int {
	return (el + 3) % 5
}

// =============================================================================
// Ten Gods (十神)
// =============================================================================

// tenGod returns the Ten God name for a stem relative to the Day Master stem.
func tenGod(dayStem, otherStem int) string {
	dmEl := elementID(stemElement[dayStem])
	otherEl := elementID(stemElement[otherStem])
	dmYang := stemYinYang(dayStem)
	otherYang := stemYinYang(otherStem)
	samePolarity := dmYang == otherYang

	if dmEl == otherEl {
		if samePolarity {
			return "Peer"
		}
		return "Rob Wealth"
	}

	if productiveCycle(dmEl) == otherEl {
		if samePolarity {
			return "Eating God"
		}
		return "Hurting Officer"
	}

	if productiveCycle(otherEl) == dmEl {
		if samePolarity {
			return "Direct Resource"
		}
		return "Indirect Resource"
	}

	if controllingCycle(dmEl) == otherEl {
		if samePolarity {
			return "Direct Wealth"
		}
		return "Indirect Wealth"
	}

	if controllingCycle(otherEl) == dmEl {
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

// seasonalStrength returns a score from 0-4 for how strong an element is in a given season.
// 4 = Peak (旺), 3 = Rising (相), 2 = Flat (休), 1 = Trapped (囚), 0 = Dying (死)
func seasonalStrength(seasonPeak, element string) int {
	peakID := elementID(seasonPeak)
	elID := elementID(element)

	// Same element → Peak
	if peakID == elID {
		return 4
	}
	// Produced by peak → Rising (peak produces this)
	if productiveCycle(peakID) == elID {
		return 3
	}
	// Produces the peak → Flat (this produces peak)
	if productiveCycle(elID) == peakID {
		return 2
	}
	// Controls the peak → Trapped (this controls peak but is out of season)
	if controllingCycle(elID) == peakID {
		return 1
	}
	// Controlled by peak → Dying (peak controls this)
	return 0
}

// =============================================================================
// Day Master Strength Analysis
// =============================================================================

func determineDMStrength(dmElement string, seasonPeak string, stems, branches []int) (string, string) {
	dmID := elementID(dmElement)

	seasonScore := seasonalStrength(seasonPeak, dmElement)

	allyCount := 0
	enemyCount := 0

	for _, s := range stems {
		elID := elementID(stemElement[s])
		if elID == dmID || productiveCycle(elID) == dmID {
			allyCount++
		} else if controllingCycle(elID) == dmID {
			enemyCount++
		}
	}
	for _, b := range branches {
		elID := elementID(branchElement[b])
		if elID == dmID || productiveCycle(elID) == dmID {
			allyCount++
		} else if controllingCycle(elID) == dmID {
			enemyCount++
		}
		for _, hs := range hiddenStems[b] {
			hsID := elementID(stemElement[hs])
			if hsID == dmID || productiveCycle(hsID) == dmID {
				allyCount++
			} else if controllingCycle(hsID) == dmID {
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

func determineYongShen(dmElement, dmStrength string) (YongShenInfo, []string, []string) {
	dmID := elementID(dmElement)

	var yongShenEl string
	var reason string

	switch dmStrength {
	case "strong":
		// Controlling element or output element
		yongShenEl = elementOrder[controllingCycle(dmID)]
		reason = "Strong " + dmElement + " DM needs " + yongShenEl + " to control and balance the excess energy."
	case "weak":
		// Producing element (resource)
		yongShenEl = elementOrder[controlledBy(dmID)]
		reason = "Weak " + dmElement + " DM needs " + yongShenEl + " to produce and nourish it."
	default:
		// Balanced — Yong Shen is the element that produces the DM
		yongShenEl = elementOrder[controlledBy(dmID)]
		reason = "Balanced " + dmElement + " DM benefits from " + yongShenEl + " to maintain harmony."
	}

	ysID := elementID(yongShenEl)

	var favorable []string
	var unfavorable []string

	// Favorable: Yong Shen itself + elements that produce it
	favorable = append(favorable, yongShenEl)
	producerOfYS := elementOrder[controlledBy(ysID)]
	if producerOfYS != yongShenEl {
		favorable = append(favorable, producerOfYS)
	}

	// Unfavorable: elements that control or drain the Yong Shen
	controllerOfYS := elementOrder[controllingCycle(ysID)]
	unfavorable = append(unfavorable, controllerOfYS)
	drainOfYS := elementOrder[productiveCycle(ysID)]
	if drainOfYS != controllerOfYS {
		unfavorable = append(unfavorable, drainOfYS)
	}

	return YongShenInfo{Element: yongShenEl, Reason: reason}, favorable, unfavorable
}

// =============================================================================
// Luck Cycles (大运 / Da Yun)
// =============================================================================

// luckDirection returns true if Luck Pillars move forward through the sexagenary cycle.
// Yang male + Yin female → forward (true)
// Yin male + Yang female → backward (false)
func luckDirection(yearStem int, gender string) bool {
	stemIsYang := yearStem%2 == 0
	male := gender == "male"
	return (stemIsYang && male) || (!stemIsYang && !male)
}

// luckStartAge computes the starting age for the first Luck Pillar.
// Days between birth and next/previous jie term, divided by 3.
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

func computeLuckCycles(monthStem, monthBranch int, forward bool, startAge int) []LuckCycle {
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
			Stem:     heavenlyStems[stem],
			Branch:   earthlyBranches[branch],
			Element:  stemElement[stem],
		})

		age = endAge + 1
		if age > 120 {
			break
		}
	}
	return cycles
}

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
// Build enriched pillar data
// =============================================================================

func buildPillar(stem, branch int) Pillar {
	hs := make([]string, len(hiddenStems[branch]))
	for i, s := range hiddenStems[branch] {
		hs[i] = heavenlyStems[s]
	}
	return Pillar{
		Stem:          heavenlyStems[stem],
		StemElement:   stemElement[stem],
		StemYinYang:   stemYinYang(stem),
		Branch:        earthlyBranches[branch],
		BranchElement: branchElement[branch],
		HiddenStems:   hs,
	}
}

func buildTenGods(dayStem int, yStem, mStem, dStem, hStem int) map[string]TenGodEntry {
	return map[string]TenGodEntry{
		"yearStem":  {Element: stemElement[yStem], God: tenGod(dayStem, yStem)},
		"monthStem": {Element: stemElement[mStem], God: tenGod(dayStem, mStem)},
		"dayStem":   {Element: stemElement[dStem], God: "Day Master"},
		"hourStem":  {Element: stemElement[hStem], God: tenGod(dayStem, hStem)},
	}
}

// =============================================================================
// Main API
// =============================================================================

func computeProfile(birthdate, birthtime, gender string) Profile {
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

	// --- Life Path Number ---
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
		for _, el := range elementOrder {
			pct := int(math.Round(float64(raw[el]) * 100.0 / float64(total)))
			balance[el] = pct
			allocated += pct
		}
		balance[dominant] += 100 - allocated
	} else {
		balance = map[string]int{"wood": 20, "fire": 20, "earth": 20, "metal": 20, "water": 20}
	}

	// --- Pillar data ---
	pillars := map[string]Pillar{
		"year":  buildPillar(yStem, yBranch),
		"month": buildPillar(mStem, mBranch),
		"day":   buildPillar(dStem, dBranch),
		"hour":  buildPillar(hStem, hBranch),
	}

	// --- Day Master info ---
	dayMaster := DayMaster{
		Stem:    heavenlyStems[dStem],
		Element: dominant,
		YinYang: stemYinYang(dStem),
	}

	// --- Season ---
	seasonInfo := seasonTable[mBranch]
	season := SeasonInfo{
		Name:    seasonInfo.season,
		Branch:  earthlyBranches[mBranch],
		Element: seasonInfo.peak,
	}

	// --- Day Master Strength ---
	stems := []int{yStem, mStem, dStem, hStem}
	branches := []int{yBranch, mBranch, dBranch, hBranch}
	dmStrength, _ := determineDMStrength(dominant, seasonInfo.peak, stems, branches)

	// --- Yong Shen ---
	yongShen, favorable, unfavorable := determineYongShen(dominant, dmStrength)

	// --- Ten Gods ---
	tenGods := buildTenGods(dStem, yStem, mStem, dStem, hStem)

	// --- Luck Cycles ---
	forward := luckDirection(yStem, gender)
	startAge := luckStartAge(birthJDN, mBranch, forward, solarYear)
	luckCycles := computeLuckCycles(mStem, mBranch, forward, startAge)

	return Profile{
		Dominant:           dominant,
		Lp:                 lp,
		Balance:            balance,
		Descriptions:       elementDescriptions(),
		Pillars:            pillars,
		DayMaster:          dayMaster,
		DmStrength:         dmStrength,
		Season:             season,
		YongShen:           yongShen,
		FavorableElements:  favorable,
		UnfavorableElements: unfavorable,
		TenGods:            tenGods,
		LuckCycles:         luckCycles,
	}
}

func fallbackProfile() Profile {
	return Profile{
		Dominant:           "earth",
		Lp:                 0,
		Balance:            map[string]int{"wood": 20, "fire": 20, "earth": 20, "metal": 20, "water": 20},
		Descriptions:       elementDescriptions(),
		Pillars:            map[string]Pillar{},
		DayMaster:          DayMaster{},
		DmStrength:         "unknown",
		Season:             SeasonInfo{},
		YongShen:           YongShenInfo{},
		FavorableElements:  []string{},
		UnfavorableElements: []string{},
		TenGods:            map[string]TenGodEntry{},
		LuckCycles:         []LuckCycle{},
	}
}
