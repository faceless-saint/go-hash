package hash

import (
	"fmt"
	"hash"
	"io/ioutil"
	"regexp"
	"strings"
)

// Checksum describes a copmuted checksum and the hash algorithm used.
type Checksum struct {
	Value string
	Hash  *hash.Hash
}

// IsValid returns true iff the Checksum is valid. Note that this checks
// the formatting of the checksum string itself, not related to a file.
func (this *Checksum) IsValid() bool {
	if this.Hash == nil {
		return false
	}
	if len(this.Value) != (*this.Hash).Size()*2 {
		return false
	}
	test, _ := regexp.MatchString("^[a-fA-F0-9]+$", this.Value)
	return test
}

// ChecksumFromString attempts to parse a string using the following
// format into a Checksum object - "{hash}:{value}". If the hash prefix
// is omitted then the default algorithm is used.
func ChecksumFromString(str string) (*Checksum, error) {
	raw := strings.SplitN(str, ":", 2)
	if len(raw) < 2 {
		return &Checksum{raw[0], &Default}, nil
	}
	h, err := New(raw[0])
	if err != nil {
		return nil, err
	}
	return &Checksum{raw[1], &h}, nil

}

// ByteChecksum computes the checksum for the given bytes and hash
func ByteChecksum(data []byte, h hash.Hash) *Checksum {
	if h == nil {
		h = Default
	}
	h.Write(data)
	checksum := fmt.Sprintf("%x", h.Sum(nil))
	h.Reset()
	return &Checksum{checksum, &h}
}

// StringChecksum computes the checksum for the given string and hash
func StringChecksum(str string, h hash.Hash) *Checksum {
	return ByteChecksum([]byte(str), h)
}

// FileChecksum computes the checksum for the given file
func FileChecksum(file string, h hash.Hash) (*Checksum, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return ByteChecksum(data, h), nil
}

// VerifyBytes compares the checksum of the given byte array against the
// reference checksum, returning true iff they match.
func VerifyBytes(data []byte, checksum Checksum) bool {
	return ByteChecksum(data, *checksum.Hash).Value == checksum.Value
}

// VerifyBytes compares the checksum of the given string against the
// reference checksum, returning true iff they match.
func VerifyString(data string, checksum Checksum) bool {
	return StringChecksum(data, *checksum.Hash).Value == checksum.Value
}

// VerifyFile compares the checksum of the given file against the
// reference checksum, returning true iff they match.
func VerifyFile(file string, checksum Checksum) (bool, error) {
	sum, err := FileChecksum(file, *checksum.Hash)
	if err != nil {
		return false, err
	}
	return sum.Value == checksum.Value, nil
}

// Digest returns a short digest consisting of the first 'n' characters
// of the given string's checksum. The library default hash is used.
func Digest(str string, n int) string {
	return StringChecksum(str, Default).Value[:n]
}
