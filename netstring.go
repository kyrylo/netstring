package netstring

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
)

// We use Semantic Versioning v2.0.0
// More information: http://semver.org/
const Version = "1.0.0"

const (
	// Prefix that comes after total string size
	prefixCh = ':'

	// Suffix that denotes end of netstring
	suffixCh = ','
)

func Parse(r *bufio.Reader) ([]byte, error) {
	strLen, err := parseLen(r)
	if err != nil {
		return []byte{}, err
	}

	b, err := parseStr(r, strLen)
	if err != nil {
		return []byte{}, err
	}

	if err = stripSuffix(r); err != nil {
		return []byte{}, err
	}

	return b, nil
}

func parseLen(r *bufio.Reader) (int, error) {
	var b byte
	var err error
	var strLen int
	var bytesRead int

	for {
		b, err = r.ReadByte()
		if err != nil {
			return 0, err
		}
		bytesRead++

		if b == '0' && bytesRead == 1 {
			if !assertPrefixAhead(r) {
				return 0, errors.New("leading zeros at the front of length are prohibited")
			}
		} else if b == prefixCh {
			break
		} else if b < '0' || b > '9' {
			return 0, fmt.Errorf(
				"length number %d is not in the range of 0-9",
				int(b-'0'),
			)
		}

		strLen = strLen*10 + int(b-'0')
	}

	return strLen, err
}

func assertPrefixAhead(r *bufio.Reader) bool {
	peek, err := r.Peek(1)
	if err != nil {
		return false
	}

	return peek[0] == prefixCh
}

func parseStr(r *bufio.Reader, len int) ([]byte, error) {
	b, err := io.ReadAll(io.LimitReader(r, int64(len)))
	if err != nil {
		return []byte{}, err
	}

	return b, nil
}

func stripSuffix(r *bufio.Reader) error {
	suffix, err := r.ReadByte()
	if err != nil {
		return err
	}
	if suffix != suffixCh {
		return fmt.Errorf(
			"unexpected suffix %c, wanted %c",
			suffix, suffixCh,
		)
	}

	return nil
}

func Pack(str []byte) ([]byte, error) {
	var buf bytes.Buffer

	buf.Write([]byte(strconv.FormatInt(int64(len(str)), 10)))
	buf.WriteByte(prefixCh)
	buf.Write([]byte(str))
	buf.WriteByte(suffixCh)

	return buf.Bytes(), nil
}
