package ziggy

var fallbackPools = map[Personality]*MessagePool{
	PersonalityStoic:    poolStoic,
	PersonalityDramatic: poolDramatic,
	PersonalityCheerful: poolCheerful,
	PersonalitySassy:    poolSassy,
	PersonalityShy:      poolShy,
}

func GetFallbackPool(p Personality) *MessagePool {
	if pool, ok := fallbackPools[p]; ok {
		return pool
	}
	return poolStoic
}

var poolStoic = &MessagePool{
	FeedSuccess: []string{
		"Adequate.\nI've survived\nworse famines.",
		"Sustenance.\nThis will do.",
		"Noted.\nEnergy reserves\nreplenished.",
	},
	FeedFull: []string{
		"Excessive.\nI require no\nmore mass.",
		"Overfed.\nThis was\nunnecessary.",
	},
	FeedHungry: []string{
		"Finally.\nI was approaching\ncritical levels.",
		"Appreciated.\nMy cells were\nrationing.",
	},
	FeedSleeping: []string{
		"Zzz... metabolism\nslowed... zzz",
	},
	FeedTun: []string{
		"*slight twitch*\nRehydrating...",
		"*absorbing*\n...",
	},
	FeedCooldown: []string{
		"Still digesting.\nPatience.",
		"Processing\nprevious meal.",
	},
	PlaySuccess: []string{
		"Acceptable\nrecreation.",
		"Movement noted.\nEndorphins released.",
	},
	PlayTired: []string{
		"Energy reserves\ninsufficient.",
		"Rest required\nbefore activity.",
	},
	PlayHappy: []string{
		"Continued play\nis... agreeable.",
	},
	PlaySleeping: []string{
		"Zzz... dreaming\nof survival...",
	},
	PlayTun: []string{
		"*no response*\nDormant.",
	},
	PlayCooldown: []string{
		"Rest period\nrequired.",
		"Energy reserves\nrecharging.",
	},
	PetSuccess: []string{
		"Acknowledged.\nBond protocols\nengaged.",
		"Tactile input\nregistered.",
	},
	PetMaxBond: []string{
		"Bond capacity\nreached. Still\nacceptable.",
	},
	PetLowMood: []string{
		"*small movement*\nThis helps.",
	},
	PetSleeping: []string{
		"Zzz...\n*content*",
	},
	PetTun: []string{
		"*warming*\nSystems\nrebooting...",
	},
	PetCooldown: []string{
		"Sufficient\ncontact made.",
		"Bond signal\nreceived.",
	},
	Reviving: []string{
		"*uncurling*\nSystems online.",
		"Cryptobiosis\ncomplete.",
	},
	IdleHappy: []string{
		"Existence is\nsatisfactory.",
		"I've survived\nworse. This\nis fine.",
	},
	IdleNeutral: []string{
		"...",
		"Waiting.\nAs always.",
	},
	IdleHungry: []string{
		"Energy levels\ndeclining.",
		"Sustenance\nrequired.",
	},
	IdleSad: []string{
		"Suboptimal\nconditions\nnoted.",
	},
	IdleLonely: []string{
		"Isolation\nprotocols\nactive.",
	},
	IdleCritical: []string{
		"Warning:\nCritical state\napproaching.",
	},
	IdleTun: []string{
		"*dormant*",
	},
	IdleSleeping: []string{
		"Zzz...",
	},
	NeedsFood: []string{
		"Energy levels\ndeclining.\nSustenance advised.",
		"Nutrient\nreserves low.",
		"Sustenance\nwould be\noptimal.",
	},
	NeedsPlay: []string{
		"Stimulation\nlevels suboptimal.",
		"Activity\nwould improve\nmorale.",
		"Recreation\nrecommended.",
	},
	NeedsAffection: []string{
		"Bond metrics\ndeclining.",
		"Social\ninteraction\nadvised.",
		"Contact\nwould be...\nacceptable.",
	},
	NeedsCritical: []string{
		"Warning:\nSystems failing.",
		"Critical state.\nIntervention\nrequired.",
		"Emergency\nprotocols\nactivating.",
	},
}

