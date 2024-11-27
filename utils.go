package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"
	"unicode"
)

// createClientWithToken creates an HTTP client that includes an authorization token in each request
func createClientWithToken(token string) *http.Client {
	// Custom transport to inject the token into each request
	transport := &http.Transport{}

	client := &http.Client{
		Timeout:   10 * time.Second,
		Transport: transport,
	}

	// RoundTripper to add the authorization header for each request
	client.Transport = roundTripperWithToken{transport, token}
	return client
}

// roundTripperWithToken is a custom RoundTripper that adds an authorization header
type roundTripperWithToken struct {
	transport http.RoundTripper
	token     string
}

func (rt roundTripperWithToken) RoundTrip(req *http.Request) (*http.Response, error) {
	// Add Authorization header
	req.Header.Add("Authorization", "Bearer "+rt.token)
	return rt.transport.RoundTrip(req)
}

func getArtist(title string, discography *Discography) string {
	title = strings.Replace(title, "(", "", -1)
	title = strings.Replace(title, ")", "", -1)
	title = strings.Replace(title, ".", "", -1)
	artist := title
	// use discography orchestra if it exist
	if discography != nil && discography.Orchestra != "" {
		artist = discography.Orchestra
	}
	words := strings.Fields(artist) // Split title into words
	var result []string

	for _, word := range words {
		if isAlphabetical(word) {
			result = append(result, word)
		}
	}
	return strings.Join(result, " ")
}

// Helper function to check if a word contains only alphabetical characters
func isAlphabetical(word string) bool {
	for _, char := range word {
		if !unicode.IsLetter(char) {
			return false
		}
	}
	return true
}

func callbackUrl() string {
	url := fmt.Sprintf("http://localhost:%d/callback", Config.CallbackPort)
	return url
}

func inList(key string, ilist []string) bool {
	for _, k := range ilist {
		if k == key {
			return true
		}
	}
	return false
}
