# GoPlaylist

[![Go CI build](https://github.com/vkuznet/goplaylist/actions/workflows/go.yml/badge.svg)](https://github.com/vkuznet/goplaylist/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/vkuznet/goplaylist)](https://goreportcard.com/report/github.com/vkuznet/goplaylist)

**GoPlaylist** is a simple tool to build Spotify or YouTube playlists from provided XML/CSV files.

## Getting Started

To use this tool, you need the [Go programming language](https://go.dev/doc/install)
installed on your system. Alternatively,
you can download a pre-built static executable from the
[Releases](https://github.com/vkuznet/GoPlaylist/releases) section.

### Spotify Setup
To use Spotify, obtain client credentials from their developer site:

1. Create a new app.
2. Set it as a Web Application.
3. Configure the callback URL as `http://localhost:8888/callback`.
Ensure the port (e.g., 8888) is available on your system and matches the `callback_port` value in `config.json`.

### YouTube Setup
For YouTube integration:

1. Visit the [Google Console](https://console.cloud.google.com)
2. Navigate to `API & Services -> OAuth consent screen`.
3. Create a new app, set it as a Web Application, and configure the callback URL as `http://localhost:8888/callback`.
Ensure the port matches the `callback_port` value in `config.json`.


### Configuration and Usage

#### Compilation
To compile the tool, use:
```
go build
```
This will create the playlist executable.

#### Example `config.json`:
```
{
    "client_id": "blablabla",
    "client_secret": "secretbla",
    "callback_port": 8888,
    "verbose": 0,
    "service": "spotify"
}
```
The `service` value can be either **spotify** or **youtube** depending
on service you want to use.

#### Running the Tool
To parse a playlist and print tracks:
```
# parse given file and print out found tracks
./goplaylist -config config.json -file=testplaylist.xml -tracks
# the same operation using csv file
./goplaylist -config config.json -file=testplaylist.csv -tracks
```

To upload a playlist to Spotify or YouTube:
```
# upload my testplaylist to Spotify, i.e. ensure your config.json specifies
# "service": "spotify" option:
./goplaylist -config config.json -file=testplaylist.xml
...
# it will provide you an URL to click on and you'll go through verification
step and your playlist will be build in corresponding service
```
The tool will generate a URL for you to complete the authentication process. Once authorized, your playlist will be created.


#### Example XML playlist
```
<?xml version="1.0" encoding="UTF-8"?>
<discography>
    <track name="A la luz del candil" vocal="Jorge Durán" year="1956-09-27" genre="tango" orchestra="Carlos Di Sarli"/>
    <track name="Sin rumbo fijo" vocal="Ángel Vargas" year="1938-04-18" genre="vals" orchestra="Orquesta Tipica Victor"/>
</discography>
```

#### Example CSV playlist
```
Carlos Di Sarli,1956-09-27,A la luz del candil
Orquesta Tipica Victor,1938-04-18,Sin rumbo fijo
```


### Limitations

#### Youtube limitations
- API Quota: Limited to 10,000 units/day per client. Each search query consumes 100 units.
- Playlist Size: Maximum 5,000 videos.
- Large playlists (>100 tracks) may require multiple runs due to daily quotas. The tool skips existing tracks during reruns.

#### Spotify limitations

- Playlist Size: Maximum 10,000 tracks per playlist.
- Search Issue: Older tracks may lack proper year metadata. The tool excludes
  year: queries by default. To enable this feature, modify `spotify.go`.

### History and motivation
I created this tool as a Tango DJ to upload playlists or discographies of
specific Tango orchestras to Spotify and YouTube. The tool is optimized for
Tango music attributes, including Name, Year, Genre, Artist, Vocal, and
Orchestra. For additional attributes, submit a pull request with proper unit
tests.


### References:

- [Youtube API](https://developers.google.com/youtube/v3/getting-started)
- [Youtube Go](https://developers.google.com/youtube/v3/quickstart/go)
- [Google Console](https://console.cloud.google.com)
- [Spotify Dashboard](https://developer.spotify.com/dashboard/applications)
