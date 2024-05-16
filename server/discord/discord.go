package discord

import "github.com/bwmarrin/discordgo"

type DiscordClient struct {
	*discordgo.Session
}

func NewDiscordClient(token string) (*DiscordClient, error) {
	client, err := discordgo.New(token)
	if err != nil {
		return nil, err
	}

	client.Identify.Intents = makeIntentFlags(
		discordgo.IntentsGuildMembers,
		discordgo.IntentsGuildMessages,
		discordgo.IntentsMessageContent,
	)

	return &DiscordClient{client}, nil
}
