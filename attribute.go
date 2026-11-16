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

// Record describes the components of a record class (decoding TODO).
type Record struct{}

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
				code.Code = append(code.Code, decodeInstruction(reader, cp, start))

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
		case "Record": // TODO
			attr = append(attr, Record{})
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
	if r.err != nil {
		return
	}
	if err := binary.Read(r.Reader, binary.BigEndian, data); err != nil {
		r.err = fmt.Errorf("%w: exceeds attribute length: %d", ErrMalformed, r.Size())
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
