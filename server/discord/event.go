package discord

import "github.com/bwmarrin/discordgo"

func messageCreateEvent(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	s.ChannelMessageSend(m.ChannelID, m.Content)
}

func archiveReleaseHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.ChannelID != "" {
		return
	}

	if reference := m.Reference(); reference != nil {
		referenceMsg, err := s.ChannelMessage(reference.ChannelID, reference.MessageID)
		if err == nil {
			s.ChannelMessageDelete(referenceMsg.ChannelID, referenceMsg.ID)
			return
		}
	}

	// content := m.Content
}
