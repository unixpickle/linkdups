package main

import (
	"flag"
	"fmt"
	"github.com/unixpickle/linkdups"
	"os"
	"path/filepath"
)

func main() {
	hardlink := flag.Bool("hard", false, "Use hard links")
	absolute := flag.Bool("absolute", false, "Use absolute symbolic links")
	follow := flag.Bool("follow", false, "Follow symlinks when scanning")
	flag.Parse()

	// Get the path argument
	if len(flag.Args()) != 1 {
		fmt.Fprintln(os.Stderr, "Usage: linkdups [flags] <directory path>")
		os.Exit(1)
	}
	absPath, err := filepath.Abs(flag.Arg(0))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Generate the sums
	sums := linkdups.NewSumsSHA256()
	sums.FollowLinks = *follow
	linker := linkdups.Linker{!*hardlink, !*absolute}
	files, err := sums.Compute(absPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Create the links
	if err := linker.LinkDuplicates(files); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
