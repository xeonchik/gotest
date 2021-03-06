package engine

// Code generated by colf(1); DO NOT EDIT.
// The compiler used schema file colf\datarecord.colf.

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

var intconv = binary.BigEndian

// Colfer configuration attributes
var (
	// ColferSizeMax is the upper limit for serial byte sizes.
	ColferSizeMax = 16 * 1024 * 1024
	// ColferListMax is the upper limit for the number of elements in a list.
	ColferListMax = 64 * 1024
)

// ColferMax signals an upper limit breach.
type ColferMax string

// Error honors the error interface.
func (m ColferMax) Error() string { return string(m) }

// ColferError signals a data mismatch as as a byte index.
type ColferError int

// Error honors the error interface.
func (i ColferError) Error() string {
	return fmt.Sprintf("colfer: unknown header at byte %d", i)
}

// ColferTail signals data continuation as a byte index.
type ColferTail int

// Error honors the error interface.
func (i ColferTail) Error() string {
	return fmt.Sprintf("colfer: data continuation at byte %d", i)
}

// Course is the grounds where the game of golf is played.
type DataRecord struct {
	ID uint64

	Data string

	Sectors []*Sector

	Cities []*City

	Active bool

	Sort float32
}

// MarshalTo encodes o as Colfer into buf and returns the number of bytes written.
// If the buffer is too small, MarshalTo will panic.
// All nil entries in o.Sectors will be replaced with a new value.
// All nil entries in o.Cities will be replaced with a new value.
func (o *DataRecord) MarshalTo(buf []byte) int {
	var i int

	if x := o.ID; x >= 1<<49 {
		buf[i] = 0 | 0x80
		intconv.PutUint64(buf[i+1:], x)
		i += 9
	} else if x != 0 {
		buf[i] = 0
		i++
		for x >= 0x80 {
			buf[i] = byte(x | 0x80)
			x >>= 7
			i++
		}
		buf[i] = byte(x)
		i++
	}

	if l := len(o.Data); l != 0 {
		buf[i] = 1
		i++
		x := uint(l)
		for x >= 0x80 {
			buf[i] = byte(x | 0x80)
			x >>= 7
			i++
		}
		buf[i] = byte(x)
		i++
		i += copy(buf[i:], o.Data)
	}

	if l := len(o.Sectors); l != 0 {
		buf[i] = 2
		i++
		x := uint(l)
		for x >= 0x80 {
			buf[i] = byte(x | 0x80)
			x >>= 7
			i++
		}
		buf[i] = byte(x)
		i++
		for vi, v := range o.Sectors {
			if v == nil {
				v = new(Sector)
				o.Sectors[vi] = v
			}
			i += v.MarshalTo(buf[i:])
		}
	}

	if l := len(o.Cities); l != 0 {
		buf[i] = 3
		i++
		x := uint(l)
		for x >= 0x80 {
			buf[i] = byte(x | 0x80)
			x >>= 7
			i++
		}
		buf[i] = byte(x)
		i++
		for vi, v := range o.Cities {
			if v == nil {
				v = new(City)
				o.Cities[vi] = v
			}
			i += v.MarshalTo(buf[i:])
		}
	}

	if o.Active {
		buf[i] = 4
		i++
	}

	if v := o.Sort; v != 0 {
		buf[i] = 5
		intconv.PutUint32(buf[i+1:], math.Float32bits(v))
		i += 5
	}

	buf[i] = 0x7f
	i++
	return i
}

// MarshalLen returns the Colfer serial byte size.
// The error return option is engine.ColferMax.
func (o *DataRecord) MarshalLen() (int, error) {
	l := 1

	if x := o.ID; x >= 1<<49 {
		l += 9
	} else if x != 0 {
		for l += 2; x >= 0x80; l++ {
			x >>= 7
		}
	}

	if x := len(o.Data); x != 0 {
		if x > ColferSizeMax {
			return 0, ColferMax(fmt.Sprintf("colfer: field engine.DataRecord.Data exceeds %d bytes", ColferSizeMax))
		}
		for l += x + 2; x >= 0x80; l++ {
			x >>= 7
		}
	}

	if x := len(o.Sectors); x != 0 {
		if x > ColferListMax {
			return 0, ColferMax(fmt.Sprintf("colfer: field engine.DataRecord.Sectors exceeds %d elements", ColferListMax))
		}
		for l += 2; x >= 0x80; l++ {
			x >>= 7
		}
		for _, v := range o.Sectors {
			if v == nil {
				l++
				continue
			}
			vl, err := v.MarshalLen()
			if err != nil {
				return 0, err
			}
			l += vl
		}
		if x > ColferSizeMax {
			return 0, ColferMax(fmt.Sprintf("colfer: struct engine.DataRecord size exceeds %d bytes", ColferSizeMax))
		}
	}

	if x := len(o.Cities); x != 0 {
		if x > ColferListMax {
			return 0, ColferMax(fmt.Sprintf("colfer: field engine.DataRecord.Cities exceeds %d elements", ColferListMax))
		}
		for l += 2; x >= 0x80; l++ {
			x >>= 7
		}
		for _, v := range o.Cities {
			if v == nil {
				l++
				continue
			}
			vl, err := v.MarshalLen()
			if err != nil {
				return 0, err
			}
			l += vl
		}
		if x > ColferSizeMax {
			return 0, ColferMax(fmt.Sprintf("colfer: struct engine.DataRecord size exceeds %d bytes", ColferSizeMax))
		}
	}

	if o.Active {
		l++
	}

	if o.Sort != 0 {
		l += 5
	}

	if l > ColferSizeMax {
		return l, ColferMax(fmt.Sprintf("colfer: struct engine.DataRecord exceeds %d bytes", ColferSizeMax))
	}
	return l, nil
}

