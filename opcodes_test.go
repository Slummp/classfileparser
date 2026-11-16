package classfileparser

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func opcodeConstantPool() ConstantPool {
	return ConstantPool{
		1: Utf8("Code"),
		2: Class("Foo"),
		3: Fieldref{Class: "Foo", Name: "bar", Type: "I"},
		4: Methodref{Class: "Foo", Name: "run", Type: "()V"},
		5: InterfaceMethodref{Class: "Bar", Name: "run", Type: "()V"},
		6: InvokeDynamic{},
		7: int32(42),
		8: int64(7),
	}
}

func opcodeCases() []struct {
	bytecode []byte
	want     interface{}
} {
	return []struct {
		bytecode []byte
		want     interface{}
	}{
		{[]byte{0x00}, Nop{}},
		{[]byte{0x01}, AconstNull{}},
		{[]byte{0x02}, IconstM1{}},
		{[]byte{0x03}, Iconst0{}},
		{[]byte{0x04}, Iconst1{}},
		{[]byte{0x05}, Iconst2{}},
		{[]byte{0x06}, Iconst3{}},
		{[]byte{0x07}, Iconst4{}},
		{[]byte{0x08}, Iconst5{}},
		{[]byte{0x09}, Lconst0{}},
		{[]byte{0x0A}, Lconst1{}},
		{[]byte{0x0B}, Fconst0{}},
		{[]byte{0x0C}, Fconst1{}},
		{[]byte{0x0D}, Fconst2{}},
		{[]byte{0x0E}, Dconst0{}},
		{[]byte{0x0F}, Dconst1{}},
		{[]byte{0x1A}, Iload0{}},
		{[]byte{0x1B}, Iload1{}},
		{[]byte{0x1C}, Iload2{}},
		{[]byte{0x1D}, Iload3{}},
		{[]byte{0x1E}, Lload0{}},
		{[]byte{0x1F}, Lload1{}},
		{[]byte{0x20}, Lload2{}},
		{[]byte{0x21}, Lload3{}},
		{[]byte{0x22}, Fload0{}},
		{[]byte{0x23}, Fload1{}},
		{[]byte{0x24}, Fload2{}},
		{[]byte{0x25}, Fload3{}},
		{[]byte{0x26}, Dload0{}},
		{[]byte{0x27}, Dload1{}},
		{[]byte{0x28}, Dload2{}},
		{[]byte{0x29}, Dload3{}},
		{[]byte{0x2A}, Aload0{}},
		{[]byte{0x2B}, Aload1{}},
		{[]byte{0x2C}, Aload2{}},
		{[]byte{0x2D}, Aload3{}},
		{[]byte{0x2E}, Iaload{}},
		{[]byte{0x2F}, Laload{}},
		{[]byte{0x30}, Faload{}},
		{[]byte{0x31}, Daload{}},
		{[]byte{0x32}, Aaload{}},
		{[]byte{0x33}, Baload{}},
		{[]byte{0x34}, Caload{}},
		{[]byte{0x35}, Saload{}},
		{[]byte{0x3B}, Istore0{}},
		{[]byte{0x3C}, Istore1{}},
		{[]byte{0x3D}, Istore2{}},
		{[]byte{0x3E}, Istore3{}},
		{[]byte{0x3F}, Lstore0{}},
		{[]byte{0x40}, Lstore1{}},
		{[]byte{0x41}, Lstore2{}},
		{[]byte{0x42}, Lstore3{}},
		{[]byte{0x43}, Fstore0{}},
		{[]byte{0x44}, Fstore1{}},
		{[]byte{0x45}, Fstore2{}},
		{[]byte{0x46}, Fstore3{}},
		{[]byte{0x47}, Dstore0{}},
		{[]byte{0x48}, Dstore1{}},
		{[]byte{0x49}, Dstore2{}},
		{[]byte{0x4A}, Dstore3{}},
		{[]byte{0x4B}, Astore0{}},
		{[]byte{0x4C}, Astore1{}},
		{[]byte{0x4D}, Astore2{}},
		{[]byte{0x4E}, Astore3{}},
		{[]byte{0x4F}, Iastore{}},
		{[]byte{0x50}, Lastore{}},
		{[]byte{0x51}, Fastore{}},
		{[]byte{0x52}, Dastore{}},
		{[]byte{0x53}, Aastore{}},
		{[]byte{0x54}, Bastore{}},
		{[]byte{0x55}, Castore{}},
		{[]byte{0x56}, Sastore{}},
		{[]byte{0x57}, Pop{}},
		{[]byte{0x58}, Pop2{}},
		{[]byte{0x59}, Dup{}},
		{[]byte{0x5A}, DupX1{}},
		{[]byte{0x5B}, DupX2{}},
		{[]byte{0x5C}, Dup2{}},
		{[]byte{0x5D}, Dup2X1{}},
		{[]byte{0x5E}, Dup2X2{}},
		{[]byte{0x5F}, Swap{}},
		{[]byte{0x60}, Iadd{}},
		{[]byte{0x61}, Ladd{}},
		{[]byte{0x62}, Fadd{}},
		{[]byte{0x63}, Dadd{}},
		{[]byte{0x64}, Isub{}},
		{[]byte{0x65}, Lsub{}},
		{[]byte{0x66}, Fsub{}},
		{[]byte{0x67}, Dsub{}},
		{[]byte{0x68}, Imul{}},
		{[]byte{0x69}, Lmul{}},
		{[]byte{0x6A}, Fmul{}},
		{[]byte{0x6B}, Dmul{}},
		{[]byte{0x6C}, Idiv{}},
		{[]byte{0x6D}, Ldiv{}},
		{[]byte{0x6E}, Fdiv{}},
		{[]byte{0x6F}, Ddiv{}},
		{[]byte{0x70}, Irem{}},
		{[]byte{0x71}, Lrem{}},
		{[]byte{0x72}, Frem{}},
		{[]byte{0x73}, Drem{}},
		{[]byte{0x74}, Ineg{}},
		{[]byte{0x75}, Lneg{}},
		{[]byte{0x76}, Fneg{}},
		{[]byte{0x77}, Dneg{}},
		{[]byte{0x78}, Ishl{}},
		{[]byte{0x79}, Lshl{}},
		{[]byte{0x7A}, Ishr{}},
		{[]byte{0x7B}, Lshr{}},
		{[]byte{0x7C}, Iushr{}},
		{[]byte{0x7D}, Lushr{}},
		{[]byte{0x7E}, Iand{}},
		{[]byte{0x7F}, Land{}},
		{[]byte{0x80}, Ior{}},
		{[]byte{0x81}, Lor{}},
		{[]byte{0x82}, Ixor{}},
		{[]byte{0x83}, Lxor{}},
		{[]byte{0x85}, I2l{}},
		{[]byte{0x86}, I2f{}},
		{[]byte{0x87}, I2d{}},
		{[]byte{0x88}, L2i{}},
		{[]byte{0x89}, L2f{}},
		{[]byte{0x8A}, L2d{}},
		{[]byte{0x8B}, F2i{}},
		{[]byte{0x8C}, F2l{}},
		{[]byte{0x8D}, F2d{}},
		{[]byte{0x8E}, D2i{}},
		{[]byte{0x8F}, D2l{}},
		{[]byte{0x90}, D2f{}},
		{[]byte{0x91}, I2b{}},
		{[]byte{0x92}, I2c{}},
		{[]byte{0x93}, I2s{}},
		{[]byte{0x94}, Lcmp{}},
		{[]byte{0x95}, Fcmpl{}},
		{[]byte{0x96}, Fcmpg{}},
		{[]byte{0x97}, Dcmpl{}},
		{[]byte{0x98}, Dcmpg{}},
		{[]byte{0xAC}, Ireturn{}},
		{[]byte{0xAD}, Lreturn{}},
		{[]byte{0xAE}, Freturn{}},
		{[]byte{0xAF}, Dreturn{}},
		{[]byte{0xB0}, Areturn{}},
		{[]byte{0xB1}, Return{}},
		{[]byte{0xBE}, Arraylength{}},
		{[]byte{0xBF}, Athrow{}},
		{[]byte{0xC2}, Monitorenter{}},
		{[]byte{0xC3}, Monitorexit{}},
		{[]byte{0x10, 0xFF}, Bipush{Byte: 0xFF}},
		{[]byte{0x11, 0xFF, 0xFE}, Sipush{Short: -2}},
		{[]byte{0x12, 0x07}, Ldc(int32(42))},
		{[]byte{0x13, 0x00, 0x07}, LdcW(int32(42))},
		{[]byte{0x14, 0x00, 0x08}, Ldc2W(int64(7))},
		{[]byte{0x15, 0x05}, Iload{LocalIndex: 5}},
		{[]byte{0x16, 0x05}, Lload{LocalIndex: 5}},
		{[]byte{0x17, 0x05}, Fload{LocalIndex: 5}},
		{[]byte{0x18, 0x05}, Dload{LocalIndex: 5}},
		{[]byte{0x19, 0x05}, Aload{LocalIndex: 5}},
		{[]byte{0x36, 0x05}, Istore{LocalIndex: 5}},
		{[]byte{0x37, 0x05}, Lstore{LocalIndex: 5}},
		{[]byte{0x38, 0x05}, Fstore{LocalIndex: 5}},
		{[]byte{0x39, 0x05}, Dstore{LocalIndex: 5}},
		{[]byte{0x3A, 0x05}, Astore{LocalIndex: 5}},
		{[]byte{0xA9, 0x05}, Ret{LocalIndex: 5}},
		{[]byte{0x84, 0x05, 0xFF}, Iinc{LocalIndex: 5, Const: -1}},
		{[]byte{0x99, 0xFF, 0xFE}, Ifeq{Offset: -2}},
		{[]byte{0x9A, 0xFF, 0xFE}, Ifne{Offset: -2}},
		{[]byte{0x9B, 0xFF, 0xFE}, Iflt{Offset: -2}},
		{[]byte{0x9C, 0xFF, 0xFE}, Ifge{Offset: -2}},
		{[]byte{0x9D, 0xFF, 0xFE}, Ifgt{Offset: -2}},
		{[]byte{0x9E, 0xFF, 0xFE}, Ifle{Offset: -2}},
		{[]byte{0x9F, 0xFF, 0xFE}, IfIcmpeq{Offset: -2}},
		{[]byte{0xA0, 0xFF, 0xFE}, IfIcmpne{Offset: -2}},
		{[]byte{0xA1, 0xFF, 0xFE}, IfIcmplt{Offset: -2}},
		{[]byte{0xA2, 0xFF, 0xFE}, IfIcmpge{Offset: -2}},
		{[]byte{0xA3, 0xFF, 0xFE}, IfIcmpgt{Offset: -2}},
		{[]byte{0xA4, 0xFF, 0xFE}, IfIcmple{Offset: -2}},
		{[]byte{0xA5, 0xFF, 0xFE}, IfAcmpeq{Offset: -2}},
		{[]byte{0xA6, 0xFF, 0xFE}, IfAcmpne{Offset: -2}},
		{[]byte{0xA7, 0xFF, 0xFE}, Goto{Offset: -2}},
		{[]byte{0xA8, 0xFF, 0xFE}, Jsr{Offset: -2}},
		{[]byte{0xC6, 0xFF, 0xFE}, Ifnull{Offset: -2}},
		{[]byte{0xC7, 0xFF, 0xFE}, Ifnonnull{Offset: -2}},
		{[]byte{0xAA, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2}, Tableswitch{Default: 1, Low: 0, High: 0, Offsets: []int32{2}}},
		{[]byte{0xAB, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 5, 0, 0, 0, 2}, Lookupswitch{Default: 1, Pairs: []struct{ Match, Offset int32 }{{Match: 5, Offset: 2}}}},
		{[]byte{0xB2, 0x00, 0x03}, Getstatic{Fieldref: Fieldref{Class: "Foo", Name: "bar", Type: "I"}}},
		{[]byte{0xB3, 0x00, 0x03}, Putstatic{Fieldref: Fieldref{Class: "Foo", Name: "bar", Type: "I"}}},
		{[]byte{0xB4, 0x00, 0x03}, Getfield{Fieldref: Fieldref{Class: "Foo", Name: "bar", Type: "I"}}},
		{[]byte{0xB5, 0x00, 0x03}, Putfield{Fieldref: Fieldref{Class: "Foo", Name: "bar", Type: "I"}}},
		{[]byte{0xB6, 0x00, 0x04}, Invokevirtual{Methodref: Methodref{Class: "Foo", Name: "run", Type: "()V"}}},
		{[]byte{0xB7, 0x00, 0x04}, Invokespecial{Methodref: Methodref{Class: "Foo", Name: "run", Type: "()V"}}},
		{[]byte{0xB7, 0x00, 0x05}, Invokespecial{Methodref: Methodref{Class: "Bar", Name: "run", Type: "()V"}, IsInterface: true}},
		{[]byte{0xB8, 0x00, 0x04}, Invokestatic{Methodref: Methodref{Class: "Foo", Name: "run", Type: "()V"}}},
		{[]byte{0xB8, 0x00, 0x05}, Invokestatic{Methodref: Methodref{Class: "Bar", Name: "run", Type: "()V"}, IsInterface: true}},
		{[]byte{0xB9, 0x00, 0x05, 0x02, 0x00}, Invokeinterface{InterfaceMethodref: InterfaceMethodref{Class: "Bar", Name: "run", Type: "()V"}, Count: 2}},
		{[]byte{0xBA, 0x00, 0x06, 0x00, 0x00}, Invokedynamic{InvokeDynamic: InvokeDynamic{}}},
		{[]byte{0xBB, 0x00, 0x02}, New("Foo")},
		{[]byte{0xBC, 0x0A}, Newarray{Type: 10}},
		{[]byte{0xBD, 0x00, 0x02}, Anewarray("Foo")},
		{[]byte{0xC0, 0x00, 0x02}, Checkcast("Foo")},
		{[]byte{0xC1, 0x00, 0x02}, Instanceof("Foo")},
		{[]byte{0xC4, 0x15, 0x01, 0x00}, Wide{OpCode: 0x15, LocalIndex: 256}},
		{[]byte{0xC4, 0x84, 0x01, 0x00, 0xFF, 0xFF}, Wide{OpCode: 0x84, LocalIndex: 256, Const: 0xFFFF}},
		{[]byte{0xC5, 0x00, 0x02, 0x03}, Multianewarray{Class: "Foo", Dimension: 3}},
		{[]byte{0xC8, 0xFF, 0xFF, 0xFF, 0xFE}, GotoW{Offset: -2}},
		{[]byte{0xC9, 0xFF, 0xFF, 0xFF, 0xFE}, JsrW{Offset: -2}},
	}
}

