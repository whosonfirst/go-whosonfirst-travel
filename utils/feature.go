package utils

import (
	"github.com/whosonfirst/go-whosonfirst-geojson-v2"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/feature"
	"github.com/whosonfirst/go-whosonfirst-readwrite/reader"
	"github.com/whosonfirst/go-whosonfirst-uri"
	"strconv"
)

// put this in go-whosonfirst-readwrite? would it ever be used by anything
// but this... ?

func LoadFeatureFromString(r reader.Reader, str_id string) (geojson.Feature, error) {

	id, err := strconv.ParseInt(str_id, 10, 64)

	if err != nil {
		return nil, err
	}

	return LoadFeature(r, id)
}

func LoadFeature(r reader.Reader, id int64) (geojson.Feature, error) {

	uri, err := uri.Id2RelPath(id)

	if err != nil {
		return nil, err
	}

	fh, err := r.Read(uri)

	if err != nil {
		return nil, err
	}

	return feature.LoadWOFFeatureFromReader(fh)
}
