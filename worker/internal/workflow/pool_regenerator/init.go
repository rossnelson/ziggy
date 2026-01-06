package pool_regenerator

import (
	"fmt"

	"ziggy/internal/registry"
)

func init() {
	registry.RegisterWorkflow(registry.Definition{
		Name:     "PoolRegeneratorWorkflow",
		Workflow: Workflow,
		IDPattern: func(owner string) string {
			return fmt.Sprintf("ziggy-pool-%s", owner)
		},
		NewInput: func(owner, ziggyID, _ string) any {
			return Input{ZiggyWorkflowID: ziggyID}
		},
		AutoStart: true,
	})
}
