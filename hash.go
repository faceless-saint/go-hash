/* hash is a convenience wrapper around many common hashing operations
 * and algorithms. Supported operations include computing and comparing
 * checksums and digests, including git-style checksums.
 */
package hash

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash"
	"math"
	"strings"
)

// Hash exposes the underlying hash.Hash definition.
type Hash hash.Hash

// Default is a shortcut to the defualt SHA256 hashing algorithm.
var Default hash.Hash = sha256.New()

// New returns a new hash implementing the chosen algorithm.
// Supported algorithms: "sha512", "sha256", "sha1", "md5", "git"
func New(h string) (hash.Hash, error) {
	switch h {
	case "sha512":
		return sha512.New(), nil
	case "sha256":
		return sha256.New(), nil
	case "sha1":
		return sha1.New(), nil
	case "md5":
		return md5.New(), nil
	case "git":
		return &git{kind: "blob"}, nil
	case strings.Contains(g, "git-"):
		return &git{kind: strings.Replace(h, "git-", "", 1)}, nil
	case "":
		return Default, nil
	default:
		return nil, fmt.Errorf("error: unsupported hash type %s", h)
	}
}

// ByteCountToString returns a human readable representation of the
// given raw byte count using SI units to compactly display the size.
func ByteCountToString(bytes uint64) string {
	var unit uint64
	unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	exp := uint64(math.Log2(float64(bytes)) / math.Log2(float64(unit)))
	char := string("kMGTPE"[exp-1])
	return fmt.Sprintf("%7.2f %sB",
		float64(bytes)/math.Pow(float64(unit),
			float64(exp)), char)
}
