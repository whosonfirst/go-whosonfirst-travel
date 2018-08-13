package main

import (
	"flag"
	"github.com/whosonfirst/go-whosonfirst-cli/flags"
	"github.com/whosonfirst/go-whosonfirst-travel/traveler"
	"log"
)

func main() {

	var belongs_to flags.MultiInt64
	flag.Var(&belongs_to, "belongs-to", "...")

	mode := flag.String("mode", "repo", "...")

	flag.Parse()

	t, err := traveler.NewDefaultBelongsToTraveler()
	t.Mode = *mode
	t.BelongsTo = belongs_to

	paths := flag.Args()
	err = t.Travel(paths...)

	if err != nil {
		log.Fatal(err)
	}
}
