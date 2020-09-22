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

	member, err := s.State.Member(m.GuildID, m.Author.ID)
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
		return
	}

	st, err := s.GuildRoles(m.GuildID)
	if err != nil {
		log.Println(err)
		return
	}

	roles := discordgo.Roles(st)

	sort.Sort(roles)

	mask := 0
	if m.Content == "!um" {
		mask = discordgo.PermissionVoiceSpeak
	}

	for _, r := range roles[1:] {
		if _, err := s.GuildRoleEdit(
			m.GuildID,
			r.ID,
			r.Name,
			r.Color,
			r.Hoist,
			r.Permissions|mask,
			r.Mentionable,
		); err != nil {
			log.Println(err)
		}
	}
}
