package classfileparser

import (
	"bytes"
	"encoding/binary"
	"io"
)

type Attribute interface{}

// Code : Contient le bytecode et les informations relatives à la méthode.
type Code struct {
	MaxStack       uint16
	MaxLocals      uint16
	Code           []interface{}
	ExceptionTable []ExceptionTableEntry
	Attributes     []Attribute
}

// ConstantValue : Attribut d'un champ constant
type ConstantValue struct {
	ConstantValueIndex uint16
}

// Deprecated : Attribut indiquant qu'une méthode est obsolète
type Deprecated struct{}

// Exceptions : Attribut indiquant les exceptions vérifiées lancées par une méthode
type Exceptions struct {
	NumberOfExceptions  uint16
	ExceptionIndexTable []uint16
}

// InnerClasses : Attribut contenant des informations sur les classes internes
type InnerClasses struct {
	NumberOfClasses uint16
	InnerClassInfo  []InnerClassInfo
}

// LineNumberTable : Attribut de correspondance entre les lignes source et le bytecode
type LineNumberTable struct {
	LineNumberTableLength uint16
	LineNumberTable       []LineNumberTableEntry
}

// LocalVariableTable : Attribut contenant des informations sur les variables locales
type LocalVariableTable struct {
	LocalVariableTableLength uint16
	LocalVariableTable       []LocalVariableEntry
}

// LocalVariableTypeTable : Attribut contenant des informations sur les types des variables locales
type LocalVariableTypeTable struct {
	LocalVariableTypeTableLength uint16
	LocalVariableTypeTable       []LocalVariableTypeEntry
}

// MethodParameters : Attribut contenant des informations sur les paramètres de la méthode
type MethodParameters struct {
	MethodParametersCount uint16
	MethodParameters      []MethodParameter
}

// RuntimeVisibleAnnotations : Attribut contenant des annotations visibles à l'exécution
type RuntimeVisibleAnnotations struct {
	NumAnnotations uint16
	Annotations    []Annotation
}

// RuntimeInvisibleAnnotations : Attribut contenant des annotations invisibles à l'exécution
type RuntimeInvisibleAnnotations struct {
	NumAnnotations uint16
	Annotations    []Annotation
}

// RuntimeVisibleParameterAnnotations : Attribut contenant des annotations visibles pour les paramètres
type RuntimeVisibleParameterAnnotations struct {
	NumParameters        uint16
	ParameterAnnotations []ParameterAnnotation
}

// RuntimeInvisibleParameterAnnotations : Attribut contenant des annotations invisibles pour les paramètres
type RuntimeInvisibleParameterAnnotations struct {
	NumParameters        uint16
	ParameterAnnotations []ParameterAnnotation
}

// SourceFile : Attribut contenant le nom du fichier source
type SourceFile struct {
	SourcefileIndex uint16
}

// SourceDebugExtension : Attribut contenant des données de débogage source
type SourceDebugExtension struct {
	DebugExtension []byte
}

// Signature : Attribut contenant la signature d'une classe, méthode ou champ
type Signature struct {
	SignatureIndex uint16
}

// StackMapTable : Attribut contenant des informations sur la pile et les variables locales
type StackMapTable struct {
	NumberOfEntries uint16
	Entries         []StackMapFrame
}

// Synthetic : Attribut indiquant que l'élément est synthétique
type Synthetic struct{}

// EnclosingMethod : Attribut indiquant la méthode englobante
type EnclosingMethod struct {
	ClassIndex  uint16
	MethodIndex uint16
}

// BootstrapMethods : Attribut contenant des méthodes de bootstrap
type BootstrapMethods struct {
	NumBootstrapMethods uint16
	BootstrapMethods    []BootstrapMethod
}

// ModuleInfo : Attribut contenant des informations sur un module
type ModuleInfo struct {
	NameIndex uint16
	Flags     uint16
	Requires  []ModuleRequire
}

