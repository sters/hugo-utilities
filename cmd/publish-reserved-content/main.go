package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/morikuni/failure"
	"github.com/sters/hugo-utilities/publish"
)

func abs(p string) string {
	f, _ := filepath.Abs(p)
	return f
}

func dirwalk(dir string) ([]string, error) {
	dir = abs(dir)
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, failure.Wrap(err)
	}

	paths := make([]string, 0, len(files))
	for _, file := range files {
		if file.IsDir() {
			child, err := dirwalk(filepath.Join(dir, file.Name()))
			if err != nil {
				return nil, failure.Wrap(err)
			}
			paths = append(paths, child...)
			continue
		}

		p, err := filepath.Abs(filepath.Join(dir, file.Name()))
		if err != nil {
			fmt.Fprintf(os.Stderr, "%+v", err)
			continue
		}
		paths = append(paths, p)
	}

	return paths, nil
}

func main() {
	var (
		reservedKey        string
		draftKey           string
		basePath           string
		imageReplaceFormat string
	)
	flag.StringVar(&reservedKey, "reservedKey", "reserved", "hugo content's reservation bool key")
	flag.StringVar(&draftKey, "draftKey", "draft", "hugo content's draft bool key")
	flag.StringVar(&basePath, "basePath", "", "hugo content's root directory")
	flag.StringVar(&imageReplaceFormat, "imageReplaceFormat", "![]($1)", "how to replace image url")
	flag.Parse()

	if reservedKey == "" {
		fmt.Fprintf(os.Stderr, "reservedKey is required.\n")
		os.Exit(1)
	}
	if draftKey == "" {
		fmt.Fprintf(os.Stderr, "draftKey is required.\n")
		os.Exit(1)
	}
	if basePath == "" {
		fmt.Fprintf(os.Stderr, "basePath is required.\n")
		os.Exit(1)
	}
	if imageReplaceFormat == "" {
		fmt.Fprintf(os.Stderr, "imageReplaceFormat is required.\n")
		os.Exit(1)
	}

	dirs, err := dirwalk(basePath)
	if err != nil {
		log.Fatal(err)
	}

	p := publish.New(reservedKey, draftKey, imageReplaceFormat)
	for _, filepath := range dirs {
		err := p.CheckReservedAndPublish(filepath)
		c, ok := failure.CodeOf(err)
		if !ok {
			// = no error
			fmt.Fprintf(os.Stdout, "%s is published.\n", filepath)
			continue
		}

		switch c {
		case publish.ErrContentIsReservedButNotDraft:
			fmt.Fprintf(os.Stderr, "%s is reserved but not draft.\n", filepath)
		case publish.ErrFileContentMismatch:
			fmt.Fprintf(os.Stderr, "%s is maybe breaking content.\n", filepath)
		case publish.ErrContentIsNotTheTimeYet:
			fmt.Fprintf(os.Stderr, "%s is still waiting.\n", filepath)
		}
	}
}
