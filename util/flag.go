package util

import (
	"flag"
	"fmt"
	"os"
)

type BasicFlag struct {
	BaseDir string
}

func ParseFlag() *BasicFlag {
	dir := flag.String("dir", "", "")
	flag.Parse()

	if *dir == "" {
		fmt.Fprintf(
			os.Stderr,
			"required dir argument\n\n",
		)
		flag.Usage()

		os.Exit(1)
	}

	return &BasicFlag{
		BaseDir: *dir,
	}
}
