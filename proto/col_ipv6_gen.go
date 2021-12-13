// Code generated by ./cmd/ch-gen-int, DO NOT EDIT.

package proto

import (
	"github.com/go-faster/errors"
)

// ColIPv6 represents IPv6 column.
type ColIPv6 []IPv6

// Compile-time assertions for ColIPv6.
var (
	_ Input  = ColIPv6{}
	_ Result = (*ColIPv6)(nil)
)

// Type returns ColumnType of IPv6.
func (ColIPv6) Type() ColumnType {
	return ColumnTypeIPv6
}

// Rows returns count of rows in column.
func (c ColIPv6) Rows() int {
	return len(c)
}

// Reset resets data in row, preserving capacity for efficiency.
func (c *ColIPv6) Reset() {
	*c = (*c)[:0]
}

// NewArrIPv6 returns new Array(IPv6).
func NewArrIPv6() *ColArr {
	return &ColArr{
		Data: new(ColIPv6),
	}
}

// AppendIPv6 appends slice of IPv6 to Array(IPv6).
func (c *ColArr) AppendIPv6(data []IPv6) {
	d := c.Data.(*ColIPv6)
	*d = append(*d, data...)
	c.Offsets = append(c.Offsets, uint64(len(*d)))
}

// EncodeColumn encodes IPv6 rows to *Buffer.
func (c ColIPv6) EncodeColumn(b *Buffer) {
	const size = 128 / 8
	offset := len(b.Buf)
	b.Buf = append(b.Buf, make([]byte, size*len(c))...)
	for _, v := range c {
		binPutIPv6(
			b.Buf[offset:offset+size],
			v,
		)
		offset += size
	}
}

// DecodeColumn decodes IPv6 rows from *Reader.
func (c *ColIPv6) DecodeColumn(r *Reader, rows int) error {
	const size = 128 / 8
	data, err := r.ReadRaw(rows * size)
	if err != nil {
		return errors.Wrap(err, "read")
	}
	v := *c
	for i := 0; i < len(data); i += size {
		v = append(v,
			binIPv6(data[i:i+size]),
		)
	}
	*c = v
	return nil
}
