package main

import (
	"fmt"
	"log"

	"github.com/go-vgo/robotgo"
)

type Check struct {
	X      int
	Y      int
	Colors []string
}

func checkColor(c Check, debug bool) bool {
	color := robotgo.GetPixelColor(c.X+LEFT, c.Y+TOP)
	for i := range c.Colors {
		if c.Colors[i] == color {
			return true
		}
	}
	if debug {
		snagRect(RECT{
			Left:   int32(c.X - 20),
			Top:    int32(c.Y - 20),
			Right:  int32(c.X + 20),
			Bottom: int32(c.Y + 20),
		}, "color-debug.png")
		log.Println(color, "is not", c.Colors)
	}
	return false
}

func checkColorInverse(c Check) bool {
	color := robotgo.GetPixelColor(c.X+LEFT, c.Y+TOP)
	for i := range c.Colors {
		if c.Colors[i] == color {
			snagRect(RECT{
				Left:   int32(c.X - 20),
				Top:    int32(c.Y - 20),
				Right:  int32(c.X + 20),
				Bottom: int32(c.Y + 20),
			}, fmt.Sprintf("color-debug-%d.png", FAILCOUNT))
			FAILCOUNT++
			return true
		}
	}
	return false
}

var EnemyHealth = Check{X: 736, Y: 430, Colors: []string{"d93030", "eb3434", "db3131", "da3030", "e83333", "0f0303"}}
var RegAttack = Check{X: 430, Y: 137}
var RegAttackUnused = Check{X: 471, Y: 120, Colors: []string{"f89b9b"}}
var RegAttackUsed = Check{X: 471, Y: 120, Colors: []string{"7c4e4e", "c27a7a", "c47b7b", "d48484", "d58585", "b06e6e", "c37a7a", "d68686", "ba7474"}}
var EnterITOPOD = Check{X: 370, Y: 250}
var ITOPODStartBox = Check{X: 613, Y: 219}
var ITOPODEngage = Check{X: 625, Y: 322}
var ITOPODHeader = Check{X: 445, Y: 45}
var ITOPODOptimal = Check{X: 708, Y: 233}
var ITOPODBoxOpen = Check{X: 605, Y: 225, Colors: []string{"ffffff"}}

var IdleModeAbility = Check{X: 330, Y: 137}
var IdleModeOn = Check{X: 316, Y: 114, Colors: []string{"ffeb04"}}

var OCR_ITOPOD_START_BOX = RECT{Left: 598, Top: 216, Right: 655, Bottom: 236}
var OCR_AP_KILL_COUNT_2LINE = RECT{Left: 470, Top: 127, Right: 740, Bottom: 155}
var OCR_AP_KILL_COUNT_1LINE = RECT{Left: 470, Top: 127, Right: 740, Bottom: 140}
var EXP_VALIDATION_LINE = RECT{Left: 319, Top: 595, Right: 600, Bottom: 596}

var MenuMoneyPit = Check{X: 231, Y: 100}
var MenuAdventure = Check{X: 231, Y: 131}
var MenuBloodMagic = Check{X: 231, Y: 265}
var MenuYggdrasil = Check{X: 231, Y: 345}

// Money Pit Related
var DailySpin = Check{X: 822, Y: 263}
var DailySpinNoBS = Check{X: 717, Y: 587}

// Blood Magic Related
var CastSpells = Check{X: 392, Y: 142}
var IronPill = Check{X: 737, Y: 244}

// Yggdrasil Related
var EatAllFruit = Check{X: 827, Y: 517}
var EatMaxFruit = Check{X: 835, Y: 479}
