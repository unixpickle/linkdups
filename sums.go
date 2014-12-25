package linkdups

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	pathlib "path"
)

// SumFunc is a function which synchronously computes a hash of a file.
type SumFunc func(path string) (string, error)

// Sums is a configuration for finding the sums of files in a directory.
type Sums struct {
	FollowLinks bool
	SumFunc     SumFunc
}

// NewSumsSHA256 creates a new Sums which uses SHA256 to hash files.
func NewSumsSHA256() *Sums {
	return &Sums{SumFunc: SHA256File}
}

// Compute finds the sum of every file in a named file or directory.
// The returned map will have hashes as keys and slices of paths as values.
func (s *Sums) Compute(name string) (map[string][]string, error) {
	if isDir, err := s.IsDirectory(name); err != nil {
		return nil, err
	} else if !isDir {
		if sum, err := s.SumFunc(name); err != nil {
			return nil, err
		} else {
			return map[string][]string{sum: []string{name}}, nil
		}
	} else {
		return s.listDirectory(name)
	}
}

// IsDirectory returns whether a named path is counted as a directory.
// If name is a symbolic link to a directory, s.FollowLinks is returned.
func (s *Sums) IsDirectory(name string) (bool, error) {
	info, err := os.Stat(name)
	if err != nil {
		return false, err
	}
	if info.IsDir() {
		// Make sure the file isn't a link unless we are allowed to follow
		// links.
		return ((info.Mode()&os.ModeSymlink) == 0 || s.FollowLinks), nil
	}
	return false, nil
}

func (s *Sums) listDirectory(name string) (map[string][]string, error) {
	// List the directory
	dir, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer dir.Close()
	res := map[string][]string{}

	for {
		// Get the next file
		nextName, err := dir.Readdirnames(1)
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		// Read the hashes in the file
		fullName := pathlib.Join(name, nextName[0])
		subRes, err := s.Compute(fullName)
		if err != nil {
			return nil, err
		}

		// Add the hashes to the result
		for key, val := range subRes {
			if list, ok := res[key]; ok {
				res[key] = append(list, val...)
			} else {
				res[key] = val
			}
		}
	}

	return res, nil
}

// SHA256File computes the SHA256 hash of a file located at a given path.
func SHA256File(path string) (string, error) {
	hasher := sha256.New()
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(hasher.Sum([]byte{})), nil
}
