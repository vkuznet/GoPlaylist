package main

import (
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/dhowden/tag"
)

//go:embed index.html
var indexHTML embed.FS
var musicDir string

type Track struct {
	Index    string
	Filename string
	Title    string
	Artist   string
	Album    string
	Genre    string
	Year     string
	Duration string
}

func main() {
	var port string
	flag.StringVar(&musicDir, "musicDir", "./music", "Directory to load music files from")
	flag.StringVar(&port, "port", "8080", "DJ server port number")
	flag.Parse()

	http.HandleFunc("/", serveIndex)
	http.HandleFunc("/playlist", playlistHandler)
	http.Handle("/music/", http.StripPrefix("/music/", http.FileServer(http.Dir(musicDir))))

	fmt.Println("Server running on :8080")
	http.ListenAndServe(":"+port, nil)
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	data, _ := indexHTML.ReadFile("index.html")
	w.Write(data)
}

func playlistHandler(w http.ResponseWriter, r *http.Request) {
	files, err := os.ReadDir(musicDir)
	if err != nil {
		http.Error(w, "Failed to read music directory", http.StatusInternalServerError)
		return
	}

	var tracks []Track
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		ext := strings.ToLower(filepath.Ext(f.Name()))
		if ext != ".mp3" && ext != ".mp4" && ext != ".aif" {
			continue
		}

		path := filepath.Join(musicDir, f.Name())
		track := Track{Filename: f.Name()}

		// extract index from file Name
		var filenameRegex = regexp.MustCompile(`^(\d{3})\s*-\s*`)
		matches := filenameRegex.FindStringSubmatch(f.Name())
		if len(matches) > 1 {
			track.Index = matches[1]
		}

		file, err := os.Open(path)
		if err == nil {
			metadata, err := tag.ReadFrom(file)
			if err == nil {
				track.Title = metadata.Title()
				track.Artist = metadata.Artist()
				track.Album = metadata.Album()
				track.Genre = metadata.Genre()
				if metadata.Year() > 0 {
					track.Year = fmt.Sprintf("%d", metadata.Year())
				}
				track.Duration = ""
			}
			file.Close()
		}
		tracks = append(tracks, track)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tracks)
}
