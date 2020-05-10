package netstring

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"testing"
)

func TestParse(t *testing.T) {
	cases := []struct {
		desc   string
		stream []byte
		err    error
		bytes  []byte
	}{
		{
			"Netstring length is less than 4 bytes",
			[]byte{0x32, prefixCh, 0x61, 0x62, 0x63, 0x2c, suffixCh},
			errors.New("got unexpected netstring prefix c, wanted :"),
			[]byte{},
		},
		{
			"Netstring length is more than 4 bytes",
			[]byte{
				0x32, 0x32, 0x32, 0x32, 0x01, prefixCh, 0x61,
				0x62, 0x63, 0x2c, suffixCh,
			},
			errors.New("got unexpected netstring prefix \x01, wanted :"),
			[]byte{},
		},
		{
			"Correct netstring (4-byte length, prefix : and suffix ,)",
			[]byte{
				0x05, 0x00, 0x00, 0x00, prefixCh, 0x68, 0x65,
				0x6c, 0x6c, 0x6f, suffixCh,
			},
			nil,
			[]byte{0x68, 0x65, 0x6c, 0x6c, 0x6f},
		},
		{
			"Correct netstring (4-byte length, prefix : and suffix ,)",
			[]byte{
				0x05, 0x00, 0x00, 0x00, prefixCh, 0x68, 0x65,
				0x6c, 0x6c, 0x6f, suffixCh,
			},
			nil,
			[]byte{0x68, 0x65, 0x6c, 0x6c, 0x6f},
		},
		{
			"Netstring without suffix ,",
			[]byte{
				0x05, 0x00, 0x00, 0x00, prefixCh, 0x68, 0x65,
				0x6c, 0x6c, 0x6f,
			},
			io.EOF,
			[]byte{},
		},
		{
			"Netstring longer than specified length",
			[]byte{
				0x05, 0x00, 0x00, 0x00, prefixCh, 0x68, 0x65,
				0x6c, 0x6c, 0x6f, 0x6f, 0x6f, suffixCh,
			},
			errors.New("got unexpected netstring suffix o, wanted ,"),
			[]byte{},
		},
		{
			"Netstring shorter than specified length",
			[]byte{
				0x05, 0x00, 0x00, 0x00, prefixCh, 0x68, 0x65,
				0x6c, 0x6c, suffixCh,
			},
			io.EOF,
			[]byte{},
		},
	}

	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			r := bytes.NewReader(c.stream)
			buf := bufio.NewReader(r)

			b, err := Parse(buf)
			if c.err == nil && err != nil {
				t.Fatalf("unexpected error: '%q'", err)
				return
			}
			if c.err != nil && err != nil {
				if c.err.Error() != err.Error() {
					t.Fatalf("expected error %q, want %q", err, c.err)
				}
				return
			}
			if c.err != nil && err == nil {
				t.Errorf("expected error %q, got nothing", c.err)
			}
			if !bytes.Equal(b, c.bytes) {
				t.Errorf("bytes: got %c, want %c", b, c.bytes)
			}
		})
	}
}

func TestPack(t *testing.T) {
	cases := []struct {
		desc   string
		input  []byte
		result []byte
	}{
		{
			`Packs string "sunshine" into a netstring`,
			[]byte("sunshine"),
			[]byte{
				0x08, 0x00, 0x00, 0x00, prefixCh, 0x73, 0x75,
				0x6e, 0x73, 0x68, 0x69, 0x6e, 0x65, suffixCh,
			},
		},
		{
			`Packs string "саншайн" into a netstring`,
			[]byte("саншайн"),
			[]byte{
				0x0e, 0x00, 0x00, 0x00, prefixCh, 0xd1, 0x81,
				0xd0, 0xb0, 0xd0, 0xbd, 0xd1, 0x88, 0xd0, 0xb0,
				0xd0, 0xb9, 0xd0, 0xbd, suffixCh,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			b := Pack(c.input)
			if !bytes.Equal(b, c.result) {
				t.Errorf("bytes: got %c, want %c", b, c.result)
			}
		})
	}
}

func TestParseAndPackInteropability(t *testing.T) {
	b := []byte{
		0x05, 0x00, 0x00, 0x00, prefixCh, 0x68, 0x65,
		0x6c, 0x6c, 0x6f, suffixCh,
	}
	buf := bufio.NewReader(bytes.NewReader(b))

	str, err := Parse(buf)
	if err != nil {
		t.Fatalf("unexpected error: '%q'", err)
		return
	}

	parsed := Pack(str)
	if !bytes.Equal(b, parsed) {
		t.Errorf("wanted %c to be the same as %c", b, parsed)
	}
}

func TestPackAndParseInteropability(t *testing.T) {
	// sunshine
	b := []byte{0x73, 0x75, 0x6e, 0x73, 0x68, 0x69, 0x6e, 0x65}

	packed := Pack(b)
	str, err := Parse(bufio.NewReader(bytes.NewReader(packed)))
	if err != nil {
		t.Fatalf("unexpected error: '%q'", err)
		return
	}

	if !bytes.Equal(b, str) {
		t.Errorf("wanted %c to be the same as %c", b, str)
	}
}
