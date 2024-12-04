# iTunesXML2CSV
This area contains code which transforms iTunes (macOS Music app) playlist from
XML format to CSV.

To compile code just use `go build` and then you can run it as following:
```
# build code
go build

# save your favorite playlist from iTunes in XML format to some destination

# run code
./itunesXML -xmlInput /path/playlist.xml -csvOutput playlist.csv
```
