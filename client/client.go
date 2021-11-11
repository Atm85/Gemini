package client

import (
	"errors"
	"github.com/iAtomPlaza/gemini/config"
	"log"

	"github.com/bwmarrin/discordgo"
)

var client *Client

type Client struct {
	token  string
	Prefix string

	User    *discordgo.User
	Session *discordgo.Session
	Config  *config.Config
}

func Get() *Client {
	return client
}

func New(config *config.Config) (*Client, error) {

	discord, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		return nil, err
	}

	user, err := discord.User("@me")
	if err != nil {
		return nil, err
	}

	client = &Client{
		token:   config.Token,
		Prefix:  config.Prefix,
		User:    user,
		Session: discord,
		Config:  config,
	}

	return client, nil
}

func (client *Client) Start() {

	log.Printf("starting...")

	err := client.Session.Open()
	if err != nil {
		log.Printf("%s", err.Error())
		return
	}

	log.Printf("%s is online!", client.User.Username)
	<-make(chan struct{})
}

func (client *Client) Role(guildID, roleID string) (*discordgo.Role, error) {

	roles, _ := client.Session.GuildRoles(guildID)
	for _, role := range roles {
		if role.ID == roleID {
			return role, nil
		}
	}

	return nil, errors.New("roleID not found")
}
