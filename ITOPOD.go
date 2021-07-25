package main

import (
	"fmt"
	"log"
	"math"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const MIN_ITOPOD_TIER = 1
const PP_BASE = 200

var baseEXP = map[int]int{
	1:  1,
	2:  2,
	3:  4,
	4:  8,
	5:  14,
	6:  22,
	7:  32,
	8:  44,
	9:  58,
	10: 74,
	11: 92,
	12: 112,
	13: 134,
	14: 158,
	15: 184,
	16: 212,
	17: 242,
	18: 274,
	19: 308,
	20: 344,
}

func IdleITOPOD(dur time.Duration) {

	expBonus := 1.0
	ppBonus := 1.0

	defer writeStatusLog()

	clickCheckWait(MenuAdventure)

	// Turn off idling if it was on
	if checkColor(IdleModeOn, false) {
		clickCheckWait(IdleModeAbility)
	}

	// Do the hover OCR thing to get current kill counts
	floorMap, killMap := parseKillMap()

	killingStarted := time.Now()
	var kills int
	var ap int
	var exp int
	var pp int

	var currentTier int
	targetTier := pickTier(killMap, 0)
	for {
		// TimeCheck
		loopEnd := killingStarted.Add(dur)
		if time.Now().After(loopEnd) {
			log.Println("Loop Completed")
			return
		}
		targetTier = pickTier(killMap, targetTier)
		if targetTier != currentTier {
			//fmt.Println("\nTransitioning to tier", targetTier, "floor", floorMap[targetTier], "kills to go", killMap[targetTier])
			currentTier = targetTier
			clickCheckWait(EnterITOPOD) // open tower interface
			for {
				if checkColor(ITOPODBoxOpen, false) {
					break
				}
			}
			clickCheckWait(ITOPODStartBox) // click into the box
			typeWait(fmt.Sprintf("%d", floorMap[targetTier]))
			clickCheckWait(ITOPODEngage) // trigger the level
		}
		// remove the tool tip so we can see the level
		clickCheckWait(ITOPODEngage)
		// Wait for enemy to spawn
		for {
			if checkColor(EnemyHealth, false) {
				break
			}
		}
		matches := clickCheckValidate(RegAttackUnused, false)
		if matches {
			fmt.Printf("\nAttack did not go through\n")
			clickCheckWait(RegAttack)
		}
		if killMap[currentTier] == 1 {
			ap++
			expGained := int(math.Floor(float64(baseEXP[targetTier]) * expBonus))
			exp += expGained
			fmt.Println("")
			//fmt.Println("Gained", expGained, "exp and 1 AP")
		}
		killMap = updateKillMap(killMap)
		ppGained := int(math.Floor(float64(floorMap[targetTier]+PP_BASE) * ppBonus))
		pp += ppGained
		//fmt.Println("Gained", ppGained, "pp")
		kills++
		STATUS = genStat(killingStarted, kills, ap, exp, pp)
		fmt.Printf("\r%s", STATUS)
	}
}

func genStat(start time.Time, kills int, ap int, exp int, pp int) string {
	now := time.Now()
	minutes := now.Sub(start).Minutes()
	hours := now.Sub(start).Hours()
	killStats := fmt.Sprintf("Kills/KPM: %d/%.2f", kills, float64(kills)/minutes)
	expStats := fmt.Sprintf("EXP/EPM: %d/%.2f", exp, float64(exp)/minutes)
	apStats := fmt.Sprintf("AP/APM/KPA: %d/%.2f/%.2f", ap, float64(ap)/minutes, float64(kills)/float64(ap))
	ppStats := fmt.Sprintf("PP/PPPH: %.1f/%.2f", float64(pp)/float64(1000000), float64(float64(pp)/float64(1000000))/hours)
	statBlock := fmt.Sprintf("Hours: %.2f %s %s %s %s", hours, killStats, expStats, apStats, ppStats)
	return statBlock
}

func updateKillMap(killMap map[int]int) map[int]int {

	for tier, kills := range killMap {
		kills--
		if kills < 1 {
			kills = 40 - tier
			killMap[tier] = kills
		} else {
			killMap[tier] = kills
		}
	}
	return killMap
}

func parseKillMap() (floorMap map[int]int, killMap map[int]int) {

	for {
		err, floorMap, killMap := parseKillMapRetryable()
		if err != nil {
			log.Println(err)
			time.Sleep(time.Second)
			clickCheckWait(RegAttack)
		} else {
			return floorMap, killMap
		}
	}
}

func parseKillMapRetryable() (err error, floorMap map[int]int, killMap map[int]int) {

	clickCheckWait(EnterITOPOD)
	clickCheckWait(ITOPODOptimal)

	killMap = make(map[int]int)
	floorMap = make(map[int]int)

	snagRect(OCR_ITOPOD_START_BOX, "ocr.png")
	cmd := exec.Command("tesseract", "ocr.png", "stdout", "--dpi", "109")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("tesseract error %v", err), floorMap, killMap
	}
	optimal, err := strconv.ParseInt(numbers(output), 10, 32)
	if err != nil {
		return fmt.Errorf("optimal level parse error %v", err), floorMap, killMap
	}
	fmt.Println("Optimal ITOPOD is", optimal)

	max_tier := math.Floor(float64(optimal)/50) + 1

	for i := MIN_ITOPOD_TIER; i <= int(max_tier); i++ {
		idleFloor := (i * 50) - 5
		tier := i
		floorMap[tier] = idleFloor
	}
	floorMap[int(max_tier)] = int(optimal)
	clickCheckWait(ITOPODEngage)
	for i := MIN_ITOPOD_TIER; i <= int(max_tier); i++ {
		clickCheckWait(EnterITOPOD)

		clickCheckWait(ITOPODStartBox)
		typeWait(fmt.Sprintf("%d", floorMap[i]))
		clickCheckWait(ITOPODEngage)
		moveCheck(ITOPODHeader)
		snagRect(OCR_AP_KILL_COUNT_1LINE, "ocr.png")
		cmd = exec.Command("tesseract", "ocr.png", "stdout", "--dpi", "109")
		output, err = cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("tesseract error %v", err), floorMap, killMap
		}
		text := readable(output)
		if !strings.Contains(text, "kills") {
			snagRect(OCR_AP_KILL_COUNT_2LINE, "ocr.png")
			cmd := exec.Command("tesseract", "ocr.png", "stdout", "--dpi", "109")
			output, err := cmd.CombinedOutput()
			if err != nil {
				return fmt.Errorf("tesseract error %v", err), floorMap, killMap
			}
			text := readable(output)
			raw, counter, err := getKills(text)
			if err != nil {
				return err, floorMap, killMap
			}
			fmt.Println("Tier", i, "Raw:", raw)
			killMap[i] = counter
		} else {
			raw, counter, err := getKills(text)
			if err != nil {
				return err, floorMap, killMap
			}
			fmt.Println("Tier", i, "Raw:", raw)
			killMap[i] = counter
		}
	}
	return nil, floorMap, killMap
}

func getKills(input string) (string, int, error) {

	index := strings.Index(input, ".")
	var output string
	if index == -1 {
		output = input
	} else {
		output = input[:index]
	}
	num := numbers([]byte(output))
	kills, err := strconv.ParseInt(num, 10, 64)
	if err != nil {
		return output, int(kills), fmt.Errorf("unable to parse kills from %v", input)
	}
	return output, int(kills), nil
}

func pickTier(killMap map[int]int, currentTier int) int {
	maxTier := 0
	desiredTier := 0
	lowTierCount := 99
	for tier, kills := range killMap {
		if tier > maxTier {
			maxTier = tier
		}
		if kills < lowTierCount {
			lowTierCount = kills
			desiredTier = tier
		}
		if kills == lowTierCount && tier > desiredTier {
			desiredTier = tier
			lowTierCount = kills
		}
	}

	if lowTierCount == 2 || lowTierCount == 3 {
		return int(math.Max(float64(currentTier), float64(desiredTier)))
	}
	if lowTierCount == 1 {
		return desiredTier
	}
	return maxTier
}
