package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/uhthomas/adya/internal"
)

func Main(ctx context.Context) error {
	token, ok := os.LookupEnv("TOKEN")
	if !ok {
		return errors.New("missing token")
	}

	s, err := discordgo.New(token)
	if err != nil {
		return fmt.Errorf("new session: %w", err)
	}

	s.AddHandler(internal.Handle)

	s.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)

	if err := s.Open(); err != nil {
		return fmt.Errorf("open: %w", err)
	}
	defer s.Close()

	<-ctx.Done()
	return nil
}

func main() {
	if err := Main(context.Background()); err != nil {
		log.Fatal(err)
	}
}
