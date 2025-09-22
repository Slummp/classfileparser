package classfileparser

import (
	"encoding/binary"
	"math"
)

type ConstantPool map[uint16]interface{}

type Utf8 string

type Class string

type MethodType string

type Module string

type Package string

type Fieldref struct {
	Class string
	Name  string
	Type  string
}

type Methodref struct {
	Class string
	Name  string
	Type  string
}

type InterfaceMethodref struct {
	Class string
	Name  string
	Type  string
}

type NameAndType struct {
	Name string
	Type string
}

type MethodHandle struct {
	Kind  string
	Class string
	Name  string
	Type  string
}

type Dynamic struct {
	Name           string
	Type           string
	BootstrapIndex uint16
}

type InvokeDynamic struct {
	BootstrapIndex uint16
	Name           string
	Type           string
}

func (cf *ClassFile) GetConstantPool() (ConstantPool, error) {
	cp := ConstantPool{}
	for i, cpItem := range cf.ConstantPool {
		switch cpItem.Tag {
		case 1:
			cp[uint16(i+1)] = Utf8(cpItem.Info)
			break
		case 3:
			cp[uint16(i+1)] = int32(binary.BigEndian.Uint32(cpItem.Info))
			break
		case 4:
			cp[uint16(i+1)] = math.Float32frombits(binary.BigEndian.Uint32(cpItem.Info))
			break
		case 5:
			cp[uint16(i+1)] = int64(binary.BigEndian.Uint64(cpItem.Info))
			break
		case 6:
			cp[uint16(i+1)] = math.Float64frombits(binary.BigEndian.Uint64(cpItem.Info))
			break
		case 7:
			cp[uint16(i+1)] = Class(getString(cpItem.Info, cf.ConstantPool))
			break
		case 8:
			cp[uint16(i+1)] = getString(cpItem.Info, cf.ConstantPool)
			break
		case 9:
			cp[uint16(i+1)] = Fieldref(getClassNameType(cpItem.Info, cf.ConstantPool))
			break
		case 10:
			cp[uint16(i+1)] = Methodref(getClassNameType(cpItem.Info, cf.ConstantPool))
			break
		case 11:
			cp[uint16(i+1)] = InterfaceMethodref(getClassNameType(cpItem.Info, cf.ConstantPool))
			break
		case 12:
			cp[uint16(i+1)] =
				getNameType(cpItem.Info, cf.ConstantPool)
			break
		case 15:
			item := cp[uint16(binary.BigEndian.Uint16(cpItem.Info[1:3]))]
			var cnt struct {
				Class string
				Name  string
				Type  string
			}
			switch itemTyped := item.(type) {
			case Fieldref:
				cnt = itemTyped
			case Methodref:
				cnt = itemTyped
			case InterfaceMethodref:
				cnt = itemTyped
			}
			cp[uint16(i+1)] = MethodHandle{
				Kind: map[byte]string{
					1: "getField",
					2: "getStatic",
					3: "putField",
					4: "putStatic",
					5: "invokeVirtual",
					6: "invokeStatic",
					7: "invokeSpecial",
					8: "newInvokeSpecial",
					9: "invokeInterface",
				}[cpItem.Info[0]],
				Class: cnt.Class,
				Name:  cnt.Name,
				Type:  cnt.Type,
			}
			break
		case 16:
			cp[uint16(i+1)] = MethodType(getString(cpItem.Info, cf.ConstantPool))
			break
		case 17:
			cp[uint16(i+1)] = Dynamic{} // TODO
			break
		case 18:
			cp[uint16(i+1)] = InvokeDynamic{} // TODO
			break
		case 19:
			cp[uint16(i+1)] = Module(getString(cpItem.Info, cf.ConstantPool))
			break
		case 20:
			cp[uint16(i+1)] = Package(getString(cpItem.Info, cf.ConstantPool))
			break
		}
	}
	return cp, nil
}

func getString(i []byte, cp []CpInfo) string {
	index := int16(binary.BigEndian.Uint16(i))
	return string(cp[index-1].Info)
}

func getNameType(i []byte, cp []CpInfo) NameAndType {
	return NameAndType{
		Name: getString(i[0:2], cp),
		Type: getString(i[2:4], cp),
	}
}

func getClassNameType(i []byte, cp []CpInfo) struct {
	Class string
	Name  string
	Type  string
} {
	indexNameAndType := cp[int16(binary.BigEndian.Uint16(i[2:4]))-1].Info
	return struct {
		Class string
		Name  string
		Type  string
	}{
		Class: getString(cp[int16(binary.BigEndian.Uint16(i[0:2]))-1].Info, cp),
		Name:  getString(indexNameAndType[0:2], cp),
		Type:  getString(indexNameAndType[2:4], cp),
	}
}
