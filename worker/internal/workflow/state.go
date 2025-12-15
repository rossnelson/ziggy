package workflow

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
)

type ZiggyState struct {
	Fullness       float64   `json:"fullness"`
	Happiness      float64   `json:"happiness"`
	Bond           float64   `json:"bond"`
	HP             float64   `json:"hp"`
	Stage          Stage     `json:"stage"`
	TimeOfDay      TimeOfDay `json:"timeOfDay"`
	Sleeping       bool      `json:"sleeping"`
	Message        string    `json:"message"`
	LastAction     Action    `json:"lastAction,omitempty"`
	LastActionTime int64     `json:"lastActionTime,omitempty"`
	Age            float64   `json:"age"`
	Generation     int       `json:"generation"`
}

func NewZiggyState() ZiggyState {
	return ZiggyState{
		Fullness:   70,
		Happiness:  70,
		Bond:       50,
		HP:         100,
		Stage:      StageAdult,
		TimeOfDay:  TimeDay,
		Sleeping:   false,
		Message:    "I've survived\nworse than this.\nBut barely.",
		Age:        0,
		Generation: 1,
	}
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
