package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/disintegration/imaging"
)

func BossSnipe(zoneName string) {

	clickCheckWait(MenuAdventure)
	start := time.Now()
	killCounter := 0

	currentZone := "Unknown"

	// Turn off idling if it was on
	if checkColor(IdleModeOn, false) {
		clickCheckWait(IdleModeAbility)
	}

	for {
		if !checkColor(MyHealthFull, false) {
			clickCheckRight(AdventureLeft)
			currentZone = "SafeZone"
			//log.Println("Healing")
			for {
				if checkColor(MyHealthFull, false) {
					break
				}
			}
		}
		//log.Println("Going to Zone", zoneName)
		if currentZone != zoneName {
			clickCheckRight(AdventureRight)
			for {
				snagRect(OCR_ADVENTURE_ZONE_NAME, "ocr.png")
				srcImage, _ := imaging.Open("ocr.png")
				// make it bigger
				dstImage := imaging.Resize(srcImage, 400, 0, imaging.Lanczos)
				// then sharpen for the best chance of tesseracting
				dstImage = imaging.Sharpen(dstImage, 5)
				imaging.Save(dstImage, "ocr.png")
				cmd := exec.Command("tesseract", "ocr.png", "stdout", "--dpi", "109")
				output, err := cmd.CombinedOutput()
				if err != nil {
					log.Println(err)
					return
				}
				currentZone = readable(output)
				currentZone = strings.TrimSpace(currentZone)
				if currentZone == zoneName {
					break
				} else {
					clickCheckWait(AdventureLeft)
				}
			}
		}

		//log.Println("Finding a boss")
		skills := SkillStatus{}
		for {
			for {
				if checkColor(EnemyHealth, false) {
					break
				}
			}
			if checkColor(BossCrown, false) {
				break
			} else {
				clickCheckWait(AdventureLeft)
				clickCheckWait(AdventureRight)
			}
		}
		//log.Println("start fighting")
		for {
			if checkColor(EnemyHealthEmpty, false) {
				if !checkColor(EnemyStatsVisible, false) {
					killCounter++
					now := time.Now()
					minutes := now.Sub(start).Minutes()
					killStats := fmt.Sprintf("Kills/KPM: %d/%.2f", killCounter, float64(killCounter)/minutes)
					log.Println("Enemy is Dead woo!", killCounter, killStats)
					break
				} else {
					log.Println("Enemy is playing dead!")
				}
			}
			skills.BlockTillBestAttack()
		}
	}
}

type SkillStatus struct {
	// Attacks
	Regular  bool
	Strong   bool
	Piercing bool
	Ultimate bool
	// Defensive
	Parry   bool
	Block   bool
	DefBuff bool
	// Healing
	Heal  bool
	Regen bool
	// Offensive Buffs
	OffBuff  bool
	UltBuff  bool
	Charge   bool
	MegaBuff bool
	// Misc
	Paralyze bool
}

func (s *SkillStatus) GetAttackStatus() {
	s.Regular = checkColor(RegAttackUnused, false)
	s.Strong = checkColor(StrongAttackUnused, false)
	s.Piercing = checkColor(PiercingAttackUnused, false)
	s.Ultimate = checkColor(UltimateAttackUnused, false)
}
func (s *SkillStatus) GetThirdRowStatus() {
	s.Paralyze = checkColor(ParalyzeSkillUnused, false)
	s.Regen = checkColor(HyperRegenSkillUnused, false)
	s.MegaBuff = checkColor(MegaBuffSkillUnused, false)
}

func (s *SkillStatus) BlockTillBestAttack() {

	// remove the tool tip so we can see all the skills
	defer clickCheckWait(ITOPODEngage)

	for {
		s.GetAttackStatus()
		s.GetThirdRowStatus()
		if s.MegaBuff {
			clickCheckWait(MegaBuffSkillUnused)
		}
		if s.Ultimate {
			clickCheckWait(UltimateAttackUnused)
			return
		}
		if s.Piercing {
			clickCheckWait(PiercingAttackUnused)
			return
		}
		if s.Strong {
			clickCheckWait(StrongAttackUnused)
			return
		}
		if s.Regular {
			clickCheckWait(RegAttackUnused)
			return
		}
	}
}
