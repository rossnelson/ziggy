package ziggy

import (
	"testing"
)

func TestGetMostUrgentNeed(t *testing.T) {
	tests := []struct {
		name     string
		state    ZiggyState
		expected NeedType
	}{
		{
			name: "no need when stats are high",
			state: ZiggyState{
				Fullness:  70,
				Happiness: 70,
				Bond:      70,
				HP:        80,
				Sleeping:  false,
			},
			expected: NeedNone,
		},
		{
			name: "needs food when fullness is lowest",
			state: ZiggyState{
				Fullness:  50, // Below 60 threshold
				Happiness: 70,
				Bond:      70,
				HP:        60,
				Sleeping:  false,
			},
			expected: NeedFood,
		},
		{
			name: "needs play when happiness is lowest",
			state: ZiggyState{
				Fullness:  70,
				Happiness: 50, // Below 60 threshold
				Bond:      70,
				HP:        60,
				Sleeping:  false,
			},
			expected: NeedPlay,
		},
		{
			name: "needs affection when bond is lowest",
			state: ZiggyState{
				Fullness:  70,
				Happiness: 70,
				Bond:      50, // Below 60 threshold
				HP:        60,
				Sleeping:  false,
			},
			expected: NeedAffection,
		},
		{
			name: "critical takes priority over other needs",
			state: ZiggyState{
				Fullness:  10,
				Happiness: 10,
				Bond:      10,
				HP:        30, // Below 40 critical threshold
				Sleeping:  false,
			},
			expected: NeedCritical,
		},
		{
			name: "no need when sleeping",
			state: ZiggyState{
				Fullness:  10,
				Happiness: 10,
				Bond:      10,
				HP:        50,
				Sleeping:  true,
			},
			expected: NeedNone,
		},
		{
			name: "no need when in tun state",
			state: ZiggyState{
				Fullness:  10,
				Happiness: 10,
				Bond:      10,
				HP:        0,
				Sleeping:  false,
			},
			expected: NeedNone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.state.GetMostUrgentNeed()
			if got != tt.expected {
				t.Errorf("GetMostUrgentNeed() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestPoolSelectorPicksNeedMessages(t *testing.T) {
	personalities := []Personality{
		PersonalityStoic,
		PersonalityDramatic,
		PersonalityCheerful,
		PersonalitySassy,
		PersonalityShy,
	}

	needCategories := []string{
		"needsFood",
		"needsPlay",
		"needsAffection",
		"needsCritical",
	}

	for _, p := range personalities {
		pool := GetFallbackPool(p)
		selector := NewPoolSelector(nil, pool, poolStoic)

		for _, category := range needCategories {
			t.Run(string(p)+"/"+category, func(t *testing.T) {
				msg := selector.Pick(category)
				if msg == "" {
					t.Errorf("Pool %s missing message for category %s", p, category)
				}
			})
		}
	}
}
