# linkdups

**linkdups** is a simple tool for removing duplicate files from a directory.

# Installation

Install the `linkdups` command as follows:

    go install github.com/unixpickle/linkdups/linkdups

# Usage

Simply run `linkdups` on a directory to create symbolic links between duplicate files:

    linkdups /path/to/directory

Optionally, you can make `linkdups` create hard links rather than symbolic links:

    linkdups -hard /path/to/directory

# TODO

 * Use "path/filepath" for computing absolute paths
 * Support relative symlinks
