package main

// Global window offsets
var TOP int
var LEFT int

// Windows border tweaks
const (
	BAR_OFFSET_TOP  = 0
	BAR_OFFSET_LEFT = 2
)

// Trying to quit? ya I know
var QUIT bool

// Status string to print on exit
var STATUS string

// Want to pause? pfft, fine ok you can
var PAUSE Pauser

// Metrics stuff, I mean that is why you are here right?
var AppMetrics Metrics

// Amount of time to wait for NGU to act on a mouse click before sending a second
const FrameDelayMs = 25

// For uniquely identifying debug pngs
var FAILCOUNT int

// Lowest tower tier to attempt
var MIN_ITOPOD_TIER int

// Max levels above optimal to snipe
var MAX_ITOPOD_SNIPE int

// Base PP gained based on difficulty, 200/700/2000
const PP_BASE = 700

// Debugging variable needed for unique OCR attempts
var PARSE_ATTEMPTS = 0

// Set this to be true to save every freaking tesseract ocr
var EXTRASANITY bool

// The number of click/color timings to average together for the floating/instant fps measurement
const ColorTimingAvg = 10
