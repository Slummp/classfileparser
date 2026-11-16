package classfileparser

import (
	"bytes"
	"errors"
	"reflect"
	"testing"
)

func TestGetClassFile(t *testing.T) {
	c := testClass{
		minor: 3,
		major: 61,
		pool: [][]byte{
			cpUtf8("Foo"),                      // 1
			cpClass(1),                         // 2
			cpUtf8("java/lang/Object"),         // 3
			cpClass(3),                         // 4
			cpUtf8("java/lang/Runnable"),       // 5
			cpClass(5),                         // 6
			cpUtf8("count"),                    // 7
			cpUtf8("I"),                        // 8
			cpUtf8("run"),                      // 9
			cpUtf8("([[ILjava/lang/String;)V"), // 10
			cpUtf8("Signature"),                // 11
			cpUtf8("TT;"),                      // 12
		},
		access:     0x0001,
		this:       2,
		super:      4,
		interfaces: []uint16{6},
		fields:     []testMember{{access: 0x0002, name: 7, descriptor: 8, attributes: []testAttribute{{name: 11, info: u2(12)}}}},
		methods:    []testMember{{access: 0x0008, name: 9, descriptor: 10}},
		attributes: []testAttribute{{name: 11, info: u2(12)}},
	}
	cf, err := Open(bytes.NewReader(c.bytes()))
	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}

	got, err := cf.GetClassFile()
	if err != nil {
		t.Fatalf("GetClassFile failed: %v", err)
	}

	want := &ClassStruct{
		Access:     []string{"ACC_PUBLIC"},
		ThisClass:  "Foo",
		SuperClass: "java/lang/Object",
		Interfaces: []string{"java/lang/Runnable"},
		Fields: []Field{{
			Access:     []string{"ACC_PRIVATE"},
			Name:       "count",
			Type:       "I",
			Attributes: []Attribute{Signature("TT;")},
		}},
		Methods: []Method{{
			Access:      []string{"ACC_STATIC"},
			Name:        "run",
			ReturnType:  "V",
			ParamsTypes: []string{"[[I", "Ljava/lang/String;"},
		}},
		Attributes: []Attribute{Signature("TT;")},
	}
	want.Version.MinorVersion = 3
	want.Version.MajorVersion = 61
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v\nwant %#v", got, want)
	}
}

func TestGetClassFileWithoutSuperClass(t *testing.T) {
	cf, err := Open(bytes.NewReader(minimalClass().bytes()))
	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}

	got, err := cf.GetClassFile()
	if err != nil {
		t.Fatalf("GetClassFile failed: %v", err)
	}
	if got.ThisClass != "Foo" || got.SuperClass != "" {
		t.Fatalf("got this %q and super %q, want \"Foo\" and \"\"", got.ThisClass, got.SuperClass)
	}
}

func TestGetClassFileMalformed(t *testing.T) {
	invalidAttribute := []AttributeInfo{{AttributeNameIndex: 9}}

	cases := map[string]struct {
		edit    func(cf *ClassFile)
		message string
	}{
		"constant pool": {
			func(cf *ClassFile) { cf.ConstantPool = pool(cpClass(9)) },
			"failed to read constant pool entry 1: malformed class file: invalid constant pool index: 9",
		},
		"interface": {
			func(cf *ClassFile) { cf.Interfaces = []uint16{9} },
			"failed to read interface 0: malformed class file: invalid constant pool index: 9",
		},
		"field name": {
			func(cf *ClassFile) { cf.Fields = []FieldInfo{{NameIndex: 9, DescriptorIndex: 4}} },
			"failed to read field 0: malformed class file: invalid constant pool index: 9",
		},
		"field descriptor": {
			func(cf *ClassFile) { cf.Fields = []FieldInfo{{NameIndex: 3, DescriptorIndex: 9}} },
			"failed to read field count: malformed class file: invalid constant pool index: 9",
		},
		"field attribute": {
			func(cf *ClassFile) {
				cf.Fields = []FieldInfo{{NameIndex: 3, DescriptorIndex: 4, Attributes: invalidAttribute}}
			},
			"failed to read field count: failed to read attribute name: malformed class file: invalid constant pool index: 9",
		},
		"method name": {
			func(cf *ClassFile) { cf.Methods = []MethodInfo{{NameIndex: 9, DescriptorIndex: 6}} },
			"failed to read method 0: malformed class file: invalid constant pool index: 9",
		},
		"method descriptor": {
			func(cf *ClassFile) { cf.Methods = []MethodInfo{{NameIndex: 5, DescriptorIndex: 9}} },
			"failed to read method run: malformed class file: invalid constant pool index: 9",
		},
		"method attribute": {
			func(cf *ClassFile) {
				cf.Methods = []MethodInfo{{NameIndex: 5, DescriptorIndex: 6, Attributes: invalidAttribute}}
			},
			"failed to read method run: failed to read attribute name: malformed class file: invalid constant pool index: 9",
		},
		"invalid method descriptor": {
			func(cf *ClassFile) { cf.Methods = []MethodInfo{{NameIndex: 5, DescriptorIndex: 4}} },
			"failed to read method run: malformed class file: invalid method descriptor: I",
		},
		"class attribute": {
			func(cf *ClassFile) { cf.Attributes = invalidAttribute },
			"failed to read attribute name: malformed class file: invalid constant pool index: 9",
		},
		"this class": {
			func(cf *ClassFile) { cf.ThisClass = 9 },
			"failed to read this class: malformed class file: invalid constant pool index: 9",
		},
		"super class": {
			func(cf *ClassFile) { cf.SuperClass = 9 },
			"failed to read super class: malformed class file: invalid constant pool index: 9",
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			cf := &ClassFile{
				ConstantPool: pool(cpUtf8("Foo"), cpClass(1), cpUtf8("count"), cpUtf8("I"), cpUtf8("run"), cpUtf8("()V")),
				ThisClass:    2,
			}
			c.edit(cf)
			_, err := cf.GetClassFile()
			if !errors.Is(err, ErrMalformed) || err.Error() != c.message {
				t.Fatalf("got %v, want %s", err, c.message)
			}
		})
	}
}

func TestReadSignature(t *testing.T) {
	cases := map[string]struct {
		params     []string
		returnType string
	}{
		"()V":                                    {nil, "V"},
		"(I)V":                                   {[]string{"I"}, "V"},
		"(BCDFIJSZ)V":                            {[]string{"B", "C", "D", "F", "I", "J", "S", "Z"}, "V"},
		"(Ljava/lang/String;)I":                  {[]string{"Ljava/lang/String;"}, "I"},
		"([I[[J[[[Ljava/a/B;)[[I":                {[]string{"[I", "[[J", "[[[Ljava/a/B;"}, "[[I"},
		"(ILjava/util/List;Z)Ljava/lang/Object;": {[]string{"I", "Ljava/util/List;", "Z"}, "Ljava/lang/Object;"},
	}

	for descriptor, c := range cases {
		t.Run(descriptor, func(t *testing.T) {
			params, returnType, err := readSignature(descriptor)
			if err != nil || !reflect.DeepEqual(params, c.params) || returnType != c.returnType {
				t.Fatalf("got %q, %q, %v, want %q, %q", params, returnType, err, c.params, c.returnType)
			}
		})
	}

	_, _, err := readSignature("I")
	if want := "malformed class file: invalid method descriptor: I"; !errors.Is(err, ErrMalformed) || err.Error() != want {
		t.Fatalf("got %v, want %s", err, want)
	}
}
