package classfileparser

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"testing"
)

func fullClass() testClass {
	return testClass{
		minor:      3,
		major:      61,
		pool:       [][]byte{cpUtf8("Foo"), cpClass(1), cpUtf8("x"), cpLong(7)},
		access:     0x0021,
		this:       2,
		super:      2,
		interfaces: []uint16{2},
		fields:     []testMember{{access: 0x0002, name: 3, descriptor: 3, attributes: []testAttribute{{name: 3, info: []byte("f")}}}},
		methods:    []testMember{{access: 0x0001, name: 3, descriptor: 3, attributes: []testAttribute{{name: 3, info: []byte("m")}}}},
		attributes: []testAttribute{{name: 3, info: []byte("c")}},
	}
}

func TestOpen(t *testing.T) {
	cf, err := Open(bytes.NewReader(fullClass().bytes()))
	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}

	want := &ClassFile{
		Magic:             0xCAFEBABE,
		MinorVersion:      3,
		MajorVersion:      61,
		ConstantPoolCount: 6,
		ConstantPool: []CpInfo{
			{Tag: 1, Info: []byte("Foo")},
			{Tag: 7, Info: []byte{0, 1}},
			{Tag: 1, Info: []byte("x")},
			{Tag: 5, Info: []byte{0, 0, 0, 0, 0, 0, 0, 7}},
			{},
		},
		AccessFlags:     0x0021,
		ThisClass:       2,
		SuperClass:      2,
		InterfacesCount: 1,
		Interfaces:      []uint16{2},
		FieldsCount:     1,
		Fields: []FieldInfo{{
			AccessFlags:     0x0002,
			NameIndex:       3,
			DescriptorIndex: 3,
			AttributesCount: 1,
			Attributes:      []AttributeInfo{{AttributeNameIndex: 3, AttributeLength: 1, Info: []byte("f")}},
		}},
		MethodsCount: 1,
		Methods: []MethodInfo{{
			AccessFlags:     0x0001,
			NameIndex:       3,
			DescriptorIndex: 3,
			AttributesCount: 1,
			Attributes:      []AttributeInfo{{AttributeNameIndex: 3, AttributeLength: 1, Info: []byte("m")}},
		}},
		AttributesCount: 1,
		Attributes:      []AttributeInfo{{AttributeNameIndex: 3, AttributeLength: 1, Info: []byte("c")}},
	}
	if !reflect.DeepEqual(cf, want) {
		t.Fatalf("got %+v\nwant %+v", cf, want)
	}
}

func TestOpenConstantPoolEntries(t *testing.T) {
	cases := map[string][]byte{
		"Utf8":               cpUtf8("hello"),
		"empty Utf8":         cpUtf8(""),
		"Integer":            cpInteger(-1),
		"Float":              cpFloat(1.5),
		"Long":               cpLong(-1),
		"Double":             cpDouble(1.5),
		"Class":              cpClass(1),
		"String":             cpString(1),
		"Fieldref":           cpFieldref(1, 2),
		"Methodref":          cpMethodref(1, 2),
		"InterfaceMethodref": cpInterfaceMethodref(1, 2),
		"NameAndType":        cpNameAndType(1, 2),
		"MethodHandle":       cpMethodHandle(6, 1),
		"MethodType":         cpMethodType(1),
		"Dynamic":            cpDynamic(0, 1),
		"InvokeDynamic":      cpInvokeDynamic(0, 1),
		"Module":             cpModule(1),
		"Package":            cpPackage(1),
	}

	for name, entry := range cases {
		t.Run(name, func(t *testing.T) {
			c := minimalClass()
			c.pool = append(c.pool, entry)
			cf, err := Open(bytes.NewReader(c.bytes()))
			if err != nil {
				t.Fatalf("Open failed: %v", err)
			}
			want := CpInfo{Tag: entry[0], Info: entry[1:]}
			if entry[0] == 1 {
				want.Info = entry[3:]
			}
			if got := cf.ConstantPool[2]; !reflect.DeepEqual(got, want) {
				t.Fatalf("got %+v, want %+v", got, want)
			}
		})
	}
}