func parseInstructions(bytecode []byte) ([]interface{}, error) {
	attrs, err := parseAttributes([]AttributeInfo{{AttributeNameIndex: 1, Info: codeInfo(bytecode, nil)}}, opcodeConstantPool())
	if err != nil {
		return nil, err
	}
	return attrs[0].(Code).Code, nil
}

func TestOpcodes(t *testing.T) {
	seen := map[byte]bool{}
	for _, c := range opcodeCases() {
		seen[c.bytecode[0]] = true
		t.Run(fmt.Sprintf("0x%02X", c.bytecode), func(t *testing.T) {
			got, err := parseInstructions(c.bytecode)
			if err != nil {
				t.Fatalf("parse failed: %v", err)
			}
			if want := []interface{}{c.want}; !reflect.DeepEqual(got, want) {
				t.Fatalf("got %#v, want %#v", got, want)
			}
		})
	}
	if len(seen) != 202 {
		t.Fatalf("%d opcodes tested, want 202", len(seen))
	}
}

func TestOpcodesTruncated(t *testing.T) {
	for _, c := range opcodeCases() {
		if len(c.bytecode) == 1 {
			continue
		}
		t.Run(fmt.Sprintf("0x%02X", c.bytecode), func(t *testing.T) {
			_, err := parseInstructions(c.bytecode[:len(c.bytecode)-1])
			prefix := "failed to read attribute Code: failed to read instruction at offset 0: malformed class file: "
			if !errors.Is(err, ErrMalformed) || !strings.HasPrefix(err.Error(), prefix) {
				t.Fatalf("got %v, want an error starting with %s", err, prefix)
			}
		})
	}
}

