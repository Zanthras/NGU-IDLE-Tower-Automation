package main

import (
	"fmt"
	"os"
	"time"

	"github.com/go-vgo/robotgo"
)

func clickCheckWait(c Check) {
	click(c.X, c.Y, true)
}

func clickCheckRight(c Check) {
	clickRight(c.X, c.Y, true)
}

func clickCheckNoWait(c Check) {
	click(c.X, c.Y, false)
}

func moveCheck(c Check) {
	if QUIT {
		os.Exit(0)
	}
	frameDelay := int64(FrameDelayMs)

	// move slightly offset instantly
	robotgo.Move(LEFT+c.X, TOP+c.Y+10)
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
	// Conservative value that does work
	//frameDelay := int64(33 * 3)
	// Faster value that will probably work
	frameDelay := int64(FrameDelayMs)
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
	// Conservative value that does work
	//frameDelay := int64(33 * 2)
	// Faster value that will probably work
	frameDelay := int64(FrameDelayMs)
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

func clickRight(x int, y int, frameWait bool) {
	if QUIT {
		os.Exit(0)
	}
	frameDelay := int64(FrameDelayMs)
	start := time.Now()
	robotgo.Move(LEFT+x, TOP+y)
	robotgo.ClickRight()
	end := time.Now()
	duration := end.Sub(start)
	if frameWait && duration.Milliseconds() < frameDelay {
		delay := frameDelay - duration.Milliseconds()
		time.Sleep(time.Millisecond * time.Duration(delay))
	}
}

func checkColor(c Check, debug bool) bool {
	start := time.Now()
	color := robotgo.GetPixelColor(c.X+LEFT, c.Y+TOP)
	AppMetrics.RecordClick(time.Now().Sub(start))
	if debug {
		snagRect(RECT{
			Left:   int32(c.X - 20),
			Top:    int32(c.Y - 20),
			Right:  int32(c.X + 20),
			Bottom: int32(c.Y + 20),
		}, "color-debug.png")
	}
	for i := range c.Colors {
		if c.Colors[i] == color {
			return true
		}
	}
	return false
}

func getColor(c Check, debug bool) string {
	start := time.Now()
	color := robotgo.GetPixelColor(c.X+LEFT, c.Y+TOP)
	AppMetrics.RecordClick(time.Now().Sub(start))
	if debug {
		snagRect(RECT{
			Left:   int32(c.X - 20),
			Top:    int32(c.Y - 20),
			Right:  int32(c.X + 20),
			Bottom: int32(c.Y + 20),
		}, fmt.Sprintf("color-debug-%s.png", color))
	}
	return color
}

func checkColorInverse(c Check, debug bool) bool {
	color := robotgo.GetPixelColor(c.X+LEFT, c.Y+TOP)
	for i := range c.Colors {
		if c.Colors[i] == color {
			if debug {
				snagRect(RECT{
					Left:   int32(c.X - 20),
					Top:    int32(c.Y - 20),
					Right:  int32(c.X + 20),
					Bottom: int32(c.Y + 20),
				}, fmt.Sprintf("color-debug-%d.png", FAILCOUNT))
				FAILCOUNT++
			}
			return true
		}
	}
	return false
}

type CheckResult struct {
	Match bool
	Color string
}

func MultiCheckColor(checks []*Check, debug bool) map[*Check]CheckResult {
	start := time.Now()
	var Top int
	var Bottom int
	var Left int
	var Right int
	for i, check := range checks {
		if i == 0 {
			Top = check.Y
			Bottom = check.Y
			Left = check.X
			Right = check.X
			continue
		}
		if check.Y < Top {
			Top = check.Y
		}
		if check.Y > Bottom {
			Bottom = check.Y
		}
		if check.X < Left {
			Left = check.X
		}
		if check.X > Right {
			Right = check.X
		}
	}
	var x, y, w, h int
	x = Left + LEFT
	y = Top + TOP
	w = Right - Left + 1
	h = Bottom - Top + 1
	// add a 10px border buffer when debugging
	if debug {
		x = x - 10
		y = y - 10
		w = w + 20
		h = h + 20
	}
	bitmap := robotgo.CaptureScreen(x, y, w, h)
	defer robotgo.FreeBitmap(bitmap)

	results := make(map[*Check]CheckResult)
	for _, check := range checks {
		color := robotgo.GetColors(bitmap, check.X-Left, check.Y-Top)
		var match bool
		for i := range check.Colors {
			if check.Colors[i] == color {
				match = true
				break
			}
		}
		results[check] = CheckResult{Match: match, Color: color}
	}
	AppMetrics.RecordClick(time.Now().Sub(start))

	if debug {
		robotgo.SaveBitmap(bitmap, "multicheck.png")
	}
	return results
}
