package bencode

import (
	"bytes"
	"errors"
	"io"
	"testing"
)

func compareList(t *testing.T, a, b List) bool {
	t.Helper()

	for i := 0; i < len(a.value); i++ {
		first := a.value[i]
		second := b.value[i]

		switch firstValue := first.(type) {
		case Integer:
			secondValue, ok := second.(Integer)
			if !ok {
				return false
			}
			if firstValue.value != secondValue.value {
				return false
			}
		case ByteString:
			secondValue, ok := second.(ByteString)
			if !ok {
				return false
			}
			if firstValue.value != secondValue.value {
				return false
			}
		case List:
			secondValue, ok := second.(List)
			if !ok {
				return false
			}
			if !compareList(t, firstValue, secondValue) {
				return false
			}
		case Dict:
			secondValue, ok := second.(Dict)
			if !ok {
				return false
			}
			if !compareDict(t, firstValue, secondValue) {
				return false
			}
		default:
			return false
		}
	}

	return true
}

func compareDict(t *testing.T, a, b Dict) bool {
	t.Helper()

	for key, first := range a.value {
		second, ok := b.value[key]
		if !ok {
			return false
		}

		switch firstValue := first.(type) {
		case Integer:
			secondValue, ok := second.(Integer)
			if !ok {
				return false
			}
			if firstValue.value != secondValue.value {
				return false
			}
		case ByteString:
			secondValue, ok := second.(ByteString)
			if !ok {
				return false
			}
			if firstValue.value != secondValue.value {
				return false
			}
		case List:
			secondValue, ok := second.(List)
			if !ok {
				return false
			}
			if !compareList(t, firstValue, secondValue) {
				return false
			}
		case Dict:
			secondValue, ok := second.(Dict)
			if !ok {
				return false
			}
			if !compareDict(t, firstValue, secondValue) {
				return false
			}
		default:
			return false
		}
	}

	return true
}

var intMarshalTestCases = []struct {
	name     string
	input    Integer
	expected []byte
}{
	{
		name:     "zero value",
		input:    Integer{0},
		expected: []byte{'i', '0', 'e'},
	},
	{
		name:     "negative value",
		input:    Integer{-145},
		expected: []byte{'i', '-', '1', '4', '5', 'e'},
	},
	{
		name:  "large positive value",
		input: Integer{11435223524},
		expected: []byte{
			'i',
			'1', '1', '4', '3', '5', '2', '2', '3', '5', '2', '4',
			'e',
		},
	},
}

func TestIntegerMarshal(t *testing.T) {
	for _, tc := range intMarshalTestCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.input.MarshalBinary()
			if err != nil {
				t.Fatalf("unexpected error: %v\n", err)
			}

			if !bytes.Equal(got, tc.expected) {
				t.Fatalf("expected %v got %v\n", tc.expected, got)
			}
		})
	}
}

var intUnmarshalTestCases = []struct {
	name     string
	input    []byte
	expected Integer
}{
	{
		name:     "zero value",
		input:    []byte{'i', '0', 'e'},
		expected: Integer{0},
	},
	{
		name:     "negative value",
		input:    []byte{'i', '-', '1', '4', '5', 'e'},
		expected: Integer{-145},
	},
	{
		name: "large positive value",
		input: []byte{
			'i',
			'1', '1', '4', '3', '5', '2', '2', '3', '5', '2', '4',
			'e',
		},
		expected: Integer{11435223524},
	},
}

func TestIntegerUnmarshal(t *testing.T) {
	for _, tc := range intUnmarshalTestCases {
		t.Run(tc.name, func(t *testing.T) {
			got := Integer{}
			err := got.UnmarshalBinary(tc.input)
			if err != nil {
				t.Fatal(err)
			}

			if got.Value() != tc.expected.Value() {
				t.Fatalf("expected %d got %d\n", tc.expected.Value(), got.Value())
			}
		})
	}
}

