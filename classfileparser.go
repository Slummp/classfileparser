package classfileparser

import (
	"encoding/binary"
	"fmt"
	"io"
)

// ClassFile represents the entire structure of a .class file
type ClassFile struct {
	Magic             uint32          // Magic number (0xCAFEBABE)
	MinorVersion      uint16          // Minor version
	MajorVersion      uint16          // Major version
	ConstantPoolCount uint16          // Number of entries in the constant pool
	ConstantPool      []CpInfo        // Constant pool entries
	AccessFlags       uint16          // Access flags (e.g., public, final)
	ThisClass         uint16          // Index of the current class in the constant pool
	SuperClass        uint16          // Index of the superclass in the constant pool
	InterfacesCount   uint16          // Number of interfaces implemented by this class
	Interfaces        []uint16        // Interface indexes in the constant pool
	FieldsCount       uint16          // Number of fields declared by this class
	Fields            []FieldInfo     // Field structures
	MethodsCount      uint16          // Number of methods declared by this class
	Methods           []MethodInfo    // Method structures
	AttributesCount   uint16          // Number of attributes associated with this class
	Attributes        []AttributeInfo // Attribute structures
}

// CpInfo represents an entry in the constant pool
type CpInfo struct {
	Tag  uint8  // Type of the constant pool entry (e.g., Class, Utf8, Methodref, etc.)
	Info []byte // Data specific to the constant type
}

// FieldInfo represents a field in the class
type FieldInfo struct {
	AccessFlags     uint16          // Access flags for the field
	NameIndex       uint16          // Index in the constant pool pointing to the name
	DescriptorIndex uint16          // Index in the constant pool pointing to the type descriptor
	AttributesCount uint16          // Number of attributes associated with the field
	Attributes      []AttributeInfo // Field attributes
}

// MethodInfo represents a method in the class
type MethodInfo struct {
	AccessFlags     uint16          // Access flags for the method (e.g., public, static, etc.)
	NameIndex       uint16          // Index in the constant pool pointing to the method name
	DescriptorIndex uint16          // Index in the constant pool pointing to the method descriptor
	AttributesCount uint16          // Number of attributes associated with the method
	Attributes      []AttributeInfo // Method attributes
}

// AttributeInfo represents an attribute structure
type AttributeInfo struct {
	AttributeNameIndex uint16 // Index in the constant pool pointing to the attribute name
	AttributeLength    uint32 // Length of the attribute in bytes
	Info               []byte // Raw attribute data (to be parsed depending on the attribute type)
}

// CodeAttribute represents the "Code" attribute of a method
type CodeAttribute struct {
	MaxStack       uint16          // Maximum stack size
	MaxLocals      uint16          // Number of local variables
	CodeLength     uint32          // Length of the bytecode
	Code           []byte          // Bytecode instructions
	ExceptionTable []ExceptionInfo // Table of exceptions
	Attributes     []AttributeInfo // Additional attributes within the code
}

// ExceptionInfo represents an entry in the exception table within a Code attribute
type ExceptionInfo struct {
	StartPC   uint16 // Start of the code block covered by this entry
	EndPC     uint16 // End of the code block covered by this entry
	HandlerPC uint16 // Location of the exception handler
	CatchType uint16 // Index in the constant pool for the exception type (or 0 for any exception)
}

// Open creates a ClassFile by parsing the content of the provided file
func Open(file io.Reader) (*ClassFile, error) {
	cf := &ClassFile{}

	// Read and validate the magic number
	if err := binary.Read(file, binary.BigEndian, &cf.Magic); err != nil {
		return nil, fmt.Errorf("failed to read magic number: %w", err)
	}
	if cf.Magic != 0xCAFEBABE {
		return nil, fmt.Errorf("invalid magic number: 0x%X", cf.Magic)
	}

	// Read versions
	if err := binary.Read(file, binary.BigEndian, &cf.MinorVersion); err != nil {
		return nil, fmt.Errorf("failed to read minor version: %w", err)
	}
	if err := binary.Read(file, binary.BigEndian, &cf.MajorVersion); err != nil {
		return nil, fmt.Errorf("failed to read major version: %w", err)
	}

	// Read constant pool count
	if err := binary.Read(file, binary.BigEndian, &cf.ConstantPoolCount); err != nil {
		return nil, fmt.Errorf("failed to read constant pool count: %w", err)
	}

	// Read constant pool entries
	cf.ConstantPool = make([]CpInfo, cf.ConstantPoolCount-1)
	for i := range cf.ConstantPool {
		var tag uint8
		if err := binary.Read(file, binary.BigEndian, &tag); err != nil {
			return nil, fmt.Errorf("failed to read constant pool tag: %w", err)
		}
		cf.ConstantPool[i].Tag = tag

		// Read the data for the constant pool entry based on the tag
		infoLength, err := getCpInfoLength(tag, file)
		if err != nil {
			return nil, fmt.Errorf("failed to read constant pool length: %w", err)
		}
		cf.ConstantPool[i].Info = make([]byte, infoLength)
		if _, err := io.ReadFull(file, cf.ConstantPool[i].Info); err != nil {
			return nil, fmt.Errorf("failed to read constant pool info: %w", err)
		}
	}

	// Read access flags
	if err := binary.Read(file, binary.BigEndian, &cf.AccessFlags); err != nil {
		return nil, fmt.Errorf("failed to read access flags: %w", err)
	}

	// Read this class, super class, and interfaces
	if err := binary.Read(file, binary.BigEndian, &cf.ThisClass); err != nil {
		return nil, fmt.Errorf("failed to read this class index: %w", err)
	}
	if err := binary.Read(file, binary.BigEndian, &cf.SuperClass); err != nil {
		return nil, fmt.Errorf("failed to read super class index: %w", err)
	}
	if err := binary.Read(file, binary.BigEndian, &cf.InterfacesCount); err != nil {
		return nil, fmt.Errorf("failed to read interfaces count: %w", err)
	}
	cf.Interfaces = make([]uint16, cf.InterfacesCount)
	if err := binary.Read(file, binary.BigEndian, &cf.Interfaces); err != nil {
		return nil, fmt.Errorf("failed to read interfaces: %w", err)
	}

	// Read fields
	if err := binary.Read(file, binary.BigEndian, &cf.FieldsCount); err != nil {
		return nil, fmt.Errorf("failed to read fields count: %w", err)
	}
	cf.Fields = make([]FieldInfo, cf.FieldsCount)
	for i := range cf.Fields {
		if err := parseField(file, &cf.Fields[i]); err != nil {
			return nil, fmt.Errorf("failed to read field %d: %w", i, err)
		}
	}

	// Read methods
	if err := binary.Read(file, binary.BigEndian, &cf.MethodsCount); err != nil {
		return nil, fmt.Errorf("failed to read methods count: %w", err)
	}
	cf.Methods = make([]MethodInfo, cf.MethodsCount)
	for i := range cf.Methods {
		if err := parseMethod(file, &cf.Methods[i]); err != nil {
			return nil, fmt.Errorf("failed to read method %d: %w", i, err)
		}
	}

	// Read attributes
	if err := binary.Read(file, binary.BigEndian, &cf.AttributesCount); err != nil {
		return nil, fmt.Errorf("failed to read attributes count: %w", err)
	}
	cf.Attributes = make([]AttributeInfo, cf.AttributesCount)
	for i := range cf.Attributes {
		if err := parseAttribute(file, &cf.Attributes[i]); err != nil {
			return nil, fmt.Errorf("failed to read attribute %d: %w", i, err)
		}
	}

	return cf, nil
}

