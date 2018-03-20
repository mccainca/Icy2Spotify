package main

import (
	"log"

	"github.com/mccainca/Icy2Spotify/config"
	"github.com/mccainca/Icy2Spotify/shout"
	"github.com/mccainca/Icy2Spotify/spotify"
)

var appConfig = config.Config{}

const configFile = "config.json"

func init() {
	appConfig = config.LoadConfiguration(configFile)
}
func main() {
	client := authenticate()

	userID, err := client.CurrentUserId()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Logged in as %s", userID)

	createPlaylistCommand := shout.NewCreatePlaylistCommand(client)

	for {
		createPlaylistCommand.Run(appConfig.StreamURL, appConfig.PlaylistID)
	}
}

func authenticate() spotify.Client {

	authURL, clientChannel, err := spotify.Authenticate(appConfig.SpotifyID, appConfig.SpotifySecret)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Please log in to Spotify by visiting the following page in your browser: %s", authURL)
	client := <-clientChannel

	return client
}
