package workflow

import (
	"ziggy/internal/temporal"
)

func RegisterWorkflows() {
	temporal.AddWorkflow("ZiggyWorkflow", ZiggyWorkflow)
	temporal.AddWorkflow("ChatWorkflow", ChatWorkflow)
	temporal.AddWorkflow("NeedUpdaterWorkflow", NeedUpdaterWorkflow)
}

func RegisterActivities(activities *Activities) {
	temporal.AddActivity("RegeneratePool", activities.RegeneratePool)
}

func RegisterChatActivities(activities *ChatActivities) {
	temporal.AddActivity("GenerateChatResponse", activities.GenerateChatResponse)
	temporal.AddActivity("QueryZiggyState", activities.QueryZiggyState)
}
