package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/agnivade/levenshtein"
	"github.com/sters/hugo-utilities/content"
	"github.com/sters/hugo-utilities/util"
)

type similarTagRelation struct {
	a        string
	b        string
	distance int
}

type similarTagRelationHash string

func (s *similarTagRelation) Hash() similarTagRelationHash {
	return similarTagRelationHash(s.a + "-" + s.b)
}

func (s *similarTagRelation) HashRev() similarTagRelationHash {
	return similarTagRelationHash(s.b + "-" + s.a)
}

type similarTagRelations map[similarTagRelationHash]*similarTagRelation

func (s similarTagRelations) appendWithoutDuplecation(r *similarTagRelation) {
	if _, ok := s[r.Hash()]; ok {
		return
	}
	if _, ok := s[r.HashRev()]; ok {
		return
	}
	s[r.Hash()] = r
}

func getAllTags(contents []*content.MarkdownContent, field string) []string {
	tags := make([]string, 0, len(contents))
	check := map[string]struct{}{}
	for _, c := range contents {
		for _, t := range c.FrontMatter.GetStrings(field) {
			if _, ok := check[t]; ok {
				continue
			}

			tags = append(tags, t)
			check[t] = struct{}{}
		}
	}

	return tags
}

func main() {
	const defaultThreshold = 0.7

	t := flag.Float64("threshold", defaultThreshold, "")
	ff := flag.String("field", "tags", "")
	basicFlags := util.ParseFlag()
	flags := struct {
		threshold float64
		field     string
	}{
		threshold: *t,
		field:     *ff,
	}

	contents, err := util.ReadAllContents(basicFlags.BaseDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed: %s", err)
		os.Exit(1)
	}

	tags := getAllTags(contents, flags.field)
	r := similarTagRelations{}

	for i, tagA := range tags {
		for j, tagB := range tags {
			if i >= j || tagA == tagB {
				continue
			}

			distance := levenshtein.ComputeDistance(tagA, tagB)
			g := func(s string) int {
				return int(flags.threshold * float64(len([]rune(s))))
			}
			if distance >= g(tagA) || distance >= g(tagB) {
				continue
			}

			r.appendWithoutDuplecation(
				&similarTagRelation{
					a:        tagA,
					b:        tagB,
					distance: distance,
				},
			)
		}
	}

	result := make([]string, 0, len(r))
	for _, n := range r {
		result = append(
			result,
			fmt.Sprintf("%s, %s: distance = %d", n.a, n.b, n.distance),
		)
	}

	sort.Slice(result, func(i, j int) bool { return result[i] < result[j] })

	fmt.Print(strings.Join(result, "\n"))
}
