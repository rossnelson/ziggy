package need_updater

import (
	"fmt"

	"ziggy/internal/registry"
)

func init() {
	registry.RegisterWorkflow(registry.Definition{
		Name:     "NeedUpdaterWorkflow",
		Workflow: Workflow,
		IDPattern: func(owner string) string {
			return fmt.Sprintf("ziggy-needs-%s", owner)
		},
		NewInput: func(owner, ziggyID, _ string) any {
			return Input{ZiggyWorkflowID: ziggyID, Iteration: 0}
		},
		AutoStart: true,
	})
}
