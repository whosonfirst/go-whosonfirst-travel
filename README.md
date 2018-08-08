# go-whosonfirst-travel

Go package for traveling Who's On First documents and their relations

## Install

You will need to have both `Go` (specifically a version of Go more recent than 1.7 so let's just assume you need [Go 1.9](https://golang.org/dl/) or higher) and the `make` programs installed on your computer. Assuming you do just type:

```
make bin
```

All of this package's dependencies are bundled with the code in the `vendor` directory.

## Important

This should still be considered experimental.

## Tools

Don't get too attached to anything here. It might all change.

### wof-travel-id

```
./bin/wof-travel-id -source /usr/local/data/whosonfirst-data-venue-us-ca/data -supersedes 890535433
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
./bin/wof-travel-id -source /usr/local/data/whosonfirst-data-venue-us-ca/data -superseded-by 1108808531
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
./bin/wof-travel-id -source /usr/local/data/whosonfirst-data-venue-us-ca/data -source /usr/local/data/whosonfirst-data/data -supersedes -parent 890535433
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