func TestOpcodesInvalidConstant(t *testing.T) {
	for _, c := range opcodeCases() {
		op := c.bytecode[0]
		if op != 0x12 && op != 0x13 && op != 0x14 && !(op >= 0xB2 && op <= 0xBB) && op != 0xBD && op != 0xC0 && op != 0xC1 && op != 0xC5 {
			continue
		}
		t.Run(fmt.Sprintf("0x%02X", op), func(t *testing.T) {
			bytecode := append([]byte{}, c.bytecode...)
			if op == 0x12 {
				bytecode[1] = 0x09
			} else {
				bytecode[1], bytecode[2] = 0x00, 0x09
			}
			_, err := parseInstructions(bytecode)
			want := "failed to read attribute Code: failed to read instruction at offset 0: malformed class file: invalid constant pool index: 9"
			if !errors.Is(err, ErrMalformed) || err.Error() != want {
				t.Fatalf("got %v, want %s", err, want)
			}
		})
	}
}

func TestUnknownOpcodes(t *testing.T) {
	for op := 0xCA; op <= 0xFF; op++ {
		_, err := parseInstructions([]byte{byte(op)})
		want := fmt.Sprintf("failed to read attribute Code: failed to read instruction at offset 0: malformed class file: unknown opcode: 0x%X", op)
		if !errors.Is(err, ErrMalformed) || err.Error() != want {
			t.Errorf("got %v, want %s", err, want)
		}
	}
}
