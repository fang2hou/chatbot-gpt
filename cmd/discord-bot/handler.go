package main

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

var (
	// messageHandlers is a list of message handlers.
	messageHandlers []func(*discordgo.Session, *discordgo.MessageCreate) bool

	// interactionHandlers is a map of interaction handlers.
	interactionHandlers = make(
		map[string]map[string]func(*discordgo.Session, *discordgo.InteractionCreate),
	)

	// slashCommands is a list of slash commands.
	slashCommands = struct {
		ClearContext func(alias string) *discordgo.ApplicationCommand
	}{
		ClearContext: func(alias string) *discordgo.ApplicationCommand {
			return &discordgo.ApplicationCommand{
				Name:        alias,
				Description: "Clear the context for GPT Model",
				Type:        discordgo.ChatApplicationCommand,
			}
		},
	}
)

// addHandlers adds the handlers to the Discord client.
func addHandlers() {
	// Messages
	DiscordClient.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		for _, handler := range messageHandlers {
			if ok := handler(s, m); ok {
				return
			}
		}
	})

	// Slash commands
	DiscordClient.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if _, ok := interactionHandlers[i.GuildID]; !ok {
			return
		}

		if handler, ok := interactionHandlers[i.GuildID][i.ApplicationCommandData().Name]; ok {
			handler(s, i)
		}
	})
}

// initSlashCommands initializes the slash commands.
func initSlashCommands() {
	for serverID, serverConfig := range ServerConfigMap {
		commands, commandsGetErr := DiscordClient.ApplicationCommands(
			DiscordClient.State.Application.ID,
			serverID,
		)
		if commandsGetErr != nil {
			Logger.Error("failed to get slash commands", zap.Error(commandsGetErr))
			return
		}

		for _, command := range commands {
			if err := DiscordClient.ApplicationCommandDelete(DiscordClient.State.Application.ID, serverID, command.ID); err != nil {
				Logger.Error("failed to delete slash command", zap.Error(err))
				continue
			}
		}

		if serverConfig.Commands.ClearContext.Enable {
			interactionHandlers[serverID] = make(
				map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate),
			)

			for _, alias := range serverConfig.Commands.ClearContext.Aliases {
				interactionHandlers[serverID][alias] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
					Logger.Debug(
						"received interaction",
						zap.String("command", alias),
						zap.String("user", i.Member.User.Username),
					)

					MessageDatabase.Clear(i.Member.User.ID)

					if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Embeds: []*discordgo.MessageEmbed{
								{
									Title:       "✅ " + Localizer.Fetch("cleared", serverConfig.Language),
									Description: Localizer.Fetch("clear_context", serverConfig.Language),
									Timestamp:   time.Now().Format(time.RFC3339),
									Color:       0x379C6F,
								},
							},
						},
					}); err != nil {
						Logger.Error("failed to respond to interaction", zap.Error(err))
						return
					}
				}

				if _, err := DiscordClient.ApplicationCommandCreate(DiscordClient.State.User.ID, serverID, slashCommands.ClearContext(alias)); err != nil {
					Logger.Error("failed to create slash command", zap.Error(err))
					continue
				}
			}
		}
	}
}
