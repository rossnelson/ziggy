package ziggy

import (
	"fmt"

	"ziggy/internal/ai"
	"ziggy/internal/registry"
)

func Register() {
	// Register workflow (Weight 100 ensures dependent workflows start first)
	registry.RegisterWorkflow(registry.Definition{
		Name:     "ZiggyWorkflow",
		Workflow: Workflow,
		IDPattern: func(owner string) string {
			return fmt.Sprintf("ziggy-%s", owner)
		},
		NewInput: func(owner, _, tz string) any {
			return Input{Owner: owner, Timezone: tz, Generation: 1}
		},
		AutoStart: true,
		Weight:    100,
		Primary:   true,
	})

	// Register activities
	aiClient := ai.NewClient()
	activities := NewActivities(aiClient)
	registry.RegisterActivity(registry.ActivityDef{
		Name:     "ProcessAction",
		Activity: activities.ProcessAction,
	})
	registry.RegisterActivity(registry.ActivityDef{
		Name:     "RegeneratePool",
		Activity: activities.RegeneratePool,
	})
}
