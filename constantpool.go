package classfileparser

import (
	"encoding/binary"
	"fmt"
	"math"
	"slices"
)

// ConstantPool is a typed view of the constant pool, indexed by the JVM slot number
type ConstantPool map[uint16]interface{}

// Utf8 represents a CONSTANT_Utf8 entry decoded as a Go string
type Utf8 string

// Class represents a CONSTANT_Class entry resolved to its internal JVM name
type Class string

// MethodType represents a CONSTANT_MethodType signature descriptor
type MethodType string

// Module describes a CONSTANT_Module entry holding a module name
type Module string

// Package describes a CONSTANT_Package entry holding a package name
type Package string

// Fieldref bundles class/name/type info from a CONSTANT_Fieldref entry
type Fieldref struct {
	Class string
	Name  string
	Type  string
}

// Methodref bundles class/name/type info from a CONSTANT_Methodref entry
type Methodref struct {
	Class string
	Name  string
	Type  string
}

// InterfaceMethodref bundles class/name/type info from a CONSTANT_InterfaceMethodref entry
type InterfaceMethodref struct {
	Class string
	Name  string
	Type  string
}

// NameAndType represents a CONSTANT_NameAndType entry
type NameAndType struct {
	Name string
	Type string
}

// MethodHandle holds the decoded data of a CONSTANT_MethodHandle entry
type MethodHandle struct {
	Kind        string
	Class       string
	Name        string
	Type        string
	IsInterface bool // true if the target is an interface method (InterfaceMethodref)
}

// Dynamic represents a CONSTANT_Dynamic entry (bootstrap handling TODO)
type Dynamic struct {
	Name           string
	Type           string
	BootstrapIndex uint16
}

// InvokeDynamic represents a CONSTANT_InvokeDynamic entry (bootstrap handling TODO)
type InvokeDynamic struct {
	BootstrapIndex uint16
	Name           string
	Type           string
}

// GetConstantPool resolves the raw constant pool entries into typed Go values for easier access
func (cf *ClassFile) GetConstantPool() (ConstantPool, error) {
	cp := ConstantPool{}
	for i, cpItem := range cf.ConstantPool {
		if cpItem.Tag == 0 {
			continue
		}
		value, err := resolveConstant(cpItem, cf.ConstantPool)
		if err != nil {
			return nil, fmt.Errorf("failed to read constant pool entry %d: %w", i+1, err)
		}
		cp[uint16(i+1)] = value
	}
	return cp, nil
}

func resolveConstant(cpItem CpInfo, cp []CpInfo) (interface{}, error) {
	switch cpItem.Tag {
	case 1:
		return Utf8(cpItem.Info), nil
	case 3:
		return int32(binary.BigEndian.Uint32(cpItem.Info)), nil
	case 4:
		return math.Float32frombits(binary.BigEndian.Uint32(cpItem.Info)), nil
	case 5:
		return int64(binary.BigEndian.Uint64(cpItem.Info)), nil
	case 6:
		return math.Float64frombits(binary.BigEndian.Uint64(cpItem.Info)), nil
	case 7:
		name, err := getString(cpItem.Info, cp)
		return Class(name), err
	case 8:
		return getString(cpItem.Info, cp)
	case 9:
		ref, err := getClassNameType(cpItem.Info, cp)
		return Fieldref(ref), err
	case 10:
		ref, err := getClassNameType(cpItem.Info, cp)
		return Methodref(ref), err
	case 11:
		ref, err := getClassNameType(cpItem.Info, cp)
		return InterfaceMethodref(ref), err
	case 12:
		return getNameType(cpItem.Info, cp)
	case 15:
		kind, ok := map[byte]string{
			1: "getField",
			2: "getStatic",
			3: "putField",
			4: "putStatic",
			5: "invokeVirtual",
			6: "invokeStatic",
			7: "invokeSpecial",
			8: "newInvokeSpecial",
			9: "invokeInterface",
		}[cpItem.Info[0]]
		if !ok {
			return nil, fmt.Errorf("invalid method handle kind: %d", cpItem.Info[0])
		}
		ref, err := getEntry(cpItem.Info[1:3], cp, 9, 10, 11)
		if err != nil {
			return nil, err
		}
		cnt, err := getClassNameType(ref.Info, cp)
		return MethodHandle{
			Kind:        kind,
			Class:       cnt.Class,
			Name:        cnt.Name,
			Type:        cnt.Type,
			IsInterface: ref.Tag == 11,
		}, err
	case 16:
		descriptor, err := getString(cpItem.Info, cp)
		return MethodType(descriptor), err
	case 17:
		return Dynamic{}, nil // TODO
	case 18:
		return InvokeDynamic{}, nil // TODO
	case 19:
		name, err := getString(cpItem.Info, cp)
		return Module(name), err
	case 20:
		name, err := getString(cpItem.Info, cp)
		return Package(name), err
	default:
		return nil, fmt.Errorf("unknown constant pool tag: %d", cpItem.Tag)
	}
}

func getEntry(i []byte, cp []CpInfo, tags ...uint8) (CpInfo, error) {
	index := binary.BigEndian.Uint16(i)
	if index == 0 || int(index) > len(cp) {
		return CpInfo{}, fmt.Errorf("invalid constant pool index: %d", index)
	}
	entry := cp[index-1]
	if !slices.Contains(tags, entry.Tag) {
		return CpInfo{}, fmt.Errorf("invalid constant pool tag at index %d: %d", index, entry.Tag)
	}
	return entry, nil
}

func getString(i []byte, cp []CpInfo) (string, error) {
	entry, err := getEntry(i, cp, 1)
	return string(entry.Info), err
}

func getNameType(i []byte, cp []CpInfo) (NameAndType, error) {
	name, err := getString(i[0:2], cp)
	if err != nil {
		return NameAndType{}, err
	}
	typ, err := getString(i[2:4], cp)
	return NameAndType{Name: name, Type: typ}, err
}

func getClassNameType(i []byte, cp []CpInfo) (ref struct {
	Class string
	Name  string
	Type  string
}, err error) {
	class, err := getEntry(i[0:2], cp, 7)
	if err != nil {
		return ref, err
	}
	nameAndType, err := getEntry(i[2:4], cp, 12)
	if err != nil {
		return ref, err
	}
	if ref.Class, err = getString(class.Info, cp); err != nil {
		return ref, err
	}
	nt, err := getNameType(nameAndType.Info, cp)
	ref.Name, ref.Type = nt.Name, nt.Type
	return ref, err
}
