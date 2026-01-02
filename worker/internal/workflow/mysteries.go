package workflow

var educationalMysteries = []Mystery{
	{
		ID:          "signals-queries",
		Title:       "Signals & Queries",
		Description: "Learn how workflows communicate with the outside world!",
		Track:       "educational",
		Concept:     "Signals & Queries",
		Hints: []string{
			"Signals are async messages sent TO a workflow - like when you click Feed or Pet, that's a signal to me!",
			"Signals queue up and wait - even if I'm busy, I'll process them in order when I'm ready.",
			"Queries are read-only requests - they let you check my state without changing anything. The UI polls my state this way!",
		},
		Solution: "Signals send data into workflows asynchronously, while Queries read workflow state without side effects. Together they enable external communication with running workflows.",
	},
	{
		ID:          "continue-as-new",
		Title:       "Continue-As-New",
		Description: "Discover how workflows can run forever without running out of memory!",
		Track:       "educational",
		Concept:     "Continue-as-new",
		Hints: []string{
			"Temporal stores every event in workflow history - but history can't grow forever!",
			"Continue-as-new starts a fresh workflow execution while passing along important state.",
			"I use this myself! When my history gets too long, I continue-as-new with my current stats preserved.",
		},
		Solution: "Continue-as-new prevents unbounded history growth by atomically starting a new execution with carried-over state, enabling truly long-running workflows.",
	},
	{
		ID:          "activities",
		Title:       "Activities",
		Description: "Learn about the workers that do the real work outside workflows!",
		Track:       "educational",
		Concept:     "Activities",
		Hints: []string{
			"Workflows must be deterministic - they can't call APIs or databases directly.",
			"Activities are where side effects happen - API calls, database writes, sending emails.",
			"When I need to generate AI responses, that happens in an Activity, not in my workflow code!",
		},
		Solution: "Activities encapsulate non-deterministic operations. Workflows orchestrate when activities run, while activities do the actual external work.",
	},
	{
		ID:          "replay",
		Title:       "Workflow Replay",
		Description: "See how Temporal makes workflows fault-tolerant through replay!",
		Track:       "educational",
		Concept:     "Workflow replay",
		Hints: []string{
			"When a worker crashes, the workflow needs to recover - but how does it know where it was?",
			"Temporal replays the entire history - re-running your code but skipping completed activities.",
			"This is why workflows must be deterministic - the same inputs must produce the same decisions!",
		},
		Solution: "Replay re-executes workflow code against recorded history. Completed activities return cached results, so the workflow reaches its previous state and continues.",
	},
	{
		ID:          "child-workflows",
		Title:       "Child Workflows",
		Description: "Learn how to break complex workflows into smaller pieces!",
		Track:       "educational",
		Concept:     "Child workflows",
		Hints: []string{
			"Some tasks are complex enough to be their own workflow - with their own history and lifecycle.",
			"Parent workflows can spawn child workflows and wait for results or let them run independently.",
			"Child workflows can even outlive their parent! This is useful for long-running subtasks.",
		},
		Solution: "Child workflows provide modularity and independent failure domains. They're ideal for sub-processes that need their own history or may need to be managed separately.",
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
