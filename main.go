package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-vgo/robotgo"
)

var TOP int
var LEFT int
var QUIT bool
var FAILCOUNT int
var STATUS string

func main() {

	spam := flag.Bool("spam", false, "Spam the location and color")
	flag.Parse()

	evChan := robotgo.EventStart()
	go func() {
		for event := range evChan {
			if string(event.Keychar) == "`" {
				fmt.Println("")
				log.Println("Thank you for playing!")
				writeStatusLog()
				QUIT = true
			}
		}
	}()

	getNGULocation()
	log.Println("Window Found at", TOP, LEFT)
	if *spam {
		spamDetails()
	}

	// Do a daily loop forever!
	for {
		// Run for 12 hours to make the pill available
		IdleITOPOD(time.Hour * 12)
		// 12 hours in do a pill and eat just max fruit
		CastIronPill()
		EatFruit(false)
		IdleITOPOD(time.Hour * 12)
		// every 24 hours spin the wheel and eat the fruit (and pill again)
		CastIronPill()
		EatFruit(true)
		SpinTheWheel()
	}
}

func readable(b []byte) string {

	var clean []byte
	for i := range b {
		isNum := b[i] >= 48 && b[i] <= 57
		isLow := b[i] >= 97 && b[i] <= 122
		isUpp := b[i] >= 65 && b[i] <= 90
		isSym := b[i] >= 32 && b[i] <= 46
		if isNum || isLow || isUpp || isSym {
			clean = append(clean, b[i])
		}
		if b[i] == 10 {
			clean = append(clean, 32)
		}
	}
	return string(clean)
}

func numbers(b []byte) string {
	var clean []byte
	for i := range b {
		isNum := b[i] >= 48 && b[i] <= 57
		if isNum {
			clean = append(clean, b[i])
		}
	}
	return string(clean)
}

func spamDetails() {
	for {
		if QUIT {
			os.Exit(0)
		}
		start := time.Now()
		x, y := robotgo.GetMousePos()
		color := robotgo.GetPixelColor(x, y)
		end := time.Now()
		timer := end.Sub(start)
		fmt.Println("pos:", x-LEFT, y-TOP, "color---- ", color, timer)
	}
}

func clickCheckValidate(c Check, inverse bool) bool {
	if QUIT {
		os.Exit(0)
	}
	// Always input relative pixel locations based on alt printscreen x,y paint coords
	frameDelay := int64(33 * 2)
	start := time.Now()
	robotgo.Move(LEFT+c.X, TOP+c.Y)
	robotgo.Click()
	duration := time.Now().Sub(start)
	var delay int64
	if duration.Milliseconds() < frameDelay {
		delay = frameDelay - duration.Milliseconds()
		//log.Println("Sleeping for", delay, "ms")
		time.Sleep(time.Millisecond * time.Duration(delay))
	}
	//log.Println("Desired:", frameDelay, "requested:", delay, "cost", time.Now().Sub(start).Milliseconds(), "ms")
	if inverse {
		return checkColorInverse(c)
	} else {
		return checkColor(c, false)
	}
}

func clickCheckWait(c Check) {
	click(c.X, c.Y, true)
}

func clickCheckNoWait(c Check) {
	click(c.X, c.Y, false)
}

func moveCheck(c Check) {
	if QUIT {
		os.Exit(0)
	}
	frameDelay := int64(33)

	// move slightly offset instantly
	robotgo.Move(LEFT+c.X, TOP+c.Y+80)
	robotgo.Click()
	// move slowly to the real location
	robotgo.MoveMouseSmooth(LEFT+c.X, TOP+c.Y, 2.0, 10.0)
	end := time.Now()
	start := time.Now()
	duration := end.Sub(start)
	if duration.Milliseconds() < frameDelay {
		delay := frameDelay - duration.Milliseconds()
		time.Sleep(time.Millisecond * time.Duration(delay))
	}
}

func typeWait(val string) {
	if QUIT {
		os.Exit(0)
	}
	frameDelay := int64(33 * 3)
	start := time.Now()
	robotgo.TypeStr(val)
	end := time.Now()
	duration := end.Sub(start)
	if duration.Milliseconds() < frameDelay {
		delay := frameDelay - duration.Milliseconds()
		time.Sleep(time.Millisecond * time.Duration(delay))
	}
}

func click(x int, y int, frameWait bool) {
	if QUIT {
		os.Exit(0)
	}
	// Always input relative pixel locations based on alt printscreen x,y paint coords
	frameDelay := int64(33 * 2)
	start := time.Now()
	robotgo.Move(LEFT+x, TOP+y)
	robotgo.Click()
	end := time.Now()
	duration := end.Sub(start)
	if frameWait && duration.Milliseconds() < frameDelay {
		delay := frameDelay - duration.Milliseconds()
		time.Sleep(time.Millisecond * time.Duration(delay))
	}
}

func SpinTheWheel() {
	clickCheckWait(MenuMoneyPit)
	clickCheckWait(DailySpin)
	clickCheckWait(DailySpinNoBS)
}

func EatFruit(all bool) {
	clickCheckWait(MenuYggdrasil)
	if all {
		clickCheckWait(EatAllFruit)
	} else {
		clickCheckWait(EatMaxFruit)
	}
}

func CastIronPill() {
	clickCheckWait(MenuBloodMagic)
	clickCheckWait(CastSpells)
	clickCheckWait(IronPill)
}

func writeStatusLog() {
	f, err := os.OpenFile("itopod_rewards.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	if _, err := f.WriteString(STATUS + "\n"); err != nil {
		log.Println(err)
	}
}
