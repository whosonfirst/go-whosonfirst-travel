package travel

import (
	"context"
	"fmt"
	"github.com/whosonfirst/go-reader"
	"github.com/whosonfirst/go-whosonfirst-feature/properties"
	wof_reader "github.com/whosonfirst/go-whosonfirst-reader"
	"log"
	"sync"
	"time"
)

// TravelFunc is a callback function to be invoked for each `geojson.Feature` encountered during a travel session.
type TravelFunc func(context.Context, []byte, int64) error

// TravelOptions is a struct containing configuration details for a travel session.
type TravelOptions struct {
	// TravelFunc is a callback function to be invoked for each `geojson.Feature` encountered during a travel session.
	Callback TravelFunc
	// A `reader.Reader` instance used to load GeoJSON Feature data.
	Reader reader.Reader
	// A boolean flag to indicate whether to record timing information.
	Timings bool
	// A boolean flag to indcate whether or not the same record should be traveled more than once. If true then records will only be traveled once.
	Singleton bool
	// A boolean flag to indicate whether a travel session should include the records that a feature supersedes .
	Supersedes bool
	// A boolean flag to indicate whether a travel session should include the records that a feature is superseded by .
	SupersededBy bool
	// A boolean flag to indicate whether a travel session should include a feature's parent record.
	ParentID bool
	// A boolean flag to indicate whether a travel session should include the record in a feature's hierarchy.
	Hierarchy bool
}

// DefaultTravelFunc returns a TravelFunc callback function that prints the current step, the feature's ID and name as well as its inception and cessation dates.
func DefaultTravelFunc() (TravelFunc, error) {

	f := func(ctx context.Context, body []byte, step int64) error {

		id, err := properties.Id(body)

		if err != nil {
			return fmt.Errorf("Failed to derive ID, %w", err)
		}

		label, err := properties.Name(body)

		if err != nil {
			return fmt.Errorf("Failed to derive name, %w", err)
		}

		inception := properties.Inception(body)
		cessation := properties.Cessation(body)

		is_deprecated := ""

		deprecated, err := properties.IsDeprecated(body)

		if err != nil {
			return err
		}

		if deprecated.IsKnown() && deprecated.IsTrue() {
			is_deprecated = "DEPRECATED"
		}

		fmt.Printf("[%d] %s %s [%s] [%s] %s\n", step, id, label, inception, cessation, is_deprecated)
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
	}

	return &opts, nil
}

// Traveler is a struct for walking the tree of supersedes or superseded_by relations for a Who's On First record.
type Traveler struct {
	// Options is a TravelOptions struct containing configuration details for the travel session.
	Options  *TravelOptions
	mu       *sync.RWMutex
	travelog map[int64]int
	Step     int64
}

// Create a new Traveler instance.
func NewTraveler(opts *TravelOptions) (*Traveler, error) {

	travelog := make(map[int64]int)

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
func (t *Traveler) TravelFeature(ctx context.Context, f []byte) error {

	select {
	case <-ctx.Done():
		return nil
	default:
		// pass
	}

	opts := t.Options

	t.mu.RLock()

	id, err := properties.Id(f)

	if err != nil {
		return fmt.Errorf("Failed to derive ID, %w", id)
	}

	visits, visited := t.travelog[id]

	if opts.Singleton && visited {
		t.mu.RUnlock()
		return nil
	}

	t1 := time.Now()

	if opts.Timings {

		defer func() {
			log.Printf("time to travel feature ID %d %v\n", id, time.Since(t1))
		}()
	}

	t.mu.RUnlock()

	t.mu.Lock()
	t.Step += 1
	step := t.Step
	t.mu.Unlock()

	cb := opts.Callback
	err = cb(ctx, f, step)

	if err != nil {
		return err
	}

	t.mu.Lock()

	if !visited {
		visits = 1
	} else {
		visits += 1
	}

	t.travelog[id] = visits
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

	f, err := wof_reader.LoadBytes(ctx, opts.Reader, id)

	if err != nil {
		return err
	}

	return t.TravelFeature(ctx, f)
}

func (t *Traveler) travelParent(ctx context.Context, f []byte) error {

	parent_id, err := properties.ParentId(f)

	if err != nil {
		return fmt.Errorf("Failed to derive parent ID, %w", err)
	}

	return t.TravelID(ctx, parent_id)
}

func (t *Traveler) travelSupersedes(ctx context.Context, f []byte) error {

	supersedes := properties.Supersedes(f)

	for _, id := range supersedes {
		t.TravelID(ctx, id)
	}

	return nil
}

func (t *Traveler) travelSupersededBy(ctx context.Context, f []byte) error {

	superseded_by := properties.SupersededBy(f)

	for _, id := range superseded_by {
		t.TravelID(ctx, id)
	}

	return nil
}

func (t *Traveler) travelHierarchies(ctx context.Context, f []byte) error {

	hierarchies := properties.Hierarchies(f)

	for _, hier := range hierarchies {

		for _, id := range hier {
			t.TravelID(ctx, id)
		}
	}

	return nil
}
