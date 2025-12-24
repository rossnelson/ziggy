package workflow

import (
	"ziggy/internal/temporal"
)

func RegisterWorkflows() {
	temporal.AddWorkflow("ZiggyWorkflow", ZiggyWorkflow)
}

func RegisterActivities(activities *Activities) {
	temporal.AddActivity("RegeneratePool", activities.RegeneratePool)
}
