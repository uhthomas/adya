package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
)

func Main(context.Context) error {
	token, ok := os.LookupEnv("TOKEN")
	if !ok {
		return errors.New("missing token")
	}

	s, err := discordgo.New(token)
	if err != nil {
		return fmt.Errorf("new session: %w", err)
	}
	defer s.Close()

	if err := s.Open(); err != nil {
		return fmt.Errorf("open: %w", err)
	}

	return nil
}

func main() {
	if err := Main(context.Background()); err != nil {
		log.Fatal(err)
	}
}
