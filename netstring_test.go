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
			"Netstring is empty",
			[]byte{},
			io.EOF,
			[]byte{},
		},
		{
			"Netstring length is 0",
			[]byte{0x30, prefixCh, suffixCh},
			nil,
			[]byte{},
		},
		{
			"Netstring length starts with a leading 0",
			[]byte{0x30, 0x38, prefixCh, 0x73, 0x75, 0x6e, 0x73, 0x68, 0x69, 0x6e, 0x65, suffixCh},
			errors.New("leading zeros at the front of length are prohibited"),
			[]byte{},
		},
		{
			"Netstring is 0",
			[]byte{0x30},
			errors.New("leading zeros at the front of length are prohibited"),
			[]byte{},
		},
		{
			"Netstring length consists of 1 digit",
			[]byte{
				0x38, prefixCh, 0x73, 0x75, 0x6e, 0x73, 0x68, 0x69, 0x6e, 0x65, suffixCh,
			},
			nil,
			[]byte("sunshine"),
		},
		{
			"Netstring length consists of 2 digits",
			[]byte{
				0x31, 0x34, prefixCh, 0x70, 0x65, 0x72, 0x66, 0x65, 0x63, 0x74, 0x6c,
				0x69, 0x62, 0x72, 0x61, 0x72, 0x79, suffixCh,
			},
			nil,
			[]byte("perfectlibrary"),
		},
		{
			"Netstring length consists of 3 digits",
			[]byte{
				0x31, 0x30, 0x35, prefixCh, 0x41, 0x20, 0x6e, 0x65, 0x74, 0x73, 0x74,
				0x72, 0x69, 0x6e, 0x67, 0x20, 0x69, 0x73, 0x20, 0x61, 0x20, 0x73, 0x65,
				0x6c, 0x66, 0x2d, 0x64, 0x65, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x69, 0x6e,
				0x67, 0x20, 0x65, 0x6e, 0x63, 0x6f, 0x64, 0x69, 0x6e, 0x67, 0x20, 0x6f,
				0x66, 0x20, 0x61, 0x20, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x2e, 0x20,
				0x4e, 0x65, 0x74, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x73, 0x20, 0x61,
				0x72, 0x65, 0x20, 0x76, 0x65, 0x72, 0x79, 0x20, 0x65, 0x61, 0x73, 0x79,
				0x20, 0x74, 0x6f, 0x20, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65,
				0x20, 0x61, 0x6e, 0x64, 0x20, 0x74, 0x6f, 0x20, 0x70, 0x61, 0x72, 0x73,
				0x65, 0x2e, suffixCh,
			},
			nil,
			[]byte(
				"A netstring is a self-delimiting encoding of a string. Netstrings " +
					"are very easy to generate and to parse.",
			),
		},
		{
			"Netstring consists of digigts",
			[]byte{
				0x32, prefixCh, 0x31, 0x34, suffixCh,
			},
			nil,
			[]byte("14"),
		},
		{
			"Netstring length is shorter than the actual string",
			[]byte{
				0x30, prefixCh, 0x73, 0x75, 0x6e, 0x73, 0x68, 0x69, 0x6e, 0x65, suffixCh,
			},
			errors.New("unexpected suffix s, wanted ,"),
			[]byte{},
		},
		{
			"Netstring length is longer than the actual string",
			[]byte{
				0x31, 0x31, 0x31, prefixCh, 0x73, 0x75, 0x6e, 0x73, 0x68, 0x69, 0x6e, 0x65,
				0x73, 0x75, 0x6e, 0x73, 0x68, 0x69, 0x6e, 0x65, prefixCh,
			},
			io.EOF,
			[]byte{},
		},
		{
			"Netstring is missing prefix :",
			[]byte{
				0x38, 0x73, 0x75, 0x6e, 0x73, 0x68, 0x69, 0x6e, 0x65, suffixCh,
			},
			errors.New("length number 67 is not in the range of 0-9"),
			[]byte{},
		},
		{
			"Netstring is missing suffix ,",
			[]byte{
				0x38, prefixCh, 0x73, 0x75, 0x6e, 0x73, 0x68, 0x69, 0x6e, 0x65,
			},
			io.EOF,
			[]byte{},
		},
		{
			"Netstring ends with other suffix than ,",
			[]byte{
				0x38, prefixCh, 0x73, 0x75, 0x6e, 0x73, 0x68, 0x69, 0x6e, 0x65, prefixCh,
			},
			errors.New("unexpected suffix :, wanted ,"),
			[]byte{},
		},
		{
			"Netstring includes non-digits in the length field",
			[]byte{
				0x6e, 0x65, prefixCh, 0x73, 0x75, 0x6e, 0x73, 0x68, 0x69, 0x6e, 0x65,
				prefixCh,
			},
			errors.New("length number 62 is not in the range of 0-9"),
			[]byte{},
		},
		{
			"Netstring includes both digits and non-digits in the length field",
			[]byte{
				0x38, 0x65, prefixCh, 0x73, 0x75, 0x6e, 0x73, 0x68, 0x69, 0x6e, 0x65,
				prefixCh,
			},
			errors.New("length number 53 is not in the range of 0-9"),
			[]byte{},
		},
	}

	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			netstr, err := Parse(bufio.NewReader(bytes.NewReader(c.stream)))
			if c.err == nil && err != nil {
				t.Fatalf("unexpected error: '%q'", err)
				return
			}
			if c.err != nil && err != nil {
				if c.err.Error() != err.Error() {
					t.Fatalf("expected error %q, got %q", c.err, err)
				}
				return
			}
			if c.err != nil && err == nil {
				t.Errorf("expected error %q, got nothing", c.err)
			}
			if !bytes.Equal(netstr, c.bytes) {
				t.Errorf("bytes: got %c, want %c", netstr, c.bytes)
			}
		})
	}
}

