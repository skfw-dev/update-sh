package tests

import (
	"fmt"
	"testing"
	"update-sh/nextgen/cores"
)

func TestHex(t *testing.T) {
	fmt.Printf("Hex: %s\n", cores.ToHex([]byte("Hello, World!")))
}
