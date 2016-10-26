package hash

import (
    "crypto/sha1"
	"fmt"
)

// git is an implementation of hash.Hash for Git blob checksums. It
// corresponds to the "git" hash type for the New function.
type git struct {
	data []byte
}

func (this *git) Write(data []byte) (int, error) {
    this.data = append(this.data, data...)
    return len(data), nil
}
func (this *git) Sum(data []byte) []byte {
    hash := sha1.Sum(append([]byte(fmt.Sprintf("blob %d\x00", len(this.data))), this.data...))
    return append(data, hash[:]...)
}
func (this *git) Reset()         { this.data = nil }
func (this *git) Size() int      { return 20 }
func (this *git) BlockSize() int { return 64 }
