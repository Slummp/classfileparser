package classfileparser

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

// Attribute represents a generic attribute found in a JVM class file.
type Attribute interface{}

// Code holds the bytecode and the method metadata.
type Code struct {
	MaxStack       uint16
	MaxLocals      uint16
	Code           []interface{}
	ExceptionTable []ExceptionTableEntry
	Attributes     []Attribute
}

// ConstantValue is the attribute for a constant field.
type ConstantValue struct {
	ConstantValueIndex uint16
}

// Deprecated marks a method as obsolete.
type Deprecated struct{}

// Exceptions lists the checked exceptions thrown by a method.
type Exceptions struct {
	NumberOfExceptions  uint16
	ExceptionIndexTable []uint16
}

// InnerClasses stores information about nested classes.
type InnerClasses struct {
	NumberOfClasses uint16
	InnerClassInfo  []InnerClassInfo
}

// LineNumberTable links source lines to bytecode positions.
type LineNumberTable struct {
	LineNumberTableLength uint16
	LineNumberTable       []LineNumberTableEntry
}

// LocalVariableTable stores information about local variables.
type LocalVariableTable struct {
	LocalVariableTableLength uint16
	LocalVariableTable       []LocalVariableEntry
}

// LocalVariableTypeTable stores generic type information for local variables.
type LocalVariableTypeTable struct {
	LocalVariableTypeTableLength uint16
	LocalVariableTypeTable       []LocalVariableTypeEntry
}

// MethodParameters lists information about the method parameters.
type MethodParameters struct {
	MethodParametersCount uint16
	MethodParameters      []MethodParameter
}

// RuntimeVisibleAnnotations contains runtime-visible annotations.
type RuntimeVisibleAnnotations struct {
	NumAnnotations uint16
	Annotations    []Annotation
}

// RuntimeInvisibleAnnotations contains runtime-invisible annotations.
type RuntimeInvisibleAnnotations struct {
	NumAnnotations uint16
	Annotations    []Annotation
}

// RuntimeVisibleParameterAnnotations stores runtime-visible parameter annotations.
type RuntimeVisibleParameterAnnotations struct {
	NumParameters        uint16
	ParameterAnnotations []ParameterAnnotation
}

// RuntimeInvisibleParameterAnnotations stores runtime-invisible parameter annotations.
type RuntimeInvisibleParameterAnnotations struct {
	NumParameters        uint16
	ParameterAnnotations []ParameterAnnotation
}

// SourceFile holds the source file name.
type SourceFile struct {
	SourcefileIndex uint16
}

// SourceDebugExtension holds source debugging data.
type SourceDebugExtension struct {
	DebugExtension []byte
}

// Signature stores the signature of a class, method, or field.
type Signature string

// StackMapTable stores stack and local variable verification data.
type StackMapTable struct {
	NumberOfEntries uint16
	Entries         []StackMapFrame
}

// Synthetic marks an element as synthetic.
type Synthetic struct{}

// EnclosingMethod records the enclosing method.
type EnclosingMethod struct {
	ClassIndex  uint16
	MethodIndex uint16
}

// BootstrapMethods lists bootstrap methods.
type BootstrapMethods struct {
	NumBootstrapMethods uint16
	BootstrapMethods    []BootstrapMethod
}

// ModuleInfo stores module metadata.
type ModuleInfo struct {
	NameIndex uint16
	Flags     uint16
	Requires  []ModuleRequire
}

// ModulePackages stores the packages declared in a module.
type ModulePackages struct {
	NumberOfPackages uint16
	PackageIndex     []uint16
}

// NestHost stores the host class.
type NestHost struct {
	HostClassIndex uint16
}

// NestMembers lists the member classes.
type NestMembers struct {
	NumberOfMembers uint16
	ClassIndex      []uint16
}

// PermittedSubclasses lists the permitted subclasses.
type PermittedSubclasses struct {
	NumberOfSubclasses uint16
	SubclassIndex      []uint16
}