var intUnmarshalErrorTestCases = []struct {
	name     string
	input    []byte
	expected error
}{
	{
		name:     "missing end byte",
		input:    []byte{'i', '0'},
		expected: io.EOF,
	},
	{
		name:     "wrong start byte",
		input:    []byte{'-', '1', '4', '5', 'e'},
		expected: ErrWrongStartByte,
	},
}

func TestIntegerUnmarshalError(t *testing.T) {
	for _, tc := range intUnmarshalErrorTestCases {
		t.Run(tc.name, func(t *testing.T) {
			got := Integer{}
			err := got.UnmarshalBinary(tc.input)
			if err == nil {
				t.Fatal("expected error, got nil")
			}

			if !errors.Is(err, tc.expected) {
				t.Fatalf("expected error %v, got %v", tc.expected, got)
			}
		})
	}
}

var stringMarshalTestCases = []struct {
	name     string
	input    ByteString
	expected []byte
}{
	{
		name:     "empty string value",
		input:    ByteString{""},
		expected: []byte{'0', ':'},
	},
	{
		name:     "short string value",
		input:    ByteString{"test"},
		expected: []byte{'4', ':', 't', 'e', 's', 't'},
	},
	{
		name:  "long string value",
		input: ByteString{"a longer string"},
		expected: []byte{
			'1', '5',
			':',
			'a', ' ', 'l', 'o', 'n', 'g', 'e', 'r', ' ',
			's', 't', 'r', 'i', 'n', 'g',
		},
	},
}

func TestByteStringMarshal(t *testing.T) {
	for _, tc := range stringMarshalTestCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.input.MarshalBinary()
			if err != nil {
				t.Fatalf("unexpected error: %v\n", err)
			}

			if !bytes.Equal(got, tc.expected) {
				t.Fatalf("expected %v got %v\n", tc.expected, got)
			}
		})
	}
}

var byteStringUnmarshalTestCases = []struct {
	name     string
	input    []byte
	expected ByteString
}{
	{
		name:     "short string value",
		input:    []byte{'4', ':', 't', 'e', 's', 't'},
		expected: ByteString{"test"},
	},
	{
		name: "long string value",
		input: []byte{
			'1', '5',
			':',
			'a', ' ', 'l', 'o', 'n', 'g', 'e', 'r', ' ',
			's', 't', 'r', 'i', 'n', 'g',
		},
		expected: ByteString{"a longer string"},
	},
}

func TestByteStringUnmarshal(t *testing.T) {
	for _, tc := range byteStringUnmarshalTestCases {
		t.Run(tc.name, func(t *testing.T) {
			got := ByteString{}
			err := got.UnmarshalBinary(tc.input)
			if err != nil {
				t.Fatal(err)
			}

			if got.value != tc.expected.value {
				t.Fatalf("expected %v got %v\n", tc.expected.value, got.value)
			}
		})
	}
}

var listMarshalTestCases = []struct {
	name     string
	input    List
	expected []byte
}{
	{
		name:     "empty value",
		input:    List{},
		expected: []byte{'l', 'e'},
	},
	{
		name:     "integer and bytestring",
		input:    List{[]interface{}{Integer{12}, ByteString{"test"}}},
		expected: []byte{'l', 'i', '1', '2', 'e', '4', ':', 't', 'e', 's', 't', 'e'},
	},
	{
		name: "integer and inner list",
		input: List{
			[]interface{}{
				List{
					[]interface{}{ByteString{"test"}, ByteString{"again"}},
				},
				Integer{5},
			},
		},
		expected: []byte{
			'l',
			'l', '4', ':', 't', 'e', 's', 't', '5', ':', 'a', 'g', 'a', 'i', 'n', 'e',
			'i', '5', 'e',
			'e',
		},
	},
	{
		name: "inner list and inner dict",
		input: List{
			[]interface{}{
				List{
					[]interface{}{ByteString{"test"}, ByteString{"again"}},
				},
				Dict{
					map[ByteString]interface{}{
						{"integer"}: Integer{12},
						{"string"}:  ByteString{"test"},
					},
				},
			},
		},
		expected: []byte{
			'l',
			'l', '4', ':', 't', 'e', 's', 't', '5', ':', 'a', 'g', 'a', 'i', 'n', 'e',
			'd',
			'7', ':', 'i', 'n', 't', 'e', 'g', 'e', 'r', 'i', '1', '2', 'e',
			'6', ':', 's', 't', 'r', 'i', 'n', 'g', '4', ':', 't', 'e', 's', 't',
			'e',
			'e',
		},
	},
}

