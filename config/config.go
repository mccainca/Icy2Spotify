package config

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	StreamURL     string `json:"streamUrl"`
	PlaylistID    string `json:"playlistId"`
	SpotifyID     string `json:"spotifyId"`
	SpotifySecret string `json:"spotifySecret"`
}

func LoadConfiguration(file string) Config {
	var config Config
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		log.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config
}
