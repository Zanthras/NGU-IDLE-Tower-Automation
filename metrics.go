package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Metrics struct {
	EPM       prometheus.Gauge
	EXP       prometheus.Counter
	KPM       prometheus.Gauge
	Kills     prometheus.Counter
	AP        prometheus.Counter
	APM       prometheus.Gauge
	KPA       prometheus.Gauge
	PP        prometheus.Counter
	PPPH      prometheus.Gauge
	FPS       prometheus.Gauge
	FrameRate prometheus.Gauge
}

func initMetrics() {

	m := Metrics{
		EPM: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "ngu_epm",
			Help: "Experience gained per minute",
		}),
		EXP: promauto.NewCounter(prometheus.CounterOpts{
			Name: "ngu_exp",
			Help: "Total experience gained",
		}),
		KPM: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "ngu_kpm",
			Help: "Kills per minute",
		}),
		Kills: promauto.NewCounter(prometheus.CounterOpts{
			Name: "ngu_kills",
			Help: "Total tower kills",
		}),
		APM: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "ngu_apm",
			Help: "AP earned per minute",
		}),
		KPA: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "ngu_kpa",
			Help: "Kills per AP gained",
		}),
		AP: promauto.NewCounter(prometheus.CounterOpts{
			Name: "ngu_ap",
			Help: "Total arbitrary points earned",
		}),
		PP: promauto.NewCounter(prometheus.CounterOpts{
			Name: "ngu_pp",
			Help: "Total progress points gathered",
		}),
		PPPH: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "ngu_ppph",
			Help: "Perk Points gained per hour",
		}),
		FPS: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "ngu_fps",
			Help: "Average framerate",
		}),
		FrameRate: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "ngu_fps_instant",
			Help: "Framerate from the last 10 frames",
		}),
	}
	AppMetrics = m
}
