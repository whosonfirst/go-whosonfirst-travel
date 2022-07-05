package traveler

import (
	"context"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-feature/properties"
	"github.com/whosonfirst/go-whosonfirst-iterate/v2/iterator"
	"github.com/whosonfirst/go-whosonfirst-uri"
	"io"
)

// type BelongsToTravelFunc defines custom callback function to be invoked for records matching a "wof:belongs_to" condition.
type BelongsToTravelFunc func(context.Context, []byte, int64) error

// type BelongsToTraveler defines a struct
type BelongsToTraveler struct {
	// Callback is a `BelongsToTravelFunc` function to be invoked for records matching a "wof:belongs_to" condition.
	Callback BelongsToTravelFunc
	// IteratorURI is a valid `whosonfirst/go-whosonfirst-iterate/v2` URI used to crawl (iterate) records.
	IteratorURI string
	// BelongsTo list of WOF IDs that matching records should belong to.
	BelongsTo []int64
}

// NewDefaultBelongsToTravelFunc() returns a `BelongsToTravelFunc` instance that prints metadata about records
// matching a "wof:belongs_to" condition to STDOUT.
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

		fmt.Printf("%s (%d) belongs to %d\n", name, id, container_id)
		return nil
	}

	return cb, nil
}

// NewDefaultBelongsToTraveler() returns a a `BelongsToTraveler` with default values. Specifically
// one that expects to iterate documents in `repo://` mode, using a default callback function (defined
// by `NewDefaultBelongsToTravelFunc` and an empty list of the IDs that matching records should belong
// to.
func NewDefaultBelongsToTraveler() (*BelongsToTraveler, error) {

	cb, err := NewDefaultBelongsToTravelFunc()

	if err != nil {
		return nil, fmt.Errorf("Failed to create DefaultBelongsToTravelFunc, %v", err)
	}

	belongs := make([]int64, 0)

	t := BelongsToTraveler{
		Callback:    cb,
		IteratorURI: "repo://",
		BelongsTo:   belongs,
	}

	return &t, nil
}

// Travel() iterates through all the records emitted by 'uris' and invokes `t.Callback` for records
// whose "wof:belongs_to" values match one or more of the IDs defined in `t.BelongTo`.
func (t *BelongsToTraveler) Travel(ctx context.Context, uris ...string) error {

	iter_cb := func(ctx context.Context, path string, fh io.ReadSeeker, args ...interface{}) error {

		is_alt, err := uri.IsAltFile(path)

		if err != nil {
			return fmt.Errorf("Failed to determine whether %s is alt file, %w", path, err)
		}

		if is_alt {
			return nil
		}

		f, err := io.ReadAll(fh)

		if err != nil {
			return fmt.Errorf("Unable to load '%s' because %w", path, err)
		}

		belongs_to := properties.BelongsTo(f)

		for _, id := range belongs_to {

			for _, test := range t.BelongsTo {

				if test != id {
					continue
				}

				err := t.Callback(ctx, f, id)

				if err != nil {
					return fmt.Errorf("Unable to process '%s' because %w", path, err)
				}
			}
		}

		return nil
	}

	iter, err := iterator.NewIterator(ctx, t.IteratorURI, iter_cb)

	if err != nil {
		return fmt.Errorf("Failed to create new iterator, %w", err)
	}

	err = iter.IterateURIs(ctx, uris...)

	if err != nil {
		return fmt.Errorf("Failed to iterate URIs, %w", err)
	}

	return nil
}
