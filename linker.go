package linkdups

import "os"

type Linker struct {
	Symlinks bool
}

func (l *Linker) LinkDuplicates(listing map[string][]string) error {
	for _, files := range listing {
		if len(files) == 1 {
			continue
		}
		// Link the rest of the files to the original
		orig := files[0]
		for _, path := range files[1:] {
			if err := os.Remove(path); err != nil {
				return err
			}
			if err := l.Link(orig, path); err != nil {
				return err
			}
		}
	}
	return nil
}

func (l *Linker) Link(source, dest string) error {
	if !l.Symlinks {
		return os.Link(source, dest)
	} else {
		return os.Symlink(source, dest)
	}
}