// ModulePackages : Attribut contenant des informations sur les packages d'un module
type ModulePackages struct {
	NumberOfPackages uint16
	PackageIndex     []uint16
}

// NestHost : Attribut contenant la classe hôte
type NestHost struct {
	HostClassIndex uint16
}

// NestMembers : Attribut contenant des informations sur les classes membres
type NestMembers struct {
	NumberOfMembers uint16
	ClassIndex      []uint16
}

// PermittedSubclasses : Attribut contenant des informations sur les sous-classes autorisées
type PermittedSubclasses struct {
	NumberOfSubclasses uint16
	SubclassIndex      []uint16
}

// Définition des structures supplémentaires nécessaires pour les attributs
type ExceptionTableEntry struct {
	StartPc   uint16
	EndPc     uint16
	HandlerPc uint16
	CatchType uint16
}

type InnerClassInfo struct {
	InnerClassIndex       uint16
	OuterClassIndex       uint16
	InnerNameIndex        uint16
	InnerClassAccessFlags uint16
}

type LineNumberTableEntry struct {
	StartPc    uint16
	LineNumber uint16
}

type LocalVariableEntry struct {
	StartPc        uint16
	Length         uint16
	NameIndex      uint16
	SignatureIndex uint16
	Index          uint16
}

type LocalVariableTypeEntry struct {
	StartPc        uint16
	Length         uint16
	NameIndex      uint16
	SignatureIndex uint16
	Index          uint16
}

type MethodParameter struct {
	NameIndex   uint16
	AccessFlags uint16
}

type Annotation struct {
	TypeIndex            uint16
	NumElementValuePairs uint16
	ElementValuePairs    []ElementValuePair
}

type ParameterAnnotation struct {
	NumAnnotations uint16
	Annotations    []Annotation
}

type ElementValuePair struct {
	ElementNameIndex uint16
	Value            ElementValue
}

type ElementValue struct {
	Tag   uint8
	Value uint16
}

type StackMapFrame struct {
	FrameType   uint8
	OffsetDelta uint16
	StackItems  []StackItem
}

type StackItem struct {
	Tag   uint8
	Value uint16
}

type BootstrapMethod struct {
	MethodRefIndex uint16
	ArgumentsCount uint16
	Arguments      []uint16
}

type ModuleRequire struct {
	NameIndex    uint16
	Flags        uint16
	VersionIndex uint16
}

