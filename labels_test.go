package dns

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnpackLabel(t *testing.T) {
	cases := []struct {
		Input     string
		Expected  string
		Offset    int
		IsPointer bool
		Err       error
	}{
		{Input: "\x00", Err: ErrLabelEmpty},
		{Input: "\x40", Err: ErrLabelTooLong},
		{Input: "\x01", Err: io.ErrShortBuffer},
		{Input: "\x3f", Err: io.ErrShortBuffer},

		{Input: "\x01.", Err: ErrLabelInvalid},
		{Input: "\x02..", Err: ErrLabelInvalid},
		{Input: "\x07.00000A", Err: ErrLabelInvalid},
		{Input: "\x03123", Err: ErrLabelInvalid},
		{Input: "\x0asome.email", Err: ErrLabelInvalid},
		{Input: "\x01-", Err: ErrLabelInvalid},

		{Input: "\x01a", Expected: "a"},
		{Input: "\x031-a", Expected: "1-a"},

		// Pointer tests
		{Input: "\xc0\x02", Err: ErrLabelPointerIllegal, IsPointer: true},
		{Input: "\x01\xc0\x00", Offset: 1, Err: ErrLabelPointerIllegal, IsPointer: true},
		{Input: "\x06domain\xc0\x00", Expected: "domain", Offset: 7, IsPointer: true},
	}

	for _, c := range cases {
		t.Logf("Label unpacking input: %q\n", c.Input)

		var label string
		var n int
		var err error

		b := []byte(c.Input)
		if c.IsPointer {
			label, n, err = unpackLabelPointer(b, c.Offset)
		} else {
			label, n, err = unpackLabel(b, c.Offset)
		}

		assert.Equal(t, c.Err, err)
		if err == nil {
			if c.IsPointer {
				assert.Equal(t, 2, n)
			} else {
				assert.Equal(t, len(c.Input), n)
			}
		}
		assert.Equal(t, c.Expected, label)
	}
}
