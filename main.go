package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/iAtomPlaza/dgoc"
	"github.com/iAtomPlaza/gemini/client"
	"github.com/iAtomPlaza/gemini/command"
	"github.com/iAtomPlaza/gemini/config"
	"math"
	"math/rand"
	"strconv"
)

func main() {

	conf, err := config.New("./global.json")
	if err != nil {
		panic(err.Error())
	}

	bot, err := client.New(conf)
	if err != nil {
		panic(err.Error())
	}

	defer bot.Start()

	// add event listeners
	bot.Session.AddHandler(messageEvent)

	//register commands...
	commandHandler := dgoc.New(bot.Session)
	dgoc.SetPrefix(bot.Config.Prefix)
	err = commandHandler.AddCommand(
		&command.Help{},
		&command.Level{})

	if err != nil {
		fmt.Println(err)
	}

	// load database into memory
	err = config.InitDatabase(config.Path)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func messageEvent(s *discordgo.Session, m *discordgo.MessageCreate) {

	// do nothing if the user is a bot
	if m.Author.Bot {
		return
	}

	server   := config.GetServer(m.GuildID)
	user, ok := server.Users[m.Author.ID]

	// create new entry if user is not found in config
	if !ok {
		user = &config.User{
			Level:      1,
			Experience: 0,
		}

		server.Users[m.Author.ID] = user
		go server.Save()
	}

	addXp(user, server, s, m)
}

func addXp(user *config.User, server *config.Server, s *discordgo.Session, m *discordgo.MessageCreate) {

	// add a random integer between 5 and 10 to existing experience
	user.Experience += rand.Intn(10-5) + 5
	go server.Save()

	// calculate next level
	// the formula is: experience**(1/4)
	var exponent, base float64
	base = float64(user.Experience)
	exponent = float64(1) / float64(4)
	calculation := math.Pow(base, exponent)
	if user.Level < int(calculation) {
		addLevel(user, server, s, m)
	}
}

func addLevel(user *config.User, server *config.Server, s *discordgo.Session, m *discordgo.MessageCreate) {

	user.Level += 1
	go server.Save()

	_, _ = s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Description: fmt.Sprintf("gg %s! :tada: You just advanced to level %d", m.Author.Mention(), user.Level),
	})

	// role rewards
	roleID, ok := server.Rewards[strconv.Itoa(user.Level)]
	if !ok {
		return
	}

	role, err := client.Get().Role(m.GuildID, roleID)
	if err != nil {
		return
	}

	_ = s.GuildMemberRoleAdd(m.GuildID, m.Author.ID, role.ID)
}
