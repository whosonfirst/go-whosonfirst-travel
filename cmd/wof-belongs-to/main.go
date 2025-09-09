package main

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/sfomuseum/go-flags/multi"
	"github.com/whosonfirst/go-whosonfirst-feature/properties"
	"github.com/whosonfirst/go-whosonfirst-travel/traveler"
)

// please move all the BelongsToResult and BelongsToResultSet stuff
// in to a proper re-usalbe package (20180814/thisisaaronland)

type BelongsToResultSet struct {
	results []*BelongsToResult
	mu      *sync.RWMutex
}

func NewBelongsToResultSet() (*BelongsToResultSet, error) {

	results := make([]*BelongsToResult, 0)
	mu := new(sync.RWMutex)

	rs := BelongsToResultSet{
		results: results,
		mu:      mu,
	}

	return &rs, nil
}

func (rs *BelongsToResultSet) AddResult(r *BelongsToResult) error {

	rs.mu.Lock()
	defer rs.mu.Unlock()

	rs.results = append(rs.results, r)
	return nil
}

func (rs *BelongsToResultSet) Results() []*BelongsToResult {

	return rs.results
}

func (rs *BelongsToResultSet) Sort() {

	sort.Slice(rs.results, func(i, j int) bool {

		str_i := rs.results[i].Placetype
		str_j := rs.results[j].Placetype

		switch strings.Compare(str_i, str_j) {
		case -1:
			return true
		case 1:
			return false
		}

		return rs.results[i].Label < rs.results[j].Label
	})
}

func (rs *BelongsToResultSet) AsJSON(wr io.Writer) error {

	b, err := json.Marshal(rs.Results())

	if err != nil {
		return err
	}

	br := bytes.NewReader(b)
	_, err = io.Copy(wr, br)

	if err != nil {
		return err
	}

	return nil
}

func (rs *BelongsToResultSet) AsMarkdown(wr io.Writer) error {

	for i, r := range rs.Results() {

		if i == 0 {

			head, _ := r.ToCSVHeader()

			for _, col := range head {
				wr.Write([]byte(fmt.Sprintf("| %s ", col)))
			}

			wr.Write([]byte("|\n"))

			for range head {
				wr.Write([]byte("| --- "))
			}

			wr.Write([]byte("|\n"))
		}

		out, _ := r.ToCSVRow()

		for _, col := range out {
			wr.Write([]byte(fmt.Sprintf("| %s ", col)))
		}

		wr.Write([]byte("|\n"))
	}

	return nil
}

func (rs *BelongsToResultSet) AsCSV(wr io.Writer, header bool) error {

	csv_wr := csv.NewWriter(wr)

	for i, r := range rs.Results() {

		if i == 0 && header {

			head, _ := r.ToCSVHeader()
			err := csv_wr.Write(head)

			if err != nil {
				return err
			}

		}

		out, _ := r.ToCSVRow()
		err := csv_wr.Write(out)

		if err != nil {
			return err
		}
	}

	csv_wr.Flush()

	return nil
}

type BelongsToResult struct {
	BelongsToId int64  `json:"belongs_to"`
	Id          int64  `json:"id"`
	ParentId    int64  `json:"parent_id"`
	Placetype   string `json:"placetype"`
	Label       string `json:"label"`
}

func (r *BelongsToResult) ToCSVHeader() ([]string, error) {

	head := []string{
		"belongs_to",
		"id",
		"parent_id",
		"placetype",
		"label",
	}

	return head, nil
}

func (r *BelongsToResult) ToCSVRow() ([]string, error) {

	out := []string{
		strconv.FormatInt(r.BelongsToId, 10),
		strconv.FormatInt(r.Id, 10),
		strconv.FormatInt(r.ParentId, 10),
		r.Placetype,
		r.Label,
	}

	return out, nil
}

func main() {

	var belongs_to multi.MultiInt64
	flag.Var(&belongs_to, "belongs-to", "...")

	var include_placetype multi.MultiString
	flag.Var(&include_placetype, "include-placetype", "...")

	var exclude_placetype multi.MultiString
	flag.Var(&exclude_placetype, "exclude-placetype", "...")

	mode := flag.String("mode", "repo", "...")

	as_json := flag.Bool("json", false, "...")
	as_markdown := flag.Bool("markdown", false, "...")
	as_ids := flag.Bool("ids", false, "...")

	csv_header := flag.Bool("csv-header", false, "...")
	sort_rs := flag.Bool("sort", false, "...")

	flag.Parse()

	ctx := context.Background()

	rs, err := NewBelongsToResultSet()

	if err != nil {
		log.Fatal(err)
	}

	cb := func(r *BelongsToResult) error {
		return rs.AddResult(r)
	}

	// we should make this a canned TravelFunc once we figure out
	// what the method signature looks like... (20180314/thisisaaronland)

	filter_cb := func(ctx context.Context, f []byte, belongsto_id int64) error {

		pt, err := properties.Placetype(f)

		if err != nil {
			return fmt.Errorf("Faild to derive placetype, %w", err)
		}

		if len(include_placetype) > 0 {

			if !include_placetype.Contains(pt) {
				return nil
			}
		}

		if len(exclude_placetype) > 0 {

			if exclude_placetype.Contains(pt) {
				return nil
			}
		}

		id, err := properties.Id(f)

		if err != nil {
			return fmt.Errorf("Faild to derive ID, %w", err)
		}

		name, err := properties.Name(f)

		if err != nil {
			return fmt.Errorf("Faild to derive name, %w", err)
		}

		parent_id, err := properties.ParentId(f)

		if err != nil {
			return fmt.Errorf("Faild to derive parent ID, %w", err)
		}

		r := BelongsToResult{
			BelongsToId: belongsto_id,
			Id:          id,
			ParentId:    parent_id,
			Placetype:   pt,
			Label:       name, // whosonfirst.LabelOrDerived(f),
		}

		return cb(&r)
	}

	t, err := traveler.NewDefaultBelongsToTraveler()
	t.IteratorURI = *mode
	t.BelongsTo = belongs_to
	t.Callback = filter_cb

	paths := flag.Args()
	err = t.Travel(ctx, paths...)

	if err != nil {
		log.Fatal(err)
	}

	if *sort_rs {
		rs.Sort()
	}

	if *as_ids {

		for _, r := range rs.Results() {
			fmt.Println(r.Id)
		}

	} else if *as_json {

		err := rs.AsJSON(os.Stdout)

		if err != nil {
			log.Fatal(err)
		}

	} else if *as_markdown {

		err := rs.AsMarkdown(os.Stdout)

		if err != nil {
			log.Fatal(err)
		}

	} else {

		err := rs.AsCSV(os.Stdout, *csv_header)

		if err != nil {
			log.Fatal(err)
		}

	}
}
