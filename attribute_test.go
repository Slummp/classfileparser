package classfileparser

import (
	"errors"
	"reflect"
	"testing"
)

func attributeConstantPool() ConstantPool {
	return ConstantPool{
		1: Utf8("Code"),
		2: Utf8("Signature"),
		3: Utf8("TT;"),
		4: Class("Foo"),
		5: Utf8("Unknown"),
	}
}

func TestParseAttributesCode(t *testing.T) {
	nested := testAttribute{name: 2, info: u2(3)}
	exceptions := []ExceptionTableEntry{{StartPc: 0, EndPc: 1, HandlerPc: 1, CatchType: 4}}

	attrs, err := parseAttributes([]AttributeInfo{{AttributeNameIndex: 1, Info: codeInfo([]byte{0x00, 0xB1}, exceptions, nested)}}, attributeConstantPool())
	if err != nil {
		t.Fatalf("parseAttributes failed: %v", err)
	}

	want := []Attribute{Code{
		MaxStack:       1,
		MaxLocals:      1,
		Code:           []interface{}{Nop{}, Return{}},
		ExceptionTable: exceptions,
		Attributes:     []Attribute{Signature("TT;")},
	}}
	if !reflect.DeepEqual(attrs, want) {
		t.Fatalf("got %#v, want %#v", attrs, want)
	}
}

func TestParseAttributesNames(t *testing.T) {
	cases := map[string]Attribute{
		"ConstantValue":                        ConstantValue{},
		"Deprecated":                           Deprecated{},
		"Exceptions":                           Exceptions{},
		"InnerClasses":                         InnerClasses{},
		"LineNumberTable":                      LineNumberTable{},
		"LocalVariableTable":                   LocalVariableTable{},
		"LocalVariableTypeTable":               LocalVariableTypeTable{},
		"MethodParameters":                     MethodParameters{},
		"RuntimeVisibleAnnotations":            RuntimeVisibleAnnotations{},
		"RuntimeInvisibleAnnotations":          RuntimeInvisibleAnnotations{},
		"RuntimeVisibleParameterAnnotations":   RuntimeVisibleParameterAnnotations{},
		"RuntimeInvisibleParameterAnnotations": RuntimeInvisibleParameterAnnotations{},
		"SourceFile":                           SourceFile{},
		"SourceDebugExtension":                 SourceDebugExtension{},
		"Signature":                            Signature("TT;"),
		"StackMapTable":                        StackMapTable{},
		"Synthetic":                            Synthetic{},
		"EnclosingMethod":                      EnclosingMethod{},
		"BootstrapMethods":                     BootstrapMethods{},
		"Module":                               ModuleInfo{},
		"ModulePackages":                       ModulePackages{},
		"NestHost":                             NestHost{},
		"NestMembers":                          NestMembers{},
		"PermittedSubclasses":                  PermittedSubclasses{},
		"Record":                               Record{},
	}

	for name, want := range cases {
		t.Run(name, func(t *testing.T) {
			cp := attributeConstantPool()
			cp[10] = Utf8(name)
			attrs, err := parseAttributes([]AttributeInfo{{AttributeNameIndex: 10, Info: u2(3)}}, cp)
			if err != nil {
				t.Fatalf("parseAttributes failed: %v", err)
			}
			if !reflect.DeepEqual(attrs, []Attribute{want}) {
				t.Fatalf("got %#v, want %#v", attrs, want)
			}
		})
	}
}

func TestParseAttributesUnknownName(t *testing.T) {
	attrs, err := parseAttributes([]AttributeInfo{{AttributeNameIndex: 5, Info: []byte("anything")}}, attributeConstantPool())
	if err != nil || len(attrs) != 0 {
		t.Fatalf("got %#v, %v, want no attribute", attrs, err)
	}
}

func TestParseAttributesMalformed(t *testing.T) {
	cases := map[string]struct {
		attribute AttributeInfo
		message   string
	}{
		"missing name": {
			AttributeInfo{AttributeNameIndex: 9},
			"failed to read attribute name: malformed class file: invalid constant pool index: 9",
		},
		"name of wrong type": {
			AttributeInfo{AttributeNameIndex: 4},
			"failed to read attribute name: malformed class file: invalid constant pool type at index 4: classfileparser.Class",
		},
		"empty Code": {
			AttributeInfo{AttributeNameIndex: 1},
			"failed to read attribute Code: malformed class file: exceeds attribute length: 0",
		},
		"bytecode shorter than code length": {
			AttributeInfo{AttributeNameIndex: 1, Info: cat(u2(1), u2(1), u4(5))},
			"failed to read attribute Code: failed to read instruction at offset 0: malformed class file: exceeds attribute length: 8",
		},
		"truncated exception table": {
			AttributeInfo{AttributeNameIndex: 1, Info: cat(u2(1), u2(1), u4(1), u1(0xB1), u2(1), u2(0))},
			"failed to read attribute Code: malformed class file: exceeds attribute length: 13",
		},
		"nested attribute longer than Code": {
			AttributeInfo{AttributeNameIndex: 1, Info: cat(u2(1), u2(1), u4(1), u1(0xB1), u2(0), u2(1), u2(2), u4(100))},
			"failed to read attribute Code: malformed class file: invalid attribute length: 100",
		},
		"invalid nested attribute": {
			AttributeInfo{AttributeNameIndex: 1, Info: codeInfo([]byte{0xB1}, nil, testAttribute{name: 9})},
			"failed to read attribute Code: failed to read attribute name: malformed class file: invalid constant pool index: 9",
		},
		"empty Signature": {
			AttributeInfo{AttributeNameIndex: 2},
			"failed to read attribute Signature: malformed class file: exceeds attribute length: 0",
		},
		"Signature of wrong type": {
			AttributeInfo{AttributeNameIndex: 2, Info: u2(4)},
			"failed to read attribute Signature: malformed class file: invalid constant pool type at index 4: classfileparser.Class",
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := parseAttributes([]AttributeInfo{c.attribute}, attributeConstantPool())
			if !errors.Is(err, ErrMalformed) || isTruncation(err) || err.Error() != c.message {
				t.Fatalf("got %v, want %s", err, c.message)
			}
		})
	}
}
