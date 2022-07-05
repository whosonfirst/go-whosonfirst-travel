package travel

import (
	"context"
	"fmt"
	"github.com/whosonfirst/go-reader"
	"path/filepath"
	"sync/atomic"
	"testing"
)

func TestDefaultTravelFunc(t *testing.T) {

	_, err := DefaultTravelFunc()

	if err != nil {
		t.Fatalf("Failed to create default travel func, %v", err)
	}
}

func TestDefaultTravelOptions(t *testing.T) {

	_, err := DefaultTravelOptions()

	if err != nil {
		t.Fatalf("Failed to create default travel options, %v", err)
	}
}

func TestTravelID(t *testing.T) {

	ctx := context.Background()

	sfo := int64(102527513)

	count := int32(0)
	expected := int32(1)

	tr_opts, err := DefaultTravelOptions()

	if err != nil {
		t.Fatalf("Failed to create default travel options, %v", err)
	}

	rel_path := "fixtures/data"
	abs_path, err := filepath.Abs(rel_path)

	if err != nil {
		t.Fatalf("Failed to derive absolute path for %s, %v", rel_path, err)
	}

	reader_uri := fmt.Sprintf("fs://%s", abs_path)

	r, err := reader.NewReader(ctx, reader_uri)

	if err != nil {
		t.Fatalf("Failed to create new reader, %v", err)
	}

	cb := func(ctx context.Context, body []byte, id int64) error {
		atomic.AddInt32(&count, 1)
		return nil
	}

	tr_opts.Reader = r
	tr_opts.Callback = cb

	tr, err := NewTraveler(tr_opts)

	if err != nil {
		t.Fatalf("Failed to create new traveler, %v", err)
	}

	err = tr.TravelID(ctx, sfo)

	if err != nil {
		t.Fatalf("Failed to travel %d, %v", sfo, err)
	}

	if count != expected {
		t.Fatalf("Unexpected count: %d", count)
	}
}
