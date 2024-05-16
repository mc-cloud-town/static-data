package discord

import "github.com/bwmarrin/discordgo"

func makeIntentFlags(intents ...discordgo.Intent) discordgo.Intent {
	var result discordgo.Intent
	for _, intent := range intents {
		result |= intent
	}
	return result
}
