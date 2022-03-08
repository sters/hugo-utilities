package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/sters/hugo-utilities/util"
)

type tagRelationHash string

type tagRelation struct {
	a     string
	b     string
	count int
}

func (s *tagRelation) Hash() tagRelationHash {
	return tagRelationHash(s.a + "-" + s.b)
}

type tagRelations map[tagRelationHash]*tagRelation

func (t tagRelations) append(tr *tagRelation) {
	h := tr.Hash()
	if tt, ok := t[h]; ok {
		tt.count += tr.count
		return
	}

	t[h] = tr
}

func main() {
	t := flag.String("tag", "", "")
	ff := flag.String("field", "tags", "")
	basicFlags := util.ParseFlag()
	flags := struct {
		tag   string
		field string
	}{
		tag:   strings.TrimSpace(*t),
		field: strings.TrimSpace(*ff),
	}

	if flags.tag == "" {
		fmt.Fprintf(os.Stderr, "Required tag argument\n\n")
		flag.Usage()
		os.Exit(1)
	}

	contents, err := util.ReadAllContents(basicFlags.BaseDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed: %s", err)
		os.Exit(1)
	}

	trs := tagRelations{}
	for _, c := range contents {
		tags := c.FrontMatter.GetStrings(flags.field)

		for i, ta := range tags {
			for j, tb := range tags {
				if i == j {
					continue
				}
				if ta != flags.tag {
					continue
				}

				trs.append(&tagRelation{
					a:     ta,
					b:     tb,
					count: 1,
				})
			}
		}
	}

	str := make([]*tagRelation, 0, len(trs))
	for _, tr := range trs {
		str = append(str, tr)
	}

	sort.Slice(str, func(i, j int) bool {
		if str[i].count == str[j].count {
			if str[i].a == str[j].a {
				return str[i].b < str[j].b
			}
			return str[i].a < str[j].a
		}
		return str[i].count > str[j].count
	})

	for _, tr := range str {
		fmt.Printf("(%s, %s) = %d\n", tr.a, tr.b, tr.count)
	}
}
