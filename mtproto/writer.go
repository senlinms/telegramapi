package mtproto

import (
	"bytes"
	"math/big"
)

type Writer struct {
	buf bytes.Buffer
}

func NewWriter() *Writer {
	return &Writer{}
}

func NewWriterCmd(cmd uint32) *Writer {
	w := NewWriter()
	w.WriteCmd(cmd)
	return w
}

func (w *Writer) Clear() {
	w.buf.Truncate(0)
}

func (w *Writer) Bytes() []byte {
	return w.buf.Bytes()
}

func (w *Writer) WriteByte(v byte) {
	_ = w.buf.WriteByte(v)
}

func (w *Writer) WriteUint24(v uint32) {
	var b [3]byte
	b[0] = byte(v)
	b[1] = byte(v >> (8 * 1))
	b[2] = byte(v >> (8 * 2))
	_, _ = w.buf.Write(b[:])
}

func (w *Writer) WriteUint32(v uint32) {
	var b [4]byte
	b[0] = byte(v)
	b[1] = byte(v >> (8 * 1))
	b[2] = byte(v >> (8 * 2))
	b[3] = byte(v >> (8 * 3))
	_, _ = w.buf.Write(b[:])
}

func (w *Writer) WriteUint64(v uint64) {
	var b [8]byte
	b[0] = byte(v)
	b[1] = byte(v >> (8 * 1))
	b[2] = byte(v >> (8 * 2))
	b[3] = byte(v >> (8 * 3))
	b[4] = byte(v >> (8 * 4))
	b[5] = byte(v >> (8 * 5))
	b[6] = byte(v >> (8 * 6))
	b[7] = byte(v >> (8 * 7))
	_, _ = w.buf.Write(b[:])
}

func (w *Writer) WriteCmd(v uint32) {
	w.WriteUint32(v)
}
func (w *Writer) WriteInt(v int) {
	w.WriteUint32(uint32(v))
}

func (w *Writer) Write(v []byte) {
	_, _ = w.buf.Write(v)
}

func (w *Writer) WriteUint128(v []byte) {
	_ = v[15]
	w.Write(v[:16])
}

func (w *Writer) ZeroPad(n int) {
	w.buf.Grow(n)
	for i := 0; i < n; i++ {
		_ = w.buf.WriteByte(0)
	}
}

func (w *Writer) WriteStringLen(v int) int {
	if v < 0 {
		panic("negative len")
	}
	if v < 254 {
		w.WriteByte(byte(v))
		return PaddingOf(1 + v)
	} else {
		w.buf.Grow(4)
		w.WriteByte(254)
		w.WriteUint24(uint32(v))
		return PaddingOf(4 + v)
	}
}

func (w *Writer) WriteString(v []byte) {
	pad := w.WriteStringLen(len(v))
	w.Write(v)
	w.ZeroPad(pad)
}

func (w *Writer) WriteBigInt(v *big.Int) {
	b := v.Bytes()
	if len(b) == 0 {
		b = []byte{0}
	}

	w.WriteString(b)
}

func (w *Writer) WriteVectorLong(v []uint64) {
	w.WriteCmd(IDVectorLong)
	w.WriteInt(len(v))
	for _, el := range v {
		w.WriteUint64(el)
	}
}

func (w *Writer) PaddingTo(bs int) int {
	len := w.buf.Len()
	r := len % bs
	if r == 0 {
		return 0
	} else {
		return bs - r
	}
}

type WritableToWriter interface {
	WriteTo(w *Writer)
}

func BytesOf(ww WritableToWriter) []byte {
	w := NewWriter()
	ww.WriteTo(w)
	return w.Bytes()
}
