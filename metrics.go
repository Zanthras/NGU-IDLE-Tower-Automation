package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Metrics struct {
	EPM              MetricGauge
	EXP              MetricCounter
	KPM              MetricGauge
	Kills            MetricCounter
	AP               MetricCounter
	APM              MetricGauge
	KPA              MetricGauge
	PP               MetricCounter
	PPPH             MetricGauge
	FPS              MetricGauge
	FrameRate        MetricGauge
	Rescans          MetricCounter
	RescansNeeded    MetricCounter
	IdleDuration     time.Duration
	Start            time.Time
	clickCount       int
	clickDuration    time.Duration
	clickFloatingDur []time.Duration
	TierKills        MetricsMap
}

func (m *Metrics) generateSTATUS() {
	adjustedDur := time.Now().Sub(m.Start) - PAUSE.Duration() - m.IdleDuration
	minutes := adjustedDur.Minutes()
	hours := adjustedDur.Hours()
	m.KPM.Set(m.Kills.Value / minutes)
	killStats := fmt.Sprintf("Kills/KPM: %s/%.2f", prettyNum(m.Kills.Value, true), m.KPM.Value)
	m.EPM.Set(m.EXP.Value / minutes)
	expStats := fmt.Sprintf("EXP/EPM: %s/%s", prettyNum(m.EXP.Value, false), prettyNum(m.EPM.Value, false))
	m.APM.Set(m.AP.Value / minutes)
	m.KPA.Set(m.Kills.Value / m.AP.Value)
	apStats := fmt.Sprintf("AP/APM/KPA: %s/%.2f/%.2f", prettyNum(m.AP.Value, true), m.APM.Value, m.KPA.Value)
	m.PPPH.Set(m.PP.Value / float64(1000000) / hours)
	ppStats := fmt.Sprintf("PP/PPPH: %.1f/%.2f", m.PP.Value/float64(1000000), m.PPPH.Value)
	var totalTime time.Duration
	for i := range m.clickFloatingDur {
		totalTime += m.clickFloatingDur[i]
	}
	avgFPS := m.clickDuration / time.Duration(m.clickCount)
	floatingFPS := totalTime / time.Duration(len(m.clickFloatingDur))
	avgTime := fmt.Sprintf("FPS/Instant %v/%v", 1000/avgFPS.Milliseconds(), 1000/floatingFPS.Milliseconds())
	m.FPS.Set(float64(1000 / avgFPS.Milliseconds()))
	m.FrameRate.Set(float64(1000 / floatingFPS.Milliseconds()))
	var brokeCount string
	if m.Rescans.Value > 0 {
		brokeCount = fmt.Sprintf(" Rescans %d Broken: %d", int(m.Rescans.Value), int(m.RescansNeeded.Value))
	} else {
		brokeCount = ""
	}
	STATUS = fmt.Sprintf("Hours: %.2f %s %s %s %s%s %s %v    ", hours, killStats, expStats, apStats, ppStats, brokeCount, avgTime, m.IdleDuration.Round(time.Second/10))
}

func (m *Metrics) RecordClick(dur time.Duration) {
	m.clickCount++
	m.clickDuration += dur
	if len(m.clickFloatingDur) >= ColorTimingAvg {
		m.clickFloatingDur = m.clickFloatingDur[1:]
	}
	m.clickFloatingDur = append(m.clickFloatingDur, dur)
}

func (m *Metrics) Init() {
	m.EPM.InitProm()
	m.EXP.InitProm()
	m.KPM.InitProm()
	m.Kills.InitProm()
	m.AP.InitProm()
	m.APM.InitProm()
	m.KPA.InitProm()
	m.PP.InitProm()
	m.PPPH.InitProm()
	m.FPS.InitProm()
	m.FrameRate.InitProm()
	m.Rescans.InitProm()
	m.RescansNeeded.InitProm()
	m.TierKills.InitProm()
}

type MetricGauge struct {
	Value  float64
	Name   string
	Descr  string
	prom   prometheus.Gauge
	inited bool
}

func (m *MetricGauge) InitProm() {
	// Only allow init once
	if m.inited {
		return
	}
	m.prom = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "ngu_" + m.Name,
		Help: m.Descr,
	})
	m.inited = true
}

func (m *MetricGauge) Set(value float64) {
	m.Value = value
	m.prom.Set(value)
}

type MetricCounter struct {
	Value  float64
	Name   string
	Descr  string
	prom   prometheus.Counter
	inited bool
}

func (m *MetricCounter) InitProm() {
	// Only allow init once
	if m.inited {
		return
	}
	m.prom = promauto.NewCounter(prometheus.CounterOpts{
		Name: "ngu_" + m.Name,
		Help: m.Descr,
	})
	m.inited = true
}

func (m *MetricCounter) Add(value float64) {
	m.Value += value
	m.prom.Add(value)
}

func (m *MetricCounter) Inc() {
	m.Value++
	m.prom.Inc()
}

type MetricsMap struct {
	Values map[string]prometheus.Counter
	Name   string
	Prefix string
	Descr  string
	inited bool
	lk     sync.Mutex
}

func (m *MetricsMap) InitProm() {
	// Only allow init once
	if m.inited {
		return
	}
	m.Values = make(map[string]prometheus.Counter)
}

func (m *MetricsMap) Inc(key string) {
	// serialize access to the map
	m.lk.Lock()
	defer m.lk.Unlock()
	if _, found := m.Values[key]; !found {
		prom := promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("ngu_%s_%s_%s", m.Name, m.Prefix, key),
			Help: m.Descr,
		})
		m.Values[key] = prom
	}
	m.Values[key].Inc()
}

func initMetrics() {
	m := Metrics{
		EPM:           MetricGauge{Name: "epm", Descr: "Experience gained per minute"},
		EXP:           MetricCounter{Name: "exp", Descr: "Total experience gained"},
		KPM:           MetricGauge{Name: "kpm", Descr: "Kills per minute"},
		Kills:         MetricCounter{Name: "kills", Descr: "Total tower kills"},
		APM:           MetricGauge{Name: "apm", Descr: "AP earned per minute"},
		KPA:           MetricGauge{Name: "kpa", Descr: "Kills per AP gained"},
		AP:            MetricCounter{Name: "ap", Descr: "Total arbitrary points earned"},
		PP:            MetricCounter{Name: "pp", Descr: "Total progress points gathered"},
		PPPH:          MetricGauge{Name: "ppph", Descr: "Perk Points gained per hour"},
		FPS:           MetricGauge{Name: "fps", Descr: "Average framerate"},
		FrameRate:     MetricGauge{Name: "fps_instant", Descr: "Framerate from the last 10 frames"},
		Rescans:       MetricCounter{Name: "rescans", Descr: "Total kill counter rescans"},
		RescansNeeded: MetricCounter{Name: "rescans_needed", Descr: "Total rescan values that were different post scan"},
		TierKills:     MetricsMap{Name: "kills", Prefix: "tier", Descr: "Total kills per tier"},
	}
	m.Init()
	m.Start = time.Now()
	AppMetrics = m
}
