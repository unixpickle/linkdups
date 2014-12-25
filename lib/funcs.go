package linkdups

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
)

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
