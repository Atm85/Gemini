package command

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/iAtomPlaza/dgoc"
	"github.com/iAtomPlaza/gemini/client"
)


type Help struct{
	Name string
	Desc string
}

func (command *Help) Prepare() {
	command.Name = "Help"
	command.Desc = "Lists bot commands"
}

func (command *Help) Execute(ctx *dgoc.Context, args []string) {

	commands := dgoc.CommandMap
	c := client.Get()

	embedDescription := ""
	for _, cmd := range commands {
		embedDescription += fmt.Sprintf("`%s`: %s \n", c.Prefix+cmd.Name, cmd.Desc)
	}

	_, _ = ctx.Session.ChannelMessageSendComplex(ctx.Message.ChannelID, &discordgo.MessageSend{
		Embed: &discordgo.MessageEmbed{
			Description: embedDescription,
			Color:       0xd7b480,
		},
	})
}
