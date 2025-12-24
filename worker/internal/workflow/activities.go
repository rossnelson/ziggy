package workflow

import (
	"context"
	"time"

	"ziggy/internal/ai"
)

type Activities struct {
	aiClient *ai.Client
}

func NewActivities(aiClient *ai.Client) *Activities {
	return &Activities{aiClient: aiClient}
}

type PoolRegenerationInput struct {
	Personality Personality `json:"personality"`
	Stage       Stage       `json:"stage"`
	Bond        float64     `json:"bond"`
}

type PoolRegenerationOutput struct {
	Pool        *MessagePool `json:"pool"`
	GeneratedAt time.Time    `json:"generatedAt"`
}

func (a *Activities) RegeneratePool(ctx context.Context, input PoolRegenerationInput) (*PoolRegenerationOutput, error) {
	// No AI client - return empty result, workflow will use fallback pools
	if a.aiClient == nil {
		return &PoolRegenerationOutput{
			Pool:        nil,
			GeneratedAt: time.Now(),
		}, nil
	}

	bondDescription := getBondDescription(input.Bond)

	aiInput := ai.PoolGenerationInput{
		Personality:     string(input.Personality),
		Stage:           string(input.Stage),
		BondDescription: bondDescription,
	}

	aiPool, err := a.aiClient.GeneratePool(ctx, aiInput)
	if err != nil {
		// API error - return empty, use fallback pools
		return &PoolRegenerationOutput{
			Pool:        nil,
			GeneratedAt: time.Now(),
		}, nil
	}

	pool := convertAIPool(aiPool)

	return &PoolRegenerationOutput{
		Pool:        pool,
		GeneratedAt: time.Now(),
	}, nil
}

func getBondDescription(bond float64) string {
	if bond >= 80 {
		return "deeply bonded (best friends)"
	}
	if bond >= 60 {
		return "close bond (good friends)"
	}
	if bond >= 40 {
		return "developing bond (getting to know each other)"
	}
	if bond >= 20 {
		return "new acquaintance (still shy)"
	}
	return "barely met (very timid)"
}

func convertAIPool(aiPool *ai.MessagePool) *MessagePool {
	return &MessagePool{
		FeedSuccess:  aiPool.FeedSuccess,
		FeedFull:     aiPool.FeedFull,
		FeedHungry:   aiPool.FeedHungry,
		FeedSleeping: aiPool.FeedSleeping,
		FeedTun:      aiPool.FeedTun,
		FeedCooldown: aiPool.FeedCooldown,

		PlaySuccess:  aiPool.PlaySuccess,
		PlayTired:    aiPool.PlayTired,
		PlayHappy:    aiPool.PlayHappy,
		PlaySleeping: aiPool.PlaySleeping,
		PlayTun:      aiPool.PlayTun,
		PlayCooldown: aiPool.PlayCooldown,

		PetSuccess:  aiPool.PetSuccess,
		PetMaxBond:  aiPool.PetMaxBond,
		PetLowMood:  aiPool.PetLowMood,
		PetSleeping: aiPool.PetSleeping,
		PetTun:      aiPool.PetTun,
		PetCooldown: aiPool.PetCooldown,

		Reviving: aiPool.Reviving,

		IdleHappy:    aiPool.IdleHappy,
		IdleNeutral:  aiPool.IdleNeutral,
		IdleHungry:   aiPool.IdleHungry,
		IdleSad:      aiPool.IdleSad,
		IdleLonely:   aiPool.IdleLonely,
		IdleCritical: aiPool.IdleCritical,
		IdleTun:      aiPool.IdleTun,
		IdleSleeping: aiPool.IdleSleeping,
	}
}
