package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	var playlistFile string
	flag.StringVar(&playlistFile, "playlistFile", "", "input playlist file")
	var dst string
	flag.StringVar(&dst, "dst", "", "destination folder")
	flag.Parse()

	if playlistFile == "" || dst == "" {
		fmt.Println("Usage: ./copytracks -playlistFile <file.m3u> -dst <destination_folder>")
		return
	}

	// ensure destination exists
	if _, err := os.Stat(dst); os.IsNotExist(err) {
		if err := os.MkdirAll(dst, 0755); err != nil {
			fmt.Printf("Failed to create destination dir: %v\n", err)
			return
		}
	}

	content, err := os.ReadFile(playlistFile)
	if err != nil {
		fmt.Printf("Failed to read playlist file: %v\n", err)
		return
	}

	lines := strings.Split(string(content), "\r")
	var artist, track string

	re := regexp.MustCompile(`#EXTINF:\d+,(.*) - (.*)`)
	idx := 1 // initial index for tracks

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || line == "#EXTM3U" {
			continue
		}

		if strings.HasPrefix(line, "#EXTINF") {
			parts := strings.SplitN(line, ",", 2)
			if len(parts) == 2 {
				trackInfo := parts[1]
				matches := re.FindStringSubmatch("#EXTINF:0," + trackInfo)
				if len(matches) == 3 {
					track = sanitizeFileName(matches[1])
					artist = sanitizeFileName(matches[2])
				}
			}
		} else {
			srcPath := line
			ext := filepath.Ext(srcPath)
			newFileName := fmt.Sprintf("%03d - %s - %s%s", idx, artist, track, ext)
			destPath := filepath.Join(dst, newFileName)

			err := copyFile(srcPath, destPath)
			if err != nil {
				fmt.Printf("Failed to copy %s → %s: %v\n", srcPath, destPath, err)
			} else {
				fmt.Printf("Copied: %s → %s\n", srcPath, destPath)
				idx += 1
			}
		}
	}
}

func copyFile(src, dest string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	return out.Sync()
}

func sanitizeFileName(name string) string {
	return strings.Map(func(r rune) rune {
		if strings.ContainsRune(`\/:*?"<>|`, r) {
			return '_'
		}
		return r
	}, name)
}