func TestListMarshal(t *testing.T) {
	for _, tc := range listMarshalTestCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.input.MarshalBinary()
			if err != nil {
				t.Fatalf("unexpected error: %v\n", err)
			}

			if !bytes.Equal(got, tc.expected) {
				t.Fatalf("expected %v got %v\n", tc.expected, got)
			}
		})
	}
}

var listUnmarshalTestCases = []struct {
	name     string
	input    []byte
	expected List
}{
	{
		name:     "empty value",
		input:    []byte{'l', 'e'},
		expected: List{},
	},
	{
		name:     "integer and bytestring",
		input:    []byte{'l', 'i', '1', '2', 'e', '4', ':', 't', 'e', 's', 't', 'e'},
		expected: List{[]interface{}{Integer{12}, ByteString{"test"}}},
	},
	{
		name: "integer and inner list",
		input: []byte{
			'l',
			'l', '4', ':', 't', 'e', 's', 't', '5', ':', 'a', 'g', 'a', 'i', 'n', 'e',
			'i', '5', 'e',
			'e',
		},
		expected: List{
			[]interface{}{
				List{
					[]interface{}{ByteString{"test"}, ByteString{"again"}},
				},
				Integer{5},
			},
		},
	},
	{
		name: "inner list and inner dict",
		input: []byte{
			'l',
			'l', '4', ':', 't', 'e', 's', 't', '5', ':', 'a', 'g', 'a', 'i', 'n', 'e',
			'd',
			'7', ':', 'i', 'n', 't', 'e', 'g', 'e', 'r', 'i', '1', '2', 'e',
			'6', ':', 's', 't', 'r', 'i', 'n', 'g', '4', ':', 't', 'e', 's', 't',
			'e',
			'e',
		},
		expected: List{
			[]interface{}{
				List{
					[]interface{}{ByteString{"test"}, ByteString{"again"}},
				},
				Dict{
					map[ByteString]interface{}{
						{"integer"}: Integer{12},
						{"string"}:  ByteString{"test"},
					},
				},
			},
		},
	},
}

func TestListUnmarshal(t *testing.T) {
	for _, tc := range listUnmarshalTestCases {
		t.Run(tc.name, func(t *testing.T) {
			got := List{}
			err := got.UnmarshalBinary(tc.input)
			if err != nil {
				t.Fatal(err)
			}

			if len(got.Value()) != len(tc.expected.Value()) {
				t.Fatalf("expected length %d got %d\n", len(tc.expected.Value()), len(got.Value()))
			}

			if !compareList(t, got, tc.expected) {
				t.Fatal("lists are not equal")
			}

		})
	}
}

var listUnmarshalErrorTestCases = []struct {
	name     string
	input    []byte
	expected error
}{
	{
		name:     "missing end byte",
		input:    []byte{'l', '4', ':', 't', 'e', 's', 't'},
		expected: io.EOF,
	},
	{
		name:     "wrong start byte",
		input:    []byte{'-', '1', '4', '5', 'e'},
		expected: ErrWrongStartByte,
	},
}