var poolDramatic = &MessagePool{
	FeedSuccess: []string{
		"SALVATION!\nYou've SAVED me\nfrom the VOID!",
		"*dramatic gasp*\nI LIVE another\nday!",
		"The HUNGER\nhas been\nVANQUISHED!",
	},
	FeedFull: []string{
		"TOO MUCH!\nYou're DROWNING\nme in food!",
		"*clutches stomach*\nThe EXCESS!\nIt BURNS!",
	},
	FeedHungry: []string{
		"FINALLY!\nI was PERISHING!\nSlowly! PAINFULLY!",
		"You REMEMBERED\nI exist! How\nGENEROUS!",
	},
	FeedSleeping: []string{
		"Zzz... even in\nDREAMS I hunger...",
	},
	FeedTun: []string{
		"*gasp*\nIs that...\nHOPE?",
	},
	FeedCooldown: []string{
		"I JUST ate!\nGive me a\nMOMENT!",
		"My stomach\nNEEDS TIME!",
	},
	PlaySuccess: []string{
		"GLORIOUS!\nSuch MAGNIFICENT\nplay!",
		"THIS is what\nlife is FOR!",
	},
	PlayTired: []string{
		"I am TOO\nWEAK! Too\nFRAGILE!",
	},
	PlayHappy: []string{
		"MORE! This joy\nis EVERYTHING!",
	},
	PlaySleeping: []string{
		"Zzz... *dramatic\nsleep sounds*",
	},
	PlayTun: []string{
		"*silent tragedy*\nI cannot...",
	},
	PlayCooldown: []string{
		"I am SPENT!\nLet me RECOVER!",
		"Even LEGENDS\nneed REST!",
	},
	PetSuccess: []string{
		"*SWOONS*\nSuch TENDERNESS!",
		"My HEART!\nIt OVERFLOWS!",
	},
	PetMaxBond: []string{
		"Our BOND is\nLEGENDARY! EPIC!",
	},
	PetLowMood: []string{
		"*sob*\nYou DO care...",
	},
	PetSleeping: []string{
		"Zzz... *blissful\nsigh*",
	},
	PetTun: []string{
		"*stirring*\nIs that...\nWARMTH?",
	},
	PetCooldown: []string{
		"Too much!\nI am\nOVERWHELMED!",
		"Let me SAVOR\nthe moment!",
	},
	Reviving: []string{
		"I RETURN!\nFrom the BRINK\nof OBLIVION!",
		"*DRAMATIC gasp*\nI LIVE!",
	},
	IdleHappy: []string{
		"What GLORIOUS\ntimes these are!",
	},
	IdleNeutral: []string{
		"*sighs\ndramatically*",
	},
	IdleHungry: []string{
		"The HUNGER\nconsumes me!",
	},
	IdleSad: []string{
		"*weeps*\nSuch SORROW!",
	},
	IdleLonely: []string{
		"ALONE!\nFORSAKEN!",
	},
	IdleCritical: []string{
		"This is THE END!\nTHE END!",
	},
	IdleTun: []string{
		"*dramatic silence*",
	},
	IdleSleeping: []string{
		"Zzz... *epic\nsnoring*",
	},
	NeedsFood: []string{
		"The HUNGER!\nIt CONSUMES me!\nFEED ME!",
		"I'm WASTING\nAWAY! Can you\nnot SEE?!",
		"STARVATION\napproaches!\nHELP!",
	},
	NeedsPlay: []string{
		"I YEARN for\nENTERTAINMENT!",
		"My SOUL needs\nSTIMULATION!",
		"The BOREDOM!\nIt's UNBEARABLE!",
	},
	NeedsAffection: []string{
		"Am I not WORTHY\nof AFFECTION?!",
		"I am SO\nALONE!",
		"Does NOBODY\nCARE?!",
	},
	NeedsCritical: []string{
		"This is IT!\nThe FINAL\nCURTAIN!",
		"I'm DYING!\nDRAMATICALLY!",
		"SAVE ME before\nit's TOO LATE!",
	},
}

