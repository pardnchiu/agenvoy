package discord

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	discordCommand "github.com/pardnchiu/agenvoy/internal/discord/command"
	discordTypes "github.com/pardnchiu/agenvoy/internal/discord/types"
)

func interactionCreate(dcSession *discordgo.Session, dcInderactionCreate *discordgo.InteractionCreate) {
	if dcInderactionCreate.Type != discordgo.InteractionApplicationCommand {
		return
	}

	data := dcInderactionCreate.ApplicationCommandData()

	var userID, username string
	if dcInderactionCreate.Member != nil {
		userID = dcInderactionCreate.Member.User.ID
		username = dcInderactionCreate.Member.User.Username
	} else if dcInderactionCreate.User != nil {
		userID = dcInderactionCreate.User.ID
		username = dcInderactionCreate.User.Username
	}

	var params []string
	for _, opt := range data.Options {
		params = append(params, opt.StringValue())
	}

	message := &discordTypes.ReceiveMessage{
		GuildID:    dcInderactionCreate.GuildID,
		ChannelID:  dcInderactionCreate.ChannelID,
		AuthorID:   userID,
		AuthorName: username,
		Content:    fmt.Sprintf("/%s %s", data.Name, strings.Join(params, " ")),
		Cmd:        fmt.Sprintf("/%s", data.Name),
		Params:     params,
		IsChannel:  dcInderactionCreate.GuildID != "",
		IsMention:  false,
		RecievedAt: time.Now().Unix(),
	}
	ctx := context.Background()
	dcSession.InteractionRespond(dcInderactionCreate.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})

	slog.Info("command received",
		slog.String("user", message.AuthorName),
		slog.String("content", message.Content),
		slog.Bool("is_channel", message.IsChannel))

	replies := discordCommand.Handler(message)
	for _, reply := range replies {
		dcReply := &discordTypes.DiscordReply{
			Session:     dcSession,
			Interaction: dcInderactionCreate,
		}
		Reply(ctx, dcReply, reply)
	}
}
