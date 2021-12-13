// Code generated by ./cmd/ch-gen-int, DO NOT EDIT.

package proto

import (
	"github.com/go-faster/errors"
	"math"
)

// ColFloat32 represents Float32 column.
type ColFloat32 []float32

// Compile-time assertions for ColFloat32.
var (
	_ Input  = ColFloat32{}
	_ Result = (*ColFloat32)(nil)
)

// Type returns ColumnType of Float32.
func (ColFloat32) Type() ColumnType {
	return ColumnTypeFloat32
}

// Rows returns count of rows in column.
func (c ColFloat32) Rows() int {
	return len(c)
}

// Reset resets data in row, preserving capacity for efficiency.
func (c *ColFloat32) Reset() {
	*c = (*c)[:0]
}

// NewArrFloat32 returns new Array(Float32).
func NewArrFloat32() *ColArr {
	return &ColArr{
		Data: new(ColFloat32),
	}
}

// AppendFloat32 appends slice of float32 to Array(Float32).
func (c *ColArr) AppendFloat32(data []float32) {
	d := c.Data.(*ColFloat32)
	*d = append(*d, data...)
	c.Offsets = append(c.Offsets, uint64(len(*d)))
}

// EncodeColumn encodes Float32 rows to *Buffer.
func (c ColFloat32) EncodeColumn(b *Buffer) {
	const size = 32 / 8
	offset := len(b.Buf)
	b.Buf = append(b.Buf, make([]byte, size*len(c))...)
	for _, v := range c {
		bin.PutUint32(
			b.Buf[offset:offset+size],
			math.Float32bits(v),
		)
		offset += size
	}
}

// DecodeColumn decodes Float32 rows from *Reader.
func (c *ColFloat32) DecodeColumn(r *Reader, rows int) error {
	const size = 32 / 8
	data, err := r.ReadRaw(rows * size)
	if err != nil {
		return errors.Wrap(err, "read")
	}
	v := *c
	for i := 0; i < len(data); i += size {
		v = append(v,
			math.Float32frombits(bin.Uint32(data[i:i+size])),
		)
	}
	*c = v
	return nil
}