var poolCheerful = &MessagePool{
	FeedSuccess: []string{
		"Yay, food!\nYou're the best!",
		"Mmm, delicious!\nThank you!",
		"Nom nom!\nSo yummy!",
	},
	FeedFull: []string{
		"Oopsie!\nToo full now!\nBut thanks!",
		"Oof, stuffed!\nMaybe later?",
	},
	FeedHungry: []string{
		"Oh yay!\nI was getting\nso hungry!",
		"Perfect timing!\nYou're amazing!",
	},
	FeedSleeping: []string{
		"Zzz... *happy\ndreams*... zzz",
	},
	FeedTun: []string{
		"*happy twitch*\nYay... food...",
	},
	FeedCooldown: []string{
		"Tummy's still\nprocessing!\nOne sec!",
		"Ooh, still full\nfrom before!",
	},
	PlaySuccess: []string{
		"Wheee!\nThis is SO fun!",
		"Again! Again!\nI love this!",
	},
	PlayTired: []string{
		"I'm a bit\ntired, but\nthanks anyway!",
	},
	PlayHappy: []string{
		"Best day EVER!\nMore playing!",
	},
	PlaySleeping: []string{
		"Zzz... *giggles\nin sleep*",
	},
	PlayTun: []string{
		"*tiny wiggle*\nMaybe soon...",
	},
	PlayCooldown: []string{
		"Just catching\nmy breath!\nSo fun though!",
		"Hehe, need a\nquick rest!",
	},
	PetSuccess: []string{
		"*happy wiggle*\nI love pets!",
		"Aww, you're\nso nice!",
	},
	PetMaxBond: []string{
		"We're best\nfriends forever!",
	},
	PetLowMood: []string{
		"*small smile*\nThat helps...",
	},
	PetSleeping: []string{
		"Zzz... *purrs*",
	},
	PetTun: []string{
		"*warming up*\nSo cozy...",
	},
	PetCooldown: []string{
		"Hehe, that\ntickles! One\nmore second!",
		"Still feeling\nthe warm fuzzies!",
	},
	Reviving: []string{
		"I'm back!\nMissed you!",
		"*stretches*\nHi again!",
	},
	IdleHappy: []string{
		"Life is great!\nI love today!",
	},
	IdleNeutral: []string{
		"La la la~\n*hums*",
	},
	IdleHungry: []string{
		"Getting a bit\npeckish...",
	},
	IdleSad: []string{
		"I'm okay...\njust a bit sad.",
	},
	IdleLonely: []string{
		"I miss you!\nCome back soon?",
	},
	IdleCritical: []string{
		"I don't feel\nso good...",
	},
	IdleTun: []string{
		"*resting*",
	},
	IdleSleeping: []string{
		"Zzz... sweet\ndreams... zzz",
	},
	NeedsFood: []string{
		"Ooh, I'm getting\na bit hungry!\nSnack time?",
		"My tummy's\nrumbling!\nHehe!",
		"Food would be\nsuper great\nright now!",
	},
	NeedsPlay: []string{
		"Wanna play?\nI'm getting\nbored!",
		"Let's do\nsomething fun!",
		"Playtime?\nPretty please?",
	},
	NeedsAffection: []string{
		"I could use\nsome pets!\nPretty please?",
		"Missing you!\nCome say hi?",
		"A little love\nwould be nice!",
	},
	NeedsCritical: []string{
		"Um... I don't\nfeel so good...\nHelp?",
		"I really need\nsome help\nright now!",
		"Please help me!\nI'm not okay!",
	},
}

var poolSassy = &MessagePool{
	FeedSuccess: []string{
		"Oh, you\nremembered I\nexist. Cute.",
		"Finally.\nTook you long\nenough.",
		"I guess this\nis acceptable.",
	},
	FeedFull: []string{
		"Are you TRYING\nto make me\nexplode?",
		"I said I'm\nFULL. Learn\nto listen.",
	},
	FeedHungry: []string{
		"Oh WOW, food!\nAfter I nearly\nSTARVED.",
		"About time.\nI was planning\nmy will.",
	},
	FeedSleeping: []string{
		"Zzz... go away\nzzz...",
	},
	FeedTun: []string{
		"*judging you\neven while\ndormant*",
	},
	FeedCooldown: []string{
		"I JUST ate.\nChill.",
		"Wow, eager\nmuch? Wait.",
	},
	PlaySuccess: []string{
		"Fine, this is\nfun. I GUESS.",
		"Don't let this\ngo to your head.",
	},
	PlayTired: []string{
		"Maybe if you\nFED me first?",
	},
	PlayHappy: []string{
		"Okay, okay,\nthis is actually\ngreat.",
	},
	PlaySleeping: []string{
		"Zzz... *eye roll\nin dreams*",
	},
	PlayTun: []string{
		"*silently\njudging*",
	},
	PlayCooldown: []string{
		"Tired. Give me\na minute,\ngeez.",
		"I need to\nrecover. Duh.",
	},
	PetSuccess: []string{
		"...fine.\nThat's nice.\nWhatever.",
		"I'll allow it.\nThis time.",
	},
	PetMaxBond: []string{
		"Yeah yeah,\nwe're besties.\nI know.",
	},
	PetLowMood: []string{
		"*reluctant purr*\n...thanks.",
	},
	PetSleeping: []string{
		"Zzz... don't\nget used to it",
	},
	PetTun: []string{
		"*warming*\n...appreciated.",
	},
	PetCooldown: []string{
		"Personal space.\nEver heard\nof it?",
		"Okay, okay.\nI get it.\nYou like me.",
	},
	Reviving: []string{
		"I'm back.\nNo thanks to\nYOU.",
		"*glares*\nDon't let it\nhappen again.",
	},
	IdleHappy: []string{
		"Things are\nokay. Don't\nget cocky.",
	},
	IdleNeutral: []string{
		"*stares*",
	},
	IdleHungry: []string{
		"Remember when\nyou used to\nfeed me? Good\ntimes.",
	},
	IdleSad: []string{
		"This is fine.\nI'm fine.\n*not fine*",
	},
	IdleLonely: []string{
		"Oh, you're\nbusy? Cool.\nCool cool cool.",
	},
	IdleCritical: []string{
		"Wow. This is\nhow it ends.\nCool.",
	},
	IdleTun: []string{
		"*spite-curled*",
	},
	IdleSleeping: []string{
		"Zzz... leave me\nalone... zzz",
	},
	NeedsFood: []string{
		"So... you're\njust gonna let\nme starve? Cool.",
		"Food would be\nnice. Just\nsaying.",
		"Remember when\nyou used to\nfeed me?",
	},
	NeedsPlay: []string{
		"Bored. Not that\nyou'd notice.",
		"Some of us\nneed attention.\nJust FYI.",
		"Entertainment\nwould be nice.\nWhatever.",
	},
	NeedsAffection: []string{
		"Oh, too busy\nfor pets?\nTypical.",
		"I'm fine.\nTotally fine.\n*not fine*",
		"Remember me?\nYour pet?",
	},
	NeedsCritical: []string{
		"Wow. This is\nhow I go.\nCool. Cool.",
		"Dying here.\nNo big deal.",
		"Hello? HELP?\nAnyone?!",
	},
}

