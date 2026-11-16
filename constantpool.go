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

type constantResolver func(info []byte, cp []CpInfo) (interface{}, error)

var methodHandleKinds = map[byte]string{
	1: "getField",
	2: "getStatic",
	3: "putField",
	4: "putStatic",
	5: "invokeVirtual",
	6: "invokeStatic",
	7: "invokeSpecial",
	8: "newInvokeSpecial",
	9: "invokeInterface",
}

var constantResolvers = map[uint8]constantResolver{
	1: func(info []byte, _ []CpInfo) (interface{}, error) {
		return Utf8(info), nil
	},
	3: func(info []byte, _ []CpInfo) (interface{}, error) {
		return int32(binary.BigEndian.Uint32(info)), nil
	},
	4: func(info []byte, _ []CpInfo) (interface{}, error) {
		return math.Float32frombits(binary.BigEndian.Uint32(info)), nil
	},
	5: func(info []byte, _ []CpInfo) (interface{}, error) {
		return int64(binary.BigEndian.Uint64(info)), nil
	},
	6: func(info []byte, _ []CpInfo) (interface{}, error) {
		return math.Float64frombits(binary.BigEndian.Uint64(info)), nil
	},
	7: func(info []byte, cp []CpInfo) (interface{}, error) {
		name, err := getString(info, cp)
		return Class(name), err
	},
	8: func(info []byte, cp []CpInfo) (interface{}, error) {
		return getString(info, cp)
	},
	9: func(info []byte, cp []CpInfo) (interface{}, error) {
		ref, err := getClassNameType(info, cp)
		return Fieldref(ref), err
	},
	10: func(info []byte, cp []CpInfo) (interface{}, error) {
		ref, err := getClassNameType(info, cp)
		return Methodref(ref), err
	},
	11: func(info []byte, cp []CpInfo) (interface{}, error) {
		ref, err := getClassNameType(info, cp)
		return InterfaceMethodref(ref), err
	},
	12: func(info []byte, cp []CpInfo) (interface{}, error) {
		return getNameType(info, cp)
	},
	15: func(info []byte, cp []CpInfo) (interface{}, error) {
		kind, ok := methodHandleKinds[info[0]]
		if !ok {
			return nil, fmt.Errorf("%w: invalid method handle kind: %d", ErrMalformed, info[0])
		}
		ref, err := getEntry(info[1:3], cp, 9, 10, 11)
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
	},
	16: func(info []byte, cp []CpInfo) (interface{}, error) {
		descriptor, err := getString(info, cp)
		return MethodType(descriptor), err
	},
	17: func([]byte, []CpInfo) (interface{}, error) {
		return Dynamic{}, nil // TODO
	},
	18: func([]byte, []CpInfo) (interface{}, error) {
		return InvokeDynamic{}, nil // TODO
	},
	19: func(info []byte, cp []CpInfo) (interface{}, error) {
		name, err := getString(info, cp)
		return Module(name), err
	},
	20: func(info []byte, cp []CpInfo) (interface{}, error) {
		name, err := getString(info, cp)
		return Package(name), err
	},
}

func resolveConstant(cpItem CpInfo, cp []CpInfo) (interface{}, error) {
	resolve, ok := constantResolvers[cpItem.Tag]
	if !ok {
		return nil, fmt.Errorf("%w: unknown constant pool tag: %d", ErrMalformed, cpItem.Tag)
	}
	return resolve(cpItem.Info, cp)
}

func getEntry(i []byte, cp []CpInfo, tags ...uint8) (CpInfo, error) {
	index := binary.BigEndian.Uint16(i)
	if index == 0 || int(index) > len(cp) {
		return CpInfo{}, fmt.Errorf("%w: invalid constant pool index: %d", ErrMalformed, index)
	}
	entry := cp[index-1]
	if !slices.Contains(tags, entry.Tag) {
		return CpInfo{}, fmt.Errorf("%w: invalid constant pool tag at index %d: %d", ErrMalformed, index, entry.Tag)
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

func getConstant[T any](cp ConstantPool, index uint16) (T, error) {
	var value T
	entry, ok := cp[index]
	if !ok {
		return value, fmt.Errorf("%w: invalid constant pool index: %d", ErrMalformed, index)
	}
	value, ok = entry.(T)
	if !ok {
		return value, fmt.Errorf("%w: invalid constant pool type at index %d: %T", ErrMalformed, index, entry)
	}
	return value, nil
}
