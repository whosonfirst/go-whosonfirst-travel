package main

import (
	"flag"
	"github.com/whosonfirst/go-whosonfirst-cli/flags"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2"
	"github.com/whosonfirst/go-whosonfirst-readwrite/reader"
	"github.com/whosonfirst/go-whosonfirst-travel"
	"github.com/whosonfirst/go-whosonfirst-travel/utils"
	"log"
)

func Travel(f geojson.Feature, opts *travel.TravelOptions) error {

	err := opts.Callback(f)

	if err != nil {
		return err
	}

	return nil
}

func main() {

	var sources flags.MultiString
	flag.Var(&sources, "source", "One or more filesystem based sources to use to read WOF ID data, which may or may not be part of the sources to graph. This is work in progress.")

	var follow flags.MultiString
	flag.Var(&follow, "follow", "...")

	flag.Parse()

	r, err := reader.NewMultiReaderFromStrings(sources...)

	if err != nil {
		log.Fatal(err)
	}

	for _, str_id := range flag.Args() {

		f, err := utils.LoadFeatureFromString(r, str_id)

		if err != nil {
			log.Fatal(err)
		}

		opts, err := travel.DefaultTravelOptions()

		if err != nil {
			log.Fatal(err)
		}

		err = Travel(f, opts)

		if err != nil {
			log.Fatal(err)
		}
	}

}