// Helper functions for parsing
func getCpInfoLength(tag uint8, file io.Reader) (int, error) {
	switch tag {
	case 1: // CONSTANT_Utf8
		// Lire la longueur du champ Utf8 (2 octets)
		var length uint16
		if err := binary.Read(file, binary.BigEndian, &length); err != nil {
			return 0, fmt.Errorf("failed to read CONSTANT_Utf8 length: %v", err)
		}
		return int(length), nil

	case 3, 4: // CONSTANT_Integer, CONSTANT_Float
		return 4, nil

	case 5, 6: // CONSTANT_Long, CONSTANT_Double
		return 8, nil

	case 7, 8: // CONSTANT_Class, CONSTANT_String
		return 2, nil

	case 9, 10, 11: // CONSTANT_Fieldref, CONSTANT_Methodref, CONSTANT_InterfaceMethodref
		return 4, nil

	case 12: // CONSTANT_NameAndType
		return 4, nil

	case 15: // CONSTANT_MethodHandle
		return 3, nil

	case 16: // CONSTANT_MethodType
		return 2, nil

	case 18: // CONSTANT_InvokeDynamic
		return 4, nil

	case 19, 20: // CONSTANT_Module, CONSTANT_Package (Java 9+)
		return 2, nil

	default:
		return 0, fmt.Errorf("unknown constant pool tag: %d", tag)
	}
}

func parseField(file io.Reader, field *FieldInfo) error {
	if err := binary.Read(file, binary.BigEndian, &field.AccessFlags); err != nil {
		return err
	}
	if err := binary.Read(file, binary.BigEndian, &field.NameIndex); err != nil {
		return err
	}
	if err := binary.Read(file, binary.BigEndian, &field.DescriptorIndex); err != nil {
		return err
	}
	if err := binary.Read(file, binary.BigEndian, &field.AttributesCount); err != nil {
		return err
	}
	field.Attributes = make([]AttributeInfo, field.AttributesCount)
	for i := range field.Attributes {
		if err := parseAttribute(file, &field.Attributes[i]); err != nil {
			return err
		}
	}
	return nil
}

func parseMethod(file io.Reader, method *MethodInfo) error {
	if err := binary.Read(file, binary.BigEndian, &method.AccessFlags); err != nil {
		return err
	}
	if err := binary.Read(file, binary.BigEndian, &method.NameIndex); err != nil {
		return err
	}
	if err := binary.Read(file, binary.BigEndian, &method.DescriptorIndex); err != nil {
		return err
	}
	if err := binary.Read(file, binary.BigEndian, &method.AttributesCount); err != nil {
		return err
	}
	method.Attributes = make([]AttributeInfo, method.AttributesCount)
	for i := range method.Attributes {
		if err := parseAttribute(file, &method.Attributes[i]); err != nil {
			return err
		}
	}
	return nil
}

func parseAttribute(file io.Reader, attr *AttributeInfo) error {
	if err := binary.Read(file, binary.BigEndian, &attr.AttributeNameIndex); err != nil {
		return err
	}
	if err := binary.Read(file, binary.BigEndian, &attr.AttributeLength); err != nil {
		return err
	}
	attr.Info = make([]byte, attr.AttributeLength)
	if _, err := io.ReadFull(file, attr.Info); err != nil {
		return err
	}
	return nil
}
