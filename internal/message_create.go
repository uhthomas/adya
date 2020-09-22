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
			g, err := s.Guild(m.GuildID)
			if err != nil {
				return fmt.Errorf("guild: %w", err)
			}
			for _, vs := range g.VoiceStates {
				if err := s.GuildMemberMute(g.ID, vs.UserID, true); err != nil {
					log.Println(err)
				}
			}
			return nil
		}, true); err != nil {
			log.Printf("mute: %v", err)
			return
		}
	case "!unmute":
		if err := h.messageCreate(s, m, func() error {
			g, err := s.Guild(m.GuildID)
			if err != nil {
				return fmt.Errorf("guild: %w", err)
			}
			for _, vs := range g.VoiceStates {
				if err := s.GuildMemberMute(g.ID, vs.UserID, false); err != nil {
					log.Println(err)
				}
			}
			return nil
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
	defer s.MessageReactionRemove(m.ChannelID, m.ID, "🔄", s.State.User.ID)
	if err := f(); err != nil {
		if err := s.MessageReactionAdd(m.ChannelID, m.ID, "❌"); err != nil {
			return fmt.Errorf("add error emoji: %w", err)
		}
		return err
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
