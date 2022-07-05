package traveler

import (
	"context"
	"sync/atomic"
	"testing"
)

func TestBelongsToTraveler(t *testing.T) {

	ctx := context.Background()

	count := int32(0)
	expected := int32(3)

	cb := func(ctx context.Context, body []byte, id int64) error {
		atomic.AddInt32(&count, 1)
		return nil
	}

	california := int64(85688637)

	belongs := []int64{california}

	tr := BelongsToTraveler{
		Callback:  cb,
		Mode:      "repo://",
		BelongsTo: belongs,
	}

	err := tr.Travel(ctx, "../fixtures")

	if err != nil {
		t.Fatalf("Failed to determine belongs to, %v", err)
	}

	if count != expected {
		t.Fatalf("Unexpected count, %d", count)
	}
}