// ExceptionTableEntry describes the additional entries needed for attributes.
type ExceptionTableEntry struct {
	StartPc   uint16
	EndPc     uint16
	HandlerPc uint16
	CatchType uint16
}

// InnerClassInfo describes an InnerClasses entry.
type InnerClassInfo struct {
	InnerClassIndex       uint16
	OuterClassIndex       uint16
	InnerNameIndex        uint16
	InnerClassAccessFlags uint16
}

// LineNumberTableEntry links a bytecode position to a source line.
type LineNumberTableEntry struct {
	StartPc    uint16
	LineNumber uint16
}

// LocalVariableEntry details a local variable in the dedicated table.
type LocalVariableEntry struct {
	StartPc        uint16
	Length         uint16
	NameIndex      uint16
	SignatureIndex uint16
	Index          uint16
}

// LocalVariableTypeEntry specifies the generic type of a local variable.
type LocalVariableTypeEntry struct {
	StartPc        uint16
	Length         uint16
	NameIndex      uint16
	SignatureIndex uint16
	Index          uint16
}

// MethodParameter describes a method parameter.
type MethodParameter struct {
	NameIndex   uint16
	AccessFlags uint16
}

// Annotation represents a JVM annotation.
type Annotation struct {
	TypeIndex            uint16
	NumElementValuePairs uint16
	ElementValuePairs    []ElementValuePair
}

// ParameterAnnotation groups the annotations of a method parameter.
type ParameterAnnotation struct {
	NumAnnotations uint16
	Annotations    []Annotation
}

// ElementValuePair associates an element name with an annotation value.
type ElementValuePair struct {
	ElementNameIndex uint16
	Value            ElementValue
}

// ElementValue holds an annotation value.
type ElementValue struct {
	Tag   uint8
	Value uint16
}

// StackMapFrame describes a frame of the StackMapTable.
type StackMapFrame struct {
	FrameType   uint8
	OffsetDelta uint16
	StackItems  []StackItem
}

// StackItem represents a stack element inside the StackMapTable.
type StackItem struct {
	Tag   uint8
	Value uint16
}

// BootstrapMethod references a bootstrap method and its arguments.
type BootstrapMethod struct {
	MethodRefIndex uint16
	ArgumentsCount uint16
	Arguments      []uint16
}

// ModuleRequire describes a requires entry of a ModuleInfo.
type ModuleRequire struct {
	NameIndex    uint16
	Flags        uint16
	VersionIndex uint16
}

