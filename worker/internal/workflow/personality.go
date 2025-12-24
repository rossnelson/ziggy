package workflow

import "time"

type Personality string

const (
	PersonalityStoic    Personality = "stoic"
	PersonalityDramatic Personality = "dramatic"
	PersonalityCheerful Personality = "cheerful"
	PersonalitySassy    Personality = "sassy"
	PersonalityShy      Personality = "shy"
)

type CareMetrics struct {
	TotalInteractions int       `json:"totalInteractions"`
	LastInteractionAt time.Time `json:"lastInteractionAt"`
	AvgFullness       float64   `json:"avgFullness"`
	AvgBond           float64   `json:"avgBond"`
}

func (m *CareMetrics) RecordInteraction(fullness, bond float64, now time.Time) {
	m.TotalInteractions++
	m.LastInteractionAt = now

	if m.TotalInteractions == 1 {
		m.AvgFullness = fullness
		m.AvgBond = bond
		return
	}

	alpha := 0.1
	m.AvgFullness = alpha*fullness + (1-alpha)*m.AvgFullness
	m.AvgBond = alpha*bond + (1-alpha)*m.AvgBond
}

func DerivePersonality(metrics CareMetrics, bond float64, now time.Time) Personality {
	if metrics.TotalInteractions == 0 {
		return PersonalityShy
	}

	timeSinceInteraction := now.Sub(metrics.LastInteractionAt)
	neglected := timeSinceInteraction > 2*time.Hour || metrics.TotalInteractions < 10

	if neglected && bond < 40 {
		return PersonalitySassy
	}
	if neglected {
		return PersonalityDramatic
	}
	if bond > 70 && metrics.AvgFullness > 60 {
		return PersonalityCheerful
	}
	if bond < 30 {
		return PersonalityShy
	}
	return PersonalityStoic
}
