package internal

import (
	"errors"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

func (h *Handler) MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	switch m.Content {
	case "!mute":
		if err := h.messageCreate(s, m, func() error {
			return members(s, m.GuildID, func(m *discordgo.Member) error {
				return s.GuildMemberMute(m.GuildID, m.User.ID, true)
			})
		}, true); err != nil {
			log.Printf("mute: %v", err)
			return
		}
	case "!unmute":
		if err := h.messageCreate(s, m, func() error {
			return members(s, m.GuildID, func(m *discordgo.Member) error {
				return s.GuildMemberMute(m.GuildID, m.User.ID, false)
			})
		}, true); err != nil {
			log.Printf("unmute: %v", err)
			return
		}
	}
}

func (h *Handler) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate, f func() error, admin bool) error {
	if admin && !h.Admin(m.Author.ID) {
		return errors.New("not an admin")
	}
	if err := s.MessageReactionAdd(m.ChannelID, m.ID, "🔄"); err != nil {
		return fmt.Errorf("add in-progress emoji: %w", err)
	}
	if err := f(); err != nil {
		return err
	}
	if err := s.MessageReactionRemove(m.ChannelID, m.ID, "🔄", s.State.User.ID); err != nil {
		return fmt.Errorf("remove in-progress emoji: %w", err)
	}
	if err := s.MessageReactionAdd(m.ChannelID, m.ID, "✅"); err != nil {
		return fmt.Errorf("add done emoji: %w", err)
	}
	return nil
}

func members(s *discordgo.Session, guild string, f func(*discordgo.Member) error) error {
	var after string
	for {
		ms, err := s.GuildMembers(guild, after, 1000)
		if err != nil {
			return fmt.Errorf("members: %w", err)
		}
		if len(ms) == 0 {
			return nil
		}
		for _, m := range ms {
			if err := f(m); err != nil {
				return err
			}
		}
		after = ms[len(ms)-1].User.ID
	}
}
