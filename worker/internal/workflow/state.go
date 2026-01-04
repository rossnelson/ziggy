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

	// Action cooldowns
	CooldownFeed = 30 * time.Second
	CooldownPlay = 60 * time.Second
	CooldownPet  = 10 * time.Second
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

	// Cooldown tracking
	LastFeedTime time.Time `json:"lastFeedTime,omitempty"`
	LastPlayTime time.Time `json:"lastPlayTime,omitempty"`
	LastPetTime  time.Time `json:"lastPetTime,omitempty"`
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

	// Cooldown remaining in seconds (0 = ready)
	FeedCooldown float64 `json:"feedCooldown"`
	PlayCooldown float64 `json:"playCooldown"`
	PetCooldown  float64 `json:"petCooldown"`
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

	if current.HP == 0 {
		current.LastUpdateTime = now
		return current
	}

	// Note: Sleep state is controlled by workflow signals (wake) and time-of-day
	// transitions in the workflow loop, not overridden here

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

	// During egg stage, only bond decays (no fullness/happiness decay)
	age := s.LastUpdateTime.Sub(s.CreatedAt).Seconds()
	isEgg := GetStageForAge(age) == StageEgg

	bondProtection := 0.0
	if s.Bond > 50 {
		bondProtection = (s.Bond - 50) / 100
	}

	if s.Sleeping {
		if !isEgg {
			s.Fullness -= DecayFullnessSleep
			s.Happiness += RecoverHappinessSleep
		}
	} else {
		if !isEgg {
			s.Fullness -= DecayFullnessAwake * (1 - bondProtection)
			s.Happiness -= DecayHappinessAwake * (1 - bondProtection)
		}
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

	// During egg stage, only bond decays (no fullness/happiness decay)
	age := s.LastUpdateTime.Sub(s.CreatedAt).Seconds()
	isEgg := GetStageForAge(age) == StageEgg

	bondProtection := 0.0
	if s.Bond > 50 {
		bondProtection = (s.Bond - 50) / 100
	}

	if s.Sleeping {
		if !isEgg {
			s.Fullness -= DecayFullnessSleep * fraction
			s.Happiness += RecoverHappinessSleep * fraction
		}
	} else {
		if !isEgg {
			s.Fullness -= DecayFullnessAwake * (1 - bondProtection) * fraction
			s.Happiness -= DecayHappinessAwake * (1 - bondProtection) * fraction
		}
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
		Fullness:     s.Fullness,
		Happiness:    s.Happiness,
		Bond:         s.Bond,
		HP:           s.HP,
		Stage:        GetStageForAge(age),
		TimeOfDay:    GetTimeOfDay(now, s.Timezone),
		Sleeping:     s.Sleeping,
		Personality:  s.Personality,
		Message:      s.Message,
		LastAction:   s.LastAction,
		Age:          age,
		Generation:   s.Generation,
		FeedCooldown: cooldownRemaining(s.LastFeedTime, s.GetEffectiveCooldown(ActionFeed), now),
		PlayCooldown: cooldownRemaining(s.LastPlayTime, s.GetEffectiveCooldown(ActionPlay), now),
		PetCooldown:  cooldownRemaining(s.LastPetTime, s.GetEffectiveCooldown(ActionPet), now),
	}
}

func cooldownRemaining(lastTime time.Time, cooldown time.Duration, now time.Time) float64 {
	if lastTime.IsZero() {
		return 0
	}
	elapsed := now.Sub(lastTime)
	if elapsed >= cooldown {
		return 0
	}
	return (cooldown - elapsed).Seconds()
}

func cooldownMultiplier(stat float64) float64 {
	if stat < 20 {
		return 0.25
	}
	if stat < 40 {
		return 0.5
	}
	if stat < 60 {
		return 0.75
	}
	return 1.0
}

func (s *ZiggyState) GetEffectiveCooldown(action Action) time.Duration {
	switch action {
	case ActionFeed:
		return time.Duration(float64(CooldownFeed) * cooldownMultiplier(s.Fullness))
	case ActionPlay:
		return time.Duration(float64(CooldownPlay) * cooldownMultiplier(s.Happiness))
	case ActionPet:
		return time.Duration(float64(CooldownPet) * cooldownMultiplier(s.Bond))
	default:
		return 0
	}
}

// NeedType represents what Ziggy needs most urgently
type NeedType string

const (
	NeedNone      NeedType = ""
	NeedFood      NeedType = "needsFood"
	NeedPlay      NeedType = "needsPlay"
	NeedAffection NeedType = "needsAffection"
	NeedCritical  NeedType = "needsCritical"
)

// GetMostUrgentNeed returns what Ziggy needs most based on current stats
func (s *ZiggyState) GetMostUrgentNeed() NeedType {
	// Don't show needs while sleeping or in tun state
	if s.Sleeping || s.HP == 0 {
		return NeedNone
	}

	// Critical HP takes priority
	if s.HP < 40 {
		return NeedCritical
	}

	// Find the lowest stat that's below threshold
	const threshold = 60.0

	if s.Fullness < threshold && s.Fullness <= s.Happiness && s.Fullness <= s.Bond {
		return NeedFood
	}
	if s.Happiness < threshold && s.Happiness <= s.Bond {
		return NeedPlay
	}
	if s.Bond < threshold {
		return NeedAffection
	}

	return NeedNone
}

// GetMostRecentActionTime returns the most recent time an action was performed
func (s *ZiggyState) GetMostRecentActionTime() time.Time {
	latest := s.LastFeedTime
	if s.LastPlayTime.After(latest) {
		latest = s.LastPlayTime
	}
	if s.LastPetTime.After(latest) {
		latest = s.LastPetTime
	}
	return latest
}