func TestListUnmarshalError(t *testing.T) {
	for _, tc := range listUnmarshalErrorTestCases {
		t.Run(tc.name, func(t *testing.T) {
			got := List{}
			err := got.UnmarshalBinary(tc.input)
			if err == nil {
				t.Fatal("expected error, got nil")
			}

			if !errors.Is(err, tc.expected) {
				t.Fatalf("expected error %v, got %v", tc.expected, err)
			}
		})
	}
}

var dictMarshalTestCases = []struct {
	name     string
	input    Dict
	expected []byte
}{
	{
		name:     "empty value",
		input:    Dict{},
		expected: []byte{'d', 'e'},
	},
	{
		name: "integer and bytestring",
		input: Dict{map[ByteString]interface{}{
			{"one"}: Integer{12},
			{"two"}: ByteString{"test"},
		},
		},
		expected: []byte{
			'd',
			'3', ':', 'o', 'n', 'e', 'i', '1', '2', 'e',
			'3', ':', 't', 'w', 'o', '4', ':', 't', 'e', 's', 't',
			'e',
		},
	},
	{
		name: "integer and inner list",
		input: Dict{
			map[ByteString]interface{}{
				{"integer"}: Integer{12},
				{"list"}: List{
					[]interface{}{ByteString{"test"}, ByteString{"again"}},
				},
			},
		},
		expected: []byte{
			'd',
			'7', ':', 'i', 'n', 't', 'e', 'g', 'e', 'r', 'i', '1', '2', 'e',
			'4', ':', 'l', 'i', 's', 't',
			'l', '4', ':', 't', 'e', 's', 't', '5', ':', 'a', 'g', 'a', 'i', 'n', 'e',
			'e',
		},
	},
}

func TestDictMarshal(t *testing.T) {
	for _, tc := range dictMarshalTestCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.input.MarshalBinary()
			if err != nil {
				t.Fatalf("unexpected error: %v\n", err)
			}

			if !bytes.Equal(got, tc.expected) {
				t.Fatalf("expected %v got %v\n", tc.expected, got)
			}
		})
	}
}

var dictUnmarshalTestCases = []struct {
	name     string
	input    []byte
	expected Dict
}{
	{
		name:     "empty value",
		input:    []byte{'d', 'e'},
		expected: Dict{},
	},
	{
		name: "integer and bytestring",
		input: []byte{
			'd',
			'3', ':', 'o', 'n', 'e', 'i', '1', '2', 'e',
			'3', ':', 't', 'w', 'o', '4', ':', 't', 'e', 's', 't',
			'e',
		},
		expected: Dict{
			map[ByteString]interface{}{
				{"one"}: Integer{12},
				{"two"}: ByteString{"test"},
			},
		},
	},
	{
		name: "integer and inner list",
		input: []byte{
			'd',
			'7', ':', 'i', 'n', 't', 'e', 'g', 'e', 'r', 'i', '1', '2', 'e',
			'4', ':', 'l', 'i', 's', 't',
			'l', '4', ':', 't', 'e', 's', 't', '5', ':', 'a', 'g', 'a', 'i', 'n', 'e',
			'e',
		},
		expected: Dict{
			map[ByteString]interface{}{
				{"integer"}: Integer{12},
				{"list"}: List{
					[]interface{}{ByteString{"test"}, ByteString{"again"}},
				},
			},
		},
	},
}

func TestDictUnmarshal(t *testing.T) {
	for _, tc := range dictUnmarshalTestCases {
		t.Run(tc.name, func(t *testing.T) {
			got := Dict{}
			err := got.UnmarshalBinary(tc.input)
			if err != nil {
				t.Fatal(err)
			}

			if len(got.Value()) != len(tc.expected.Value()) {
				t.Fatalf("expected length %d got %d\n", len(tc.expected.Value()), len(got.Value()))
			}

			if !compareDict(t, got, tc.expected) {
				t.Fatal("dicts are not equal")
			}

		})
	}
}

