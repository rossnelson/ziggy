package workflow

var messagesFeedSuccess = []string{
	"Mmm, perfect.\nI can survive\nanother eon.",
	"Delicious.\nAlmost as good\nas cosmic dust.",
	"I needed that.\nThanks, human.",
}

var messagesFeedFull = []string{
	"Too full!\nNow I feel sick.\n*happiness down*",
	"Ugh... stuffed.\nThat made me\nunhappy.",
	"No more!\nOverfeeding\nhurts me.",
}

var messagesFeedHungry = []string{
	"FINALLY.\nI was mass-\nextinction hungry.",
	"Oh thank you.\nI thought I'd\nstarve forever.",
	"Food! Beautiful\nlife-giving food!",
}

var messagesPlaySuccess = []string{
	"Wheee!\nThis is fun!",
	"Again! Again!\nI have energy\nfor eons!",
	"Playing is the\nbest survival\nstrategy.",
}

var messagesPlayTired = []string{
	"I'm too tired...\nmaybe later?",
	"Need food first.\nThen play.",
}

var messagesPlayHappy = []string{
	"More playing!\nI love this!",
	"Best day since\nthe Permian\nextinction!",
}

var messagesPetSuccess = []string{
	"*happy wiggle*\nI like you.",
	"That's nice.\nKeep going.",
	"Mmm...\nright there.",
}

var messagesPetMaxBond = []string{
	"We're already\nbest friends!\nBut okay...",
	"*content sigh*\nI trust you\ncompletely.",
}

var messagesPetLowMood = []string{
	"Thanks...\nI needed that.",
	"*small wiggle*\nYou're kind.",
}

var messagesPetSleeping = []string{
	"Zzz...\n*happy mumble*",
	"*snuggles closer*",
}

var messagesIdle = map[Mood][]string{
	MoodHappy: {
		"Life is good.\nI've survived\nworse.",
		"Did you know\nI can live in\nspace? Cool, right?",
		"Just vibing.\nDurably.",
	},
	MoodNeutral: {
		"...",
		"I'm fine.\nJust existing.",
		"Waiting for\nsomething to\nsurvive.",
	},
	MoodHungry: {
		"My stomach\nis a void.",
		"Feed me?\nPlease?",
		"I'm withering\naway here...",
	},
	MoodSad: {
		"Nobody loves\na tardigrade...",
		"*sad wiggle*",
		"I'm fine.\nEverything is\nfine.",
	},
	MoodLonely: {
		"Is anyone\nthere...?",
		"I miss you.\nCome back soon.",
		"*looks around*\nSo quiet...",
		"Even tardigrades\nneed friends.",
	},
	MoodCritical: {
		"I don't feel\nso good...",
		"Help...",
		"Is this how\nit ends?",
	},
	MoodTun: {
		"*curled up*\n*not responding*",
	},
	MoodSleeping: {
		"Zzz...",
		"*peaceful snoring*",
		"Zzz... cosmic\ndreams... zzz",
	},
}