func parseAttributes(attributes []AttributeInfo, cp ConstantPool) ([]Attribute, error) {
	var attr []Attribute
	for _, a := range attributes {
		reader := &binaryReader{Reader: bytes.NewReader(a.Info)}
		name, err := getConstant[Utf8](cp, a.AttributeNameIndex)
		if err != nil {
			return nil, fmt.Errorf("failed to read attribute name: %w", err)
		}
		switch name {
		case "Code":
			var code Code

			reader.read(&code.MaxStack)
			reader.read(&code.MaxLocals)

			var codeLength uint32
			reader.read(&codeLength)
			start := reader.Len()
			for start-reader.Len() < int(codeLength) {
				pc := start - reader.Len()
				var opcode uint8
				reader.read(&opcode)

				switch opcode {
				case 0x00:
					code.Code = append(code.Code, Nop{})
				case 0x01:
					code.Code = append(code.Code, AconstNull{})
				case 0x02:
					code.Code = append(code.Code, IconstM1{})
				case 0x03:
					code.Code = append(code.Code, Iconst0{})
				case 0x04:
					code.Code = append(code.Code, Iconst1{})
				case 0x05:
					code.Code = append(code.Code, Iconst2{})
				case 0x06:
					code.Code = append(code.Code, Iconst3{})
				case 0x07:
					code.Code = append(code.Code, Iconst4{})
				case 0x08:
					code.Code = append(code.Code, Iconst5{})
				case 0x09:
					code.Code = append(code.Code, Lconst0{})
				case 0x0A:
					code.Code = append(code.Code, Lconst1{})
				case 0x0B:
					code.Code = append(code.Code, Fconst0{})
				case 0x0C:
					code.Code = append(code.Code, Fconst1{})
				case 0x0D:
					code.Code = append(code.Code, Fconst2{})
				case 0x0E:
					code.Code = append(code.Code, Dconst0{})
				case 0x0F:
					code.Code = append(code.Code, Dconst1{})
				case 0x10:
					var instr Bipush
					reader.read(&instr.Byte)
					code.Code = append(code.Code, instr)
				case 0x11:
					var instr Sipush
					reader.read(&instr.Short)
					code.Code = append(code.Code, instr)
				case 0x12:
					var cpIndex uint8
					reader.read(&cpIndex)
					value, err := getConstant[interface{}](cp, uint16(cpIndex))
					if err != nil {
						reader.fail(err)
					}
					code.Code = append(code.Code, Ldc{Value: value})
				case 0x13:
					code.Code = append(code.Code, LdcW{Value: readConstant[interface{}](reader, cp)})
				case 0x14:
					code.Code = append(code.Code, Ldc2W{Value: readConstant[interface{}](reader, cp)})
				case 0x15:
					var instr Iload
					reader.read(&instr.LocalIndex)
					code.Code = append(code.Code, instr)
				case 0x16:
					var instr Lload
					reader.read(&instr.LocalIndex)
					code.Code = append(code.Code, instr)
				case 0x17:
					var instr Fload
					reader.read(&instr.LocalIndex)
					code.Code = append(code.Code, instr)
				case 0x18:
					var instr Dload
					reader.read(&instr.LocalIndex)
					code.Code = append(code.Code, instr)
				case 0x19:
					var instr Aload
					reader.read(&instr.LocalIndex)
					code.Code = append(code.Code, instr)
				case 0x1A:
					code.Code = append(code.Code, Iload0{})
				case 0x1B:
					code.Code = append(code.Code, Iload1{})
				case 0x1C:
					code.Code = append(code.Code, Iload2{})
				case 0x1D:
					code.Code = append(code.Code, Iload3{})
				case 0x1E:
					code.Code = append(code.Code, Lload0{})
				case 0x1F:
					code.Code = append(code.Code, Lload1{})
				case 0x20:
					code.Code = append(code.Code, Lload2{})
				case 0x21:
					code.Code = append(code.Code, Lload3{})
				case 0x22:
					code.Code = append(code.Code, Fload0{})
				case 0x23:
					code.Code = append(code.Code, Fload1{})
				case 0x24:
					code.Code = append(code.Code, Fload2{})
				case 0x25:
					code.Code = append(code.Code, Fload3{})
				case 0x26:
					code.Code = append(code.Code, Dload0{})
				case 0x27:
					code.Code = append(code.Code, Dload1{})
				case 0x28:
					code.Code = append(code.Code, Dload2{})
				case 0x29:
					code.Code = append(code.Code, Dload3{})
				case 0x2A:
					code.Code = append(code.Code, Aload0{})
				case 0x2B:
					code.Code = append(code.Code, Aload1{})
				case 0x2C:
					code.Code = append(code.Code, Aload2{})
				case 0x2D:
					code.Code = append(code.Code, Aload3{})
				case 0x2E:
					code.Code = append(code.Code, Iaload{})
				case 0x2F:
					code.Code = append(code.Code, Laload{})
				case 0x30:
					code.Code = append(code.Code, Faload{})
				case 0x31:
					code.Code = append(code.Code, Daload{})
				case 0x32:
					code.Code = append(code.Code, Aaload{})
				case 0x33:
					code.Code = append(code.Code, Baload{})
				case 0x34:
					code.Code = append(code.Code, Caload{})
				case 0x35:
					code.Code = append(code.Code, Saload{})
				case 0x36:
					var instr Istore
					reader.read(&instr.LocalIndex)
					code.Code = append(code.Code, instr)
				case 0x37:
					var instr Lstore
					reader.read(&instr.LocalIndex)
					code.Code = append(code.Code, instr)
				case 0x38:
					var instr Fstore
					reader.read(&instr.LocalIndex)
					code.Code = append(code.Code, instr)
				case 0x39:
					var instr Dstore
					reader.read(&instr.LocalIndex)
					code.Code = append(code.Code, instr)
				case 0x3A:
					var instr Astore
					reader.read(&instr.LocalIndex)
					code.Code = append(code.Code, instr)
				case 0x3B:
					code.Code = append(code.Code, Istore0{})
				case 0x3C:
					code.Code = append(code.Code, Istore1{})
				case 0x3D:
					code.Code = append(code.Code, Istore2{})
				case 0x3E:
					code.Code = append(code.Code, Istore3{})
				case 0x3F:
					code.Code = append(code.Code, Lstore0{})
				case 0x40:
					code.Code = append(code.Code, Lstore1{})
				case 0x41:
					code.Code = append(code.Code, Lstore2{})
				case 0x42:
					code.Code = append(code.Code, Lstore3{})
				case 0x43:
					code.Code = append(code.Code, Fstore0{})
				case 0x44:
					code.Code = append(code.Code, Fstore1{})
				case 0x45:
					code.Code = append(code.Code, Fstore2{})
				case 0x46:
					code.Code = append(code.Code, Fstore3{})
				case 0x47:
					code.Code = append(code.Code, Dstore0{})
				case 0x48:
					code.Code = append(code.Code, Dstore1{})
				case 0x49:
					code.Code = append(code.Code, Dstore2{})
				case 0x4A:
					code.Code = append(code.Code, Dstore3{})
				case 0x4B:
					code.Code = append(code.Code, Astore0{})
				case 0x4C:
					code.Code = append(code.Code, Astore1{})
				case 0x4D:
					code.Code = append(code.Code, Astore2{})
				case 0x4E:
					code.Code = append(code.Code, Astore3{})
				case 0x4F:
					code.Code = append(code.Code, Iastore{})
				case 0x50:
					code.Code = append(code.Code, Lastore{})
				case 0x51:
					code.Code = append(code.Code, Fastore{})
				case 0x52:
					code.Code = append(code.Code, Dastore{})
				case 0x53:
					code.Code = append(code.Code, Aastore{})
				case 0x54:
					code.Code = append(code.Code, Bastore{})
				case 0x55:
					code.Code = append(code.Code, Castore{})
				case 0x56:
					code.Code = append(code.Code, Sastore{})
				case 0x57:
					code.Code = append(code.Code, Pop{})
				case 0x58:
					code.Code = append(code.Code, Pop2{})
				case 0x59:
					code.Code = append(code.Code, Dup{})
				case 0x5A:
					code.Code = append(code.Code, DupX1{})
				case 0x5B:
					code.Code = append(code.Code, DupX2{})
				case 0x5C:
					code.Code = append(code.Code, Dup2{})
				case 0x5D:
					code.Code = append(code.Code, Dup2X1{})
				case 0x5E:
					code.Code = append(code.Code, Dup2X2{})
				case 0x5F:
					code.Code = append(code.Code, Swap{})
				case 0x60:
					code.Code = append(code.Code, Iadd{})
				case 0x61:
					code.Code = append(code.Code, Ladd{})
				case 0x62:
					code.Code = append(code.Code, Fadd{})
				case 0x63:
					code.Code = append(code.Code, Dadd{})
				case 0x64:
					code.Code = append(code.Code, Isub{})
				case 0x65:
					code.Code = append(code.Code, Lsub{})
				case 0x66:
					code.Code = append(code.Code, Fsub{})
				case 0x67:
					code.Code = append(code.Code, Dsub{})
				case 0x68:
					code.Code = append(code.Code, Imul{})
				case 0x69:
					code.Code = append(code.Code, Lmul{})
				case 0x6A:
					code.Code = append(code.Code, Fmul{})
				case 0x6B:
					code.Code = append(code.Code, Dmul{})
				case 0x6C:
					code.Code = append(code.Code, Idiv{})
				case 0x6D:
					code.Code = append(code.Code, Ldiv{})
				case 0x6E:
					code.Code = append(code.Code, Fdiv{})
				case 0x6F:
					code.Code = append(code.Code, Ddiv{})
				case 0x70:
					code.Code = append(code.Code, Irem{})
				case 0x71:
					code.Code = append(code.Code, Lrem{})
				case 0x72:
					code.Code = append(code.Code, Frem{})
				case 0x73:
					code.Code = append(code.Code, Drem{})
				case 0x74:
					code.Code = append(code.Code, Ineg{})
				case 0x75:
					code.Code = append(code.Code, Lneg{})
				case 0x76:
					code.Code = append(code.Code, Fneg{})
				case 0x77:
					code.Code = append(code.Code, Dneg{})
				case 0x78:
					code.Code = append(code.Code, Ishl{})
				case 0x79:
					code.Code = append(code.Code, Lshl{})
				case 0x7A:
					code.Code = append(code.Code, Ishr{})
				case 0x7B:
					code.Code = append(code.Code, Lshr{})
				case 0x7C:
					code.Code = append(code.Code, Iushr{})
				case 0x7D:
					code.Code = append(code.Code, Lushr{})
				case 0x7E:
					code.Code = append(code.Code, Iand{})
				case 0x7F:
					code.Code = append(code.Code, Land{})
				case 0x80:
					code.Code = append(code.Code, Ior{})
				case 0x81:
					code.Code = append(code.Code, Lor{})
				case 0x82:
					code.Code = append(code.Code, Ixor{})
				case 0x83:
					code.Code = append(code.Code, Lxor{})
				case 0x84:
					var instr Iinc
					reader.read(&instr.LocalIndex)
					reader.read(&instr.Const)
					code.Code = append(code.Code, instr)
				case 0x85:
					code.Code = append(code.Code, I2l{})
				case 0x86:
					code.Code = append(code.Code, I2f{})
				case 0x87:
					code.Code = append(code.Code, I2d{})
				case 0x88:
					code.Code = append(code.Code, L2i{})
				case 0x89:
					code.Code = append(code.Code, L2f{})
				case 0x8A:
					code.Code = append(code.Code, L2d{})
				case 0x8B:
					code.Code = append(code.Code, F2i{})
				case 0x8C:
					code.Code = append(code.Code, F2l{})
				case 0x8D:
					code.Code = append(code.Code, F2d{})
				case 0x8E:
					code.Code = append(code.Code, D2i{})
				case 0x8F:
					code.Code = append(code.Code, D2l{})
				case 0x90:
					code.Code = append(code.Code, D2f{})
				case 0x91:
					code.Code = append(code.Code, I2b{})
				case 0x92:
					code.Code = append(code.Code, I2c{})
				case 0x93:
					code.Code = append(code.Code, I2s{})
				case 0x94:
					code.Code = append(code.Code, Lcmp{})
				case 0x95:
					code.Code = append(code.Code, Fcmpl{})
				case 0x96:
					code.Code = append(code.Code, Fcmpg{})
				case 0x97:
					code.Code = append(code.Code, Dcmpl{})
				case 0x98:
					code.Code = append(code.Code, Dcmpg{})
				case 0x99:
					var instr Ifeq
					reader.read(&instr.Offset)
					code.Code = append(code.Code, instr)
				case 0x9A:
					var instr Ifne
					reader.read(&instr.Offset)
					code.Code = append(code.Code, instr)
				case 0x9B:
					var instr Iflt
					reader.read(&instr.Offset)
					code.Code = append(code.Code, instr)
				case 0x9C:
					var instr Ifge
					reader.read(&instr.Offset)
					code.Code = append(code.Code, instr)
				case 0x9D:
					var instr Ifgt
					reader.read(&instr.Offset)
					code.Code = append(code.Code, instr)
				case 0x9E:
					var instr Ifle
					reader.read(&instr.Offset)
					code.Code = append(code.Code, instr)
				case 0x9F:
					var instr IfIcmpeq
					reader.read(&instr.Offset)
					code.Code = append(code.Code, instr)
				case 0xA0:
					var instr IfIcmpne
					reader.read(&instr.Offset)
					code.Code = append(code.Code, instr)
				case 0xA1:
					var instr IfIcmplt
					reader.read(&instr.Offset)
					code.Code = append(code.Code, instr)
				case 0xA2:
					var instr IfIcmpge
					reader.read(&instr.Offset)
					code.Code = append(code.Code, instr)
				case 0xA3:
					var instr IfIcmpgt
					reader.read(&instr.Offset)
					code.Code = append(code.Code, instr)
				case 0xA4:
					var instr IfIcmple
					reader.read(&instr.Offset)
					code.Code = append(code.Code, instr)
				case 0xA5:
					var instr IfAcmpeq
					reader.read(&instr.Offset)
					code.Code = append(code.Code, instr)
				case 0xA6:
					var instr IfAcmpne
					reader.read(&instr.Offset)
					code.Code = append(code.Code, instr)
				case 0xA7:
					var instr Goto
					reader.read(&instr.Offset)
					code.Code = append(code.Code, instr)
				case 0xA8:
					var instr Jsr
					reader.read(&instr.Offset)
					code.Code = append(code.Code, instr)
				case 0xA9:
					var instr Ret
					reader.read(&instr.LocalIndex)
					code.Code = append(code.Code, instr)
				case 0xAA:
					skipSwitchPadding(reader, start)
					var instr Tableswitch
					reader.read(&instr.Default)
					reader.read(&instr.Low)
					reader.read(&instr.High)
					if instr.Low > instr.High {
						reader.fail(fmt.Errorf("%w: invalid tableswitch: low (%d) > high (%d)", ErrMalformed, instr.Low, instr.High))
						break
					}
					count := int64(instr.High) - int64(instr.Low) + 1
					if count*4 > int64(reader.Len()) {
						reader.fail(fmt.Errorf("%w: invalid tableswitch: %d offsets exceed the remaining code", ErrMalformed, count))
						break
					}
					instr.Offsets = make([]int32, count)
					reader.read(instr.Offsets)
					code.Code = append(code.Code, instr)
				case 0xAB:
					skipSwitchPadding(reader, start)
					var instr Lookupswitch
					var npairs int32
					reader.read(&instr.Default)
					reader.read(&npairs)
					if npairs < 0 || int64(npairs)*8 > int64(reader.Len()) {
						reader.fail(fmt.Errorf("%w: invalid lookupswitch: bad number of pairs (%d)", ErrMalformed, npairs))
						break
					}
					instr.Pairs = make([]struct{ Match, Offset int32 }, npairs)
					reader.read(instr.Pairs)
					code.Code = append(code.Code, instr)
				case 0xAC:
					code.Code = append(code.Code, Ireturn{})
				case 0xAD:
					code.Code = append(code.Code, Lreturn{})
				case 0xAE:
					code.Code = append(code.Code, Freturn{})
				case 0xAF:
					code.Code = append(code.Code, Dreturn{})
				case 0xB0:
					code.Code = append(code.Code, Areturn{})
				case 0xB1:
					code.Code = append(code.Code, Return{})
				case 0xB2:
					code.Code = append(code.Code, Getstatic{Fieldref: readConstant[Fieldref](reader, cp)})
				case 0xB3:
					code.Code = append(code.Code, Putstatic{Fieldref: readConstant[Fieldref](reader, cp)})
				case 0xB4:
					code.Code = append(code.Code, Getfield{Fieldref: readConstant[Fieldref](reader, cp)})
				case 0xB5:
					code.Code = append(code.Code, Putfield{Fieldref: readConstant[Fieldref](reader, cp)})
				case 0xB6:
					code.Code = append(code.Code, Invokevirtual{Methodref: readConstant[Methodref](reader, cp)})
				case 0xB7:
					ref, isInterface := readMethodTarget(reader, cp)
					code.Code = append(code.Code, Invokespecial{Methodref: ref, IsInterface: isInterface})
				case 0xB8:
					ref, isInterface := readMethodTarget(reader, cp)
					code.Code = append(code.Code, Invokestatic{Methodref: ref, IsInterface: isInterface})
				case 0xB9:
					instr := Invokeinterface{
						InterfaceMethodref: readConstant[InterfaceMethodref](reader, cp),
					}
					reader.read(&instr.Count)
					var void byte
					reader.read(&void)
					code.Code = append(code.Code, instr)
				case 0xBA:
					instr := Invokedynamic{
						InvokeDynamic: readConstant[InvokeDynamic](reader, cp),
					}
					var void byte
					reader.read(&void)
					reader.read(&void)
					code.Code = append(code.Code, instr)
				case 0xBB:
					code.Code = append(code.Code, New(readConstant[Class](reader, cp)))
				case 0xBC:
					var instr Newarray
					reader.read(&instr.Type)
					code.Code = append(code.Code, instr)
				case 0xBD:
					code.Code = append(code.Code, Anewarray(readConstant[Class](reader, cp)))
				case 0xBE:
					code.Code = append(code.Code, Arraylength{})
				case 0xBF:
					code.Code = append(code.Code, Athrow{})
				case 0xC0:
					code.Code = append(code.Code, Checkcast(readConstant[Class](reader, cp)))
				case 0xC1:
					code.Code = append(code.Code, Instanceof(readConstant[Class](reader, cp)))
				case 0xC2:
					code.Code = append(code.Code, Monitorenter{})
				case 0xC3:
					code.Code = append(code.Code, Monitorexit{})
				case 0xC4:
					var instr Wide
					reader.read(&instr.OpCode)
					reader.read(&instr.LocalIndex)
					if instr.OpCode == 0x84 {
						reader.read(&instr.Const)
					}
					code.Code = append(code.Code, instr)
				case 0xC5:
					instr := Multianewarray{
						Class: string(readConstant[Class](reader, cp)),
					}
					reader.read(&instr.Dimension)
					code.Code = append(code.Code, instr)
				case 0xC6:
					var instr Ifnull
					reader.read(&instr.Offset)
					code.Code = append(code.Code, instr)
				case 0xC7:
					var instr Ifnonnull
					reader.read(&instr.Offset)
					code.Code = append(code.Code, instr)
				case 0xC8:
					var instr GotoW
					reader.read(&instr.Offset)
					code.Code = append(code.Code, instr)
				case 0xC9:
					var instr JsrW
					reader.read(&instr.Offset)
					code.Code = append(code.Code, instr)
				default:
					reader.fail(fmt.Errorf("%w: unknown opcode: 0x%X", ErrMalformed, opcode))
				}

				if start-reader.Len() > int(codeLength) {
					reader.fail(fmt.Errorf("%w: exceeds code length: %d", ErrMalformed, codeLength))
				}
				if reader.err != nil {
					reader.err = fmt.Errorf("failed to read instruction at offset %d: %w", pc, reader.err)
					break
				}
			}

			var exceptionTableLength uint16
			reader.read(&exceptionTableLength)
			code.ExceptionTable = make([]ExceptionTableEntry, exceptionTableLength)
			for i := 0; i < int(exceptionTableLength); i++ {
				exception := &ExceptionTableEntry{}
				reader.read(&exception.StartPc)
				reader.read(&exception.EndPc)
				reader.read(&exception.HandlerPc)
				reader.read(&exception.CatchType)
				code.ExceptionTable[i] = *exception
			}

			var attributesCount uint16
			reader.read(&attributesCount)
			nestedAttributes := make([]AttributeInfo, attributesCount)
			for i := 0; i < int(attributesCount); i++ {
				attribute := &AttributeInfo{}
				reader.read(&attribute.AttributeNameIndex)
				reader.read(&attribute.AttributeLength)
				if int64(attribute.AttributeLength) > int64(reader.Len()) {
					reader.fail(fmt.Errorf("%w: invalid attribute length: %d", ErrMalformed, attribute.AttributeLength))
					break
				}
				attribute.Info = make([]byte, attribute.AttributeLength)
				reader.read(attribute.Info)
				nestedAttributes[i] = *attribute
			}
			if reader.err == nil {
				nested, err := parseAttributes(nestedAttributes, cp)
				if err != nil {
					reader.fail(err)
				}
				code.Attributes = nested
			}

			attr = append(attr, code)
		case "ConstantValue": // TODO
			attr = append(attr, ConstantValue{})
		case "Deprecated":
			attr = append(attr, Deprecated{})
		case "Exceptions": // TODO
			attr = append(attr, Exceptions{})
		case "InnerClasses": // TODO
			attr = append(attr, InnerClasses{})
		case "LineNumberTable": // TODO
			attr = append(attr, LineNumberTable{})
		case "LocalVariableTable": // TODO
			attr = append(attr, LocalVariableTable{})
		case "LocalVariableTypeTable": // TODO
			attr = append(attr, LocalVariableTypeTable{})
		case "MethodParameters": // TODO
			attr = append(attr, MethodParameters{})
		case "RuntimeVisibleAnnotations": // TODO
			attr = append(attr, RuntimeVisibleAnnotations{})
		case "RuntimeInvisibleAnnotations": // TODO
			attr = append(attr, RuntimeInvisibleAnnotations{})
		case "RuntimeVisibleParameterAnnotations": // TODO
			attr = append(attr, RuntimeVisibleParameterAnnotations{})
		case "RuntimeInvisibleParameterAnnotations": // TODO
			attr = append(attr, RuntimeInvisibleParameterAnnotations{})
		case "SourceFile": // TODO
			attr = append(attr, SourceFile{})
		case "SourceDebugExtension": // TODO
			attr = append(attr, SourceDebugExtension{})
		case "Signature":
			attr = append(attr, Signature(readConstant[Utf8](reader, cp)))
		case "StackMapTable": // TODO
			attr = append(attr, StackMapTable{})
		case "Synthetic":
			attr = append(attr, Synthetic{})
		case "EnclosingMethod": // TODO
			attr = append(attr, EnclosingMethod{})
		case "BootstrapMethods": // TODO
			attr = append(attr, BootstrapMethods{})
		case "Module": // TODO
			attr = append(attr, ModuleInfo{})
		case "ModulePackages": // TODO
			attr = append(attr, ModulePackages{})
		case "NestHost": // TODO
			attr = append(attr, NestHost{})
		case "NestMembers": // TODO
			attr = append(attr, NestMembers{})
		case "PermittedSubclasses": // TODO
			attr = append(attr, PermittedSubclasses{})
		}

		if reader.err != nil {
			return nil, fmt.Errorf("failed to read attribute %s: %w", name, reader.err)
		}
	}

	return attr, nil
}

type binaryReader struct {
	*bytes.Reader
	err error
}

func (r *binaryReader) read(data interface{}) {
	if r.err == nil {
		r.err = binary.Read(r.Reader, binary.BigEndian, data)
	}
}

func (r *binaryReader) fail(err error) {
	if r.err == nil {
		r.err = err
	}
}

func skipSwitchPadding(reader *binaryReader, codeStart int) {
	pc := codeStart - reader.Len() // address right after the opcode
	if padding := (4 - pc%4) % 4; padding > 0 {
		reader.Seek(int64(padding), io.SeekCurrent)
	}
}

func readMethodTarget(reader *binaryReader, cp ConstantPool) (Methodref, bool) {
	var index uint16
	reader.read(&index)
	if ref, ok := cp[index].(InterfaceMethodref); ok {
		return Methodref(ref), true
	}
	ref, err := getConstant[Methodref](cp, index)
	if err != nil {
		reader.fail(err)
	}
	return ref, false
}

func readConstant[T any](reader *binaryReader, cp ConstantPool) T {
	var index uint16
	reader.read(&index)
	value, err := getConstant[T](cp, index)
	if err != nil {
		reader.fail(err)
	}
	return value
}
