package discord

import (
	"context"
	"os"
	"path/filepath"

	"github.com/bwmarrin/discordgo"
	discordTypes "github.com/pardnchiu/agenvoy/internal/discord/types"
)

const (
	// * if content over 2000, split into multiple messages
	replayMax = 2000
)

func Reply(ctx context.Context, dcReply *discordTypes.DiscordReply, reply discordTypes.ReplyMessage) error {
	var embeds []*discordgo.MessageEmbed

	if reply.ImageURL != "" {
		embeds = []*discordgo.MessageEmbed{
			{
				Image: &discordgo.MessageEmbedImage{
					URL: reply.ImageURL,
				},
			},
		}
	}

	var files []*discordgo.File
	for _, path := range reply.FilePaths {
		f, err := os.Open(path)
		if err != nil {
			continue
		}
		defer f.Close()
		files = append(files, &discordgo.File{
			Name:   filepath.Base(path),
			Reader: f,
		})
	}

	if dcReply.Interaction != nil {
		_, err := dcReply.Session.FollowupMessageCreate(dcReply.Interaction.Interaction, true, &discordgo.WebhookParams{
			Content: reply.Content,
			Embeds:  embeds,
			Files:   files,
		})
		return err
	}

	chunks := split(reply.Content)
	for i, chunk := range chunks {
		var ref *discordgo.MessageReference
		if i == 0 {
			ref = dcReply.Reference
		}
		var chunkEmbeds []*discordgo.MessageEmbed
		var chunkFiles []*discordgo.File
		if i == len(chunks)-1 {
			chunkEmbeds = embeds
			chunkFiles = files
		}
		_, err := dcReply.Session.ChannelMessageSendComplex(dcReply.ChannelID, &discordgo.MessageSend{
			Content:   chunk,
			Reference: ref,
			Embeds:    chunkEmbeds,
			Files:     chunkFiles,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func split(s string) []string {
	if len(s) <= replayMax {
		return []string{s}
	}
	var chunks []string
	for len(s) > replayMax {
		cut := replayMax
		if idx := isLast(s[:cut]); idx > 0 {
			cut = idx + 1
		}
		chunks = append(chunks, s[:cut])
		s = s[cut:]
	}
	if s != "" {
		chunks = append(chunks, s)
	}
	return chunks
}

func isLast(s string) int {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == '\n' {
			return i
		}
	}
	return -1
}
