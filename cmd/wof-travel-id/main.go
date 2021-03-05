package main

import (
	"context"
	"flag"
	"fmt"
	"encoding/csv"
	"github.com/sfomuseum/go-flags/multi"
	"github.com/whosonfirst/go-reader"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/properties/whosonfirst"
	"github.com/whosonfirst/go-whosonfirst-travel"
	"github.com/whosonfirst/go-whosonfirst-travel/utils"
	"log"
	"os"
	"strings"
	"strconv"
)

func main() {

	var sources multi.MultiString
	flag.Var(&sources, "source", "One or more filesystem based sources to use to read WOF ID data, which may or may not be part of the sources to graph. This is work in progress.")

	var follow multi.MultiString
	flag.Var(&follow, "follow", "...")

	parent_id := flag.Bool("parent", false, "...")
	supersedes := flag.Bool("supersedes", false, "...")
	superseded_by := flag.Bool("superseded-by", false, "...")
	hierarchies := flag.Bool("hierarchies", false, "...")
	singleton := flag.Bool("singleton", true, "...")
	timings := flag.Bool("timings", false, "...")

	ids := flag.Bool("ids", false, "...")
	as_markdown := flag.Bool("markdown", false, "...")
	as_csv := flag.Bool("csv", false, "...")
	
	flag.Parse()

	ctx := context.Background()

	r, err := reader.NewMultiReaderFromURIs(ctx, sources...)

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
	opts.Timings = *timings

	// please move these in to travel.go or equivalent...
	// (20180815/thisisaaronland)

	if *ids {

		cb := func(ctx context.Context, f geojson.Feature, step int64) error {
			fmt.Println(f.Id())
			return nil
		}

		opts.Callback = cb
	}

	if *as_markdown {

		cb := func(ctx context.Context, f geojson.Feature, step int64) error {

			if step == 1 {

				fmt.Println("| step | id | label |")
				fmt.Println("| --- | --- | --- |")
			}

			id := f.Id()
			label := whosonfirst.LabelOrDerived(f)

			fmt.Printf("| %d | %s | %s |\n", step, id, label)
			return nil
		}

		opts.Callback = cb
	}

	if *as_csv {

		wr := csv.NewWriter(os.Stdout)
		
		cb := func(ctx context.Context, f geojson.Feature, step int64) error {

			if step == 1 {

				out := []string{
					"step",
					"id",
					"label",
					"inception",
					"cessation",
					"parent",					
					"supersedes",
					"superseded_by",
				}

				err := wr.Write(out)

				if err != nil {
					return err
				}
			}

			id := f.Id()
			label := whosonfirst.LabelOrDerived(f)

			parent := whosonfirst.ParentId(f)
			
			supersedes := whosonfirst.Supersedes(f)
			supersedes_str := make([]string, len(supersedes))
			
			for idx, id := range supersedes {
				supersedes_str[idx] = strconv.FormatInt(id, 10)
			}

			superseded_by := whosonfirst.SupersededBy(f)
			superseded_by_str := make([]string, len(superseded_by))

			for idx, id := range superseded_by {
				superseded_by_str[idx] = strconv.FormatInt(id, 10)
			}
			
			out := []string{
				strconv.FormatInt(step, 10),
				id,
				label,
				whosonfirst.Inception(f),
				whosonfirst.Cessation(f),
				strconv.FormatInt(parent, 10),								
				strings.Join(supersedes_str, ","),
				strings.Join(superseded_by_str, ","),				
			}

			err := wr.Write(out)

			if err != nil {
				return err
			}

			wr.Flush()
			return nil
		}

		opts.Callback = cb
	}
	
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
