package classfileparser

import (
	"reflect"
	"testing"
)

func TestFindFlags(t *testing.T) {
	cases := map[string]struct {
		t     Type
		value uint16
		want  []string
	}{
		"all class flags": {ClassT, 0xFFFF, []string{
			"ACC_PUBLIC", "ACC_FINAL", "ACC_SUPER", "ACC_INTERFACE", "ACC_ABSTRACT",
			"ACC_SYNTHETIC", "ACC_ANNOTATION", "ACC_ENUM", "ACC_MODULE",
		}},
		"all field flags": {FieldT, 0xFFFF, []string{
			"ACC_PUBLIC", "ACC_PRIVATE", "ACC_PROTECTED", "ACC_STATIC", "ACC_FINAL",
			"ACC_VOLATILE", "ACC_TRANSIENT", "ACC_SYNTHETIC", "ACC_ENUM",
		}},
		"all method flags": {MethodT, 0xFFFF, []string{
			"ACC_PUBLIC", "ACC_PRIVATE", "ACC_PROTECTED", "ACC_STATIC", "ACC_FINAL",
			"ACC_SYNCHRONIZED", "ACC_BRIDGE", "ACC_VARARGS", "ACC_NATIVE", "ACC_ABSTRACT",
			"ACC_STRICT", "ACC_SYNTHETIC",
		}},
		"all nested class flags": {NestedT, 0xFFFF, []string{
			"ACC_PUBLIC", "ACC_PRIVATE", "ACC_PROTECTED", "ACC_STATIC", "ACC_FINAL",
			"ACC_INTERFACE", "ACC_ABSTRACT", "ACC_SYNTHETIC", "ACC_ANNOTATION", "ACC_ENUM",
		}},
		"some class flags":        {ClassT, 0x0421, []string{"ACC_PUBLIC", "ACC_SUPER", "ACC_ABSTRACT"}},
		"same bit, other meaning": {MethodT, 0x0020, []string{"ACC_SYNCHRONIZED"}},
		"no flag":                 {MethodT, 0x0000, nil},
		"flag of another type":    {ClassT, 0x0002, nil},
		"unknown type":            {Type(42), 0xFFFF, nil},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			if got := findFlags(c.t, c.value); !reflect.DeepEqual(got, c.want) {
				t.Fatalf("got %v, want %v", got, c.want)
			}
		})
	}
}

func TestFindFlagsIsStable(t *testing.T) {
	want := findFlags(MethodT, 0xFFFF)
	for i := 0; i < 100; i++ {
		if got := findFlags(MethodT, 0xFFFF); !reflect.DeepEqual(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
	}
}
