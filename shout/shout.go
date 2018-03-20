package shout

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/mccainca/Icy2Spotify/spotify"
)

type icyparseerror struct {
	s string
}

type CreatePlaylistCommand struct {
	client spotify.Client
}

func NewCreatePlaylistCommand(client spotify.Client) CreatePlaylistCommand {
	return CreatePlaylistCommand{client: client}
}

func (ipe *icyparseerror) Error() string {
	return ipe.s
}

func parseIcy(rdr *bufio.Reader, c byte) (string, error) {
	numbytes := int(c) * 16
	bytes := make([]byte, numbytes)
	n, err := io.ReadFull(rdr, bytes)
	if err != nil {
		log.Panic(err)
	}
	if n != numbytes {
		return "", &icyparseerror{"didn't get enough data"} // may be invalid
	}
	return strings.Split(strings.Split(string(bytes), "=")[1], ";")[0], nil
}

func extractMetadata(rdr io.Reader, skip int) <-chan string {
	ch := make(chan string)
	go func() {
		bufrdr := bufio.NewReaderSize(rdr, skip)
		for {
			skipbytes := make([]byte, skip)

			_, err := io.ReadFull(bufrdr, skipbytes)
			if err != nil {
				log.Printf("Failed: %v\n", err)
				close(ch)
				break
			}
			c, err := bufrdr.ReadByte()
			if err != nil {
				log.Panic(err)
			}
			if c > 0 {
				meta, err := parseIcy(bufrdr, c)
				if err != nil {
					log.Panic(err)
				}
				ch <- meta
			}
		}
	}()
	return ch
}

func (c CreatePlaylistCommand) Run(url string, playlist string) {
	client := &http.Client{}

	log.Printf("Getting stream from : %s\n", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Error: %s", err)
		return
	}

	req.Header.Add("Icy-MetaData", "1")
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error: %s", err)
		return
	}

	amount := 0
	if _, err = fmt.Sscan(resp.Header.Get("Icy-Metaint"), &amount); err != nil {
		fmt.Printf("Error: %s", err)
		return
	}

	metaChan := extractMetadata(resp.Body, amount)
	for meta := range metaChan {
		artist, track := splitTrackInfo(meta)

		spTrack, err := c.client.FindTrack(fmt.Sprintf("%s %s", artist, track))
		if err != nil {
			log.Printf("Unable to add track: %s by %s to playlist : %s", track, artist, err)
		} else {
			c.client.AddTrackToPlaylist(playlist, spTrack)
			log.Printf("Added %s by %s to playlist", track, artist)
		}
	}
}

func splitTrackInfo(metaData string) (string, string) {

	s := strings.Replace(metaData, "'", "", -1)
	t := strings.Split(s, " - ")
	return t[0], t[1]
}
