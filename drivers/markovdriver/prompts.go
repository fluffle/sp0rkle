package markovdriver

import "math/rand"

var insults = []string{
	"insult me", "roast me", "burn me", "give me your best shot",
	"tell me why i'm terrible", "hurt my feelings", "give me a burn",
	"give me a roast", "destroy me", "roast my existence",
	"hit me with your best insult", "tear me apart",
	"mock my intelligence", "call me stupid", "roast my logic",
	"tell me i'm a moron", "critique my brain", "call me dim",
	"mock my wit", "tell me i'm slow", "roast my iq",
	"mock my thoughts", "tell me i'm an idiot", "roast my brain cells",
	"mock my outfit", "insult my style", "roast my look",
	"mock my haircut", "critique my fashion", "roast my clothes",
	"mock my grooming", "critique my appearance", "tell me i look bad",
	"roast my wardrobe", "mock my hair", "insult my look",
	"mock my charisma", "call me awkward", "roast my social skills",
	"mock my vibes", "tell me why i'm uncool", "roast my personality",
	"call me weird", "mock my charm", "roast my social life",
	"call me a loser", "mock my presence", "tell me i'm awkward",
	"mock my skill", "tell me i'm bad at this", "roast my performance",
	"mock my gaming", "critique my work", "roast my talent",
	"tell me i'm incompetent", "mock my effort", "roast my results",
	"call me a failure", "tell me i suck", "critique my play",
	"mock my voice", "insult my singing", "roast my accent",
	"tell me i sound bad", "critique my voice", "roast my speaking",
	"mock my laugh", "tell me i'm noisy", "roast my tone",
	"mock my pronunciation", "insult my singing voice", "mock my audio",
	"mock my aura", "insult my soul", "roast my energy",
	"critique my existence", "roast my spirit", "mock my essence",
	"critique my vibe", "roast my cosmic presence",
	"mock my life's purpose", "roast my destiny",
	"insult my being", "critique my soul",
}

func randomPrompt() string {
	return insults[rand.Intn(len(insults))]
}