func TestPack(t *testing.T) {
	cases := []struct {
		desc   string
		input  string
		err    error
		result []byte
	}{
		{
			"Packs an empty string",
			"",
			nil,
			[]byte{0x30, prefixCh, suffixCh},
		},
		{
			"Packs a string with 1-digit length into a netstring",
			"sunshine",
			nil,
			[]byte{
				0x38, prefixCh, 0x73, 0x75, 0x6e, 0x73, 0x68, 0x69, 0x6e, 0x65,
				suffixCh,
			},
		},
		{
			"Packs a string with 2-digit length into a netstring",
			"perfectlibrary",
			nil,
			[]byte{
				0x31, 0x34, prefixCh, 0x70, 0x65, 0x72, 0x66, 0x65, 0x63, 0x74, 0x6c,
				0x69, 0x62, 0x72, 0x61, 0x72, 0x79, suffixCh,
			},
		},
		{
			"Packs a string with 3-digit length into a netstring",
			"A netstring is a self-delimiting encoding of a string. Netstrings " +
				"are very easy to generate and to parse.",
			nil,
			[]byte{
				0x31, 0x30, 0x35, prefixCh, 0x41, 0x20, 0x6e, 0x65, 0x74, 0x73, 0x74,
				0x72, 0x69, 0x6e, 0x67, 0x20, 0x69, 0x73, 0x20, 0x61, 0x20, 0x73, 0x65,
				0x6c, 0x66, 0x2d, 0x64, 0x65, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x69, 0x6e,
				0x67, 0x20, 0x65, 0x6e, 0x63, 0x6f, 0x64, 0x69, 0x6e, 0x67, 0x20, 0x6f,
				0x66, 0x20, 0x61, 0x20, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x2e, 0x20,
				0x4e, 0x65, 0x74, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x73, 0x20, 0x61,
				0x72, 0x65, 0x20, 0x76, 0x65, 0x72, 0x79, 0x20, 0x65, 0x61, 0x73, 0x79,
				0x20, 0x74, 0x6f, 0x20, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65,
				0x20, 0x61, 0x6e, 0x64, 0x20, 0x74, 0x6f, 0x20, 0x70, 0x61, 0x72, 0x73,
				0x65, 0x2e, suffixCh,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			b, err := Pack([]byte(c.input))
			if c.err == nil && err != nil {
				t.Fatalf("unexpected error: '%q'", err)
				return
			}
			if c.err != nil && err != nil {
				if c.err.Error() != err.Error() {
					t.Fatalf("expected error %q, got %q", c.err, err)
				}
				return
			}
			if c.err != nil && err == nil {
				t.Errorf("expected error %q, got nothing", c.err)
			}
			if !bytes.Equal(b, c.result) {
				t.Errorf("bytes: got %c, want %c", b, c.result)
			}
		})
	}
}

func TestParseAndPackInteropability(t *testing.T) {
	// 5:hello,
	netstr := []byte{0x35, prefixCh, 0x68, 0x65, 0x6c, 0x6c, 0x6f, suffixCh}

	str, err := Parse(bufio.NewReader(bytes.NewReader(netstr)))
	if err != nil {
		t.Fatalf("unexpected error: '%q'", err)
		return
	}

	packed, err := Pack(str)

	if err != nil {
		t.Fatalf("unexpected error: '%q'", err)
		return
	}

	if !bytes.Equal(netstr, packed) {
		t.Errorf("wanted %c to be the same as %c", packed, netstr)
	}
}

func TestPackAndParseInteropability(t *testing.T) {
	str := []byte("sunshine")

	packed, err := Pack(str)
	if err != nil {
		t.Fatalf("unexpected error: '%q'", err)
		return
	}

	netstr, err := Parse(bufio.NewReader(bytes.NewReader(packed)))
	if err != nil {
		t.Fatalf("unexpected error: '%q'", err)
		return
	}

	if !bytes.Equal(netstr, str) {
		t.Errorf("wanted %c to be the same as %c", netstr, str)
	}
}