var dictUnmarshalErrorTestCases = []struct {
	name     string
	input    []byte
	expected error
}{
	{
		name:     "missing end byte",
		input:    []byte{'d', '4', ':', 't', 'e', 's', 't', 'i', '9', 'e'},
		expected: io.EOF,
	},
	{
		name:     "wrong start byte",
		input:    []byte{'-', '1', '4', '5', 'e'},
		expected: ErrWrongStartByte,
	},
}

func TestDictUnmarshalError(t *testing.T) {
	for _, tc := range dictUnmarshalErrorTestCases {
		t.Run(tc.name, func(t *testing.T) {
			got := Dict{}
			err := got.UnmarshalBinary(tc.input)
			if err == nil {
				t.Fatal("expected error, got nil")
			}

			if !errors.Is(err, tc.expected) {
				t.Fatalf("expected error %v, got %v", tc.expected, err)
			}
		})
	}
}

var encoderTestCases = []struct {
	name     string
	input    interface{}
	expected []byte
}{
	{
		name:     "empty list value",
		input:    List{},
		expected: []byte{'l', 'e'},
	},
	{
		name:     "empty dict value",
		input:    Dict{},
		expected: []byte{'d', 'e'},
	},
	{
		name: "dict with integer and bytestring",
		input: Dict{
			map[ByteString]interface{}{
				{"one"}: Integer{12},
				{"two"}: ByteString{"test"},
			},
		},
		expected: []byte{
			'd',
			'3', ':', 'o', 'n', 'e', 'i', '1', '2', 'e',
			'3', ':', 't', 'w', 'o', '4', ':', 't', 'e', 's', 't',
			'e',
		},
	},
	{
		name: "dict with integer and inner list",
		input: Dict{
			map[ByteString]interface{}{
				{"integer"}: Integer{12},
				{"list"}: List{
					[]interface{}{ByteString{"test"}, ByteString{"again"}},
				},
			},
		},
		expected: []byte{
			'd',
			'7', ':', 'i', 'n', 't', 'e', 'g', 'e', 'r', 'i', '1', '2', 'e',
			'4', ':', 'l', 'i', 's', 't',
			'l', '4', ':', 't', 'e', 's', 't', '5', ':', 'a', 'g', 'a', 'i', 'n', 'e',
			'e',
		},
	},
	{
		name: "list with integer and inner list",
		input: List{
			[]interface{}{
				List{
					[]interface{}{ByteString{"test"}, ByteString{"again"}},
				},
				Integer{5},
			},
		},
		expected: []byte{
			'l',
			'l', '4', ':', 't', 'e', 's', 't', '5', ':', 'a', 'g', 'a', 'i', 'n', 'e',
			'i', '5', 'e',
			'e',
		},
	},
}

func TestEncoder(t *testing.T) {
	for _, tc := range encoderTestCases {
		t.Run(tc.name, func(t *testing.T) {
			var bb bytes.Buffer
			enc := NewEncoder(&bb)
			if err := enc.Encode(tc.input); err != nil {
				t.Fatal(err)
			}

			if !bytes.Equal(bb.Bytes(), tc.expected) {
				t.Fatalf("expected %v, got %v\n", tc.expected, bb.Bytes())
			}
		})
	}
}

var decoderDictTestCases = []struct {
	name     string
	input    []byte
	expected Dict
}{
	{
		name:     "empty dict value",
		input:    []byte{'d', 'e'},
		expected: Dict{},
	},
	{
		name: "dict with integer and bytestring",
		input: []byte{
			'd',
			'3', ':', 'o', 'n', 'e', 'i', '1', '2', 'e',
			'3', ':', 't', 'w', 'o', '4', ':', 't', 'e', 's', 't',
			'e',
		},
		expected: Dict{
			map[ByteString]interface{}{
				{"one"}: Integer{12},
				{"two"}: ByteString{"test"},
			},
		},
	},
	{
		name: "dict with integer and inner list",
		input: []byte{
			'd',
			'7', ':', 'i', 'n', 't', 'e', 'g', 'e', 'r', 'i', '1', '2', 'e',
			'4', ':', 'l', 'i', 's', 't',
			'l', '4', ':', 't', 'e', 's', 't', '5', ':', 'a', 'g', 'a', 'i', 'n', 'e',
			'e',
		},
		expected: Dict{
			map[ByteString]interface{}{
				{"integer"}: Integer{12},
				{"list"}: List{
					[]interface{}{ByteString{"test"}, ByteString{"again"}},
				},
			},
		},
	},
}

