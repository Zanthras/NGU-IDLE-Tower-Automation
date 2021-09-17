package main

import (
	"time"
)

var FAILCOUNT int
var CLICKDURATION time.Duration
var CLICKCOUNT int

type Check struct {
	X      int
	Y      int
	Colors []string
}

var EnemyHealth = Check{X: 736, Y: 430, Colors: []string{"d93030", "eb3434", "db3131", "da3030", "e83333", "0f0303"}}
var EnemyHealthEmpty = Check{X: 736, Y: 430, Colors: []string{"fafafa"}}
var EnemyStatsVisible = Check{X: 734, Y: 350, Colors: []string{"000000"}}
var MyHealthFull = Check{X: 514, Y: 447, Colors: []string{"ec3434"}}
var EnterITOPOD = Check{X: 370, Y: 250}
var ITOPODStartBox = Check{X: 613, Y: 219}
var ITOPODEngage = Check{X: 625, Y: 322}
var ITOPODHeader = Check{X: 445, Y: 45}
var ITOPODOptimal = Check{X: 708, Y: 233}
var ITOPODBoxOpen = Check{X: 605, Y: 225, Colors: []string{"ffffff"}}
var AdventureLeft = Check{X: 731, Y: 236}
var AdventureRight = Check{X: 944, Y: 236}
var BossCrown = Check{X: 741, Y: 307, Colors: []string{"f7ef29"}}

// Adventure Skills
var RegAttackUnused = Check{X: 471, Y: 120, Colors: []string{"f89b9b"}}
var RegAttackUsed = Check{X: 471, Y: 120, Colors: []string{"7c4e4e"}}
var StrongAttackUnused = Check{X: 581, Y: 120, Colors: []string{"f89b9b"}}
var PiercingAttackUnused = Check{X: 781, Y: 120, Colors: []string{"f89b9b"}}
var UltimateAttackUnused = Check{X: 881, Y: 120, Colors: []string{"f89b9b"}}
var ParalyzeSkillUnused = Check{X: 371, Y: 193, Colors: []string{"c39494"}}
var HyperRegenSkillUnused = Check{X: 471, Y: 193, Colors: []string{"c39494"}}
var MegaBuffSkillUnused = Check{X: 671, Y: 193, Colors: []string{"c39494"}}

var IdleModeAbility = Check{X: 330, Y: 137}
var IdleModeOn = Check{X: 316, Y: 114, Colors: []string{"ffeb04"}}

var OCR_ITOPOD_START_BOX = RECT{Left: 598, Top: 216, Right: 655, Bottom: 236}
var OCR_AP_KILL_COUNT_2LINE = RECT{Left: 470, Top: 127, Right: 740, Bottom: 155}
var OCR_AP_KILL_COUNT_1LINE = RECT{Left: 470, Top: 127, Right: 740, Bottom: 140}
var EXP_VALIDATION_LINE = RECT{Left: 319, Top: 595, Right: 600, Bottom: 596}
var EXP_VALIDATION_LINE2 = RECT{Left: 319, Top: 580, Right: 600, Bottom: 581}
var OCR_ADVENTURE_ZONE_NAME = RECT{Left: 753, Top: 226, Right: 907, Bottom: 247}

// Input Selections

// Menu Selections
var MenuMoneyPit = Check{X: 231, Y: 100}
var MenuAdventure = Check{X: 231, Y: 131}
var MenuInventory = Check{X: 231, Y: 155}
var MenuBloodMagic = Check{X: 231, Y: 265}
var MenuYggdrasil = Check{X: 231, Y: 345}

// Inventory Related
var Loadout1 = Check{X: 332, Y: 284}
var Loadout2 = Check{X: 362, Y: 284}
var Loadout3 = Check{X: 392, Y: 284}

// Money Pit Related
var DailySpin = Check{X: 822, Y: 263}
var DailySpinNoBS = Check{X: 717, Y: 587}

// Blood Magic Related
var CapRituals = Check{X: 822, Y: 135}
var CastSpells = Check{X: 392, Y: 142}
var IronPill = Check{X: 737, Y: 244}

// NGU Related
var NGUMagicEnergySwitch = Check{X: 373, Y: 136}
var CapNGU = Check{X: 625, Y: 186}

// Yggdrasil Related
var EatAllFruit = Check{X: 827, Y: 517}
var EatMaxFruit = Check{X: 835, Y: 479}

var EnNGUSkill1Up = Check{X: 518, Y: 265}
var EnNGUSkill2Up = Check{X: 518, Y: 300}
var EnNGUSkill3Up = Check{X: 518, Y: 335}
var EnNGUSkill4Up = Check{X: 518, Y: 370}
var EnNGUSkill5Up = Check{X: 518, Y: 405}
var EnNGUSkill6Up = Check{X: 518, Y: 440}
var EnNGUSkill7Up = Check{X: 518, Y: 475}
var EnNGUSkill8Up = Check{X: 518, Y: 510}
var EnNGUSkill9Up = Check{X: 518, Y: 545}
