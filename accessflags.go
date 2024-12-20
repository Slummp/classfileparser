package classfileparser

type Type uint8

const (
	ClassT Type = iota
	FieldT
	MethodT
	NestedT
)

var flags map[Type]map[int]string = map[Type]map[int]string{
	ClassT: {
		0x0001: "ACC_PUBLIC",
		0x0010: "ACC_FINAL",
		0x0020: "ACC_SUPER",
		0x0200: "ACC_INTERFACE",
		0x0400: "ACC_ABSTRACT",
		0x1000: "ACC_SYNTHETIC",
		0x2000: "ACC_ANNOTATION",
		0x4000: "ACC_ENUM",
		0x8000: "ACC_MODULE",
	},
	FieldT: {
		0x0001: "ACC_PUBLIC",
		0x0002: "ACC_PRIVATE",
		0x0004: "ACC_PROTECTED",
		0x0008: "ACC_STATIC",
		0x0010: "ACC_FINAL",
		0x0040: "ACC_VOLATILE",
		0x0080: "ACC_TRANSIENT",
		0x1000: "ACC_SYNTHETIC",
		0x4000: "ACC_ENUM",
	},
	MethodT: {
		0x0001: "ACC_PUBLIC",
		0x0002: "ACC_PRIVATE",
		0x0004: "ACC_PROTECTED",
		0x0008: "ACC_STATIC",
		0x0010: "ACC_FINAL",
		0x0020: "ACC_SYNCHRONIZED",
		0x0040: "ACC_BRIDGE",
		0x0080: "ACC_VARARGS",
		0x0100: "ACC_NATIVE",
		0x0400: "ACC_ABSTRACT",
		0x0800: "ACC_STRICT",
		0x1000: "ACC_SYNTHETIC",
	},
	NestedT: {
		0x0001: "ACC_PUBLIC",
		0x0002: "ACC_PRIVATE",
		0x0004: "ACC_PROTECTED",
		0x0008: "ACC_STATIC",
		0x0010: "ACC_FINAL",
		0x0200: "ACC_INTERFACE",
		0x0400: "ACC_ABSTRACT",
		0x1000: "ACC_SYNTHETIC",
		0x2000: "ACC_ANNOTATION",
		0x4000: "ACC_ENUM",
	},
}

func findFlags(t Type, value uint16) []string {
	var ret []string
	for flag, name := range flags[t] {
		if value&uint16(flag) != 0 {
			ret = append(ret, name)
		}
	}
	return ret
}
