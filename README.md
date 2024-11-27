# GoPlaylist

[![Go CI build](https://github.com/dmwm/goplaylist/actions/workflows/go-ci.yml/badge.svg)](https://github.com/dmwm/goplaylist/actions/workflows/go-ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/dmwm/goplaylist)](https://goreportcard.com/report/github.com/dmwm/goplaylist)


a simple tool to build either Spotify or Youtube playlists from provided
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
