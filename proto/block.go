package proto

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/go-faster/errors"
)

// BlockInfo describes block.
type BlockInfo struct {
	Overflows bool
	BucketNum int
}

func (i BlockInfo) String() string {
	return fmt.Sprintf("overflows: %v, buckets: %d", i.Overflows, i.BucketNum)
}

const endField = 0 // end of field pairs

// fields of BlockInfo.
const (
	blockInfoOverflows = 1
	blockInfoBucketNum = 2
)

// Encode to Buffer.
func (i BlockInfo) Encode(b *Buffer) {
	b.PutUVarInt(blockInfoOverflows)
	b.PutBool(i.Overflows)

	b.PutUVarInt(blockInfoBucketNum)
	b.PutInt32(int32(i.BucketNum))

	b.PutUVarInt(endField)
}

func (i *BlockInfo) Decode(r *Reader) error {
	for {
		f, err := r.UVarInt()
		if err != nil {
			return errors.Wrap(err, "field id")
		}
		switch f {
		case blockInfoOverflows:
			v, err := r.Bool()
			if err != nil {
				return errors.Wrap(err, "overflows")
			}
			i.Overflows = v
		case blockInfoBucketNum:
			v, err := r.Int32()
			if err != nil {
				return errors.Wrap(err, "bucket number")
			}
			i.BucketNum = int(v)
		case endField:
			return nil
		default:
			return errors.Errorf("unknown field %d", f)
		}
	}
}

// Input of query.
type Input []InputColumn

// Into returns INSERT INTO table (c0, c..., cn) VALUES query.
func (i Input) Into(table string) string {
	return fmt.Sprintf("INSERT INTO %s %s VALUES", strconv.QuoteToASCII(table), i.Columns())
}

// Columns returns "(foo, bar, baz)" formatted list of Input column names.
func (i Input) Columns() string {
	var (
		b   strings.Builder
		buf [64]byte
	)

	b.WriteRune('(')
	for idx, v := range i {
		escaped := strconv.AppendQuoteToASCII(buf[:0], v.Name)
		b.Write(escaped)
		if idx != len(i)-1 {
			b.WriteRune(',')
		}
	}
	b.WriteRune(')')

	return b.String()
}

type InputColumn struct {
	Name string
	Data ColInput
}

// ResultColumn can be uses as part of Results or as single Result.
type ResultColumn struct {
	Name string    // Name of column. Inferred if not provided.
	Data ColResult // Data of column, required.
}

// DecodeResult implements Result as "single result" helper.
func (c ResultColumn) DecodeResult(r *Reader, b Block) error {
	v := Results{c}
	return v.DecodeResult(r, b)
}

// AutoResult is ResultColumn with type inference.
func AutoResult(name string) ResultColumn {
	return ResultColumn{
		Name: name,
		Data: &ColAuto{},
	}
}

func (c InputColumn) EncodeStart(buf *Buffer) {
	buf.PutString(c.Name)
	buf.PutString(string(c.Data.Type()))
}

type Block struct {
	Info    BlockInfo
	Columns int
	Rows    int
}

func (b Block) EncodeAware(buf *Buffer, version int) {
	if FeatureBlockInfo.In(version) {
		b.Info.Encode(buf)
	}

	buf.PutInt(b.Columns)
	buf.PutInt(b.Rows)
}

func (b Block) EncodeBlock(buf *Buffer, version int, input []InputColumn) error {
	if FeatureBlockInfo.In(version) {
		b.Info.Encode(buf)
	}
	if err := b.EncodeRawBlock(buf, input); err != nil {
		return errors.Wrap(err, "raw block")
	}
	return nil
}

func (b Block) EncodeRawBlock(buf *Buffer, input []InputColumn) error {
	buf.PutInt(b.Columns)
	buf.PutInt(b.Rows)
	for _, col := range input {
		if r := col.Data.Rows(); r != b.Rows {
			return errors.Errorf("%q has %d rows, expected %d", col.Name, r, b.Rows)
		}
		col.EncodeStart(buf)
		if v, ok := col.Data.(Preparable); ok {
			if err := v.Prepare(); err != nil {
				return errors.Wrap(err, "prepare")
			}
		}
		if col.Data.Rows() == 0 {
			continue
		}
		if v, ok := col.Data.(StateEncoder); ok {
			v.EncodeState(buf)
		}
		col.Data.EncodeColumn(buf)
	}
	return nil
}

const (
	maxColumnsInBlock = 1_000_000
	maxRowsInBLock    = 1_000_000_000
)

func checkRows(n int) error {
	if n < 0 {
		return errors.New("negative")
	}
	if n > maxRowsInBLock {
		return errors.Errorf("%d bigger than maximum %d (possible OOM)", n, maxRowsInBLock)
	}
	return nil
}

func (b *Block) End() bool {
	return b.Columns == 0 && b.Rows == 0
}

func (b *Block) DecodeRawBlock(r *Reader, target Result) error {
	{
		v, err := r.Int()
		if err != nil {
			return errors.Wrap(err, "columns")
		}
		if v > maxColumnsInBlock || v < 0 {
			return errors.Errorf("invalid columns number %d", v)
		}
		b.Columns = v
	}
	{
		v, err := r.Int()
		if err != nil {
			return errors.Wrap(err, "rows")
		}
		if err := checkRows(v); err != nil {
			return errors.Wrap(err, "rows count")
		}
		b.Rows = v
	}
	if b.End() {
		// End of data, special case.
		return nil
	}
	if target == nil && b.Rows > 0 {
		return errors.New("got rows without target")
	}
	if target == nil {
		// Just skipping rows and types.
		for i := 0; i < b.Columns; i++ {
			if _, err := r.Str(); err != nil {
				return errors.Wrapf(err, "column [%d] name", i)
			}
			if _, err := r.Str(); err != nil {
				return errors.Wrapf(err, "column [%d] type", i)
			}
		}
		return nil
	}
	if err := target.DecodeResult(r, *b); err != nil {
		return errors.Wrap(err, "target")
	}

	return nil
}

func (b *Block) DecodeBlock(r *Reader, version int, target Result) error {
	if FeatureBlockInfo.In(version) {
		if err := b.Info.Decode(r); err != nil {
			return errors.Wrap(err, "info")
		}
	}
	if err := b.DecodeRawBlock(r, target); err != nil {
		return err
	}

	return nil
}
