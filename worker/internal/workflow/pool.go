package workflow

import "math/rand"

type MessagePool struct {
	FeedSuccess  []string `json:"feedSuccess"`
	FeedFull     []string `json:"feedFull"`
	FeedHungry   []string `json:"feedHungry"`
	FeedSleeping []string `json:"feedSleeping"`
	FeedTun      []string `json:"feedTun"`
	FeedCooldown []string `json:"feedCooldown"`

	PlaySuccess  []string `json:"playSuccess"`
	PlayTired    []string `json:"playTired"`
	PlayHappy    []string `json:"playHappy"`
	PlaySleeping []string `json:"playSleeping"`
	PlayTun      []string `json:"playTun"`
	PlayCooldown []string `json:"playCooldown"`

	PetSuccess   []string `json:"petSuccess"`
	PetMaxBond   []string `json:"petMaxBond"`
	PetLowMood   []string `json:"petLowMood"`
	PetSleeping  []string `json:"petSleeping"`
	PetTun       []string `json:"petTun"`
	PetCooldown  []string `json:"petCooldown"`

	Reviving []string `json:"reviving"`

	IdleHappy    []string `json:"idleHappy"`
	IdleNeutral  []string `json:"idleNeutral"`
	IdleHungry   []string `json:"idleHungry"`
	IdleSad      []string `json:"idleSad"`
	IdleLonely   []string `json:"idleLonely"`
	IdleCritical []string `json:"idleCritical"`
	IdleTun      []string `json:"idleTun"`
	IdleSleeping []string `json:"idleSleeping"`

	// Need-based coaxing messages
	NeedsFood      []string `json:"needsFood"`
	NeedsPlay      []string `json:"needsPlay"`
	NeedsAffection []string `json:"needsAffection"`
	NeedsCritical  []string `json:"needsCritical"`
}

type PoolSelector struct {
	runtime  *MessagePool
	fallback *MessagePool
	generic  *MessagePool
}

func NewPoolSelector(runtime, fallback, generic *MessagePool) *PoolSelector {
	return &PoolSelector{
		runtime:  runtime,
		fallback: fallback,
		generic:  generic,
	}
}

func (s *PoolSelector) pickFrom(runtime, fallback, generic []string) string {
	if len(runtime) > 0 {
		return runtime[rand.Intn(len(runtime))]
	}
	if len(fallback) > 0 {
		return fallback[rand.Intn(len(fallback))]
	}
	if len(generic) > 0 {
		return generic[rand.Intn(len(generic))]
	}
	return ""
}

func (s *PoolSelector) get(category string) []string {
	if s.runtime != nil {
		if msgs := s.getFromPool(s.runtime, category); len(msgs) > 0 {
			return msgs
		}
	}
	if s.fallback != nil {
		if msgs := s.getFromPool(s.fallback, category); len(msgs) > 0 {
			return msgs
		}
	}
	if s.generic != nil {
		return s.getFromPool(s.generic, category)
	}
	return nil
}

func (s *PoolSelector) getFromPool(pool *MessagePool, category string) []string {
	switch category {
	case "feedSuccess":
		return pool.FeedSuccess
	case "feedFull":
		return pool.FeedFull
	case "feedHungry":
		return pool.FeedHungry
	case "feedSleeping":
		return pool.FeedSleeping
	case "feedTun":
		return pool.FeedTun
	case "feedCooldown":
		return pool.FeedCooldown
	case "playSuccess":
		return pool.PlaySuccess
	case "playTired":
		return pool.PlayTired
	case "playHappy":
		return pool.PlayHappy
	case "playSleeping":
		return pool.PlaySleeping
	case "playTun":
		return pool.PlayTun
	case "playCooldown":
		return pool.PlayCooldown
	case "petSuccess":
		return pool.PetSuccess
	case "petMaxBond":
		return pool.PetMaxBond
	case "petLowMood":
		return pool.PetLowMood
	case "petSleeping":
		return pool.PetSleeping
	case "petTun":
		return pool.PetTun
	case "petCooldown":
		return pool.PetCooldown
	case "reviving":
		return pool.Reviving
	case "idleHappy":
		return pool.IdleHappy
	case "idleNeutral":
		return pool.IdleNeutral
	case "idleHungry":
		return pool.IdleHungry
	case "idleSad":
		return pool.IdleSad
	case "idleLonely":
		return pool.IdleLonely
	case "idleCritical":
		return pool.IdleCritical
	case "idleTun":
		return pool.IdleTun
	case "idleSleeping":
		return pool.IdleSleeping
	case "needsFood":
		return pool.NeedsFood
	case "needsPlay":
		return pool.NeedsPlay
	case "needsAffection":
		return pool.NeedsAffection
	case "needsCritical":
		return pool.NeedsCritical
	default:
		return nil
	}
}

func (s *PoolSelector) Pick(category string) string {
	msgs := s.get(category)
	if len(msgs) == 0 {
		return ""
	}
	return msgs[rand.Intn(len(msgs))]
}
