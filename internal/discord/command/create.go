package discordCommand

import (
	"log/slog"
	"os"

	"github.com/bwmarrin/discordgo"
	discordTypes "github.com/pardnchiu/agenvoy/internal/discord/types"
)

func Create(dcBot *discordTypes.DiscordBot, dcSession *discordgo.Session) {
	var command []*discordgo.ApplicationCommand
	for _, cmd := range commands {
		switch cmd {
		case CmdHelp:
			command = append(command, &discordgo.ApplicationCommand{
				Name:        cmd.Text(),
				Description: "Show how to use",
			})
		case CmdRole:
			command = append(command, &discordgo.ApplicationCommand{
				Name:        cmd.Text(),
				Description: "Assign role session to handle",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "name",
						Description: "Role name",
						Required:    true,
					},
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "message",
						Description: "Message",
						Required:    true,
					},
				},
			})
		}
	}

	guildID := os.Getenv("DISCORD_GUILD_ID")
	for _, cmd := range command {
		command, err := dcSession.ApplicationCommandCreate(dcSession.State.User.ID, guildID, cmd)
		if err != nil {
			slog.Warn("failed to register command",
				slog.String("error", err.Error()))
			continue
		}
		dcBot.Commands = append(dcBot.Commands, command)
	}
}
