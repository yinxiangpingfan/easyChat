package tools

import (
	"crypto/rand"
	"fmt"
	"io"
)

func GenerateSalt() string {
	b := make([]byte, 16)
	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%x", b)
}
