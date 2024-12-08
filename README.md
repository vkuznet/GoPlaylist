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

# use sortBy option
./goplaylist -config spotify.json -file=testplaylist.xml -tracks -sortBy=year
{Orchestra:Francisco Canaro Year:1927-02-17 Name:La cumparsita Artist: Genre:Tango Vocal:Instrumental}
{Orchestra:Francisco Canaro Year:1929-04-17 Name:La cumparsita Artist: Genre:Tango Vocal:Instrumental}
{Orchestra:Francisco Canaro Year:1933-02-14 Name:La cumparsita Artist: Genre:Tango Vocal:Instrumental}
{Orchestra:Orquesta Tipica Victor Year:1938-04-18 Name:Sin rumbo fijo Artist: Genre:vals Vocal:Ángel Vargas}
{Orchestra:Angel D'Agostino Year:1945-11-02 Name:La cumparsita Artist: Genre:Tango Vocal:Ángel Vargas}
{Orchestra:Anibal Troilo Year:1951 Name:La cumparsita Artist: Genre:Tango Vocal:Instrumental}
{Orchestra:Anibal Troilo Year:1952 Name:La cumparsita Artist: Genre:Tango Vocal:Instrumental}
{Orchestra:Anibal Troilo Year:1953 Name:Vuelve la serenata Artist: Genre:Vals Vocal:Raúl Beron, Jorge Casal}
{Orchestra:Carlos Di Sarli Year:1956-09-27 Name:A la luz del candil Artist: Genre:tango Vocal:Jorge Durán}
{Orchestra:Anibal Troilo Year:1963-04-25 Name:La cumparsita Artist: Genre:Tango Vocal:Instrumental}

# use multiple keys for sortBy option
./goplaylist -config spotify.json -file=testplaylist.xml -tracks -sortBy=orchestra,year
{Orchestra:Angel D'Agostino Year:1945-11-02 Name:La cumparsita Artist: Genre:Tango Vocal:Ángel Vargas}
{Orchestra:Anibal Troilo Year:1951 Name:La cumparsita Artist: Genre:Tango Vocal:Instrumental}
{Orchestra:Anibal Troilo Year:1952 Name:La cumparsita Artist: Genre:Tango Vocal:Instrumental}
{Orchestra:Anibal Troilo Year:1953 Name:Vuelve la serenata Artist: Genre:Vals Vocal:Raúl Beron, Jorge Casal}
{Orchestra:Anibal Troilo Year:1963-04-25 Name:La cumparsita Artist: Genre:Tango Vocal:Instrumental}
{Orchestra:Carlos Di Sarli Year:1956-09-27 Name:A la luz del candil Artist: Genre:tango Vocal:Jorge Durán}
{Orchestra:Francisco Canaro Year:1927-02-17 Name:La cumparsita Artist: Genre:Tango Vocal:Instrumental}
{Orchestra:Francisco Canaro Year:1929-04-17 Name:La cumparsita Artist: Genre:Tango Vocal:Instrumental}
{Orchestra:Francisco Canaro Year:1933-02-14 Name:La cumparsita Artist: Genre:Tango Vocal:Instrumental}
{Orchestra:Orquesta Tipica Victor Year:1938-04-18 Name:Sin rumbo fijo Artist: Genre:vals Vocal:Ángel Vargas}

# matches only specific orchestra
./goplaylist -config spotify.json -file=testplaylist.xml -tracks -sortBy=year -filterBy='{"orchestra": "anibal troilo"}'
{Orchestra:Anibal Troilo Year:1951 Name:La cumparsita Artist: Genre:Tango Vocal:Instrumental}
{Orchestra:Anibal Troilo Year:1952 Name:La cumparsita Artist: Genre:Tango Vocal:Instrumental}
{Orchestra:Anibal Troilo Year:1953 Name:Vuelve la serenata Artist: Genre:Vals Vocal:Raúl Beron, Jorge Casal}
{Orchestra:Anibal Troilo Year:1963-04-25 Name:La cumparsita Artist: Genre:Tango Vocal:Instrumental}

# matches only specific genre
./goplaylist -config spotify.json -file=testplaylist.xml -tracks -sortBy=year -filterBy='{"genre":"vals"}'
{Orchestra:Orquesta Tipica Victor Year:1938-04-18 Name:Sin rumbo fijo Artist: Genre:vals Vocal:Ángel Vargas}
{Orchestra:Anibal Troilo Year:1953 Name:Vuelve la serenata Artist: Genre:Vals Vocal:Raúl Beron, Jorge Casal}

# use sort by and filter by filters, matches both genre and orchestra
./goplaylist -config spotify.json -file=testplaylist.xml -tracks -sortBy=year -filterBy='{"genre":"Vals", "orchestra": "anibal troilo"}'
{Orchestra:Anibal Troilo Year:1953 Name:Vuelve la serenata Artist: Genre:Vals Vocal:Raúl Beron, Jorge Casal}

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

You may use different options to construct precise playlist, e.g. read all Juan
D'Arienzo discography files, select vals tracks, order them by year and
construct YouTube playlist from them:

```
./goplaylist -config youtube2.json -file="/path/Juan D'Arienzo*.xml" -sortBy=year -filterBy='{"genre":"vals"}' -tracks -title "Juan D'Arienzo - Vals"
```


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
