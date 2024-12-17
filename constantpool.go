package classfileparser

import (
	"encoding/binary"
	"math"
)

type ConstantPool struct {
	Utf8               map[uint16]string             // 01
	Integer            map[uint16]int32              // 03
	Float              map[uint16]float32            // 04
	Long               map[uint16]int64              // 05
	Double             map[uint16]float64            // 06
	Class              map[uint16]string             // 07
	String             map[uint16]string             // 08
	Fieldref           map[uint16]Fieldref           // 09
	Methodref          map[uint16]Methodref          // 10
	InterfaceMethodref map[uint16]InterfaceMethodref // 11
	NameAndType        map[uint16]NameAndType        // 12
	MethodHandle       map[uint16]MethodHandle       // 15
	MethodType         map[uint16]string             // 16
	Dynamic            map[uint16]Dynamic            // 17
	InvokeDynamic      map[uint16]InvokeDynamic      // 18
	Module             map[uint16]string             // 19
	Package            map[uint16]string             // 20
}

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

func (cf *ClassFile) GetConstantPool() (*ConstantPool, error) {
	cp := &ConstantPool{
		Utf8:               make(map[uint16]string),
		Integer:            make(map[uint16]int32),
		Float:              make(map[uint16]float32),
		Long:               make(map[uint16]int64),
		Double:             make(map[uint16]float64),
		Class:              make(map[uint16]string),
		String:             make(map[uint16]string),
		Fieldref:           make(map[uint16]Fieldref),
		Methodref:          make(map[uint16]Methodref),
		InterfaceMethodref: make(map[uint16]InterfaceMethodref),
		NameAndType:        make(map[uint16]NameAndType),
		MethodHandle:       make(map[uint16]MethodHandle),
		MethodType:         make(map[uint16]string),
		Dynamic:            make(map[uint16]Dynamic),
		InvokeDynamic:      make(map[uint16]InvokeDynamic),
		Module:             make(map[uint16]string),
		Package:            make(map[uint16]string),
	}
	for i, cpItem := range cf.ConstantPool {
		switch cpItem.Tag {
		case 1:
			cp.Utf8[uint16(i+1)] = string(cpItem.Info)
			break
		case 3:
			cp.Integer[uint16(i+1)] = int32(binary.BigEndian.Uint32(cpItem.Info))
			break
		case 4:
			cp.Float[uint16(i+1)] = math.Float32frombits(binary.BigEndian.Uint32(cpItem.Info))
			break
		case 5:
			cp.Long[uint16(i+1)] = int64(binary.BigEndian.Uint64(cpItem.Info))
			break
		case 6:
			cp.Double[uint16(i+1)] = math.Float64frombits(binary.BigEndian.Uint64(cpItem.Info))
			break
		case 7:
			cp.Class[uint16(i+1)] = getString(cpItem.Info, cf.ConstantPool)
			break
		case 8:
			cp.String[uint16(i+1)] = getString(cpItem.Info, cf.ConstantPool)
			break
		case 9:
			cp.Fieldref[uint16(i+1)] = getClassNameType(cpItem.Info, cf.ConstantPool)
			break
		case 10:
			cp.Methodref[uint16(i+1)] = getClassNameType(cpItem.Info, cf.ConstantPool)
			break
		case 11:
			cp.InterfaceMethodref[uint16(i+1)] = getClassNameType(cpItem.Info, cf.ConstantPool)
			break
		case 12:
			cp.NameAndType[uint16(i+1)] =
				getNameType(cpItem.Info, cf.ConstantPool)

			break
		case 15:
			index := int8(binary.BigEndian.Uint16(cpItem.Info[0:2]))
			cnt := getClassNameType(cpItem.Info[2:4], cf.ConstantPool)
			cp.MethodHandle[uint16(i+1)] = MethodHandle{
				Kind: map[int8]string{
					1: "getField",
					2: "getStatic",
					3: "putField",
					4: "putStatic",
					5: "invokeVirtual",
					6: "invokeStatic",
					7: "invokeSpecial",
					8: "newInvokeSpecial",
					9: "invokeInterface",
				}[index],
				Class: cnt.Class,
				Name:  cnt.Name,
				Type:  cnt.Type,
			}
			break
		case 16:
			cp.MethodType[uint16(i+1)] = getString(cpItem.Info, cf.ConstantPool)
			break
		case 17:
			cp.Dynamic[uint16(i+1)] = Dynamic{} // TODO
			break
		case 18:
			cp.InvokeDynamic[uint16(i+1)] = InvokeDynamic{} // TODO
			break
		case 19:
			cp.Module[uint16(i+1)] = getString(cpItem.Info, cf.ConstantPool)
			break
		case 20:
			cp.Package[uint16(i+1)] = getString(cpItem.Info, cf.ConstantPool)
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
