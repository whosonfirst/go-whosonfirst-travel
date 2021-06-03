# go-whosonfirst-travel

Go package for traveling Who's On First documents and their relations

## Documentation

[![Go Reference](https://pkg.go.dev/badge/github.com/whosonfirst/go-whosonfirst-travel.svg)](https://pkg.go.dev/github.com/whosonfirst/go-whosonfirst-travel)

Documentation is incomplete.

## Tools

```
$> make cli
go build -mod vendor -o bin/wof-travel-id cmd/wof-travel-id/main.go
go build -mod vendor -o bin/wof-belongs-to cmd/wof-belongs-to/main.go
```

### wof-travel-id

```
$> ./bin/wof-travel-id -h
Usage of ./bin/wof-travel-id:
  -csv
    	Emit results formatted as a comma-separated values.
  -hierarchies
    	Travel the hierarchies for each ID specified.
  -ids
    	Emit results as a line-separated list of IDs (traveled).
  -markdown
    	Emit results formatted as Markdown.
  -parent
    	Travel the parent ID of each ID specified.
  -singleton
    	... (default true)
  -source value
    	One or more valid whosonfirst/go-reader URIs to use to read WOF ID data, which may or may not be part of the sources to graph. This is work in progress.
  -stdin
    	Read IDs to travel from STDIN
  -superseded-by
    	Travel records that supersede each ID specified.
  -supersedes
    	Travel records that are superseded by each ID specified.
  -timings
    	Display timing information
```

For example:

```
$> ./bin/wof-travel-id \
	-source fs:///usr/local/data/whosonfirst-data-venue-us-ca/data \
	-supersedes 890535433
	
2018/08/07 17:26:46 890535433 Rock Bar (2012 - uuuu)
2018/08/07 17:26:46 1108808495 International Club (2001 - 2011)
2018/08/07 17:26:46 1108808497 Esperanza's (1996 - 2000)
2018/08/07 17:26:46 1108808499 Nello's Place (1981 - 1991)
2018/08/07 17:26:46 1108808501 Tiffany's Lodge (1975 - 1980)
2018/08/07 17:26:46 1108808503 Joe & Grill (1968 - 1974)
2018/08/07 17:26:46 1108808507 Walt's Place (1955 - 1965)
2018/08/07 17:26:46 1108808509 Geo & Walt's Place (1953 - 1954)
2018/08/07 17:26:46 1108808511 Welsh's Place (1951 - 1953)
2018/08/07 17:26:46 1108808513 Tiffany Club (1942 - 1950)
2018/08/07 17:26:46 1108808515 Cornelius Sweeney (1937 - uuuu)
2018/08/07 17:26:46 1108808517 Tiffany Inn (1934 - 1936)
2018/08/07 17:26:46 1108808519 Martin Burke (proprietor) (1931 - uuuu)
2018/08/07 17:26:46 1108808521 M. J. Reynolds (proprietor) (1922 - uuuu)
2018/08/07 17:26:46 1108808525 Timothy J. Costello (proprietor) (1915 - 1917)
2018/08/07 17:26:46 1108808527 Michael Coody (proprietor) (1911 - 1915)
2018/08/07 17:26:46 1108808529 Mary Caulfield (proprietor) (1910 - 1910)
2018/08/07 17:26:46 1108808531 Frank P. Caulfield (proprietor) (1902 - 1909)
```

Or:

```
$> ./bin/wof-travel-id \
	-source fs:///usr/local/data/whosonfirst-data-venue-us-ca/data \
	-superseded-by 1108808531
	
2018/08/07 17:30:55 1108808531 Frank P. Caulfield (proprietor) (1902 - 1909)
2018/08/07 17:30:55 1108808529 Mary Caulfield (proprietor) (1910 - 1910)
2018/08/07 17:30:55 1108808527 Michael Coody (proprietor) (1911 - 1915)
2018/08/07 17:30:55 1108808525 Timothy J. Costello (proprietor) (1915 - 1917)
2018/08/07 17:30:55 1108808521 M. J. Reynolds (proprietor) (1922 - uuuu)
2018/08/07 17:30:55 1108808519 Martin Burke (proprietor) (1931 - uuuu)
2018/08/07 17:30:55 1108808517 Tiffany Inn (1934 - 1936)
2018/08/07 17:30:55 1108808515 Cornelius Sweeney (1937 - uuuu)
2018/08/07 17:30:55 1108808513 Tiffany Club (1942 - 1950)
2018/08/07 17:30:55 1108808511 Welsh's Place (1951 - 1953)
2018/08/07 17:30:55 1108808509 Geo & Walt's Place (1953 - 1954)
2018/08/07 17:30:55 1108808507 Walt's Place (1955 - 1965)
2018/08/07 17:30:55 1108808503 Joe & Grill (1968 - 1974)
2018/08/07 17:30:55 1108808501 Tiffany's Lodge (1975 - 1980)
2018/08/07 17:30:55 1108808499 Nello's Place (1981 - 1991)
2018/08/07 17:30:55 1108808497 Esperanza's (1996 - 2000)
2018/08/07 17:30:55 1108808495 International Club (2001 - 2011)
2018/08/07 17:30:55 890535433 Rock Bar (2012 - uuuu)
```

Or:

```
$> ./bin/wof-travel-id \
	-source fs:///usr/local/data/whosonfirst-data-venue-us-ca/data \
	-source fs:///usr/local/data/whosonfirst-data/data \
	-supersedes \
	-parent 890535433
	
2018/08/07 17:28:00 890535433 Rock Bar (2012 - uuuu)
2018/08/07 17:28:00 1108808495 International Club (2001 - 2011)
2018/08/07 17:28:00 1108808497 Esperanza's (1996 - 2000)
2018/08/07 17:28:00 102112179 La Lengua (2008-11-19 - uuuu)
2018/08/07 17:28:00 1108808499 Nello's Place (1981 - 1991)
2018/08/07 17:28:00 1108808501 Tiffany's Lodge (1975 - 1980)
2018/08/07 17:28:00 1108808503 Joe & Grill (1968 - 1974)
2018/08/07 17:28:00 1108808507 Walt's Place (1955 - 1965)
2018/08/07 17:28:00 1108808509 Geo & Walt's Place (1953 - 1954)
2018/08/07 17:28:00 1108808511 Welsh's Place (1951 - 1953)
2018/08/07 17:28:00 1108808513 Tiffany Club (1942 - 1950)
2018/08/07 17:28:00 1108808515 Cornelius Sweeney (1937 - uuuu)
2018/08/07 17:28:00 1108808517 Tiffany Inn (1934 - 1936)
2018/08/07 17:28:00 1108808519 Martin Burke (proprietor) (1931 - uuuu)
2018/08/07 17:28:00 1108808521 M. J. Reynolds (proprietor) (1922 - uuuu)
2018/08/07 17:28:00 1108808525 Timothy J. Costello (proprietor) (1915 - 1917)
2018/08/07 17:28:00 1108808527 Michael Coody (proprietor) (1911 - 1915)
2018/08/07 17:28:00 1108808529 Mary Caulfield (proprietor) (1910 - 1910)
2018/08/07 17:28:00 1108808531 Frank P. Caulfield (proprietor) (1902 - 1909)
2018/08/07 17:28:00 85887469 Transmission (uuuu - uuuu)
2018/08/07 17:28:00 85922583 San Francisco (1850-04-15 - uuuu)
2018/08/07 17:28:00 85922583 San Francisco (1850-04-15 - uuuu)
2018/08/07 17:28:00 1125285611 City of San Francisco (uuuu - uuuu)
2018/08/07 17:28:00 102087579 San Francisco County
2018/08/07 17:28:00 102087579 San Francisco County
2018/08/07 17:28:00 85688637 California (uuuu - uuuu)
2018/08/07 17:28:01 85633793 United States (uuuu - uuuu)
2018/08/07 17:28:01 102191575 North America (uuuu - uuuu)
```

_Note: There is still a condition where the same record can be processed twice (unless you explicitly enable this with the `-singleton=true` flag)_

## See also

* https://github.com/whosonfirst/go-whosonfirst-travel-image