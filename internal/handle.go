package internal

import (
	"log"
	"sort"

	"github.com/bwmarrin/discordgo"
)

func Handle(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	switch m.Content {
	case "!mu", "!um":
	default:
		return
	}

	member, err := s.GuildMember(m.GuildID, m.Author.ID)
	if err != nil {
		log.Println(err)
		return
	}

	if ok := func() bool {
		for _, r := range member.Roles {
			if r == "757735651419488347" {
				return true
			}
		}
		return false
	}(); !ok {
		log.Println("did not have adya role")
		return
	}

	switch m.Content {
	case "!mu":
		if _, err := s.ChannelMessageSend(m.ChannelID, "muting all roles"); err != nil {
			log.Println(err)
			return
		}
	case "!um":
		if _, err := s.ChannelMessageSend(m.ChannelID, "unmuting all roles"); err != nil {
			log.Println(err)
			return
		}
	}

	st, err := s.GuildRoles(m.GuildID)
	if err != nil {
		log.Println(err)
		return
	}

	roles := discordgo.Roles(st)

	sort.Sort(roles)

	for _, r := range roles[1:] {
		p := r.Permissions
		switch m.Content {
		case "!um":
			p |= discordgo.PermissionVoiceSpeak
		case "!mu":
			p &= ^discordgo.PermissionVoiceSpeak
		}
		if _, err := s.GuildRoleEdit(
			m.GuildID,
			r.ID,
			r.Name,
			r.Color,
			r.Hoist,
			p,
			r.Mentionable,
		); err != nil {
			log.Println(err)
		}
	}
}