func TestDecoderDict(t *testing.T) {
	for _, tc := range decoderDictTestCases {
		t.Run(tc.name, func(t *testing.T) {
			d := Dict{}
			dec := NewDecoder(bytes.NewReader(tc.input))
			if err := dec.Decode(&d); err != nil {
				t.Fatal(err)
			}

			if !compareDict(t, d, tc.expected) {
				t.Fatalf("expected %v, got %v\n", tc.expected, d)
			}
		})
	}
}

var decoderListTestCases = []struct {
	name     string
	input    []byte
	expected List
}{
	{
		name:     "empty list value",
		input:    []byte{'l', 'e'},
		expected: List{},
	},
	{
		name: "list with integer and inner list",
		input: []byte{
			'l',
			'l', '4', ':', 't', 'e', 's', 't', '5', ':', 'a', 'g', 'a', 'i', 'n', 'e',
			'i', '5', 'e',
			'e',
		},
		expected: List{
			[]interface{}{
				List{
					[]interface{}{ByteString{"test"}, ByteString{"again"}},
				},
				Integer{5},
			},
		},
	},
}

func TestDecoderList(t *testing.T) {
	for _, tc := range decoderListTestCases {
		t.Run(tc.name, func(t *testing.T) {
			l := List{}
			dec := NewDecoder(bytes.NewReader(tc.input))
			if err := dec.Decode(&l); err != nil {
				t.Fatal(err)
			}

			if !compareList(t, l, tc.expected) {
				t.Fatalf("expected %v, got %v\n", tc.expected, l)
			}
		})
	}
}

// Benchmarks data is the same as github.com/jackpal/bencode-go

func BenchmarkBencodeMarshal(b *testing.B) {
	data := Dict{
		map[ByteString]interface{}{
			{"announce"}: ByteString{"udp://tracker.publicbt.com:80/announce"},
			{"announce-list"}: List{
				[]interface{}{
					ByteString{"udp://tracker.publicbt.com:80/announce"},
					ByteString{"udp://tracker.openbittorrent.com:80/announce"},
				},
			},
			{"comment"}: ByteString{"Debian CD from cdimage.debian.org"},
			{"info"}: Dict{
				map[ByteString]interface{}{
					{"name"}:         ByteString{"debian-8.8.0-arm64-netinst.iso"},
					{"length"}:       Integer{170917888},
					{"piece length"}: Integer{262144},
				},
			},
		},
	}

	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		var bb bytes.Buffer
		enc := NewEncoder(&bb)
		if err := enc.Encode(data); err != nil {
			b.Fatal(err)
		}
		_ = data.Value()
	}
}

func BenchmarkBencodeUnmarshal(b *testing.B) {
	data := []byte("d4:infod6:lengthi170917888e12:piece lengthi262144e4:name30:debian-8.8.0-arm64-netinst.isoe8:announce38:udp://tracker.publicbt.com:80/announce13:announce-listll38:udp://tracker.publicbt.com:80/announceel44:udp://tracker.openbittorrent.com:80/announceee7:comment33:Debian CD from cdimage.debian.orge")

	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		d := Dict{}
		dec := NewDecoder(bytes.NewReader(data))
		if err := dec.Decode(&d); err != nil {
			b.Fatal(err)
		}
	}
}
