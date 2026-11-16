package classfileparser

import (
	"errors"
	"reflect"
	"testing"
)

func pool(entries ...[]byte) []CpInfo {
	var cp []CpInfo
	for _, entry := range entries {
		info := CpInfo{Tag: entry[0], Info: entry[1:]}
		if entry[0] == 1 {
			info.Info = entry[3:]
		}
		cp = append(cp, info)
		if entry[0] == 5 || entry[0] == 6 {
			cp = append(cp, CpInfo{})
		}
	}
	return cp
}

func TestGetConstantPool(t *testing.T) {
	cf := &ClassFile{ConstantPool: pool(
		cpUtf8("Foo"),              // 1
		cpClass(1),                 // 2
		cpUtf8("bar"),              // 3
		cpUtf8("I"),                // 4
		cpNameAndType(3, 4),        // 5
		cpFieldref(2, 5),           // 6
		cpMethodref(2, 5),          // 7
		cpInterfaceMethodref(2, 5), // 8
		cpString(3),                // 9
		cpInteger(-1),              // 10
		cpFloat(1.5),               // 11
		cpLong(-2),                 // 12-13
		cpDouble(2.5),              // 14-15
		cpMethodHandle(1, 6),       // 16
		cpMethodHandle(6, 8),       // 17
		cpMethodType(4),            // 18
		cpDynamic(0, 5),            // 19
		cpInvokeDynamic(0, 5),      // 20
		cpModule(1),                // 21
		cpPackage(3),               // 22
		cpMethodHandle(5, 24),      // 23
		cpMethodref(2, 5),          // 24
		cpUtf8(""),                 // 25
	)}

	cp, err := cf.GetConstantPool()
	if err != nil {
		t.Fatalf("GetConstantPool failed: %v", err)
	}

	want := ConstantPool{
		1:  Utf8("Foo"),
		2:  Class("Foo"),
		3:  Utf8("bar"),
		4:  Utf8("I"),
		5:  NameAndType{Name: "bar", Type: "I"},
		6:  Fieldref{Class: "Foo", Name: "bar", Type: "I"},
		7:  Methodref{Class: "Foo", Name: "bar", Type: "I"},
		8:  InterfaceMethodref{Class: "Foo", Name: "bar", Type: "I"},
		9:  "bar",
		10: int32(-1),
		11: float32(1.5),
		12: int64(-2),
		14: float64(2.5),
		16: MethodHandle{Kind: "getField", Class: "Foo", Name: "bar", Type: "I"},
		17: MethodHandle{Kind: "invokeStatic", Class: "Foo", Name: "bar", Type: "I", IsInterface: true},
		18: MethodType("I"),
		19: Dynamic{},
		20: InvokeDynamic{},
		21: Module("Foo"),
		22: Package("bar"),
		23: MethodHandle{Kind: "invokeVirtual", Class: "Foo", Name: "bar", Type: "I"},
		24: Methodref{Class: "Foo", Name: "bar", Type: "I"},
		25: Utf8(""),
	}
	if !reflect.DeepEqual(cp, want) {
		t.Fatalf("got %#v\nwant %#v", cp, want)
	}
}

func TestGetConstantPoolMethodHandleKinds(t *testing.T) {
	kinds := []string{"getField", "getStatic", "putField", "putStatic", "invokeVirtual", "invokeStatic", "invokeSpecial", "newInvokeSpecial", "invokeInterface"}

	for i, kind := range kinds {
		cf := &ClassFile{ConstantPool: pool(
			cpUtf8("Foo"), cpClass(1), cpNameAndType(1, 1), cpMethodref(2, 3), cpMethodHandle(uint8(i+1), 4),
		)}
		cp, err := cf.GetConstantPool()
		if err != nil {
			t.Fatalf("kind %d: GetConstantPool failed: %v", i+1, err)
		}
		if got := cp[5].(MethodHandle).Kind; got != kind {
			t.Errorf("kind %d: got %s, want %s", i+1, got, kind)
		}
	}
}

func TestGetConstantPoolLargeIndex(t *testing.T) {
	entries := make([][]byte, 0, 40000)
	for len(entries) < 39999 {
		entries = append(entries, cpUtf8("x"))
	}
	entries[39998] = cpUtf8("Last")
	entries = append(entries, cpClass(39999))

	cf := &ClassFile{ConstantPool: pool(entries...)}
	cp, err := cf.GetConstantPool()
	if err != nil {
		t.Fatalf("GetConstantPool failed: %v", err)
	}
	if got := cp[40000]; got != Class("Last") {
		t.Fatalf("got %#v, want Class(\"Last\")", got)
	}
}

