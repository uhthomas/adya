package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/uhthomas/adya/internal"
)

func egress(ctx context.Context, c *http.Client) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://google.com/", nil)
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}
	res, err := c.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer res.Body.Close()
	return nil
}

func Main(ctx context.Context) error {
	c := &http.Client{Timeout: 5 * time.Second}
	for i := 0; i < 5; i++ {
		if err := egress(ctx, c); err != nil {
			log.Println(err)
		} else {
			log.Println("ok!")
		}
	}

	token, ok := os.LookupEnv("TOKEN")
	if !ok {
		return errors.New("missing token")
	}

	s, err := discordgo.New(token)
	if err != nil {
		return fmt.Errorf("new session: %w", err)
	}

	s.AddHandler(internal.Handle)

	s.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)

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
