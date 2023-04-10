package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

func main() {
	DiscordClient.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		botAccount := s.State.User.Username + "#" + s.State.User.Discriminator
		Logger.Info("logged in as " + botAccount)
	})

	if err := DiscordClient.Open(); err != nil {
		Logger.Panic("failed to open discord session", zap.Error(err))
	}

	defer func() {
		for serverID := range ServerConfigMap {
			commands, commandsGetErr := DiscordClient.ApplicationCommands(DiscordClient.State.Application.ID, serverID)
			if commandsGetErr == nil {
				for _, command := range commands {
					if err := DiscordClient.ApplicationCommandDelete(DiscordClient.State.Application.ID, serverID, command.ID); err != nil {
						Logger.Error("failed to delete slash command", zap.Error(err))
						continue
					}
				}
			}

			_ = DiscordClient.Close()
		}
	}()

	initSlashCommands()
	addHandlers()

	stopBot := make(chan os.Signal, 1)
	signal.Notify(stopBot, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	Logger.Info("bot is now running.  press ctrl-c to exit.")
	<-stopBot
}
