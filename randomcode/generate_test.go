package randomcode

import (
	"fmt"
	"strings"
	"testing"
)

func TestGenerate(t *testing.T) {
	code := Generate()

	parts := strings.Split(code, "-")

	if len(parts) != 3 {
		t.Errorf("Not expected number of parts")
	}

	for _, p := range parts {
		if len(p) == 0 {
			t.Errorf("Invalid part")
		}
	}

	fmt.Println(code)
}
