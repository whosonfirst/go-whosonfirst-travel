package main

import (
	"bufio"
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/sfomuseum/go-flags/multi"
	"github.com/whosonfirst/go-reader"
	"github.com/whosonfirst/go-whosonfirst-feature/properties"
	wof_reader "github.com/whosonfirst/go-whosonfirst-reader"
	"github.com/whosonfirst/go-whosonfirst-travel"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {

	var sources multi.MultiString
	flag.Var(&sources, "source", "One or more valid whosonfirst/go-reader URIs to use to read WOF ID data, which may or may not be part of the sources to graph. This is work in progress.")

	parent_id := flag.Bool("parent", false, "Travel the parent ID of each ID specified.")
	supersedes := flag.Bool("supersedes", false, "Travel records that are superseded by each ID specified.")
	superseded_by := flag.Bool("superseded-by", false, "Travel records that supersede each ID specified.")
	hierarchies := flag.Bool("hierarchies", false, "Travel the hierarchies for each ID specified.")
	singleton := flag.Bool("singleton", true, "...")
	timings := flag.Bool("timings", false, "Display timing information")

	ids := flag.Bool("ids", false, "Emit results as a line-separated list of IDs (traveled).")
	as_markdown := flag.Bool("markdown", false, "Emit results formatted as Markdown.")
	as_csv := flag.Bool("csv", false, "Emit results formatted as a comma-separated values.")

	from_stdin := flag.Bool("stdin", false, "Read IDs to travel from STDIN")

	flag.Parse()

	to_travel := flag.Args()

	if *from_stdin {

		scanner := bufio.NewScanner(os.Stdin)

		for scanner.Scan() {
			to_travel = append(to_travel, scanner.Text())
		}

		err := scanner.Err()

		if err != nil {
			log.Fatalf("Failed to read input from STDIN, %v", err)
		}
	}

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

		cb := func(ctx context.Context, f []byte, step int64) error {

			id, err := properties.Id(f)

			if err != nil {
				return fmt.Errorf("Failed to derive ID, %w", err)
			}

			fmt.Println(id)
			return nil
		}

		opts.Callback = cb
	}

	if *as_markdown {

		cb := func(ctx context.Context, f []byte, step int64) error {

			if step == 1 {

				fmt.Println("| step | id | label |")
				fmt.Println("| --- | --- | --- |")
			}

			id, err := properties.Id(f)

			if err != nil {
				return fmt.Errorf("Failed to derive ID, %w", err)
			}

			label, err := properties.Name(f)

			if err != nil {
				return fmt.Errorf("Failed to derive name, %w", err)
			}

			fmt.Printf("| %d | %d | %s |\n", step, id, label)
			return nil
		}

		opts.Callback = cb
	}

	if *as_csv {

		wr := csv.NewWriter(os.Stdout)

		cb := func(ctx context.Context, f []byte, step int64) error {

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

			id, err := properties.Id(f)

			if err != nil {
				return fmt.Errorf("Faild to derive ID, %w", err)
			}

			label, err := properties.Name(f)

			if err != nil {
				return fmt.Errorf("Faild to derive name, %w", err)
			}

			parent_id, err := properties.ParentId(f)

			if err != nil {
				return fmt.Errorf("Faild to derive parent ID, %w", err)
			}

			inception := properties.Inception(f)
			cessation := properties.Cessation(f)

			supersedes := properties.Supersedes(f)
			supersedes_str := make([]string, len(supersedes))

			for idx, id := range supersedes {
				supersedes_str[idx] = strconv.FormatInt(id, 10)
			}

			superseded_by := properties.SupersededBy(f)
			superseded_by_str := make([]string, len(superseded_by))

			for idx, id := range superseded_by {
				superseded_by_str[idx] = strconv.FormatInt(id, 10)
			}

			out := []string{
				strconv.FormatInt(step, 10),
				strconv.FormatInt(id, 10),
				label,
				inception,
				cessation,
				strconv.FormatInt(parent_id, 10),
				strings.Join(supersedes_str, ","),
				strings.Join(superseded_by_str, ","),
			}

			err = wr.Write(out)

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

	for _, str_id := range to_travel {

		id, err := strconv.ParseInt(str_id, 10, 64)

		if err != nil {
			log.Fatal(err)
		}

		f, err := wof_reader.LoadBytes(ctx, r, id)

		if err != nil {
			log.Fatal(err)
		}

		err = tr.TravelFeature(ctx, f)

		if err != nil {
			log.Fatal(err)
		}
	}

}
