package main

import (
	"flag"
	"fmt"
	"github.com/unixpickle/linkdups"
	"os"
	pathlib "path"
)

func main() {
	hardlink := flag.Bool("hard", false, "Use hard links")
	follow := flag.Bool("follow", false, "Follow symlinks when scanning")
	flag.Parse()
	
	// Get the path argument
	if len(flag.Args()) != 1 {
		fmt.Fprintln(os.Stderr, "Usage: linkdups [flags] <directory path>")
		os.Exit(1)
	}
	absPath := flag.Arg(0)
	if !pathlib.IsAbs(absPath) {
		wd, _ := os.Getwd()
		absPath = pathlib.Clean(pathlib.Join(wd, absPath))
	}
	
	// Generate the sums
	sums := linkdups.NewSumsSHA256()
	sums.FollowLinks = *follow
	linker := linkdups.Linker{!*hardlink}
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
