package workflow

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
}
