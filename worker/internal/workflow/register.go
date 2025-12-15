package workflow

import (
	"ziggy/internal/temporal"
)

func Register() {
	temporal.AddWorkflow("ZiggyWorkflow", ZiggyWorkflow)
}
