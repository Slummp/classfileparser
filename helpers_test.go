package classfileparser

import (
	"bytes"
	"encoding/binary"
	"math"
)

func u1(v uint8) []byte {
	return []byte{v}
}

func u2(v uint16) []byte {
	return binary.BigEndian.AppendUint16(nil, v)
}

func u4(v uint32) []byte {
	return binary.BigEndian.AppendUint32(nil, v)
}

func cat(parts ...[]byte) []byte {
	return bytes.Join(parts, nil)
}

func cpUtf8(s string) []byte {
	return cat(u1(1), u2(uint16(len(s))), []byte(s))
}

func cpInteger(v int32) []byte {
	return cat(u1(3), u4(uint32(v)))
}

func cpFloat(v float32) []byte {
	return cat(u1(4), u4(math.Float32bits(v)))
}

func cpLong(v int64) []byte {
	return cat(u1(5), u4(uint32(uint64(v)>>32)), u4(uint32(v)))
}

func cpDouble(v float64) []byte {
	bits := math.Float64bits(v)
	return cat(u1(6), u4(uint32(bits>>32)), u4(uint32(bits)))
}

func cpClass(name uint16) []byte {
	return cat(u1(7), u2(name))
}

func cpString(value uint16) []byte {
	return cat(u1(8), u2(value))
}

func cpFieldref(class, nameAndType uint16) []byte {
	return cat(u1(9), u2(class), u2(nameAndType))
}

func cpMethodref(class, nameAndType uint16) []byte {
	return cat(u1(10), u2(class), u2(nameAndType))
}

func cpInterfaceMethodref(class, nameAndType uint16) []byte {
	return cat(u1(11), u2(class), u2(nameAndType))
}

func cpNameAndType(name, descriptor uint16) []byte {
	return cat(u1(12), u2(name), u2(descriptor))
}

func cpMethodHandle(kind uint8, reference uint16) []byte {
	return cat(u1(15), u1(kind), u2(reference))
}

func cpMethodType(descriptor uint16) []byte {
	return cat(u1(16), u2(descriptor))
}

func cpDynamic(bootstrap, nameAndType uint16) []byte {
	return cat(u1(17), u2(bootstrap), u2(nameAndType))
}

func cpInvokeDynamic(bootstrap, nameAndType uint16) []byte {
	return cat(u1(18), u2(bootstrap), u2(nameAndType))
}

func cpModule(name uint16) []byte {
	return cat(u1(19), u2(name))
}

func cpPackage(name uint16) []byte {
	return cat(u1(20), u2(name))
}

type testAttribute struct {
	name uint16
	info []byte
}

type testMember struct {
	access     uint16
	name       uint16
	descriptor uint16
	attributes []testAttribute
}

type testClass struct {
	minor      uint16
	major      uint16
	pool       [][]byte
	access     uint16
	this       uint16
	super      uint16
	interfaces []uint16
	fields     []testMember
	methods    []testMember
	attributes []testAttribute
}

func minimalClass() testClass {
	return testClass{
		major: 61,
		pool:  [][]byte{cpUtf8("Foo"), cpClass(1)},
		this:  2,
	}
}

func (c testClass) bytes() []byte {
	slots := 1
	for _, entry := range c.pool {
		slots++
		if entry[0] == 5 || entry[0] == 6 {
			slots++
		}
	}

	b := cat(u4(0xCAFEBABE), u2(c.minor), u2(c.major), u2(uint16(slots)))
	for _, entry := range c.pool {
		b = cat(b, entry)
	}
	b = cat(b, u2(c.access), u2(c.this), u2(c.super), u2(uint16(len(c.interfaces))))
	for _, i := range c.interfaces {
		b = cat(b, u2(i))
	}
	b = cat(b, members(c.fields), members(c.methods), attributes(c.attributes))
	return b
}

func members(ms []testMember) []byte {
	b := u2(uint16(len(ms)))
	for _, m := range ms {
		b = cat(b, u2(m.access), u2(m.name), u2(m.descriptor), attributes(m.attributes))
	}
	return b
}

func attributes(as []testAttribute) []byte {
	b := u2(uint16(len(as)))
	for _, a := range as {
		b = cat(b, u2(a.name), u4(uint32(len(a.info))), a.info)
	}
	return b
}

func codeInfo(bytecode []byte, exceptions []ExceptionTableEntry, nested ...testAttribute) []byte {
	b := cat(u2(1), u2(1), u4(uint32(len(bytecode))), bytecode, u2(uint16(len(exceptions))))
	for _, e := range exceptions {
		b = cat(b, u2(e.StartPc), u2(e.EndPc), u2(e.HandlerPc), u2(e.CatchType))
	}
	return cat(b, attributes(nested))
}
