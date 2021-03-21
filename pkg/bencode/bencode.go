// Package bencode is a library for encoding/decoding bencode data into Go data structures.
package bencode

import (
	"bytes"
	"encoding"
	"errors"
	"io"
	"sort"
	"strconv"
)

const (
	// IntegerStart is the starting delimiter for a bencoded integer.
	IntegerStart byte = 'i'
	// IntegerEnd is the ending delimiter for a bencoded integer.
	IntegerEnd byte = 'e'

	// ByteStringDelimiter is the size/string delimiter for a bencoded bytestring.
	ByteStringDelimiter byte = ':'

	// ListStart is the starting delimiter for a bencoded list.
	ListStart byte = 'l'
	// ListEnd is the ending delimiter for a bencoded list.
	ListEnd byte = 'e'

	// DictStart is the starting delimiter for a bencoded dict.
	DictStart byte = 'd'
	// DictEnd is the ending delimiter for a bencoded dict.
	DictEnd byte = 'e'
)

var (
	// ErrWrongStartByte is the error returned when a wrong start
	// delimiter is found instead of the expected one
	ErrWrongStartByte = errors.New("wrong object start byte")
	// ErrUnknownType is the error returned when a non-bencode package object
	// (Integer, ByteString, List or Dict) is passed to the API
	ErrUnknownType = errors.New("unknown object type")
)

// Integer represents the bencode integer type.
type Integer struct {
	value int64
}

// NewInteger returns a bencode Integer initialized with the given parameter.
func NewInteger(i int64) Integer {
	return Integer{i}
}

// MarshalBinary satisfies the encoding.BinaryMarshaler interface
// to marshal an Integer in binary form.
func (i Integer) MarshalBinary() ([]byte, error) {
	var bb bytes.Buffer

	bb.WriteByte(IntegerStart)
	bb.WriteString(strconv.FormatInt(i.value, 10))
	bb.WriteByte(IntegerEnd)

	return bb.Bytes(), nil
}

func (i *Integer) unmarshal(bb *bytes.Buffer) error {
	start, err := bb.ReadByte()
	if err != nil {
		return err
	}
	if start != IntegerStart {
		return ErrWrongStartByte
	}
	buf, err := bb.ReadBytes(IntegerEnd)
	if err != nil {
		return err
	}

	value, err := strconv.ParseInt(string(buf[:len(buf)-1]), 10, 64)
	if err != nil {
		return err
	}
	i.value = value

	return nil
}

// UnmarshalBinary satisfies the encoding.BinaryUnmarshaler interface
// to unmarshal an Integer from binary data.
func (i *Integer) UnmarshalBinary(data []byte) error {
	return i.unmarshal(bytes.NewBuffer(data))
}

// Value returns a representation of the Integer using Go
// standard data type int64.
func (i Integer) Value() int64 {
	return i.value
}

// ByteString represents the bencode bytestring type.
type ByteString struct {
	value string
}

// NewByteString returns a bencode ByteString initialized with the given parameter.
func NewByteString(s string) ByteString {
	return ByteString{s}
}

// MarshalBinary satisfies the encoding.BinaryMarshaler interface
// to marshal a ByteString in binary form.
func (bs ByteString) MarshalBinary() ([]byte, error) {
	var bb bytes.Buffer

	bb.WriteString(strconv.Itoa(len(bs.value)))
	bb.WriteByte(':')
	bb.WriteString(bs.value)

	return bb.Bytes(), nil
}

func (bs *ByteString) unmarshal(bb *bytes.Buffer) error {
	szBuf, err := bb.ReadBytes(ByteStringDelimiter)
	if err != nil {
		return err
	}
	sz, err := strconv.Atoi(string(szBuf[:len(szBuf)-1]))
	if err != nil {
		return err
	}
	dataBuf := make([]byte, sz)
	if _, err := io.ReadFull(bb, dataBuf); err != nil {
		return err
	}
	bs.value = string(dataBuf)

	return nil
}

