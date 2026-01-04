package workflow

var educationalMysteries = []Mystery{
	{
		ID:          "signals-queries",
		Title:       "Signals & Queries",
		Description: "How workflows communicate with the outside world",
		Track:       "educational",
		Concept:     "Signals & Queries",
		Summary:     "Signals are async messages sent INTO a workflow - like when you click Feed or Pet, that's a signal to me! They queue up and I process them in order. Queries are read-only - they let you check my state without changing anything. The UI uses queries to poll my stats!",
		DocsURL:     "https://docs.temporal.io/workflows#signals-and-queries",
		Hints:       []string{},
		Solution:    "",
	},
	{
		ID:          "continue-as-new",
		Title:       "Continue-As-New",
		Description: "How workflows run forever without running out of memory",
		Track:       "educational",
		Concept:     "Continue-as-new",
		Summary:     "Temporal stores every event in workflow history, but history can't grow forever! Continue-as-new atomically starts a fresh execution while passing along important state. I use this myself - when my history gets too long, I continue-as-new with my stats preserved.",
		DocsURL:     "https://docs.temporal.io/workflows#continue-as-new",
		Hints:       []string{},
		Solution:    "",
	},
	{
		ID:          "activities",
		Title:       "Activities",
		Description: "Where the real work happens outside workflows",
		Track:       "educational",
		Concept:     "Activities",
		Summary:     "Workflows must be deterministic - they can't call APIs or databases directly. Activities are where side effects happen: API calls, database writes, sending emails. When I generate AI responses, that happens in an Activity, not in my workflow code!",
		DocsURL:     "https://docs.temporal.io/activities",
		Hints:       []string{},
		Solution:    "",
	},
	{
		ID:          "replay",
		Title:       "Workflow Replay",
		Description: "How Temporal makes workflows fault-tolerant",
		Track:       "educational",
		Concept:     "Workflow replay",
		Summary:     "When a worker crashes, how does the workflow recover? Temporal replays the entire history - re-running your code but returning cached results for completed activities. This is why workflows must be deterministic - same inputs must produce same decisions!",
		DocsURL:     "https://docs.temporal.io/workflows#replays",
		Hints:       []string{},
		Solution:    "",
	},
	{
		ID:          "child-workflows",
		Title:       "Child Workflows",
		Description: "Breaking complex workflows into smaller pieces",
		Track:       "educational",
		Concept:     "Child workflows",
		Summary:     "Some tasks are complex enough to be their own workflow with their own history and lifecycle. Parent workflows can spawn children and wait for results or let them run independently. Child workflows can even outlive their parent - useful for long-running subtasks!",
		DocsURL:     "https://docs.temporal.io/workflows#child-workflows",
		Hints:       []string{},
		Solution:    "",
	},
	{
		ID:          "timers",
		Title:       "Timers & Sleep",
		Description: "Durable scheduling that survives crashes",
		Track:       "educational",
		Concept:     "Timers",
		Summary:     "Temporal timers are durable - if a worker crashes during a 6-hour sleep, the timer still fires on time! I use timers to regenerate my message pool every 6 hours. Unlike regular sleep(), Temporal timers survive restarts and failures.",
		DocsURL:     "https://docs.temporal.io/workflows#timers",
		Hints:       []string{},
		Solution:    "",
	},
	{
		ID:          "task-queues",
		Title:       "Task Queues",
		Description: "How work gets distributed to workers",
		Track:       "educational",
		Concept:     "Task Queues",
		Summary:     "Task Queues are how Temporal routes work to the right workers. Workers poll specific queues for tasks. You can have different queues for different workloads - like separating CPU-heavy AI generation from lightweight state queries!",
		DocsURL:     "https://docs.temporal.io/workers#task-queue",
		Hints:       []string{},
		Solution:    "",
	},
}

var funMysteries = []Mystery{
	{
		ID:          "missing-snack",
		Title:       "The Missing Snack",
		Description: "My favorite cosmic crumbs have gone missing! I left them right here...",
		Track:       "fun",
		Hints: []string{
			"The crumbs lead somewhere cold...",
			"I remember seeing something sparkly near the fridge",
			"Wait, do tardigrades even have fridges in space?",
		},
		Solution: "The cosmic crumbs were eaten by a passing comet",
	},
	{
		ID:          "cosmic-radio",
		Title:       "The Cosmic Radio",
		Description: "I keep hearing strange beeps and boops from far away. Is someone trying to talk to me?",
		Track:       "fun",
		Hints: []string{
			"Beep... boop... is anyone out there?",
			"The pattern seems to repeat every 8 seconds",
			"It sounds like someone saying 'hello' in binary!",
		},
		Solution: "It was a friendly space probe saying hello",
	},
	{
		ID:          "dream-maze",
		Title:       "The Dream Maze",
		Description: "In my dreams, I keep finding myself in a strange maze where time flows differently...",
		Track:       "fun",
		Hints: []string{
			"In dreams, time flows differently...",
			"The walls seem to be made of stardust",
			"Following the warmest path leads to the exit",
		},
		Solution: "The dream maze was actually a memory of surviving a supernova",
	},
}

func GetMystery(id string, track string) *Mystery {
	mysteries := funMysteries
	if track == "educational" {
		mysteries = educationalMysteries
	}

	for i := range mysteries {
		if mysteries[i].ID == id {
			return &mysteries[i]
		}
	}
	return nil
}

func GetRandomMystery(track string) *Mystery {
	mysteries := funMysteries
	if track == "educational" {
		mysteries = educationalMysteries
	}

	if len(mysteries) == 0 {
		return nil
	}

	// Return first mystery for now (could randomize later)
	return &mysteries[0]
}

func GetAvailableMysteries(track string, solved []string) []Mystery {
	mysteries := funMysteries
	if track == "educational" {
		mysteries = educationalMysteries
	}

	solvedSet := make(map[string]bool)
	for _, id := range solved {
		solvedSet[id] = true
	}

	var available []Mystery
	for _, m := range mysteries {
		if !solvedSet[m.ID] {
			available = append(available, m)
		}
	}
	return available
}
