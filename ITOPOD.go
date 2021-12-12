package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	"github.com/tevino/abool"
)

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
	21: 384,
	22: 422,
	23: 464,
	24: 508,
	25: 554,
	26: 602,
	27: 652,
	28: 704,
	29: 758,
	30: 814,
	31: 872,
	32: 932,
}

var OPTIMAL_LEVEL int

func IdleITOPOD(dur time.Duration) {

	log.Println("Running for", dur)

	expBonus := 1.0
	ppBonus := 1.0

	defer writeStatusLog()

	clickCheckWait(MenuAdventure)

	//move the mouse off the popup causing menu
	clickCheckWait(ITOPODBoxClose)

	// Turn off idling if it was on
	idleResults := getRectColors(IdleModeArea)
	if _, found := idleResults[IdleModeOn.Colors[0]]; found {
		clickCheckWait(IdleModeAbility)
	}

	// Do the hover OCR thing to get current kill counts
	floorMap, killMap := parseKillMap()

	killingStarted := time.Now()
	brokenEXP := abool.New()

	var currentTier int
	targetTier := pickTier(killMap, currentTier)
	for {
		// TimeCheck
		loopEnd := killingStarted.Add(dur)
		if time.Now().After(loopEnd) {
			log.Println("Loop Completed")
			return
		}
		// Check for exp validation failure, and rescan if needed
		if brokenEXP.IsSet() {
			AppMetrics.Rescans.Inc()
			fmt.Println("")
			floorMap2, killMap2 := parseKillMap()
			for tier := range killMap2 {
				if killMap2[tier] != killMap[tier] {
					AppMetrics.RescansNeeded.Inc()
					log.Println(tier, "parsed as", killMap2[tier], "expected it to be", killMap[tier])
				}
			}
			floorMap = floorMap2
			killMap = killMap2
			brokenEXP.UnSet()
		}
		targetTier = pickTier(killMap, targetTier)
		if targetTier != currentTier {
			//fmt.Println("\nTransitioning to tier", targetTier, "floor", floorMap[targetTier], "kills to go", killMap[targetTier])
			currentTier = targetTier
			for {
				clickCheckWait(EnterITOPOD) // attempt to open the tower interface, looping in case of failure to register click
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
		t1 := time.Now()
		for {
			if checkColor(EnemyHealth, false) {
				break
			}
			if time.Now().Sub(t1) > time.Second*30 {
				getColor(EnemyHealth, true)
				panic("Failed to detect enemy spawn, check color-debug-XXXXXX.png to help understand why")
			}
		}
		// Spam the attack button until the enemy is confirmed dead
		for {
			clickCheckNoWait(RegAttackUnused)
			if checkColor(EnemyHealthEmpty, false) {
				break
			}
		}
		// Record the ap/exp gains and trigger an exp validation check
		if killMap[currentTier] == 1 {
			AppMetrics.AP.Inc()
			expGained := float64(math.Floor(float64(baseEXP[targetTier]) * expBonus))
			AppMetrics.EXP.Add(expGained)
			go validateEXPLine(brokenEXP)
			fmt.Println("")
		}

		// Recalculate where the kills should be at for all the tiers
		killMap = updateKillMap(killMap)
		ppGained := float64(math.Floor(float64(floorMap[targetTier]+PP_BASE) * ppBonus))
		// Record all the per kill stats
		AppMetrics.TierKills.Inc(fmt.Sprintf("%d", currentTier))
		AppMetrics.PP.Add(ppGained)
		AppMetrics.Kills.Inc()
		AppMetrics.generateSTATUS()
		// Move the pausing from the start of the loop to the end of the loop but before the status print, this way the pause doesnt print before the final stat block
		PAUSE.Wait()
		fmt.Printf("\r%s", STATUS)
	}
}

func updateKillMap(killMap map[int]int) map[int]int {

	for tier, kills := range killMap {
		kills--
		if kills < 1 {
			if tier > 20 {
				kills = 20
				killMap[tier] = kills
			} else {
				kills = 40 - tier
				killMap[tier] = kills
			}
		} else {
			killMap[tier] = kills
		}
	}
	return killMap
}

func parseKillMap() (floorMap map[int]int, killMap map[int]int) {

	notKillingStart := time.Now()
	for {
		err, floorMap, killMap := parseKillMapRetryable()
		if err != nil {
			log.Println(err)
			time.Sleep(time.Second)
			clickCheckWait(ITOPODBoxClose)  // Close the level selection box (or click nothing)
			clickCheckWait(RegAttackUnused) // Cause a kill to force advance kills left in case of really bad parsing
		} else {
			dur := time.Now().Sub(notKillingStart)
			AppMetrics.IdleDuration += dur
			return floorMap, killMap
		}
	}
}

func captureITOPODLevel() (int, error) {

	snagRect(OCR_ITOPOD_START_BOX, "ocr.png")
	srcImage, _ := imaging.Open("ocr.png")
	// make it bigger
	dstImage := imaging.Resize(srcImage, 400, 0, imaging.Lanczos)
	// then sharpen for the best chance of tesseracting
	dstImage = imaging.Sharpen(dstImage, 5)
	imaging.Save(dstImage, "ocr.png")
	cmd := exec.Command("tesseract", "ocr.png", "stdout", "--dpi", "109")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return 0, fmt.Errorf("tesseract error %v", err)
	}
	parsed, err := strconv.ParseInt(numbers(output), 10, 32)
	if err != nil {
		return 0, fmt.Errorf("itopod level parse error %v", err)
	}
	return int(parsed), nil
}

func parseKillMapRetryable() (err error, floorMap map[int]int, killMap map[int]int) {

	PARSE_ATTEMPTS++

	killMap = make(map[int]int)
	floorMap = make(map[int]int)

	clickCheckWait(EnterITOPOD)

	// parse max level for that sweet sweet exp
	//clickCheckWait(ITOPODMax) // had a case where the optimal level was resetting after clicking... turning this off completely
	clickCheckWait(ITOPODOptimal)
	maxLevel, err := captureITOPODLevel()
	if err != nil {
		return fmt.Errorf("max %v", err), floorMap, killMap
	}
	maxTier := math.Floor(float64(maxLevel)/50) + 1

	// parse optimal level for that sweet sweet kpm (and thus apm, epm, and ppph)
	clickCheckWait(ITOPODOptimal)
	optimalLevel, err := captureITOPODLevel()
	if err != nil {
		return fmt.Errorf("max %v", err), floorMap, killMap
	}
	optimalTier := math.Floor(float64(optimalLevel)/50) + 1

	if optimalLevel < OPTIMAL_LEVEL {
		return fmt.Errorf("invalid optimal level should be greater than %d is %d", OPTIMAL_LEVEL, optimalLevel), floorMap, killMap
	}
	OPTIMAL_LEVEL = optimalLevel
	if maxTier < optimalTier {
		return fmt.Errorf("invalid max level %d, smaller than optimal level %d", maxLevel, optimalLevel), floorMap, killMap
	}

	fmt.Printf("Optimal/Max ITOPOD is %d/%d\n", optimalLevel, maxLevel)

	for i := MIN_ITOPOD_TIER; i <= int(optimalTier); i++ {
		idleFloor := (i * 50) - 5
		tier := i
		floorMap[tier] = idleFloor
	}

	// Sniping the highest tier you can is worth it for exp/ap gains
	var topTier int
	if maxTier > optimalTier && (maxLevel-optimalLevel) <= MAX_ITOPOD_SNIPE {
		fmt.Printf("Sniping tier %d at level %d\n", int(maxTier), (int(maxTier)-1)*50)
		floorMap[int(maxTier)] = (int(maxTier) - 1) * 50
		topTier = int(maxTier)
	} else {
		topTier = int(optimalTier)
	}

	// grab the optimal floor for the max tier unless its too close to the next tier
	if float64(optimalLevel) > optimalTier*50-5 {
		floorMap[int(optimalTier)] = int(optimalTier*50 - 5)
	} else {
		floorMap[int(optimalTier)] = optimalLevel
	}
	clickCheckWait(ITOPODEngage)
	for i := MIN_ITOPOD_TIER; i <= topTier; i++ {

		OCR_BOX := OCR_AP_KILL_COUNT_2LINE
		if i < 4 {
			OCR_BOX = OCR_AP_KILL_COUNT_1LINE
		}
		clickCheckWait(EnterITOPOD)

		clickCheckWait(ITOPODStartBox)
		typeWait(fmt.Sprintf("%d", floorMap[i]))
		clickCheckWait(ITOPODEngage)
		moveCheck(ITOPODHeader)
		snagRect(OCR_BOX, "ocr.png")
		srcImage, _ := imaging.Open("ocr.png")
		// make it bigger
		dstImage := imaging.Resize(srcImage, 400, 0, imaging.Lanczos)
		// then sharpen for the best chance of tesseracting
		dstImage = imaging.Sharpen(dstImage, 5)
		imaging.Save(dstImage, "ocr.png")
		cmd := exec.Command("tesseract", "ocr.png", "stdout", "--dpi", "109")
		output, err := cmd.CombinedOutput()
		if err != nil {
			if EXTRASANITY {
				os.Rename("ocr.png", fmt.Sprintf("Parse-%d_Tier-%d_Parsed-BAD1.png", PARSE_ATTEMPTS, i))
			}
			return fmt.Errorf("tesseract error %v", err), floorMap, killMap
		}
		text := readable(output)
		if !strings.Contains(text, "kills") {
			if EXTRASANITY {
				os.Rename("ocr.png", fmt.Sprintf("Parse-%d_Tier-%d_Parsed-BAD2.png", PARSE_ATTEMPTS, i))
			}
			return fmt.Errorf("badly formated tesseract input %s", text), floorMap, killMap
		}
		raw, counter, err := getKills(text)
		if err != nil {
			if EXTRASANITY {
				os.Rename("ocr.png", fmt.Sprintf("Parse-%d_Tier-%d_Parsed-BAD3.png", PARSE_ATTEMPTS, i))
			}
			return err, floorMap, killMap
		}
		if EXTRASANITY {
			os.Rename("ocr.png", fmt.Sprintf("Parse-%d_Tier-%d_Parsed-%d.png", PARSE_ATTEMPTS, i, counter))
		}
		fmt.Println("Tier", i, "Raw:", raw)
		killMap[i] = counter
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
	desiredTier := 0
	lowTierCount := 99
	for tier, kills := range killMap {
		if kills < lowTierCount {
			lowTierCount = kills
			desiredTier = tier
		}
		if kills == lowTierCount && tier > desiredTier {
			desiredTier = tier
			lowTierCount = kills
		}
	}

	// If the closest count is within 2 or 3 kills, go farm the higher of the two tiers
	if lowTierCount == 2 || lowTierCount == 3 {
		return int(math.Max(float64(currentTier), float64(desiredTier)))
	}
	// Go directly to the tier that will produce the rewards
	if lowTierCount == 1 {
		return desiredTier
	}
	// If nothing is close go farm the optimal tier
	return int(math.Floor(float64(OPTIMAL_LEVEL)/50) + 1)
}

func validateEXPLine(toggle *abool.AtomicBool) {
	time.Sleep(100 * time.Millisecond)
	colors := getRectColors(EXP_VALIDATION_LINE)
	if _, found := colors["0000ff"]; !found {
		// check the line directly above because macguffins arrive after exp thus make the exp line go up by one
		colors2 := getRectColors(EXP_VALIDATION_LINE2)
		if _, found2 := colors2["0000ff"]; !found2 {
			fat_line := EXP_VALIDATION_LINE
			fat_line.Top = fat_line.Top - 45
			fat_line.Bottom = fat_line.Bottom + 15
			snagRect(fat_line, fmt.Sprintf("exp-%d.png", PARSE_ATTEMPTS))
			toggle.Set()
		}
	}
}
