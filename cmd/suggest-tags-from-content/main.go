package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/sters/hugo-utilities/textutil/ngram"
	"github.com/sters/hugo-utilities/util"
)

const (
	defaultThreshold = 10
	defaultN         = 4
)

func main() {
	t := flag.String("target", "", "")
	n := flag.Int("n", defaultN, "")
	ff := flag.String("field", "tags", "")
	ts := flag.Int("threshold", defaultThreshold, "")
	basicFlags := util.ParseFlag()
	flags := struct {
		targetDir string
		field     string
		n         int
		threshold int
	}{
		targetDir: strings.TrimSpace(*t),
		field:     strings.TrimSpace(*ff),
		n:         *n,
		threshold: *ts,
	}

	if flags.targetDir == "" {
		fmt.Fprintf(os.Stderr, "Required target argument\n\n")
		flag.Usage()
		os.Exit(1)
	}

	contents, err := util.ReadAllContents(basicFlags.BaseDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed: %s", err)
		os.Exit(1)
	}

	targetContents, err := util.ReadAllContents(flags.targetDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed: %s", err)
		os.Exit(1)
	}

	type tag string
	type tagCount map[tag]int
	type ngramGroup string

	gramToTags := map[ngramGroup]tagCount{}

	for _, c := range contents {
		contentTags := c.FrontMatter.GetStrings(flags.field)

		for _, group := range ngram.Parse(flags.n, c.Body).Groups {
			g := ngramGroup(group)
			if _, ok := gramToTags[g]; !ok {
				gramToTags[g] = make(tagCount)
			}

			for _, contentTag := range contentTags {
				gramToTags[g][tag(contentTag)]++
			}
		}
	}

	suggestTags := tagCount{}

	for _, c := range targetContents {
		for _, group := range ngram.Parse(flags.n, c.Body).Groups {
			g := ngramGroup(group)
			if suggestTag, ok := gramToTags[g]; ok {
				for t, c := range suggestTag {
					suggestTags[t] += c
				}
			}
		}
	}

	sortedTags := []struct {
		num int
		tag tag
	}{}
	for t, n := range suggestTags {
		sortedTags = append(sortedTags, struct {
			num int
			tag tag
		}{n, t})
	}

	sort.Slice(sortedTags, func(i, j int) bool {
		return sortedTags[i].num > sortedTags[j].num
	})

	for _, tag := range sortedTags {
		if tag.num <= flags.threshold {
			break
		}
		fmt.Printf("%s (%d)\n", tag.tag, tag.num)
	}
}
