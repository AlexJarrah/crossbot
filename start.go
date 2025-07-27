package crossbot

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func (c *Config) Start(cmds *[]*Command) error {
	if err := c.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	go func() {
		if err := c.Telegram(cmds); err != nil {
			log.Println("Telegram initialization error:", err)
		}
	}()

	cancelDg, err := c.Discord(cmds)
	if err != nil {
		return fmt.Errorf("failed to initialize Discord: %w", err)
	}
	defer cancelDg()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt, syscall.SIGTERM)
	<-sc

	return nil
}