// MarshalBinary encodes o as Colfer conform encoding.BinaryMarshaler.
// All nil entries in o.Sectors will be replaced with a new value.
// All nil entries in o.Cities will be replaced with a new value.
// The error return option is engine.ColferMax.
func (o *DataRecord) MarshalBinary() (data []byte, err error) {
	l, err := o.MarshalLen()
	if err != nil {
		return nil, err
	}
	data = make([]byte, l)
	o.MarshalTo(data)
	return data, nil
}

// Unmarshal decodes data as Colfer and returns the number of bytes read.
// The error return options are io.EOF, engine.ColferError and engine.ColferMax.
func (o *DataRecord) Unmarshal(data []byte) (int, error) {
	if len(data) == 0 {
		return 0, io.EOF
	}
	header := data[0]
	i := 1

	if header == 0 {
		start := i
		i++
		if i >= len(data) {
			goto eof
		}
		x := uint64(data[start])

		if x >= 0x80 {
			x &= 0x7f
			for shift := uint(7); ; shift += 7 {
				b := uint64(data[i])
				i++
				if i >= len(data) {
					goto eof
				}

				if b < 0x80 || shift == 56 {
					x |= b << shift
					break
				}
				x |= (b & 0x7f) << shift
			}
		}
		o.ID = x

		header = data[i]
		i++
	} else if header == 0|0x80 {
		start := i
		i += 8
		if i >= len(data) {
			goto eof
		}
		o.ID = intconv.Uint64(data[start:])
		header = data[i]
		i++
	}

	if header == 1 {
		if i >= len(data) {
			goto eof
		}
		x := uint(data[i])
		i++

		if x >= 0x80 {
			x &= 0x7f
			for shift := uint(7); ; shift += 7 {
				if i >= len(data) {
					goto eof
				}
				b := uint(data[i])
				i++

				if b < 0x80 {
					x |= b << shift
					break
				}
				x |= (b & 0x7f) << shift
			}
		}

		if x > uint(ColferSizeMax) {
			return 0, ColferMax(fmt.Sprintf("colfer: engine.DataRecord.Data size %d exceeds %d bytes", x, ColferSizeMax))
		}

		start := i
		i += int(x)
		if i >= len(data) {
			goto eof
		}
		o.Data = string(data[start:i])

		header = data[i]
		i++
	}

	if header == 2 {
		if i >= len(data) {
			goto eof
		}
		x := uint(data[i])
		i++

		if x >= 0x80 {
			x &= 0x7f
			for shift := uint(7); ; shift += 7 {
				if i >= len(data) {
					goto eof
				}
				b := uint(data[i])
				i++

				if b < 0x80 {
					x |= b << shift
					break
				}
				x |= (b & 0x7f) << shift
			}
		}

		if x > uint(ColferListMax) {
			return 0, ColferMax(fmt.Sprintf("colfer: engine.DataRecord.Sectors length %d exceeds %d elements", x, ColferListMax))
		}

		l := int(x)
		a := make([]*Sector, l)
		malloc := make([]Sector, l)
		for ai := range a {
			v := &malloc[ai]
			a[ai] = v

			n, err := v.Unmarshal(data[i:])
			if err != nil {
				if err == io.EOF && len(data) >= ColferSizeMax {
					return 0, ColferMax(fmt.Sprintf("colfer: engine.DataRecord size exceeds %d bytes", ColferSizeMax))
				}
				return 0, err
			}
			i += n
		}
		o.Sectors = a

		if i >= len(data) {
			goto eof
		}
		header = data[i]
		i++
	}

	if header == 3 {
		if i >= len(data) {
			goto eof
		}
		x := uint(data[i])
		i++

		if x >= 0x80 {
			x &= 0x7f
			for shift := uint(7); ; shift += 7 {
				if i >= len(data) {
					goto eof
				}
				b := uint(data[i])
				i++

				if b < 0x80 {
					x |= b << shift
					break
				}
				x |= (b & 0x7f) << shift
			}
		}

		if x > uint(ColferListMax) {
			return 0, ColferMax(fmt.Sprintf("colfer: engine.DataRecord.Cities length %d exceeds %d elements", x, ColferListMax))
		}

		l := int(x)
		a := make([]*City, l)
		malloc := make([]City, l)
		for ai := range a {
			v := &malloc[ai]
			a[ai] = v

			n, err := v.Unmarshal(data[i:])
			if err != nil {
				if err == io.EOF && len(data) >= ColferSizeMax {
					return 0, ColferMax(fmt.Sprintf("colfer: engine.DataRecord size exceeds %d bytes", ColferSizeMax))
				}
				return 0, err
			}
			i += n
		}
		o.Cities = a

		if i >= len(data) {
			goto eof
		}
		header = data[i]
		i++
	}

	if header == 4 {
		if i >= len(data) {
			goto eof
		}
		o.Active = true
		header = data[i]
		i++
	}

	if header == 5 {
		start := i
		i += 4
		if i >= len(data) {
			goto eof
		}
		o.Sort = math.Float32frombits(intconv.Uint32(data[start:]))
		header = data[i]
		i++
	}

	if header != 0x7f {
		return 0, ColferError(i - 1)
	}
	if i < ColferSizeMax {
		return i, nil
	}
eof:
	if i >= ColferSizeMax {
		return 0, ColferMax(fmt.Sprintf("colfer: struct engine.DataRecord size exceeds %d bytes", ColferSizeMax))
	}
	return 0, io.EOF
}

