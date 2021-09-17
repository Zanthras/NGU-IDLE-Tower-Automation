package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-vgo/robotgo"
)

func measureClick() {

	clicks := 2500

	//for x := 0; x < clicks; x++ {
	//	clickVariableDelay(EnNGUSkill1Up, time.Millisecond*5*1)
	//}
	//time.Sleep(time.Second)
	//for x := 0; x < clicks; x++ {
	//	clickVariableDelay(EnNGUSkill2Up, time.Millisecond*5*2)
	//}
	//time.Sleep(time.Second)
	//for x := 0; x < clicks; x++ {
	//	clickVariableDelay(EnNGUSkill3Up, time.Millisecond*5*3)
	//}
	//time.Sleep(time.Second)
	//for x := 0; x < clicks; x++ {
	//	clickVariableDelay(EnNGUSkill4Up, time.Millisecond*5*4)
	//}
	time.Sleep(time.Second)
	for x := 0; x < clicks; x++ {
		clickVariableDelay(EnNGUSkill5Up, time.Millisecond*5*5)
	}
	//time.Sleep(time.Second)
	//for x := 0; x < clicks; x++ {
	//	clickVariableDelay(EnNGUSkill6Up, time.Millisecond*5*6)
	//}
	//time.Sleep(time.Second)
	//for x := 0; x < clicks; x++ {
	//	clickVariableDelay(EnNGUSkill7Up, time.Millisecond*5*7)
	//}
	//time.Sleep(time.Second)
	//for x := 0; x < clicks; x++ {
	//	clickVariableDelay(EnNGUSkill8Up, time.Millisecond*5*8)
	//}
	//time.Sleep(time.Second)
	//for x := 0; x < clicks; x++ {
	//	clickVariableDelay(EnNGUSkill9Up, time.Millisecond*5*9)
	//}

}

func clickVariableDelay(c Check, frameDelay time.Duration) {
	if QUIT {
		os.Exit(0)
	}
	// Always input relative pixel locations based on alt printscreen x,y paint coords
	start := time.Now()
	robotgo.MoveClick(LEFT+c.X, TOP+c.Y)
	end := time.Now()
	duration := end.Sub(start)
	if duration < frameDelay {
		delay := frameDelay - duration
		time.Sleep(delay)
	}
}

func measureAttack() {
	start := time.Now()
	clicks := int64(0)
	last := start
	first := true
	for {
		if QUIT {
			os.Exit(0)
		}
		var regCorner1 = Check{X: 428, Y: 118, Colors: []string{"f89b9b"}}
		var regCorner2 = Check{X: 515, Y: 138, Colors: []string{"f89b9b"}}
		results := MultiCheckColor([]*Check{&regCorner1, &regCorner2, &EnemyHealthEmpty}, false)

		// wait for the enemy to have health
		if !results[&EnemyHealthEmpty].Match {
			for {
				clickCheckNoWait(RegAttackUnused)
				// enemy died!
				if checkColor(EnemyHealthEmpty, false) {
					if first {
						first = false
						start = time.Now()
						last = start
						continue
					}
					clicks++
					now := time.Now()
					dur := now.Sub(start)
					interval := now.Sub(last)
					log.Println(clicks, "avg", dur.Milliseconds()/clicks, "ms", "instant", interval.Milliseconds(), "ms")
					last = now
					break
				}
			}
		}
	}
}

func measureSpawn() {
	start := time.Now()
	clicks := int64(0)
	last := start
	first := true
	for {
		if QUIT {
			os.Exit(0)
		}
		//if checkColor(EnemyHealth, false) {
		if checkColor(EnemyHealth, false) {
			if first {
				clickCheckWait(RegAttackUnused)
				first = false
				start = time.Now()
				last = start
				continue
			}
			clicks++
			now := time.Now()
			dur := now.Sub(start)
			interval := now.Sub(last)
			log.Println(clicks, "avg", dur.Milliseconds()/clicks, "ms", "instant", interval.Milliseconds(), "ms")
			last = now
			//if !checkColor(RegAttackUnused, true) {
			//	os.Exit(0)
			//}
			clickCheckWait(RegAttackUnused)
		}
	}
}

func spamDetails() {
	for {
		if QUIT {
			os.Exit(0)
		}
		PAUSE.Wait()
		start := time.Now()
		x, y := robotgo.GetMousePos()
		color := robotgo.GetPixelColor(x, y)
		end := time.Now()
		timer := end.Sub(start)
		fmt.Println("pos:", x-LEFT, y-TOP, "color---- ", color, timer)
	}
}