// UnmarshalBinary satisfies the encoding.BinaryUnmarshaler interface
// to unmarshal a ByteString from binary data.
func (bs *ByteString) UnmarshalBinary(data []byte) error {
	return bs.unmarshal(bytes.NewBuffer(data))
}

// Value returns a representation of the ByteString using Go
// standard data type string.
func (bs ByteString) Value() string {
	return bs.value
}

// List represents the bencode list type.
type List struct {
	value []interface{}
}

// NewList returns a bencode List initialized with the given parameter.
func NewList(l []interface{}) List {
	return List{l}
}

// MarshalBinary satisfies the encoding.BinaryMarshaler interface
// to marshal a List in binary form.
func (l List) MarshalBinary() ([]byte, error) {
	var bb bytes.Buffer

	bb.WriteByte(ListStart)
	for _, v := range l.value {
		value, ok := v.(encoding.BinaryMarshaler)
		if !ok {
			return []byte{}, ErrUnknownType
		}
		buf, err := value.MarshalBinary()
		if err != nil {
			return []byte{}, err
		}
		bb.Write(buf)
	}
	bb.WriteByte(ListEnd)

	return bb.Bytes(), nil
}

func (l *List) unmarshal(bb *bytes.Buffer) error {
	l.value = []interface{}{}

	start, err := bb.ReadByte()
	if err != nil {
		return err
	}
	if start != ListStart {
		return ErrWrongStartByte
	}

	for {
		cur, err := bb.ReadByte()
		if err != nil {
			return err
		}

		switch cur {
		case ListEnd:
			return nil
		case IntegerStart:
			if err := bb.UnreadByte(); err != nil {
				return err
			}
			obj := Integer{}
			if err := obj.unmarshal(bb); err != nil {
				return err
			}
			l.value = append(l.value, obj)
		case ListStart:
			if err := bb.UnreadByte(); err != nil {
				return err
			}
			obj := List{}
			if err := obj.unmarshal(bb); err != nil {
				return err
			}
			l.value = append(l.value, obj)
		case DictStart:
			if err := bb.UnreadByte(); err != nil {
				return err
			}
			obj := Dict{}
			if err := obj.unmarshal(bb); err != nil {
				return err
			}
			l.value = append(l.value, obj)
		default:
			if err := bb.UnreadByte(); err != nil {
				return err
			}
			obj := ByteString{}
			if err := obj.unmarshal(bb); err != nil {
				return err
			}
			l.value = append(l.value, obj)
		}
	}
}

// UnmarshalBinary satisfies the encoding.BinaryUnmarshaler interface
// to unmarshal a List from binary data.
func (l *List) UnmarshalBinary(data []byte) error {
	return l.unmarshal(bytes.NewBuffer(data))
}

// Value returns a representation of the List using Go
// standard data type []interface{}.
func (l List) Value() []interface{} {
	values := make([]interface{}, 0, len(l.value))

	for _, v := range l.value {
		switch value := v.(type) {
		case Integer:
			values = append(values, value.Value())
		case ByteString:
			values = append(values, value.Value())
		case List:
			values = append(values, value.Value())
		case Dict:
			values = append(values, value.Value())
		}
	}

	return values
}

// Dict represents the bencode dict type.
type Dict struct {
	value map[ByteString]interface{}
}

// NewDict returns a bencode Dict initialized with the given parameter.
func NewDict(d map[ByteString]interface{}) Dict {
	return Dict{d}
}

// MarshalBinary satisfies the encoding.BinaryMarshaler interface
// to marshal a Dict in binary form.
func (d Dict) MarshalBinary() ([]byte, error) {
	var bb bytes.Buffer

	keys := make([]ByteString, 0, len(d.value))
	for k := range d.value {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i].value < keys[j].value
	})

	bb.WriteByte(DictStart)
	for _, k := range keys {
		buf, err := k.MarshalBinary()
		if err != nil {
			return []byte{}, err
		}
		bb.Write(buf)

		value, ok := d.value[k].(encoding.BinaryMarshaler)
		if !ok {
			return []byte{}, ErrUnknownType
		}
		buf, err = value.MarshalBinary()
		if err != nil {
			return []byte{}, err
		}
		bb.Write(buf)
	}
	bb.WriteByte(DictEnd)

	return bb.Bytes(), nil
}

