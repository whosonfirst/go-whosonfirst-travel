package main

import (
	"context"
	"errors"		
	"flag"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-cli/flags"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/feature"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/properties/whosonfirst"
	"github.com/whosonfirst/go-whosonfirst-index"
	"github.com/whosonfirst/warning"
	"io"
	"log"
)

type TravelFunc func(f geojson.Feature, container_id int64) error

type Traveler struct {
	Callback  TravelFunc
	Mode      string
	BelongsTo []int64
	// IncludePlacetype []string
	// ExcludePlacetype []string
}

func NewDefaultTravelFunc() (TravelFunc, error) {

	cb := func(f geojson.Feature, container_id int64) error {

		log.Printf("%s (%s) belongs to %d\n", f.Name(), f.Id(), container_id)
		return nil
	}

	return cb, nil
}

func NewDefaultTraveler() (*Traveler, error) {

	cb, err := NewDefaultTravelFunc()

	if err != nil {
		return nil, err
	}

	belongs := make([]int64, 0)

	t := Traveler{
		Callback:  cb,
		Mode:      "repo",
		BelongsTo: belongs,
	}

	return &t, nil
}

func (t *Traveler) Travel(paths ...string) error {

	idx_cb := func(fh io.Reader, ctx context.Context, args ...interface{}) error {

		path, err := index.PathForContext(ctx)

		if err != nil {
			return err
		}		

		f, err := feature.LoadFeatureFromReader(fh)

		if err != nil && !warning.IsWarning(err) {

			msg := fmt.Sprintf("Unable to load '%s' because %s", path, err)
			return errors.New(msg)
		}

		for _, id := range whosonfirst.BelongsTo(f) {

			for _, test := range t.BelongsTo {

				if test != id {
					continue
				}

				err := t.Callback(f, id)

				if err != nil {
					msg := fmt.Sprintf("Unable to process '%s' because %s", path, err)				
					return errors.New(msg)
				}
			}
		}

		return nil
	}

	idx, err := index.NewIndexer(t.Mode, idx_cb)

	if err != nil {
		return err
	}

	for _, path := range paths {

		err := idx.IndexPath(path)

		if err != nil {
			return err
		}
	}

	return nil
}

func main() {

	var belongs_to flags.MultiInt64
	flag.Var(&belongs_to, "belongs-to", "...")

	mode := flag.String("mode", "repo", "...")

	flag.Parse()

	t, err := NewDefaultTraveler()
	t.Mode = *mode
	t.BelongsTo = belongs_to

	paths := flag.Args()
	err = t.Travel(paths...)

	if err != nil {
		log.Fatal(err)
	}
}
