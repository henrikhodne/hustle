package hustle

import (
	"crypto/sha1"
	"fmt"
	"io"
)

func sha1Sum(s string) string {
	h := sha1.New()
	io.WriteString(h, s)
	return fmt.Sprintf("%x", h.Sum(nil))
}
