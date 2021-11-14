package command

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/iAtomPlaza/dgoc"
	"github.com/iAtomPlaza/gemini/config"
	"strconv"
	"strings"
)

type Level struct {
	Name string
	Desc string
}

func (command *Level) Prepare()  {
	command.Name = "level"
	command.Desc = "Show your or another members level"
}

func (command *Level) Execute(ctx *dgoc.Context, args []string) {

	var target *discordgo.User

	if len(args) <= 0 {
		target = ctx.Message.Author
	} else {
		target, _ = ctx.Session.User(ParseID(args[0]))
	}

	server   := config.GetServer(ctx.Message.GuildID)
	user, ok := server.Users[target.ID]
	if !ok {
		_, _ = ctx.Session.ChannelMessageSendComplex(ctx.Message.ChannelID, &discordgo.MessageSend{
			Embed: &discordgo.MessageEmbed{
				Description: fmt.Sprintf("%s's Statistics...", target.Mention()),
				Fields: []*discordgo.MessageEmbedField{
					{ Name: "Level", Value: "0", Inline: true },
					{ Name: "Experience", Value: "0", Inline: true },
				},
				Color:       0xd7b480,
			},
		})

		return
	}

	_, _ = ctx.Session.ChannelMessageSendComplex(ctx.Message.ChannelID, &discordgo.MessageSend{
		Embed: &discordgo.MessageEmbed{
			Description: fmt.Sprintf("%s's Statistics...", target.Mention()),
			Fields: []*discordgo.MessageEmbedField{
				{ Name: "Level", Value: strconv.Itoa(user.Level), Inline: true },
				{ Name: "Experience", Value: strconv.Itoa(user.Experience), Inline: true },
			},
			Color:       0xd7b480,
		},
	})
}

func ParseID(input string) string {

	if strings.HasPrefix(input, "<@!") && strings.HasSuffix(input, ">") {
		return input[3:len(input)-1]
	}

	return input
}