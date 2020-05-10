package netstring

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
)

// We use Semantic Versioning v2.0.0
// More information: http://semver.org/
const Version = "1.0.0"

const (
	// A netstring carries size information. It is encoded as 4-byte uint32
	// number
	byteSize = 4

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

	if err = stripPrefix(r); err != nil {
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

func Pack(data []byte) []byte {
	var b bytes.Buffer

	strLen := make([]byte, 4)
	binary.LittleEndian.PutUint32(strLen, uint32(len(data)))

	b.Write(strLen)
	b.WriteByte(prefixCh)
	b.Write(data)
	b.WriteByte(suffixCh)

	return b.Bytes()
}

func parseLen(r *bufio.Reader) (int, error) {
	sizeBuf := make([]byte, byteSize)
	if _, err := io.ReadFull(r, sizeBuf); err != nil {
		return 0, err
	}

	len := binary.LittleEndian.Uint32(sizeBuf[0:])
	return int(len), nil
}

func stripPrefix(r *bufio.Reader) error {
	prefix, err := r.ReadByte()
	if err != nil {
		return err
	}
	if prefix != prefixCh {
		return fmt.Errorf(
			"got unexpected netstring prefix %c, wanted %c",
			prefix, prefixCh,
		)
	}

	return nil
}

func parseStr(r *bufio.Reader, len int) ([]byte, error) {
	b, err := ioutil.ReadAll(io.LimitReader(r, int64(len)))
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
			"got unexpected netstring suffix %c, wanted %c",
			suffix, suffixCh,
		)
	}

	return nil
}
