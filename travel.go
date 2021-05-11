package travel

import (
	"context"
	"fmt"
	"github.com/whosonfirst/go-reader"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/properties/whosonfirst"
	wof_reader "github.com/whosonfirst/go-whosonfirst-reader"
	"log"
	"sync"
	"time"
)

// TravelFunc is a callback function to be invoked for each `geojson.Feature` encountered during a travel session.
type TravelFunc func(context.Context, geojson.Feature, int64) error

// TravelOptions is a struct containing configuration details for a travel session.
type TravelOptions struct {
	// TravelFunc is a callback function to be invoked for each `geojson.Feature` encountered during a travel session.
	Callback TravelFunc
	// A `reader.Reader` instance used to load GeoJSON Feature data.
	Reader reader.Reader
	// A boolean flag to indicate whether to record timing information.
	Timings   bool
	Singleton bool
	// A boolean flag to indicate whether a travel session should include the records that a feature supersedes .
	Supersedes bool
	// A boolean flag to indicate whether a travel session should include the records that a feature is superseded by .
	SupersededBy bool
	// A boolean flag to indicate whether a travel session should include a feature's parent record.
	ParentID bool
	// A boolean flag to indicate whether a travel session should include the record in a feature's hierarchy.
	Hierarchy bool
	Depth     int
}

// DefaultTravelFunc returns a TravelFunc callback function that prints the current step, GeoJSON Feature's label string and ID.
func DefaultTravelFunc() (TravelFunc, error) {

	f := func(ctx context.Context, f geojson.Feature, step int64) error {

		id := f.Id()
		label := whosonfirst.LabelOrDerived(f)

		fmt.Printf("[%d] %s %s\n", step, id, label)
		return nil
	}

	return f, nil
}

// DefaultTravelOptions returns a TravelOptions struct configured as a singleton and to use the DefaultTravelFunc callback, a `null://` reader.
func DefaultTravelOptions() (*TravelOptions, error) {

	cb, err := DefaultTravelFunc()

	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	r, err := reader.NewReader(ctx, "null://")

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

// Traveler is a struct for walking the tree of supersedes or superseded_by relations for a Who's On First record.
type Traveler struct {
	// Options is a TravelOptions struct containing configuration details for the travel session.
	Options  *TravelOptions
	mu       *sync.RWMutex
	travelog map[string]int
	Step     int64
}

// Create a new Traveler instance.
func NewTraveler(opts *TravelOptions) (*Traveler, error) {

	travelog := make(map[string]int)

	mu := new(sync.RWMutex)

	t := Traveler{
		Options:  opts,
		mu:       mu,
		travelog: travelog,
		Step:     0,
	}

	return &t, nil
}

// Travel the relationships for a `geojson.Feature` instance.
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

	t1 := time.Now()

	if opts.Timings {

		defer func() {
			log.Printf("time to travel feature ID %s %v\n", str_id, time.Since(t1))
		}()
	}

	t.mu.RUnlock()

	t.mu.Lock()
	t.Step += 1
	step := t.Step
	t.mu.Unlock()

	cb := opts.Callback
	err := cb(ctx, f, step)

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
			t.travelParent(ctx, f)
		}()
	}

	if opts.Supersedes {

		wg.Add(1)

		go func() {
			defer wg.Done()
			t.travelSupersedes(ctx, f)
		}()
	}

	if opts.SupersededBy {

		wg.Add(1)

		go func() {
			defer wg.Done()
			t.travelSupersededBy(ctx, f)
		}()
	}

	if opts.Hierarchy {

		wg.Add(1)

		go func() {
			defer wg.Done()
			t.travelHierarchies(ctx, f)
		}()
	}

	wg.Wait()

	return nil
}

// Travel the relationships for a Who's On First ID.
// The ID must be able to be read by the Traveler's `reader.Reader` instance defined in the `TravelOptions`.
func (t *Traveler) TravelID(ctx context.Context, id int64) error {

	select {
	case <-ctx.Done():
		return nil
	default:
		// pass
	}

	opts := t.Options

	f, err := wof_reader.LoadFeatureFromID(ctx, opts.Reader, id)

	if err != nil {
		return err
	}

	return t.TravelFeature(ctx, f)
}

func (t *Traveler) travelParent(ctx context.Context, f geojson.Feature) error {

	id := whosonfirst.ParentId(f)
	return t.TravelID(ctx, id)
}

func (t *Traveler) travelSupersedes(ctx context.Context, f geojson.Feature) error {

	for _, id := range whosonfirst.Supersedes(f) {
		t.TravelID(ctx, id)
	}

	return nil
}

func (t *Traveler) travelSupersededBy(ctx context.Context, f geojson.Feature) error {

	for _, id := range whosonfirst.SupersededBy(f) {
		t.TravelID(ctx, id)
	}

	return nil
}

func (t *Traveler) travelHierarchies(ctx context.Context, f geojson.Feature) error {

	for _, hier := range whosonfirst.Hierarchies(f) {

		for _, id := range hier {
			t.TravelID(ctx, id)
		}
	}

	return nil
}
