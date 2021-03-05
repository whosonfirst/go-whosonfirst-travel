package traveler

import (
	"context"
	"errors"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/feature"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/properties/whosonfirst"
	"github.com/whosonfirst/go-whosonfirst-iterate/emitter"
	"github.com/whosonfirst/go-whosonfirst-iterate/iterator"
	"github.com/whosonfirst/go-whosonfirst-uri"
	"github.com/whosonfirst/warning"
	"io"
	"log"
)

type BelongsToTravelFunc func(context.Context, geojson.Feature, int64) error

type BelongsToTraveler struct {
	Callback  BelongsToTravelFunc
	Mode      string
	BelongsTo []int64
}

func NewDefaultBelongsToTravelFunc() (BelongsToTravelFunc, error) {

	cb := func(ctx context.Context, f geojson.Feature, container_id int64) error {

		log.Printf("%s (%s) belongs to %d\n", f.Name(), f.Id(), container_id)
		return nil
	}

	return cb, nil
}

func NewDefaultBelongsToTraveler() (*BelongsToTraveler, error) {

	cb, err := NewDefaultBelongsToTravelFunc()

	if err != nil {
		return nil, err
	}

	belongs := make([]int64, 0)

	t := BelongsToTraveler{
		Callback:  cb,
		Mode:      "repo",
		BelongsTo: belongs,
	}

	return &t, nil
}

func (t *BelongsToTraveler) Travel(ctx context.Context, uris ...string) error {

	iter_cb := func(ctx context.Context, fh io.ReadSeeker, args ...interface{}) error {

		path, err := emitter.PathForContext(ctx)

		if err != nil {
			return err
		}

		is_wof, err := uri.IsWOFFile(path)

		if err != nil {
			return err
		}

		if !is_wof {
			return nil
		}

		is_alt, err := uri.IsAltFile(path)

		if err != nil {
			return err
		}

		if is_alt {
			return nil
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

				err := t.Callback(ctx, f, id)

				if err != nil {
					msg := fmt.Sprintf("Unable to process '%s' because %s", path, err)
					return errors.New(msg)
				}
			}
		}

		return nil
	}

	iter, err := iterator.NewIterator(ctx, t.Mode, iter_cb)

	if err != nil {
		return err
	}

	return iter.IterateURIs(ctx, uris...)
}
