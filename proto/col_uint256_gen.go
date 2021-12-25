// Code generated by ./cmd/ch-gen-int, DO NOT EDIT.

package proto

import (
	"encoding/binary"
	"github.com/go-faster/errors"
)

// ClickHouse uses LittleEndian.
var _ = binary.LittleEndian

// ColUInt256 represents UInt256 column.
type ColUInt256 []UInt256

// Compile-time assertions for ColUInt256.
var (
	_ ColInput  = ColUInt256{}
	_ ColResult = (*ColUInt256)(nil)
	_ Column    = (*ColUInt256)(nil)
)

// Type returns ColumnType of UInt256.
func (ColUInt256) Type() ColumnType {
	return ColumnTypeUInt256
}

// Rows returns count of rows in column.
func (c ColUInt256) Rows() int {
	return len(c)
}

// Reset resets data in row, preserving capacity for efficiency.
func (c *ColUInt256) Reset() {
	*c = (*c)[:0]
}

// NewArrUInt256 returns new Array(UInt256).
func NewArrUInt256() *ColArr {
	return &ColArr{
		Data: new(ColUInt256),
	}
}

// AppendUInt256 appends slice of UInt256 to Array(UInt256).
func (c *ColArr) AppendUInt256(data []UInt256) {
	d := c.Data.(*ColUInt256)
	*d = append(*d, data...)
	c.Offsets = append(c.Offsets, uint64(len(*d)))
}

// EncodeColumn encodes UInt256 rows to *Buffer.
func (c ColUInt256) EncodeColumn(b *Buffer) {
	const size = 256 / 8
	offset := len(b.Buf)
	b.Buf = append(b.Buf, make([]byte, size*len(c))...)
	for _, v := range c {
		binPutUInt256(
			b.Buf[offset:offset+size],
			v,
		)
		offset += size
	}
}

// DecodeColumn decodes UInt256 rows from *Reader.
func (c *ColUInt256) DecodeColumn(r *Reader, rows int) error {
	const size = 256 / 8
	data, err := r.ReadRaw(rows * size)
	if err != nil {
		return errors.Wrap(err, "read")
	}
	v := *c
	// Move bound check out of loop.
	//
	// See https://github.com/golang/go/issues/30945.
	_ = data[len(data)-size]
	for i := 0; i <= len(data)-size; i += size {
		v = append(v,
			binUInt256(data[i:i+size]),
		)
	}
	*c = v
	return nil
}