// UnmarshalBinary decodes data as Colfer conform encoding.BinaryUnmarshaler.
// The error return options are io.EOF, engine.ColferError, engine.ColferTail and engine.ColferMax.
func (o *DataRecord) UnmarshalBinary(data []byte) error {
	i, err := o.Unmarshal(data)
	if i < len(data) && err == nil {
		return ColferTail(i)
	}
	return err
}

type City struct {
	Value int32
}

// MarshalTo encodes o as Colfer into buf and returns the number of bytes written.
// If the buffer is too small, MarshalTo will panic.
func (o *City) MarshalTo(buf []byte) int {
	var i int

	if v := o.Value; v != 0 {
		x := uint32(v)
		if v >= 0 {
			buf[i] = 0
		} else {
			x = ^x + 1
			buf[i] = 0 | 0x80
		}
		i++
		for x >= 0x80 {
			buf[i] = byte(x | 0x80)
			x >>= 7
			i++
		}
		buf[i] = byte(x)
		i++
	}

	buf[i] = 0x7f
	i++
	return i
}

// MarshalLen returns the Colfer serial byte size.
// The error return option is engine.ColferMax.
func (o *City) MarshalLen() (int, error) {
	l := 1

	if v := o.Value; v != 0 {
		x := uint32(v)
		if v < 0 {
			x = ^x + 1
		}
		for l += 2; x >= 0x80; l++ {
			x >>= 7
		}
	}

	if l > ColferSizeMax {
		return l, ColferMax(fmt.Sprintf("colfer: struct engine.City exceeds %d bytes", ColferSizeMax))
	}
	return l, nil
}

// MarshalBinary encodes o as Colfer conform encoding.BinaryMarshaler.
// The error return option is engine.ColferMax.
func (o *City) MarshalBinary() (data []byte, err error) {
	l, err := o.MarshalLen()
	if err != nil {
		return nil, err
	}
	data = make([]byte, l)
	o.MarshalTo(data)
	return data, nil
}

// Unmarshal decodes data as Colfer and returns the number of bytes read.
// The error return options are io.EOF, engine.ColferError and engine.ColferMax.
func (o *City) Unmarshal(data []byte) (int, error) {
	if len(data) == 0 {
		return 0, io.EOF
	}
	header := data[0]
	i := 1

	if header == 0 {
		if i+1 >= len(data) {
			i++
			goto eof
		}
		x := uint32(data[i])
		i++

		if x >= 0x80 {
			x &= 0x7f
			for shift := uint(7); ; shift += 7 {
				b := uint32(data[i])
				i++
				if i >= len(data) {
					goto eof
				}

				if b < 0x80 {
					x |= b << shift
					break
				}
				x |= (b & 0x7f) << shift
			}
		}
		o.Value = int32(x)

		header = data[i]
		i++
	} else if header == 0|0x80 {
		if i+1 >= len(data) {
			i++
			goto eof
		}
		x := uint32(data[i])
		i++

		if x >= 0x80 {
			x &= 0x7f
			for shift := uint(7); ; shift += 7 {
				b := uint32(data[i])
				i++
				if i >= len(data) {
					goto eof
				}

				if b < 0x80 {
					x |= b << shift
					break
				}
				x |= (b & 0x7f) << shift
			}
		}
		o.Value = int32(^x + 1)

		header = data[i]
		i++
	}

	if header != 0x7f {
		return 0, ColferError(i - 1)
	}
	if i < ColferSizeMax {
		return i, nil
	}
eof:
	if i >= ColferSizeMax {
		return 0, ColferMax(fmt.Sprintf("colfer: struct engine.City size exceeds %d bytes", ColferSizeMax))
	}
	return 0, io.EOF
}

