package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var Path    = "_datastore/"
var servers = make(map[string]*Server)

type User struct {
	Level      int    `json:"level"`
	Experience int    `json:"experience"`
}

type Server struct {
	ID      string            `json:"id"`
	Users   map[string]*User  `json:"users"`
	Rewards map[string]string `json:"rewards"`
}

type Config struct {
	Token  string `json:"token"`
	Prefix string `json:"prefix"`
}

func New(filepath string) (*Config, error) {

	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	config := Config{}
	err = json.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func InitDatabase(path string) error {

	files, err := ioutil.ReadDir(path)
	if err != nil {
		err = os.Mkdir(path, 0777)
		if err != nil {
			return err
		}
	}

	for _, file := range files {

		bytes, err := ioutil.ReadFile(path+file.Name())
		if err != nil {
			return err
		}

		server := Server{}
		err = json.Unmarshal(bytes, &server)

		servers[strings.TrimSuffix(file.Name(), ".json")] = &server
	}

	return nil
}

func GetServers() map[string]*Server {
	return servers
}

func GetServer(id string) *Server {

	server, ok := servers[id]
	if !ok {

		// if server does not yet exist, attempt to create it automatically
		return NewServer(id)
	}

	return server
}

func NewServer(id string) *Server {

	server := Server{
		ID:      id,
		Users:   make(map[string]*User),
		Rewards: make(map[string]string),
	}

	bytes, err := json.MarshalIndent(server, "", "    ")
	if err != nil {
		fmt.Println(err)
		return nil
	}

	err = ioutil.WriteFile(Path+id+".json", bytes, 0666)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	servers[id] = &server
	return &server
}

func (s *Server) Save() {

	bytes, err := json.MarshalIndent(s, "", "    ")
	if err != nil {
		fmt.Println(err)
		return
	}

	_ = ioutil.WriteFile(Path+s.ID+".json", bytes, 0666)
}