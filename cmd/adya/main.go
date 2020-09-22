package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/uhthomas/adya/internal"
)

func sleep(ctx context.Context, d time.Duration) {
	t := time.NewTimer(d)
	defer t.Stop()

	select {
	case <-t.C:
	case <-ctx.Done():
	}
}

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

	for i := time.Duration(0); ; i++ {
		if err := s.Open(); err != nil {
			if i <= 5 {
				sleep(ctx, i*2*time.Second)
				continue
			}
			return fmt.Errorf("open: %w", err)
		}
		defer s.Close()
		break
	}

	<-ctx.Done()
	return nil
}

func main() {
	if err := Main(context.Background()); err != nil {
		log.Fatal(err)
	}
}
