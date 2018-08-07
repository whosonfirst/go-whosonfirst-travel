package travel

import (
	"github.com/whosonfirst/go-whosonfirst-geojson-v2"
	"log"
)

type TravelFunc func(f geojson.Feature) error

type TravelOptions struct {
	Callback     TravelFunc
	Supersedes   bool
	SupersededBy bool
	ParentID     bool
	Hierarchy    bool
	Depth        int
}

func DefaultTravelFunc() (TravelFunc, error) {

	f := func(f geojson.Feature) error {

		log.Println(f.Name())
		return nil
	}

	return f, nil
}

func DefaultTravelOptions() (*TravelOptions, error) {

	cb, err := DefaultTravelFunc()

	if err != nil {
		return nil, err
	}

	opts := TravelOptions{
		Callback:     cb,
		Supersedes:   false,
		SupersededBy: false,
		ParentID:     false,
		Hierarchy:    false,
		Depth:        0,
	}

	return &opts, nil
}
