package workflow

var educationalMysteries = []Mystery{
	{
		ID:          "vanishing-signal",
		Title:       "The Vanishing Signal",
		Description: "Something strange is happening - messages seem to arrive even when I'm busy doing other things!",
		Track:       "educational",
		Concept:     "Signals & Queries",
		Hints: []string{
			"Messages can arrive even when I'm busy...",
			"It's like having a mailbox that collects letters while you're away",
			"In Temporal, these are called Signals - they queue up and wait",
		},
		Solution: "Signals are async messages that queue in the workflow",
	},
	{
		ID:          "eternal-loop",
		Title:       "The Eternal Loop",
		Description: "I've been around for so long, but somehow I feel fresh and new. How is that possible?",
		Track:       "educational",
		Concept:     "Continue-as-new",
		Hints: []string{
			"Sometimes to go forward, you start fresh...",
			"Like a snake shedding its skin, I keep my memories but get a new body",
			"Temporal calls this continue-as-new - same workflow, fresh history",
		},
		Solution: "Continue-as-new lets long-running workflows reset their history while preserving state",
	},
	{
		ID:          "sleeping-worker",
		Title:       "The Sleeping Worker",
		Description: "Even when I'm sleeping, something keeps checking if I'm okay. What could it be?",
		Track:       "educational",
		Concept:     "Activity heartbeats",
		Hints: []string{
			"Even sleeping, I send signs of life...",
			"It's like a heartbeat monitor at a hospital",
			"Activities can send heartbeats to show they're still working",
		},
		Solution: "Activity heartbeats let Temporal know long-running activities are still alive",
	},
	{
		ID:          "time-traveler",
		Title:       "The Time Traveler",
		Description: "I had the strangest dream - I relived the same moment over and over, exactly the same way each time!",
		Track:       "educational",
		Concept:     "Workflow replay",
		Hints: []string{
			"What if every moment could be replayed exactly?",
			"Temporal records everything like a perfect movie",
			"When I wake up after a crash, I replay my history to get back to where I was",
		},
		Solution: "Workflow replay reconstructs state by re-executing the workflow history deterministically",
	},
	{
		ID:          "parallel-paths",
		Title:       "The Parallel Paths",
		Description: "Sometimes I feel like I'm in multiple places at once, doing different things simultaneously!",
		Track:       "educational",
		Concept:     "Child workflows",
		Hints: []string{
			"I can be in many places at once...",
			"Like having little helpers that do tasks for me",
			"Parent workflows can spawn children that run independently",
		},
		Solution: "Child workflows let you break complex work into independent, manageable pieces",
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
