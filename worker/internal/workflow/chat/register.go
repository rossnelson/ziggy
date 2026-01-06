package chat

import (
	"fmt"

	"ziggy/internal/ai"
	"ziggy/internal/registry"
)

func Register() {
	registry.RegisterWorkflow(registry.Definition{
		Name:     "ChatWorkflow",
		Workflow: Workflow,
		IDPattern: func(owner string) string {
			return fmt.Sprintf("ziggy-chat-%s", owner)
		},
		NewInput: func(owner, ziggyID, _ string) any {
			return Input{Owner: owner, ZiggyID: ziggyID, Track: "fun"}
		},
		AutoStart: true,
	})

	aiClient := ai.NewClient()
	activities := NewActivities(aiClient)
	registry.RegisterActivity(registry.ActivityDef{
		Name:     "ProcessChatMessage",
		Activity: activities.ProcessChatMessage,
	})
	registry.RegisterActivity(registry.ActivityDef{
		Name:     "QueryZiggyState",
		Activity: activities.QueryZiggyState,
	})
}
