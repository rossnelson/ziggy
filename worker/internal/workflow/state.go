package workflow

import (
	"math"
	"time"
)

type Stage string

const (
	StageEgg   Stage = "egg"
	StageBaby  Stage = "baby"
	StageTeen  Stage = "teen"
	StageAdult Stage = "adult"
	StageElder Stage = "elder"
)

type Mood string

const (
	MoodHappy    Mood = "happy"
	MoodNeutral  Mood = "neutral"
	MoodHungry   Mood = "hungry"
	MoodSad      Mood = "sad"
	MoodLonely   Mood = "lonely"
	MoodSleeping Mood = "sleeping"
	MoodCritical Mood = "critical"
	MoodTun      Mood = "tun"
)

type TimeOfDay string

const (
	TimeNight TimeOfDay = "night"
	TimeDawn  TimeOfDay = "dawn"
	TimeDay   TimeOfDay = "day"
	TimeDusk  TimeOfDay = "dusk"
)

type Action string

const (
	ActionFeed Action = "feed"
	ActionPlay Action = "play"
	ActionPet  Action = "pet"
	ActionWake Action = "wake"
)

const (
	DecayIntervalSeconds = 10.0

	DecayFullnessAwake  = 2.0
	DecayHappinessAwake = 1.0
	DecayBond           = 0.5

	DecayFullnessSleep    = 1.0
	RecoverHappinessSleep = 0.5

	HPDecayRate   = 2.0
	HPRecoverRate = 1.0

	// Evolution age thresholds (in seconds)
	AgeEggTooBaby   = 60      // 1 minute
	AgeBabyToTeen   = 300     // 5 minutes
	AgeTeenToAdult  = 900     // 15 minutes
	AgeAdultToElder = 3600    // 1 hour
)

type ZiggyState struct {
	Fullness  float64 `json:"fullness"`
	Happiness float64 `json:"happiness"`
	Bond      float64 `json:"bond"`
	HP        float64 `json:"hp"`

	LastUpdateTime time.Time `json:"lastUpdateTime"`
	CreatedAt      time.Time `json:"createdAt"`

	Sleeping bool  `json:"sleeping"`
	Stage    Stage `json:"stage"`

	Message    string `json:"message"`
	LastAction Action `json:"lastAction,omitempty"`

	Timezone   string `json:"timezone"`
	Generation int    `json:"generation"`

	Personality     Personality  `json:"personality"`
	CareMetrics     CareMetrics  `json:"careMetrics"`
	RuntimePool     *MessagePool `json:"runtimePool,omitempty"`
	PoolGeneratedAt time.Time    `json:"poolGeneratedAt,omitempty"`
}

type ZiggyStateResponse struct {
	Fullness  float64 `json:"fullness"`
	Happiness float64 `json:"happiness"`
	Bond      float64 `json:"bond"`
	HP        float64 `json:"hp"`

	Stage       Stage       `json:"stage"`
	TimeOfDay   TimeOfDay   `json:"timeOfDay"`
	Sleeping    bool        `json:"sleeping"`
	Personality Personality `json:"personality"`

	Message    string `json:"message"`
	LastAction Action `json:"lastAction,omitempty"`

	Age        float64 `json:"age"`
	Generation int     `json:"generation"`
}

func NewZiggyState(timezone string) ZiggyState {
	now := time.Now()
	timeOfDay := GetTimeOfDay(now, timezone)

	return ZiggyState{
		Fullness:       70,
		Happiness:      70,
		Bond:           50,
		HP:             100,
		LastUpdateTime: now,
		CreatedAt:      now,
		Sleeping:       timeOfDay == TimeNight,
		Stage:          StageEgg,
		Message:        "*wiggle*\n*wiggle*",
		Timezone:       timezone,
		Generation:     1,
		Personality:    PersonalityShy,
		CareMetrics: CareMetrics{
			TotalInteractions: 0,
			LastInteractionAt: now,
			AvgFullness:       70,
			AvgBond:           50,
		},
	}
}

func GetTimeOfDay(t time.Time, timezone string) TimeOfDay {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		loc = time.UTC
	}
	hour := t.In(loc).Hour()

	if hour >= 22 || hour < 5 {
		return TimeNight
	}
	if hour >= 5 && hour < 8 {
		return TimeDawn
	}
	if hour >= 8 && hour < 18 {
		return TimeDay
	}
	return TimeDusk
}

