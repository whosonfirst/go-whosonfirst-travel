package travel

import (
	"context"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/properties/whosonfirst"
	"github.com/whosonfirst/go-whosonfirst-readwrite/reader"
	"github.com/whosonfirst/go-whosonfirst-travel/utils"
	"log"
	"sync"
)

type TravelFunc func(f geojson.Feature) error

type TravelOptions struct {
	Callback     TravelFunc
	Reader       reader.Reader
	Singleton    bool
	Supersedes   bool
	SupersededBy bool
	ParentID     bool
	Hierarchy    bool
	Depth        int
}

func DefaultTravelFunc() (TravelFunc, error) {

	f := func(f geojson.Feature) error {

		id := f.Id()
		name := f.Name()
		
		label := whosonfirst.Label(f)

		if label == "" {
			label = name
		}
		
		log.Printf("%s %s\n", id, label)
		return nil
	}

	return f, nil
}

func DefaultTravelOptions() (*TravelOptions, error) {

	cb, err := DefaultTravelFunc()

	if err != nil {
		return nil, err
	}

	r, err := reader.NewNullReader()

	if err != nil {
		return nil, err
	}

	opts := TravelOptions{
		Callback:     cb,
		Reader:       r,
		Singleton:    true,
		Supersedes:   false,
		SupersededBy: false,
		ParentID:     false,
		Hierarchy:    false,
		Depth:        0,
	}

	return &opts, nil
}

type Traveler struct {
	Options  *TravelOptions
	mu       *sync.RWMutex
	travelog map[string]int
}

func NewTraveler(opts *TravelOptions) (*Traveler, error) {

	travelog := make(map[string]int)

	mu := new(sync.RWMutex)

	t := Traveler{
		Options:  opts,
		mu:       mu,
		travelog: travelog,
	}

	return &t, nil
}

func (t *Traveler) TravelFeature(ctx context.Context, f geojson.Feature) error {

	select {
	case <-ctx.Done():
		return nil
	default:
		// pass
	}

	opts := t.Options

	t.mu.RLock()

	str_id := f.Id()
	visits, visited := t.travelog[str_id]

	if opts.Singleton && visited {
		t.mu.RUnlock()
		return nil
	}

	t.mu.RUnlock()

	cb := opts.Callback
	err := cb(f)

	if err != nil {
		return err
	}

	t.mu.Lock()

	if !visited {
		visits = 1
	} else {
		visits += 1
	}

	t.travelog[str_id] = visits
	t.mu.Unlock()

	wg := new(sync.WaitGroup)

	if opts.ParentID {

		wg.Add(1)

		go func() {
			defer wg.Done()
			pid := whosonfirst.ParentId(f)
			t.TravelID(ctx, pid)
		}()

	}

	if opts.Supersedes {

		wg.Add(1)

		go func() {

			defer wg.Done()
			for _, id := range whosonfirst.Supersedes(f) {
				t.TravelID(ctx, id)
			}
		}()

	}

	if opts.SupersededBy {

		wg.Add(1)

		go func() {

			defer wg.Done()

			for _, id := range whosonfirst.SupersededBy(f) {
				t.TravelID(ctx, id)
			}

		}()
	}

	if opts.Hierarchy {

		wg.Add(1)

		go func() {

			defer wg.Done()

			for _, hier := range whosonfirst.Hierarchies(f) {

				for _, id := range hier {
					t.TravelID(ctx, id)
				}
			}

		}()
	}

	wg.Wait()

	return nil
}

func (t *Traveler) TravelID(ctx context.Context, id int64) error {

	select {
	case <-ctx.Done():
		return nil
	default:
		// pass
	}

	opts := t.Options

	f, err := utils.LoadFeature(opts.Reader, id)

	if err != nil {
		return err
	}

	return t.TravelFeature(ctx, f)
}
