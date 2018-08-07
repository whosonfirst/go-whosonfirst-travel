package main

import (
	"context"
	"flag"
	"github.com/whosonfirst/go-whosonfirst-cli/flags"
	"github.com/whosonfirst/go-whosonfirst-readwrite/reader"
	"github.com/whosonfirst/go-whosonfirst-travel"
	"github.com/whosonfirst/go-whosonfirst-travel/utils"
	"log"
)

func main() {

	var sources flags.MultiString
	flag.Var(&sources, "source", "One or more filesystem based sources to use to read WOF ID data, which may or may not be part of the sources to graph. This is work in progress.")

	var follow flags.MultiString
	flag.Var(&follow, "follow", "...")

	parent_id := flag.Bool("parent", false, "...")
	supersedes := flag.Bool("supersedes", false, "...")
	superseded_by := flag.Bool("superseded-by", false, "...")
	hierarchies := flag.Bool("hierarchies", false, "...")
	singleton := flag.Bool("singleton", true, "...")

	flag.Parse()

	r, err := reader.NewMultiReaderFromStrings(sources...)

	if err != nil {
		log.Fatal(err)
	}

	opts, err := travel.DefaultTravelOptions()

	if err != nil {
		log.Fatal(err)
	}

	opts.Reader = r

	opts.Singleton = *singleton
	opts.ParentID = *parent_id
	opts.Hierarchy = *hierarchies
	opts.Supersedes = *supersedes
	opts.SupersededBy = *superseded_by

	tr, err := travel.NewTraveler(opts)

	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for _, str_id := range flag.Args() {

		f, err := utils.LoadFeatureFromString(r, str_id)

		if err != nil {
			log.Fatal(err)
		}

		err = tr.TravelFeature(ctx, f)

		if err != nil {
			log.Fatal(err)
		}
	}

}
