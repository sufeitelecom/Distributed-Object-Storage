package tools

import (
	"io"
	"crypto/sha256"
	"encoding/base64"
)

func CalculateHash(r io.Reader) string {
	h := sha256.New()
	io.Copy(h, r)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