func TestGetConstantPoolMalformed(t *testing.T) {
	cases := map[string]struct {
		entries [][]byte
		message string
	}{
		"index 0": {
			[][]byte{cpClass(0)},
			"failed to read constant pool entry 1: malformed class file: invalid constant pool index: 0",
		},
		"index out of range": {
			[][]byte{cpClass(9)},
			"failed to read constant pool entry 1: malformed class file: invalid constant pool index: 9",
		},
		"Class to non-Utf8": {
			[][]byte{cpInteger(1), cpClass(1)},
			"failed to read constant pool entry 2: malformed class file: invalid constant pool tag at index 1: 3",
		},
		"String to invalid index": {
			[][]byte{cpString(9)},
			"failed to read constant pool entry 1: malformed class file: invalid constant pool index: 9",
		},
		"Fieldref to non-Class": {
			[][]byte{cpUtf8("Foo"), cpFieldref(1, 1)},
			"failed to read constant pool entry 2: malformed class file: invalid constant pool tag at index 1: 1",
		},
		"Methodref to non-NameAndType": {
			[][]byte{cpUtf8("Foo"), cpClass(1), cpMethodref(2, 1)},
			"failed to read constant pool entry 3: malformed class file: invalid constant pool tag at index 1: 1",
		},
		"InterfaceMethodref with invalid class name": {
			[][]byte{cpInterfaceMethodref(2, 3), cpClass(9), cpNameAndType(4, 4), cpUtf8("I")},
			"failed to read constant pool entry 1: malformed class file: invalid constant pool index: 9",
		},
		"Fieldref with invalid name": {
			[][]byte{cpFieldref(2, 3), cpClass(4), cpNameAndType(9, 4), cpUtf8("I")},
			"failed to read constant pool entry 1: malformed class file: invalid constant pool index: 9",
		},
		"NameAndType with invalid name": {
			[][]byte{cpNameAndType(9, 2), cpUtf8("I")},
			"failed to read constant pool entry 1: malformed class file: invalid constant pool index: 9",
		},
		"NameAndType with invalid descriptor": {
			[][]byte{cpNameAndType(2, 9), cpUtf8("I")},
			"failed to read constant pool entry 1: malformed class file: invalid constant pool index: 9",
		},
		"MethodHandle kind 0": {
			[][]byte{cpMethodHandle(0, 1)},
			"failed to read constant pool entry 1: malformed class file: invalid method handle kind: 0",
		},
		"MethodHandle kind 10": {
			[][]byte{cpMethodHandle(10, 1)},
			"failed to read constant pool entry 1: malformed class file: invalid method handle kind: 10",
		},
		"MethodHandle to non-reference": {
			[][]byte{cpUtf8("Foo"), cpMethodHandle(6, 1)},
			"failed to read constant pool entry 2: malformed class file: invalid constant pool tag at index 1: 1",
		},
		"MethodHandle to invalid reference": {
			[][]byte{cpMethodHandle(6, 2), cpMethodref(9, 9)},
			"failed to read constant pool entry 1: malformed class file: invalid constant pool index: 9",
		},
		"MethodType to invalid index": {
			[][]byte{cpMethodType(9)},
			"failed to read constant pool entry 1: malformed class file: invalid constant pool index: 9",
		},
		"Module to invalid index": {
			[][]byte{cpModule(9)},
			"failed to read constant pool entry 1: malformed class file: invalid constant pool index: 9",
		},
		"Package to invalid index": {
			[][]byte{cpPackage(9)},
			"failed to read constant pool entry 1: malformed class file: invalid constant pool index: 9",
		},
		"unknown tag": {
			[][]byte{{2}},
			"failed to read constant pool entry 1: malformed class file: unknown constant pool tag: 2",
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			cf := &ClassFile{ConstantPool: pool(c.entries...)}
			_, err := cf.GetConstantPool()
			if !errors.Is(err, ErrMalformed) || err.Error() != c.message {
				t.Fatalf("got %v, want %s", err, c.message)
			}
		})
	}
}

func TestGetConstant(t *testing.T) {
	cp := ConstantPool{1: Utf8("Foo"), 2: Class("Foo")}

	if got, err := getConstant[Class](cp, 2); err != nil || got != "Foo" {
		t.Fatalf("got %#v, %v", got, err)
	}

	cases := map[string]struct {
		index   uint16
		message string
	}{
		"missing index": {9, "malformed class file: invalid constant pool index: 9"},
		"wrong type":    {1, "malformed class file: invalid constant pool type at index 1: classfileparser.Utf8"},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := getConstant[Class](cp, c.index)
			if !errors.Is(err, ErrMalformed) || err.Error() != c.message {
				t.Fatalf("got %v, want %s", err, c.message)
			}
		})
	}
}