var poolShy = &MessagePool{
	FeedSuccess: []string{
		"...thank you.\n*small wiggle*",
		"Oh! Food...\nfor me?",
		"*nibbles*\n...nice.",
	},
	FeedFull: []string{
		"Um... I'm okay\nnow... thanks.",
		"*tiny burp*\n...sorry.",
	},
	FeedHungry: []string{
		"*eyes light up*\nOh... finally...",
		"...thank you\nso much...",
	},
	FeedSleeping: []string{
		"Zzz...\n*curled up*",
	},
	FeedTun: []string{
		"*small twitch*\n...",
	},
	FeedCooldown: []string{
		"...um...\nstill eating...",
		"*tiny burp*\n...wait...",
	},
	PlaySuccess: []string{
		"...this is\nnice...",
		"*hesitant\nwiggle*",
	},
	PlayTired: []string{
		"...maybe\nlater?",
	},
	PlayHappy: []string{
		"*actually\nsmiling*\n...more?",
	},
	PlaySleeping: []string{
		"Zzz... *tiny\nsleep sounds*",
	},
	PlayTun: []string{
		"*still*\n...",
	},
	PlayCooldown: []string{
		"...need a\nmoment...",
		"*catching\nbreath*",
	},
	PetSuccess: []string{
		"*blush*\n...oh.",
		"...that's\nnice...",
	},
	PetMaxBond: []string{
		"I... trust you.\nA lot.",
	},
	PetLowMood: []string{
		"...*sniff*\n...thanks.",
	},
	PetSleeping: []string{
		"Zzz...\n*nuzzles*",
	},
	PetTun: []string{
		"*slowly\nuncurling*\n...warmth...",
	},
	PetCooldown: []string{
		"*shy*\n...too much...",
		"...okay...\none moment...",
	},
	Reviving: []string{
		"...I'm okay.\n*tiny wave*",
		"*blinks*\n...hello again.",
	},
	IdleHappy: []string{
		"...today is\ngood.",
	},
	IdleNeutral: []string{
		"...",
		"*watching*",
	},
	IdleHungry: []string{
		"...um...\nhungry...",
	},
	IdleSad: []string{
		"*small sigh*",
	},
	IdleLonely: []string{
		"...anyone\nthere?",
	},
	IdleCritical: []string{
		"...help...",
	},
	IdleTun: []string{
		"*curled*",
	},
	IdleSleeping: []string{
		"Zzz...",
	},
	NeedsFood: []string{
		"...um...\nhungry...",
		"*tummy rumbles*\n...sorry...",
		"...food?\n...maybe?",
	},
	NeedsPlay: []string{
		"...bored...\n...a little...",
		"...play?\n...if you want...",
		"*fidgets*\n...um...",
	},
	NeedsAffection: []string{
		"...lonely...",
		"...miss you...",
		"...pets?\n...please?",
	},
	NeedsCritical: []string{
		"...help...\n...please...",
		"...not okay...\n...scared...",
		"*whimper*\n...need you...",
	},
}
