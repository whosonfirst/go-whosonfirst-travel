package traveler

import (
	"context"
	"errors"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-feature/properties"
	"github.com/whosonfirst/go-whosonfirst-iterate/v2/iterator"
	"github.com/whosonfirst/go-whosonfirst-uri"
	"io"
	"log"
)

type BelongsToTravelFunc func(context.Context, []byte, int64) error

type BelongsToTraveler struct {
	Callback  BelongsToTravelFunc
	Mode      string
	BelongsTo []int64
}

func NewDefaultBelongsToTravelFunc() (BelongsToTravelFunc, error) {

	cb := func(ctx context.Context, f []byte, container_id int64) error {

		id, err := properties.Id(f)

		if err != nil {
			return fmt.Errorf("Failed to derive ID, %w", err)
		}

		name, err := properties.Name(f)

		if err != nil {
			return fmt.Errorf("Failed to derive name, %w", err)
		}

		log.Printf("%s (%d) belongs to %d\n", name, id, container_id)
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

	iter_cb := func(ctx context.Context, path string, fh io.ReadSeeker, args ...interface{}) error {

		is_alt, err := uri.IsAltFile(path)

		if err != nil {
			return err
		}

		if is_alt {
			return nil
		}

		f, err := io.ReadAll(fh)

		if err != nil {
			msg := fmt.Sprintf("Unable to load '%s' because %s", path, err)
			return errors.New(msg)
		}

		belongs_to := properties.BelongsTo(f)

		for _, id := range belongs_to {

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
