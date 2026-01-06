package pool_regenerator

import (
	"fmt"

	"ziggy/internal/registry"
)

func Register() {
	registry.RegisterWorkflow(registry.Definition{
		Name:     "PoolRegeneratorWorkflow",
		Workflow: Workflow,
		IDPattern: func(owner string) string {
			return fmt.Sprintf("ziggy-%s-pool-regenerator", owner)
		},
		NewInput: func(owner, ziggyID, _ string) any {
			return Input{ZiggyWorkflowID: ziggyID}
		},
		AutoStart: true,
	})
}
