# GoPlaylist

[![Go CI build](https://github.com/vkuznet/goplaylist/actions/workflows/go.yml/badge.svg)](https://github.com/vkuznet/goplaylist/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/vkuznet/goplaylist)](https://goreportcard.com/report/github.com/vkuznet/goplaylist)


A simple tool to build either Spotify or Youtube playlists from provided
XML/CSV files.

To build the tool you need [Go](https://go.dev/doc/install) language
to be installed on your system or use pre-build static executable from
[releases](https://github.com/vkuznet/GoPlaylist/releases) area.
To build the executable just run
```
go build
```


### Spotify setup
In order to use spotify you must obtain client's credentials from
their [developer site](https://developer.spotify.com/dashboard/applications)
Please create a new app, use Web Application and setup callback function
to `http://localhost:8888/callback`. Here the port number (8888) should
be available on your computer and you should use it in your `config.json`,
see configuration section.

### Youtube setup
For youtube integration please visit
[Google Console](https://console.cloud.google.com) and then follow
`Api & Services -> OAuth consent screen`. Over there you may setup
your new app, check that you'll use web application, setup your
callback URL to `http://localhost:8888/callback` and match your port
number (8888) with your `config.json` configuration.

### Configuration and run procedure
To compile the tool use `go build` and you'll get back `playlist` executable.

Here is example of your `config.json`:
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

To run the tool obtain your favorite playlist in XML/CSV data-format and run it
as following:

```
# parse given file and print out found tracks
./goplaylist -config config.json -file=testplaylist.xml -tracks
# the same operation using csv file
./goplaylist -config config.json -file=testplaylist.csv -tracks

# upload my testplaylist to Spotify, i.e. ensure your config.json specifies
# "service": "spotify" option:
./goplaylist -config config.json -file=testplaylist.xml
...
# it will provide you an URL to click on and you'll go through verification
step and your playlist will be build in corresponding service
```

Here is an example of `testplaylist.xml` file:
```
<?xml version="1.0" encoding="UTF-8"?>
<discography>
    <track name="A la luz del candil" vocal="Jorge Durán" year="1956-09-27" genre="tango" orchestra="Carlos Di Sarli"/>
    <track name="Sin rumbo fijo" vocal="Ángel Vargas" year="1938-04-18" genre="vals" orchestra="Orquesta Tipica Victor"/>
</discography>
```
And, similar example of using CSV data-format
```
Carlos Di Sarli,1956-09-27,A la luz del candil
Orquesta Tipica Victor,1938-04-18,Sin rumbo fijo
```


### Limitations
You should be aware of limitations coming from Spotify/Youtube providers. Here
we briefly summarize them:

#### Youtube limitations
Youtube APIs are limited to 10,000 units per day per client where search query
accounts for 100 units. Therefore, you are limited to 100 search tracks per
single day. And, your playlist may not exceed 5,000 videos. If you will
provide xml/csv file with more than 100 tracks it means that your playlist
will not be built completely (if you are using free Google plan) and you
need to run it again next day. The `goplaylist` provides ability to
re-run existing workflow, i.e. upload your playlist again and again,
and it will skip existing tracks and add only new ones to existing
playlist.

#### Spotify limitations
Spotify limits user to have not more than 10,000 tracks per playlist.
During development process we also found that spotify does not properly
search for track year. It allows to specify query like:
```
track:Bla artist:Firt Last name year:1999
```
But for older tracks it does not keep track of it and therefore
the code skips the `year:XXXX` part of the query. You may easily
enable it and recompile the code if you want to, see `spotify.go`
file and find `year:` in it.

### History and motivation
I built this tool as Tango DJ to allow upload to Youtube or Spotify either full
discography of tracks of a specific Tango Orquestra or concrete milonga
playlist. Therefore, the data structure (see `data.go`) as well provided
xml/csv test playlists are oriented to Tango tracks. For instance, I defined
`Track` to use only Name, Year, Genre, Artist, Vocal, and Orchestra attributes
which may or may not fit modern music tracks. If you'll need more attributes
make sure you'll provide proper PR with unit tests.

### References:

- [Youtube API](https://developers.google.com/youtube/v3/getting-started)
- [Youtube Go](https://developers.google.com/youtube/v3/quickstart/go)
- [Google Console](https://console.cloud.google.com)
- [Spotify Dashboard](https://developer.spotify.com/dashboard/applications)
