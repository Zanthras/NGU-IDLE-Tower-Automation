package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/go-vgo/robotgo"
)

func main() {

	snipe := flag.String("snipe", "", "Attempt to boss snipe the targeted zone")
	tower := flag.Bool("tower", false, "Run the ITOPOD for AP/EXP")
	tier := flag.Int("tier", 1, "Set the Tower minimum tier for better non AP gains")
	levelSkip := flag.Int("skip", 10, "Maximum ITOPOD levels to snipe higher for extra exp/ap")
	debug := flag.String("debug", "", "execute debug function")
	duration := flag.Duration("time", time.Hour*12, "time to run before first iron pill")
	secondHalf := flag.Bool("second", false, "Set the script to assume its the second half of as 24h cycle")
	flag.Parse()

	MIN_ITOPOD_TIER = *tier
	MAX_ITOPOD_SNIPE = *levelSkip

	//getNGULocation()
	log.Println("Window Found at", TOP, LEFT)

	// Setup the bailout key
	go watchKeyBindings("`", "~")

	// Feature selector
	switch {
	case *debug != "":
		switch *debug {
		case "attack":
			measureAttack()
		case "spawn":
			measureSpawn()
		case "spam":
			spamDetails()
		case "timing":
			measureClick()
		case "click":
			speedClick()
		case "itopodlevel":
			log.Println(captureITOPODLevel())
		}
	case *snipe != "":
		BossSnipe(*snipe)
	case *tower:

		// Export metrics
		http.Handle("/metrics", promhttp.Handler())
		go http.ListenAndServe(":9100", nil)

		// Setup a panic if the game freezes so to more easily debug after the fact
		go FreezePanic()
		initMetrics()
		firstRun := true

		// One off second half of the day run
		if *secondHalf {
			IdleITOPOD(*duration)
			CastIronPill()
			EatFruit(true)
			firstRun = false
		}

		// Do a daily loop forever!
		for {
			if firstRun {
				// customized first run time because you almost never start idling at exactly 00:00 for blood
				IdleITOPOD(*duration)
			} else {
				// Run for 12 hours to make the pill available
				IdleITOPOD(time.Hour * 12)
			}
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
}

func watchKeyBindings(quitKey string, pauseKey string) {
	evChan := robotgo.EventStart()
	for event := range evChan {
		switch {
		case string(event.Keychar) == quitKey:
			fmt.Println("")
			log.Println("Thank you for playing!")
			writeStatusLog()
			// Trigger a clean quit
			QUIT = true
			time.Sleep(time.Second)
			// Force quit if none of the normal quit locations executed
			os.Exit(0)
		case string(event.Keychar) == pauseKey:
			if PAUSE.IsLocked() {
				PAUSE.Unlock()
			} else {
				PAUSE.Lock()
			}
		}
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

func SpinTheWheel() {
	notKillingStart := time.Now()
	clickCheckWait(MenuMoneyPit)
	clickCheckWait(DailySpin)
	clickCheckWait(DailySpinNoBS)
	dur := time.Now().Sub(notKillingStart)
	AppMetrics.IdleDuration += dur
}

func EatFruit(all bool) {
	notKillingStart := time.Now()
	clickCheckWait(MenuYggdrasil)
	if all {
		clickCheckWait(EatAllFruit)
	} else {
		clickCheckWait(EatMaxFruit)
	}
	dur := time.Now().Sub(notKillingStart)
	AppMetrics.IdleDuration += dur
}

func CastIronPill() {
	notKillingStart := time.Now()
	clickCheckWait(MenuBloodMagic)
	clickCheckWait(CastSpells)
	clickCheckWait(IronPill)
	dur := time.Now().Sub(notKillingStart)
	AppMetrics.IdleDuration += dur
}

func writeStatusLog() {
	f, err := os.OpenFile("itopod_rewards.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	if STATUS == "" {
		return
	}
	if _, err := f.WriteString(STATUS + "\n"); err != nil {
		log.Println(err)
	}
}

func prettyNum(num float64, notFloat bool) string {

	if num < 1000 {
		if notFloat {
			return fmt.Sprintf("%d", int64(num))
		}
		return fmt.Sprintf("%.2f", num)
	}
	if num > math.MaxInt64 {
		return fmt.Sprintf("%.2f", num)
	}
	v := int64(num)

	sign := ""

	// Min int64 can't be negated to a usable value, so it has to be special cased.
	if v == math.MinInt64 {
		return "-9,223,372,036,854,775,808"
	}

	if v < 0 {
		sign = "-"
		v = 0 - v
	}

	parts := []string{"", "", "", "", "", "", ""}
	j := len(parts) - 1

	for v > 999 {
		parts[j] = strconv.FormatInt(v%1000, 10)
		switch len(parts[j]) {
		case 2:
			parts[j] = "0" + parts[j]
		case 1:
			parts[j] = "00" + parts[j]
		}
		v = v / 1000
		j--
	}
	parts[j] = strconv.Itoa(int(v))
	return sign + strings.Join(parts[j:], ",")
}

type Pauser struct {
	mtx    sync.Mutex
	locked bool
	start  time.Time
	total  time.Duration
}

func (p *Pauser) Lock() {
	fmt.Println("\nPaused")
	p.locked = true
	p.start = time.Now()
	p.mtx.Lock()
}

func (p *Pauser) Unlock() {
	fmt.Println("\nUnpaused")
	p.mtx.Unlock()
	p.locked = false
	p.total += time.Now().Sub(p.start)
}

func (p *Pauser) IsLocked() bool {
	return p.locked
}

func (p *Pauser) Wait() {
	p.mtx.Lock()
	p.mtx.Unlock()
}

func (p *Pauser) Duration() time.Duration {
	return p.total
}

func noop(i interface{}) {}

func FreezePanic() {
	// Known spots will panic at 30 seconds, this unknown one needs to kick in after that timeframe
	timer := time.NewTicker(time.Minute)
	var last float64
	var pauseDur time.Duration
	for range timer.C {
		// If we paused at any point during the last freeze check, the duration will have gone up and thus we can ignore
		// any potential stalling. (since it could just be a pause).
		currentPause := PAUSE.Duration()
		if currentPause > pauseDur || PAUSE.IsLocked() {
			continue
		}
		pauseDur = currentPause
		current := AppMetrics.Kills.Value
		if current == last {
			PrintStack()
			panic("Game has stalled, please debug the current state")
		}
		last = current
	}
}

func PrintStack() {
	os.Stderr.Write(Stack())
}

// Copied from the debug package, just with all set to true
func Stack() []byte {
	buf := make([]byte, 1024)
	for {
		n := runtime.Stack(buf, true)
		if n < len(buf) {
			return buf[:n]
		}
		buf = make([]byte, 2*len(buf))
	}
}
