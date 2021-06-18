package util

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/morikuni/failure"
	hugocontent "github.com/sters/simple-hugo-content-parse"
)

func ReadAllContents(dir string) ([]*hugocontent.MarkdownContent, error) {
	dirs, err := dirwalk(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed: %s", err)
		os.Exit(1)
	}

	contents := make([]*hugocontent.MarkdownContent, 0, len(dirs))
	for _, filepath := range dirs {
		if !strings.HasSuffix(filepath, ".md") || strings.HasSuffix(filepath, "_index.md") {
			continue
		}

		f, err := os.Open(filepath)
		if err != nil {
			return nil, failure.Wrap(err)
		}

		content, err := hugocontent.ParseMarkdownWithYaml(f)
		f.Close()
		if err != nil {
			return nil, failure.Wrap(err)
		}

		contents = append(contents, content)
	}

	return contents, nil
}

func dirwalk(dir string) ([]string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, failure.Wrap(err)
	}

	var paths []string
	for _, file := range files {
		if file.IsDir() {
			childFiles, err := dirwalk(filepath.Join(dir, file.Name()))
			if err != nil {
				return nil, failure.Wrap(err)
			}
			paths = append(paths, childFiles...)
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