// UnmarshalBinary decodes data as Colfer conform encoding.BinaryUnmarshaler.
// The error return options are io.EOF, engine.ColferError, engine.ColferTail and engine.ColferMax.
func (o *City) UnmarshalBinary(data []byte) error {
	i, err := o.Unmarshal(data)
	if i < len(data) && err == nil {
		return ColferTail(i)
	}
	return err
}

type Sector struct {
	Value int32
}

// MarshalTo encodes o as Colfer into buf and returns the number of bytes written.
// If the buffer is too small, MarshalTo will panic.
func (o *Sector) MarshalTo(buf []byte) int {
	var i int

	if v := o.Value; v != 0 {
		x := uint32(v)
		if v >= 0 {
			buf[i] = 0
		} else {
			x = ^x + 1
			buf[i] = 0 | 0x80
		}
		i++
		for x >= 0x80 {
			buf[i] = byte(x | 0x80)
			x >>= 7
			i++
		}
		buf[i] = byte(x)
		i++
	}

	buf[i] = 0x7f
	i++
	return i
}

// MarshalLen returns the Colfer serial byte size.
// The error return option is engine.ColferMax.
func (o *Sector) MarshalLen() (int, error) {
	l := 1

	if v := o.Value; v != 0 {
		x := uint32(v)
		if v < 0 {
			x = ^x + 1
		}
		for l += 2; x >= 0x80; l++ {
			x >>= 7
		}
	}

	if l > ColferSizeMax {
		return l, ColferMax(fmt.Sprintf("colfer: struct engine.Sector exceeds %d bytes", ColferSizeMax))
	}
	return l, nil
}

// MarshalBinary encodes o as Colfer conform encoding.BinaryMarshaler.
// The error return option is engine.ColferMax.
func (o *Sector) MarshalBinary() (data []byte, err error) {
	l, err := o.MarshalLen()
	if err != nil {
		return nil, err
	}
	data = make([]byte, l)
	o.MarshalTo(data)
	return data, nil
}

// Unmarshal decodes data as Colfer and returns the number of bytes read.
// The error return options are io.EOF, engine.ColferError and engine.ColferMax.
func (o *Sector) Unmarshal(data []byte) (int, error) {
	if len(data) == 0 {
		return 0, io.EOF
	}
	header := data[0]
	i := 1

	if header == 0 {
		if i+1 >= len(data) {
			i++
			goto eof
		}
		x := uint32(data[i])
		i++

		if x >= 0x80 {
			x &= 0x7f
			for shift := uint(7); ; shift += 7 {
				b := uint32(data[i])
				i++
				if i >= len(data) {
					goto eof
				}

				if b < 0x80 {
					x |= b << shift
					break
				}
				x |= (b & 0x7f) << shift
			}
		}
		o.Value = int32(x)

		header = data[i]
		i++
	} else if header == 0|0x80 {
		if i+1 >= len(data) {
			i++
			goto eof
		}
		x := uint32(data[i])
		i++

		if x >= 0x80 {
			x &= 0x7f
			for shift := uint(7); ; shift += 7 {
				b := uint32(data[i])
				i++
				if i >= len(data) {
					goto eof
				}

				if b < 0x80 {
					x |= b << shift
					break
				}
				x |= (b & 0x7f) << shift
			}
		}
		o.Value = int32(^x + 1)

		header = data[i]
		i++
	}

	if header != 0x7f {
		return 0, ColferError(i - 1)
	}
	if i < ColferSizeMax {
		return i, nil
	}
eof:
	if i >= ColferSizeMax {
		return 0, ColferMax(fmt.Sprintf("colfer: struct engine.Sector size exceeds %d bytes", ColferSizeMax))
	}
	return 0, io.EOF
}

// UnmarshalBinary decodes data as Colfer conform encoding.BinaryUnmarshaler.
// The error return options are io.EOF, engine.ColferError, engine.ColferTail and engine.ColferMax.
func (o *Sector) UnmarshalBinary(data []byte) error {
	i, err := o.Unmarshal(data)
	if i < len(data) && err == nil {
		return ColferTail(i)
	}
	return err
}
