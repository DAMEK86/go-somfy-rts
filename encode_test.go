package somfy

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncodeByte(t *testing.T) {
	testdata := []struct {
		in                byte
		expectedHighPulse int
	}{
		{0x1, 0},
		{0x2, 1},
		{0x4, 2},
		{0x8, 3},
		{0x10, 4},
		{0x20, 5},
		{0x40, 6},
		{0x80, 7},
	}

	e := &encoder{}
	for _, tt := range testdata {
		t.Run(fmt.Sprintf("bit 0x%x", tt.in), func(t *testing.T) {
			encodeByte := e.encodeByte(tt.in)
			assert.True(t, encodeByte[tt.expectedHighPulse].IsHigh)
		})
	}
}
