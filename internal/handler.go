package internal

import (
	"github.com/bwmarrin/discordgo"
)

type Handler struct{ admins []string }

func NewHandler(admins []string) *Handler { return &Handler{admins: admins} }

func (h *Handler) Handle(s *discordgo.Session) {
	s.AddHandler(h.MessageCreate)
}

func (h *Handler) Admin(id string) bool {
	for _, v := range h.admins {
		if id == v {
			return true
		}
	}
	return false
}