func parseAttributes(attributes []AttributeInfo, cp ConstantPool) []Attribute {
	var attr []Attribute
	for _, a := range attributes {
		reader := bytes.NewReader(a.Info)
		switch cp[a.AttributeNameIndex].(Utf8) {
		case "Code": // TODO
			var code Code

			binary.Read(reader, binary.BigEndian, &code.MaxStack)
			binary.Read(reader, binary.BigEndian, &code.MaxLocals)

			var codeLength uint32
			binary.Read(reader, binary.BigEndian, &codeLength)
			start := reader.Len()
			for {
				if reader.Len() == start-int(codeLength) {
					break
				}
				var opcode uint8
				err := binary.Read(reader, binary.BigEndian, &opcode)
				if err != nil {
					panic(err)
				}

				switch opcode {
				case 0x00:
					code.Code = append(code.Code, Nop{})
					break
				case 0x01:
					code.Code = append(code.Code, AconstNull{})
					break
				case 0x02:
					code.Code = append(code.Code, IconstM1{})
					break
				case 0x03:
					code.Code = append(code.Code, Iconst0{})
					break
				case 0x04:
					code.Code = append(code.Code, Iconst1{})
					break
				case 0x05:
					code.Code = append(code.Code, Iconst2{})
					break
				case 0x06:
					code.Code = append(code.Code, Iconst3{})
					break
				case 0x07:
					code.Code = append(code.Code, Iconst4{})
					break
				case 0x08:
					code.Code = append(code.Code, Iconst5{})
					break
				case 0x09:
					code.Code = append(code.Code, Lconst0{})
					break
				case 0x0A:
					code.Code = append(code.Code, Lconst1{})
					break
				case 0x0B:
					code.Code = append(code.Code, Fconst0{})
					break
				case 0x0C:
					code.Code = append(code.Code, Fconst1{})
					break
				case 0x0D:
					code.Code = append(code.Code, Fconst2{})
					break
				case 0x0E:
					code.Code = append(code.Code, Dconst0{})
					break
				case 0x0F:
					code.Code = append(code.Code, Dconst1{})
					break
				case 0x10:
					var instr Bipush
					binary.Read(reader, binary.BigEndian, &instr.Byte)
					code.Code = append(code.Code, instr)
					break
				case 0x11:
					var instr Sipush
					binary.Read(reader, binary.BigEndian, &instr.Short)
					code.Code = append(code.Code, instr)
					break
				// case 0x12:
				// 	var cpIndex uint16
				// 	binary.Read(reader, binary.BigEndian, &cpIndex)
				// 	code.Code = append(code.Code, Ldc{})
				// 	break
				// case 0x13:
				// 	var cpIndex uint16
				// 	binary.Read(reader, binary.BigEndian, &cpIndex)
				// 	code.Code = append(code.Code, LdcW{})
				// 	break
				// case 0x14:
				// 	var cpIndex uint16
				// 	binary.Read(reader, binary.BigEndian, &cpIndex)
				// 	code.Code = append(code.Code, Ldc2W{})
				// 	break
				case 0x15:
					var instr Iload
					binary.Read(reader, binary.BigEndian, &instr.LocalIndex)
					code.Code = append(code.Code, instr)
					break
				case 0x16:
					var instr Lload
					binary.Read(reader, binary.BigEndian, &instr.LocalIndex)
					code.Code = append(code.Code, instr)
					break
				case 0x17:
					var instr Fload
					binary.Read(reader, binary.BigEndian, &instr.LocalIndex)
					code.Code = append(code.Code, instr)
					break
				case 0x18:
					var instr Dload
					binary.Read(reader, binary.BigEndian, &instr.LocalIndex)
					code.Code = append(code.Code, instr)
					break
				case 0x19:
					var instr Aload
					binary.Read(reader, binary.BigEndian, &instr.LocalIndex)
					code.Code = append(code.Code, instr)
					break
				case 0x1A:
					code.Code = append(code.Code, Iload0{})
					break
				case 0x1B:
					code.Code = append(code.Code, Iload1{})
					break
				case 0x1C:
					code.Code = append(code.Code, Iload2{})
					break
				case 0x1D:
					code.Code = append(code.Code, Iload3{})
					break
				case 0x1E:
					code.Code = append(code.Code, Lload0{})
					break
				case 0x1F:
					code.Code = append(code.Code, Lload1{})
					break
				case 0x20:
					code.Code = append(code.Code, Lload2{})
					break
				case 0x21:
					code.Code = append(code.Code, Lload3{})
					break
				case 0x22:
					code.Code = append(code.Code, Fload0{})
					break
				case 0x23:
					code.Code = append(code.Code, Fload1{})
					break
				case 0x24:
					code.Code = append(code.Code, Fload2{})
					break
				case 0x25:
					code.Code = append(code.Code, Fload3{})
					break
				case 0x26:
					code.Code = append(code.Code, Dload0{})
					break
				case 0x27:
					code.Code = append(code.Code, Dload1{})
					break
				case 0x28:
					code.Code = append(code.Code, Dload2{})
					break
				case 0x29:
					code.Code = append(code.Code, Dload3{})
					break
				case 0x2A:
					code.Code = append(code.Code, Aload0{})
					break
				case 0x2B:
					code.Code = append(code.Code, Aload1{})
					break
				case 0x2C:
					code.Code = append(code.Code, Aload2{})
					break
				case 0x2D:
					code.Code = append(code.Code, Aload3{})
					break
				case 0x2E:
					code.Code = append(code.Code, Iaload{})
					break
				case 0x2F:
					code.Code = append(code.Code, Laload{})
					break
				case 0x30:
					code.Code = append(code.Code, Faload{})
					break
				case 0x31:
					code.Code = append(code.Code, Daload{})
					break
				case 0x32:
					code.Code = append(code.Code, Aaload{})
					break
				case 0x33:
					code.Code = append(code.Code, Baload{})
					break
				case 0x34:
					code.Code = append(code.Code, Caload{})
					break
				case 0x35:
					code.Code = append(code.Code, Saload{})
					break
				case 0x36:
					var instr Istore
					binary.Read(reader, binary.BigEndian, &instr.LocalIndex)
					code.Code = append(code.Code, instr)
					break
				case 0x37:
					var instr Lstore
					binary.Read(reader, binary.BigEndian, &instr.LocalIndex)
					code.Code = append(code.Code, instr)
					break
				case 0x38:
					var instr Fstore
					binary.Read(reader, binary.BigEndian, &instr.LocalIndex)
					code.Code = append(code.Code, instr)
					break
				case 0x39:
					var instr Dstore
					binary.Read(reader, binary.BigEndian, &instr.LocalIndex)
					code.Code = append(code.Code, instr)
					break
				case 0x3A:
					var instr Astore
					binary.Read(reader, binary.BigEndian, &instr.LocalIndex)
					code.Code = append(code.Code, instr)
					break
				case 0x3B:
					code.Code = append(code.Code, Istore0{})
					break
				case 0x3C:
					code.Code = append(code.Code, Istore1{})
					break
				case 0x3D:
					code.Code = append(code.Code, Istore2{})
					break
				case 0x3E:
					code.Code = append(code.Code, Istore3{})
					break
				case 0x3F:
					code.Code = append(code.Code, Lstore0{})
					break
				case 0x40:
					code.Code = append(code.Code, Lstore1{})
					break
				case 0x41:
					code.Code = append(code.Code, Lstore2{})
					break
				case 0x42:
					code.Code = append(code.Code, Lstore3{})
					break
				case 0x43:
					code.Code = append(code.Code, Fstore0{})
					break
				case 0x44:
					code.Code = append(code.Code, Fstore1{})
					break
				case 0x45:
					code.Code = append(code.Code, Fstore2{})
					break
				case 0x46:
					code.Code = append(code.Code, Fstore3{})
					break
				case 0x47:
					code.Code = append(code.Code, Dstore0{})
					break
				case 0x48:
					code.Code = append(code.Code, Dstore1{})
					break
				case 0x49:
					code.Code = append(code.Code, Dstore2{})
					break
				case 0x4A:
					code.Code = append(code.Code, Dstore3{})
					break
				case 0x4B:
					code.Code = append(code.Code, Astore0{})
					break
				case 0x4C:
					code.Code = append(code.Code, Astore1{})
					break
				case 0x4D:
					code.Code = append(code.Code, Astore2{})
					break
				case 0x4E:
					code.Code = append(code.Code, Astore3{})
					break
				case 0x4F:
					code.Code = append(code.Code, Iastore{})
					break
				case 0x50:
					code.Code = append(code.Code, Lastore{})
					break
				case 0x51:
					code.Code = append(code.Code, Fastore{})
					break
				case 0x52:
					code.Code = append(code.Code, Dastore{})
					break
				case 0x53:
					code.Code = append(code.Code, Aastore{})
					break
				case 0x54:
					code.Code = append(code.Code, Bastore{})
					break
				case 0x55:
					code.Code = append(code.Code, Castore{})
					break
				case 0x56:
					code.Code = append(code.Code, Sastore{})
					break
				case 0x57:
					code.Code = append(code.Code, Pop{})
					break
				case 0x58:
					code.Code = append(code.Code, Pop2{})
					break
				case 0x59:
					code.Code = append(code.Code, Dup{})
					break
				case 0x5A:
					code.Code = append(code.Code, DupX1{})
					break
				case 0x5B:
					code.Code = append(code.Code, DupX2{})
					break
				case 0x5C:
					code.Code = append(code.Code, Dup2{})
					break
				case 0x5D:
					code.Code = append(code.Code, Dup2X1{})
					break
				case 0x5E:
					code.Code = append(code.Code, Dup2X2{})
					break
				case 0x5F:
					code.Code = append(code.Code, Swap{})
					break
				case 0x60:
					code.Code = append(code.Code, Iadd{})
					break
				case 0x61:
					code.Code = append(code.Code, Ladd{})
					break
				case 0x62:
					code.Code = append(code.Code, Fadd{})
					break
				case 0x63:
					code.Code = append(code.Code, Dadd{})
					break
				case 0x64:
					code.Code = append(code.Code, Isub{})
					break
				case 0x65:
					code.Code = append(code.Code, Lsub{})
					break
				case 0x66:
					code.Code = append(code.Code, Fsub{})
					break
				case 0x67:
					code.Code = append(code.Code, Dsub{})
					break
				case 0x68:
					code.Code = append(code.Code, Imul{})
					break
				case 0x69:
					code.Code = append(code.Code, Lmul{})
					break
				case 0x6A:
					code.Code = append(code.Code, Fmul{})
					break
				case 0x6B:
					code.Code = append(code.Code, Dmul{})
					break
				case 0x6C:
					code.Code = append(code.Code, Idiv{})
					break
				case 0x6D:
					code.Code = append(code.Code, Ldiv{})
					break
				case 0x6E:
					code.Code = append(code.Code, Fdiv{})
					break
				case 0x6F:
					code.Code = append(code.Code, Ddiv{})
					break
				case 0x70:
					code.Code = append(code.Code, Irem{})
					break
				case 0x71:
					code.Code = append(code.Code, Lrem{})
					break
				case 0x72:
					code.Code = append(code.Code, Frem{})
					break
				case 0x73:
					code.Code = append(code.Code, Drem{})
					break
				case 0x74:
					code.Code = append(code.Code, Ineg{})
					break
				case 0x75:
					code.Code = append(code.Code, Lneg{})
					break
				case 0x76:
					code.Code = append(code.Code, Fneg{})
					break
				case 0x77:
					code.Code = append(code.Code, Dneg{})
					break
				case 0x78:
					code.Code = append(code.Code, Ishl{})
					break
				case 0x79:
					code.Code = append(code.Code, Lshl{})
					break
				case 0x7A:
					code.Code = append(code.Code, Ishr{})
					break
				case 0x7B:
					code.Code = append(code.Code, Lshr{})
					break
				case 0x7C:
					code.Code = append(code.Code, Iushr{})
					break
				case 0x7D:
					code.Code = append(code.Code, Lushr{})
					break
				case 0x7E:
					code.Code = append(code.Code, Iand{})
					break
				case 0x7F:
					code.Code = append(code.Code, Land{})
					break
				case 0x80:
					code.Code = append(code.Code, Ior{})
					break
				case 0x81:
					code.Code = append(code.Code, Lor{})
					break
				case 0x82:
					code.Code = append(code.Code, Ixor{})
					break
				case 0x83:
					code.Code = append(code.Code, Lxor{})
					break
				// case 0x84:
				// 	code.Code = append(code.Code, Iinc{})
				// 	break
				case 0x85:
					code.Code = append(code.Code, I2l{})
					break
				case 0x86:
					code.Code = append(code.Code, I2f{})
					break
				case 0x87:
					code.Code = append(code.Code, I2d{})
					break
				case 0x88:
					code.Code = append(code.Code, L2i{})
					break
				case 0x89:
					code.Code = append(code.Code, L2f{})
					break
				case 0x8A:
					code.Code = append(code.Code, L2d{})
					break
				case 0x8B:
					code.Code = append(code.Code, F2i{})
					break
				case 0x8C:
					code.Code = append(code.Code, F2l{})
					break
				case 0x8D:
					code.Code = append(code.Code, F2d{})
					break
				case 0x8E:
					code.Code = append(code.Code, D2i{})
					break
				case 0x8F:
					code.Code = append(code.Code, D2l{})
					break
				case 0x90:
					code.Code = append(code.Code, D2f{})
					break
				case 0x91:
					code.Code = append(code.Code, I2b{})
					break
				case 0x92:
					code.Code = append(code.Code, I2c{})
					break
				case 0x93:
					code.Code = append(code.Code, I2s{})
					break
				case 0x94:
					code.Code = append(code.Code, Lcmp{})
					break
				case 0x95:
					code.Code = append(code.Code, Fcmpl{})
					break
				case 0x96:
					code.Code = append(code.Code, Fcmpg{})
					break
				case 0x97:
					code.Code = append(code.Code, Dcmpl{})
					break
				case 0x98:
					code.Code = append(code.Code, Dcmpg{})
					break
				case 0x99:
					var instr Ifeq
					binary.Read(reader, binary.BigEndian, &instr.Offset)
					code.Code = append(code.Code, instr)
					break
				case 0x9A:
					var instr Ifne
					binary.Read(reader, binary.BigEndian, &instr.Offset)
					code.Code = append(code.Code, instr)
					break
				case 0x9B:
					var instr Iflt
					binary.Read(reader, binary.BigEndian, &instr.Offset)
					code.Code = append(code.Code, instr)
					break
				case 0x9C:
					var instr Ifge
					binary.Read(reader, binary.BigEndian, &instr.Offset)
					code.Code = append(code.Code, instr)
					break
				case 0x9D:
					var instr Ifgt
					binary.Read(reader, binary.BigEndian, &instr.Offset)
					code.Code = append(code.Code, instr)
					break
				case 0x9E:
					var instr Ifle
					binary.Read(reader, binary.BigEndian, &instr.Offset)
					code.Code = append(code.Code, instr)
					break
				case 0x9F:
					var instr IfIcmpeq
					binary.Read(reader, binary.BigEndian, &instr.Offset)
					code.Code = append(code.Code, instr)
					break
				case 0xA0:
					var instr IfIcmpne
					binary.Read(reader, binary.BigEndian, &instr.Offset)
					code.Code = append(code.Code, instr)
					break
				case 0xA1:
					var instr IfIcmplt
					binary.Read(reader, binary.BigEndian, &instr.Offset)
					code.Code = append(code.Code, instr)
					break
				case 0xA2:
					var instr IfIcmpge
					binary.Read(reader, binary.BigEndian, &instr.Offset)
					code.Code = append(code.Code, instr)
					break
				case 0xA3:
					var instr IfIcmpgt
					binary.Read(reader, binary.BigEndian, &instr.Offset)
					code.Code = append(code.Code, instr)
					break
				case 0xA4:
					var instr IfIcmple
					binary.Read(reader, binary.BigEndian, &instr.Offset)
					code.Code = append(code.Code, instr)
					break
				case 0xA5:
					var instr IfAcmpeq
					binary.Read(reader, binary.BigEndian, &instr.Offset)
					code.Code = append(code.Code, instr)
					break
				case 0xA6:
					var instr IfAcmpne
					binary.Read(reader, binary.BigEndian, &instr.Offset)
					code.Code = append(code.Code, instr)
					break
				case 0xA7:
					var instr Goto
					binary.Read(reader, binary.BigEndian, &instr.Offset)
					code.Code = append(code.Code, instr)
					break
				case 0xA8:
					var instr Jsr
					binary.Read(reader, binary.BigEndian, &instr.Offset)
					code.Code = append(code.Code, instr)
					break
				// case 0xA9:
				// 	code.Code = append(code.Code, Ret{})
				// 	break
				case 0xAC:
					code.Code = append(code.Code, Ireturn{})
					break
				case 0xAD:
					code.Code = append(code.Code, Lreturn{})
					break
				case 0xAE:
					code.Code = append(code.Code, Freturn{})
					break
				case 0xAF:
					code.Code = append(code.Code, Dreturn{})
					break
				case 0xB0:
					code.Code = append(code.Code, Areturn{})
					break
				case 0xB1:
					code.Code = append(code.Code, Return{})
					break
				case 0xB2:
					var cpIndex uint16
					binary.Read(reader, binary.BigEndian, &cpIndex)
					code.Code = append(code.Code, Getstatic(cp[cpIndex].(Fieldref)))
					break
				case 0xB3:
					var cpIndex uint16
					binary.Read(reader, binary.BigEndian, &cpIndex)
					code.Code = append(code.Code, Putstatic(cp[cpIndex].(Fieldref)))
					break
				case 0xB4:
					var cpIndex uint16
					binary.Read(reader, binary.BigEndian, &cpIndex)
					code.Code = append(code.Code, Getfield(cp[cpIndex].(Fieldref)))
					break
				case 0xB5:
					var cpIndex uint16
					binary.Read(reader, binary.BigEndian, &cpIndex)
					code.Code = append(code.Code, Putfield(cp[cpIndex].(Fieldref)))
					break
				case 0xB6:
					var cpIndex uint16
					binary.Read(reader, binary.BigEndian, &cpIndex)
					code.Code = append(code.Code, Invokevirtual(cp[cpIndex].(Methodref)))
					break
				case 0xB7:
					var cpIndex uint16
					binary.Read(reader, binary.BigEndian, &cpIndex)
					code.Code = append(code.Code, Invokespecial(cp[cpIndex].(Methodref)))
					break
				case 0xB8:
					var cpIndex uint16
					binary.Read(reader, binary.BigEndian, &cpIndex)
					code.Code = append(code.Code, Invokestatic(cp[cpIndex].(Methodref)))
					break
				case 0xBB:
					var cpIndex uint16
					binary.Read(reader, binary.BigEndian, &cpIndex)
					code.Code = append(code.Code, New(cp[cpIndex].(Class)))
					break
				case 0xBD:
					var cpIndex uint16
					binary.Read(reader, binary.BigEndian, &cpIndex)
					code.Code = append(code.Code, Anewarray(cp[cpIndex].(Class)))
					break
				case 0xBE:
					code.Code = append(code.Code, Arraylength{})
					break
				case 0xBF:
					code.Code = append(code.Code, Athrow{})
					break
				case 0xC0:
					var cpIndex uint16
					binary.Read(reader, binary.BigEndian, &cpIndex)
					code.Code = append(code.Code, Checkcast(cp[cpIndex].(Class)))
					break
				case 0xC1:
					var cpIndex uint16
					binary.Read(reader, binary.BigEndian, &cpIndex)
					code.Code = append(code.Code, Instanceof(cp[cpIndex].(Class)))
					break
				case 0xC2:
					code.Code = append(code.Code, Monitorenter{})
					break
				case 0xC3:
					code.Code = append(code.Code, Monitorexit{})
					break
				// case 0xC5:
				// 	code.Code = append(code.Code, Multianewarray{})
				// 	break
				case 0xC6:
					var instr Ifnull
					binary.Read(reader, binary.BigEndian, &instr.Offset)
					code.Code = append(code.Code, instr)
					break
				case 0xC7:
					var instr Ifnonnull
					binary.Read(reader, binary.BigEndian, &instr.Offset)
					code.Code = append(code.Code, instr)
					break
				case 0xC8:
					var instr GotoW
					binary.Read(reader, binary.BigEndian, &instr.Offset)
					code.Code = append(code.Code, instr)
					break
				case 0xC9:
					var instr JsrW
					binary.Read(reader, binary.BigEndian, &instr.Offset)
					code.Code = append(code.Code, instr)
					break
				}
			}

			var exceptionTableLength uint16
			binary.Read(reader, binary.BigEndian, &exceptionTableLength)
			code.ExceptionTable = make([]ExceptionTableEntry, exceptionTableLength)
			for i := 0; i < int(exceptionTableLength); i++ {
				exception := &ExceptionTableEntry{}
				binary.Read(reader, binary.BigEndian, &exception.StartPc)
				binary.Read(reader, binary.BigEndian, &exception.EndPc)
				binary.Read(reader, binary.BigEndian, &exception.HandlerPc)
				binary.Read(reader, binary.BigEndian, &exception.CatchType)
				code.ExceptionTable[i] = *exception
			}

			var attributesCount uint16
			binary.Read(reader, binary.BigEndian, &attributesCount)
			nestedAttributes := make([]AttributeInfo, attributesCount)
			for i := 0; i < int(attributesCount); i++ {
				attribute := &AttributeInfo{}
				binary.Read(reader, binary.BigEndian, &attribute.AttributeNameIndex)
				binary.Read(reader, binary.BigEndian, &attribute.AttributeLength)
				attribute.Info = make([]byte, attribute.AttributeLength)
				io.ReadFull(reader, attribute.Info)
				nestedAttributes[i] = *attribute
			}
			code.Attributes = parseAttributes(nestedAttributes, cp)

			attr = append(attr, code)
			break
		case "ConstantValue": // TODO
			attr = append(attr, ConstantValue{})
			break
		case "Deprecated": // TODO
			attr = append(attr, Deprecated{})
			break
		case "Exceptions": // TODO
			attr = append(attr, Exceptions{})
			break
		case "InnerClasses": // TODO
			attr = append(attr, InnerClasses{})
			break
		case "LineNumberTable": // TODO
			attr = append(attr, LineNumberTable{})
			break
		case "LocalVariableTable": // TODO
			attr = append(attr, LocalVariableTable{})
			break
		case "LocalVariableTypeTable": // TODO
			attr = append(attr, LocalVariableTypeTable{})
			break
		case "MethodParameters": // TODO
			attr = append(attr, MethodParameters{})
			break
		case "RuntimeVisibleAnnotations": // TODO
			attr = append(attr, RuntimeVisibleAnnotations{})
			break
		case "RuntimeInvisibleAnnotations": // TODO
			attr = append(attr, RuntimeInvisibleAnnotations{})
			break
		case "RuntimeVisibleParameterAnnotations": // TODO
			attr = append(attr, RuntimeVisibleParameterAnnotations{})
			break
		case "RuntimeInvisibleParameterAnnotations": // TODO
			attr = append(attr, RuntimeInvisibleParameterAnnotations{})
			break
		case "SourceFile": // TODO
			attr = append(attr, SourceFile{})
			break
		case "SourceDebugExtension": // TODO
			attr = append(attr, SourceDebugExtension{})
			break
		case "Signature": // TODO
			attr = append(attr, Signature{})
			break
		case "StackMapTable": // TODO
			attr = append(attr, StackMapTable{})
			break
		case "Synthetic": // TODO
			attr = append(attr, Synthetic{})
			break
		case "EnclosingMethod": // TODO
			attr = append(attr, EnclosingMethod{})
			break
		case "BootstrapMethods": // TODO
			attr = append(attr, BootstrapMethods{})
			break
		case "Module": // TODO
			attr = append(attr, ModuleInfo{})
			break
		case "ModulePackages": // TODO
			attr = append(attr, ModulePackages{})
			break
		case "NestHost": // TODO
			attr = append(attr, NestHost{})
			break
		case "NestMembers": // TODO
			attr = append(attr, NestMembers{})
			break
		case "PermittedSubclasses": // TODO
			attr = append(attr, PermittedSubclasses{})
			break
		}
	}

	return attr
}