func (s *ZiggyState) GetMood() Mood {
	if s.HP == 0 {
		return MoodTun
	}
	if s.Sleeping {
		return MoodSleeping
	}
	if s.HP < 20 {
		return MoodCritical
	}
	if s.Fullness < 20 {
		return MoodHungry
	}
	if s.Happiness < 20 {
		return MoodSad
	}
	if s.Bond < 20 {
		return MoodLonely
	}
	if s.Happiness > 70 && s.Fullness > 50 {
		return MoodHappy
	}
	return MoodNeutral
}

func GetStageForAge(ageSeconds float64) Stage {
	if ageSeconds < AgeEggTooBaby {
		return StageEgg
	}
	if ageSeconds < AgeBabyToTeen {
		return StageBaby
	}
	if ageSeconds < AgeTeenToAdult {
		return StageTeen
	}
	if ageSeconds < AgeAdultToElder {
		return StageAdult
	}
	return StageElder
}

func (s *ZiggyState) CalculateCurrentState(now time.Time) ZiggyState {
	current := *s

	elapsed := now.Sub(s.LastUpdateTime).Seconds()
	if elapsed <= 0 {
		return current
	}

	ticks := elapsed / DecayIntervalSeconds
	currentTimeOfDay := GetTimeOfDay(now, s.Timezone)
	shouldBeSleeping := currentTimeOfDay == TimeNight

	if current.HP == 0 {
		current.LastUpdateTime = now
		return current
	}

	if shouldBeSleeping != s.Sleeping {
		current.Sleeping = shouldBeSleeping
	}

	for i := 0.0; i < ticks; i++ {
		current.applyDecayTick()
	}

	fractionalTick := ticks - math.Floor(ticks)
	if fractionalTick > 0 {
		current.applyPartialDecayTick(fractionalTick)
	}

	current.LastUpdateTime = now
	return current
}

func (s *ZiggyState) applyDecayTick() {
	if s.HP == 0 {
		return
	}

	bondProtection := 0.0
	if s.Bond > 50 {
		bondProtection = (s.Bond - 50) / 100
	}

	if s.Sleeping {
		s.Fullness -= DecayFullnessSleep
		s.Happiness += RecoverHappinessSleep
	} else {
		s.Fullness -= DecayFullnessAwake * (1 - bondProtection)
		s.Happiness -= DecayHappinessAwake * (1 - bondProtection)
		s.Bond -= DecayBond
	}

	targetHP := (s.Fullness + s.Happiness + s.Bond) / 3
	if s.HP > targetHP {
		s.HP -= HPDecayRate
	} else if s.HP < targetHP {
		if s.Sleeping {
			s.HP += HPRecoverRate * 1.5
		} else {
			s.HP += HPRecoverRate
		}
	}

	s.Clamp()
}

func (s *ZiggyState) applyPartialDecayTick(fraction float64) {
	if s.HP == 0 {
		return
	}

	bondProtection := 0.0
	if s.Bond > 50 {
		bondProtection = (s.Bond - 50) / 100
	}

	if s.Sleeping {
		s.Fullness -= DecayFullnessSleep * fraction
		s.Happiness += RecoverHappinessSleep * fraction
	} else {
		s.Fullness -= DecayFullnessAwake * (1 - bondProtection) * fraction
		s.Happiness -= DecayHappinessAwake * (1 - bondProtection) * fraction
		s.Bond -= DecayBond * fraction
	}

	s.Clamp()
}

func (s *ZiggyState) Clamp() {
	s.Fullness = clamp(s.Fullness, 0, 100)
	s.Happiness = clamp(s.Happiness, 0, 100)
	s.Bond = clamp(s.Bond, 0, 100)
	s.HP = clamp(s.HP, 0, 100)
}

func clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func (s *ZiggyState) ToResponse(now time.Time) ZiggyStateResponse {
	age := now.Sub(s.CreatedAt).Seconds()
	return ZiggyStateResponse{
		Fullness:    s.Fullness,
		Happiness:   s.Happiness,
		Bond:        s.Bond,
		HP:          s.HP,
		Stage:       GetStageForAge(age),
		TimeOfDay:   GetTimeOfDay(now, s.Timezone),
		Sleeping:    s.Sleeping,
		Personality: s.Personality,
		Message:     s.Message,
		LastAction:  s.LastAction,
		Age:         age,
		Generation:  s.Generation,
	}
}
