package main

import (
	"flag"
	"github.com/whosonfirst/go-whosonfirst-cli/flags"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2"
	"github.com/whosonfirst/go-whosonfirst-travel"
	"github.com/whosonfirst/go-whosonfirst-travel/traveler"
	"log"
)

func main() {

	var belongs_to flags.MultiInt64
	flag.Var(&belongs_to, "belongs-to", "...")

	var include_placetype flags.MultiString
	flag.Var(&include_placetype, "include-placetype", "...")

	var exclude_placetype flags.MultiString
	flag.Var(&exclude_placetype, "exclude-placetype", "...")

	mode := flag.String("mode", "repo", "...")

	flag.Parse()

	default_cb, err := travel.DefaultTravelFunc()

	if err != nil {
		log.Fatal(err)
	}

	// we should make this a canned TravelFunc once we figure out
	// what the method signature looks like... (20180314/thisisaaronland)

	cb := func(f geojson.Feature, step int64) error {

		pt := f.Placetype()

		if len(include_placetype) > 0 {

			if !include_placetype.Contains(pt) {
				return nil
			}
		}

		if len(exclude_placetype) > 0 {

			if exclude_placetype.Contains(pt) {
				return nil
			}
		}

		return default_cb(f, step)
	}

	t, err := traveler.NewDefaultBelongsToTraveler()
	t.Mode = *mode
	t.BelongsTo = belongs_to
	t.Callback = cb

	paths := flag.Args()
	err = t.Travel(paths...)

	if err != nil {
		log.Fatal(err)
	}
}