func TestOpenMalformed(t *testing.T) {
	unknownTag := minimalClass()
	unknownTag.pool = append(unknownTag.pool, u1(2))

	cases := map[string]struct {
		data    []byte
		message string
	}{
		"invalid magic number": {
			cat(u4(0xCAFEBABF), u2(0), u2(61)),
			"malformed class file: invalid magic number: 0xCAFEBABF",
		},
		"invalid constant pool count": {
			cat(u4(0xCAFEBABE), u2(0), u2(61), u2(0)),
			"malformed class file: invalid constant pool count: 0",
		},
		"unknown constant pool tag": {
			unknownTag.bytes(),
			"failed to read constant pool length: malformed class file: unknown constant pool tag: 2",
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := Open(bytes.NewReader(c.data))
			if !errors.Is(err, ErrMalformed) || err.Error() != c.message {
				t.Fatalf("got %v, want %s", err, c.message)
			}
		})
	}
}

func TestOpenTruncated(t *testing.T) {
	data := fullClass().bytes()

	cases := map[int]string{
		0:  "failed to read magic number: EOF",
		2:  "failed to read magic number: unexpected EOF",
		4:  "failed to read minor version: EOF",
		6:  "failed to read major version: EOF",
		8:  "failed to read constant pool count: EOF",
		10: "failed to read constant pool tag: EOF",
		11: "failed to read constant pool length: failed to read CONSTANT_Utf8 length: EOF",
		13: "failed to read constant pool info: EOF",
		14: "failed to read constant pool info: unexpected EOF",
		32: "failed to read access flags: EOF",
		34: "failed to read this class index: EOF",
		36: "failed to read super class index: EOF",
		38: "failed to read interfaces count: EOF",
		40: "failed to read interfaces: EOF",
		42: "failed to read fields count: EOF",
		44: "failed to read field 0: EOF",
		46: "failed to read field 0: EOF",
		48: "failed to read field 0: EOF",
		50: "failed to read field 0: EOF",
		52: "failed to read field 0: EOF",
		54: "failed to read field 0: EOF",
		58: "failed to read field 0: unexpected EOF",
		59: "failed to read methods count: EOF",
		61: "failed to read method 0: EOF",
		63: "failed to read method 0: EOF",
		65: "failed to read method 0: EOF",
		67: "failed to read method 0: EOF",
		69: "failed to read method 0: EOF",
		71: "failed to read method 0: EOF",
		75: "failed to read method 0: unexpected EOF",
		76: "failed to read attributes count: EOF",
		78: "failed to read attribute 0: EOF",
		80: "failed to read attribute 0: EOF",
		84: "failed to read attribute 0: unexpected EOF",
	}
	if len(data) != 85 {
		t.Fatalf("unexpected test class length: %d", len(data))
	}

	for cut, message := range cases {
		_, err := Open(bytes.NewReader(data[:cut]))
		if !isTruncation(err) || err.Error() != message {
			t.Errorf("cut at %d: got %v, want %s", cut, err, message)
		}
	}

	for cut := range data {
		_, err := Open(bytes.NewReader(data[:cut]))
		if !isTruncation(err) || errors.Is(err, ErrMalformed) {
			t.Errorf("cut at %d: got %v, want a truncation error", cut, err)
		}
	}
}

func TestOpenHugeAttributeLength(t *testing.T) {
	minimal := minimalClass().bytes()
	data := cat(minimal[:len(minimal)-2], u2(1), u2(1), u4(0xFFFFFFFF), []byte("xyz"))

	_, err := Open(bytes.NewReader(data))
	if want := "failed to read attribute 0: unexpected EOF"; !isTruncation(err) || err.Error() != want {
		t.Fatalf("got %v, want %s", err, want)
	}
}

var errBoom = errors.New("boom")

type failingReader struct{}

func (failingReader) Read([]byte) (int, error) {
	return 0, errBoom
}

func TestOpenReaderError(t *testing.T) {
	data := fullClass().bytes()

	cases := map[int]string{
		0:  "failed to read magic number: boom",
		84: "failed to read attribute 0: boom",
	}

	for cut, message := range cases {
		_, err := Open(io.MultiReader(bytes.NewReader(data[:cut]), failingReader{}))
		if !errors.Is(err, errBoom) || err.Error() != message {
			t.Errorf("cut at %d: got %v, want %s", cut, err, message)
		}
	}
}