func (d *Dict) unmarshal(bb *bytes.Buffer) error {
	d.value = map[ByteString]interface{}{}

	start, err := bb.ReadByte()
	if err != nil {
		return err
	}
	if start != DictStart {
		return ErrWrongStartByte
	}

	for {
		cur, err := bb.ReadByte()
		if err != nil {
			return err
		}
		if cur == DictEnd {
			break
		}
		if err := bb.UnreadByte(); err != nil {
			return err
		}

		key := ByteString{}
		if err := key.unmarshal(bb); err != nil {
			return err
		}

		cur, err = bb.ReadByte()
		if err != nil {
			return err
		}
		if err := bb.UnreadByte(); err != nil {
			return err
		}

		switch cur {
		case IntegerStart:
			obj := Integer{}
			if err := obj.unmarshal(bb); err != nil {
				return err
			}
			d.value[key] = obj
		case ListStart:
			obj := List{}
			if err := obj.unmarshal(bb); err != nil {
				return err
			}
			d.value[key] = obj
		case DictStart:
			obj := Dict{}
			if err := obj.unmarshal(bb); err != nil {
				return err
			}
			d.value[key] = obj
		default:
			obj := ByteString{}
			if err := obj.unmarshal(bb); err != nil {
				return err
			}
			d.value[key] = obj
		}
	}

	return nil
}

// UnmarshalBinary satisfies the encoding.BinaryUnmarshaler interface
// to unmarshal a Dict from binary data.
func (d *Dict) UnmarshalBinary(data []byte) error {
	return d.unmarshal(bytes.NewBuffer(data))
}

// Value returns a representation of the Dict using Go
// standard data type map[string]interface{}.
func (d Dict) Value() map[string]interface{} {
	values := make(map[string]interface{}, len(d.value))
	for k, v := range d.value {
		switch value := v.(type) {
		case Integer:
			values[k.Value()] = value.Value()
		case ByteString:
			values[k.Value()] = value.Value()
		case List:
			values[k.Value()] = value.Value()
		case Dict:
			values[k.Value()] = value.Value()
		}
	}
	return values
}

// Encoder writes bencode values to an output stream.
type Encoder struct {
	w io.Writer
}

// NewEncoder returns a new encoder that writes to w.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w}
}

// Encode writes the bencode encoding of v to the stream.
func (e *Encoder) Encode(v interface{}) error {
	switch value := v.(type) {
	case encoding.BinaryMarshaler:
		buf, err := value.MarshalBinary()
		if err != nil {
			return err
		}
		if _, err := e.w.Write(buf); err != nil {
			return err
		}
	default:
		return ErrUnknownType
	}

	return nil
}

// A Decoder reads and decodes bencode values from an input stream.
type Decoder struct {
	r io.Reader
}

// NewDecoder returns a new decoder that reads from r.
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r}
}

// Decode reads the next bencode-encoded value from its input and
// stores it in the value pointed to by v.
func (d *Decoder) Decode(v interface{}) error {
	var bb bytes.Buffer

	if _, err := bb.ReadFrom(d.r); err != nil {
		return err
	}

	switch value := v.(type) {
	case *Integer:
		if err := value.unmarshal(&bb); err != nil {
			return err
		}
	case *ByteString:
		if err := value.unmarshal(&bb); err != nil {
			return err
		}
	case *List:
		if err := value.unmarshal(&bb); err != nil {
			return err
		}
	case *Dict:
		if err := value.unmarshal(&bb); err != nil {
			return err
		}
	default:
		return ErrUnknownType
	}

	return nil
}
