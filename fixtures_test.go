package classfileparser

import (
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

const fixturesRoot = "testdata/java/class"

func openFixture(t *testing.T, name string) (*ClassFile, *ClassStruct) {
	t.Helper()

	path := filepath.Join(fixturesRoot, name+".class")
	fh, err := os.Open(path)
	if err != nil {
		t.Fatalf("failed to open %s – compile the fixtures with testdata/java/build.sh: %v", path, err)
	}
	defer fh.Close()

	cf, err := Open(fh)
	if err != nil {
		t.Fatalf("Open %s failed: %v", name, err)
	}
	cs, err := cf.GetClassFile()
	if err != nil {
		t.Fatalf("GetClassFile %s failed: %v", name, err)
	}
	return cf, cs
}

func fixture(t *testing.T, name string) *ClassStruct {
	t.Helper()
	_, cs := openFixture(t, name)
	return cs
}

func findMethod(t *testing.T, cs *ClassStruct, name string) Method {
	t.Helper()
	for _, m := range cs.Methods {
		if m.Name == name {
			return m
		}
	}
	t.Fatalf("method %s not found in %s", name, cs.ThisClass)
	return Method{}
}

func findField(t *testing.T, cs *ClassStruct, name string) Field {
	t.Helper()
	for _, f := range cs.Fields {
		if f.Name == name {
			return f
		}
	}
	t.Fatalf("field %s not found in %s", name, cs.ThisClass)
	return Field{}
}

func codeOf(t *testing.T, m Method) Code {
	t.Helper()
	for _, attr := range m.Attributes {
		if code, ok := attr.(Code); ok {
			return code
		}
	}
	t.Fatalf("method %s has no Code attribute", m.Name)
	return Code{}
}

func hasAttribute[T any](attrs []Attribute) bool {
	for _, attr := range attrs {
		if _, ok := attr.(T); ok {
			return true
		}
	}
	return false
}

func signatureOf(attrs []Attribute) Signature {
	for _, attr := range attrs {
		if s, ok := attr.(Signature); ok {
			return s
		}
	}
	return ""
}

func checkInstructions(t *testing.T, cs *ClassStruct, want map[string][]interface{}) {
	t.Helper()
	for name, instructions := range want {
		t.Run(name, func(t *testing.T) {
			if got := codeOf(t, findMethod(t, cs, name)).Code; !reflect.DeepEqual(got, instructions) {
				t.Fatalf("got %#v\nwant %#v", got, instructions)
			}
		})
	}
}

func TestFixturesParse(t *testing.T) {
	count := 0
	err := filepath.WalkDir(fixturesRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".class") {
			return err
		}
		count++
		openFixture(t, strings.TrimSuffix(strings.TrimPrefix(path, fixturesRoot+string(filepath.Separator)), ".class"))
		return nil
	})
	if err != nil {
		t.Fatalf("failed to walk %s – compile the fixtures with testdata/java/build.sh: %v", fixturesRoot, err)
	}
	if count < 27 {
		t.Fatalf("found %d compiled fixtures, want at least 27", count)
	}
}

func TestFixtureMembers(t *testing.T) {
	cs := fixture(t, "testcases/Members")

	if cs.Version.MajorVersion != 61 || cs.ThisClass != "testcases/Members" || cs.SuperClass != "java/lang/Object" || len(cs.Interfaces) != 0 {
		t.Fatalf("unexpected header: %+v", cs)
	}
	if want := []string{"ACC_PUBLIC", "ACC_SUPER"}; !reflect.DeepEqual(cs.Access, want) {
		t.Fatalf("got class flags %v, want %v", cs.Access, want)
	}
	if !hasAttribute[SourceFile](cs.Attributes) {
		t.Fatalf("missing SourceFile attribute")
	}

	fields := map[string]struct {
		access []string
		typ    string
	}{
		"ANSWER":  {[]string{"ACC_PUBLIC", "ACC_STATIC", "ACC_FINAL"}, "I"},
		"name":    {[]string{"ACC_PRIVATE"}, "Ljava/lang/String;"},
		"counter": {[]string{"ACC_PROTECTED", "ACC_VOLATILE"}, "J"},
		"flag":    {[]string{"ACC_TRANSIENT"}, "Z"},
	}
	for name, want := range fields {
		f := findField(t, cs, name)
		if !reflect.DeepEqual(f.Access, want.access) || f.Type != want.typ {
			t.Errorf("field %s: got %v %s, want %v %s", name, f.Access, f.Type, want.access, want.typ)
		}
	}
	if !hasAttribute[ConstantValue](findField(t, cs, "ANSWER").Attributes) {
		t.Errorf("missing ConstantValue attribute on ANSWER")
	}

	methods := map[string]struct {
		access     []string
		params     []string
		returnType string
	}{
		"<init>":     {[]string{"ACC_PUBLIC"}, []string{"Ljava/lang/String;"}, "V"},
		"greet":      {[]string{"ACC_PUBLIC"}, []string{"Ljava/lang/String;"}, "Ljava/lang/String;"},
		"count":      {[]string{"ACC_PUBLIC", "ACC_STATIC", "ACC_VARARGS"}, []string{"[Ljava/lang/String;"}, "I"},
		"lock":       {[]string{"ACC_PRIVATE", "ACC_SYNCHRONIZED"}, nil, "V"},
		"nativeCall": {[]string{"ACC_PROTECTED", "ACC_NATIVE"}, nil, "V"},
		"compute":    {[]string{"ACC_PUBLIC", "ACC_FINAL"}, []string{"I", "J", "F", "D", "C", "B", "S", "Z", "Ljava/lang/Object;"}, "D"},
	}
	for name, want := range methods {
		m := findMethod(t, cs, name)
		if !reflect.DeepEqual(m.Access, want.access) || !reflect.DeepEqual(m.ParamsTypes, want.params) || m.ReturnType != want.returnType {
			t.Errorf("method %s: got %v %v %s, want %v %v %s", name, m.Access, m.ParamsTypes, m.ReturnType, want.access, want.params, want.returnType)
		}
	}
	if !hasAttribute[MethodParameters](findMethod(t, cs, "greet").Attributes) {
		t.Errorf("missing MethodParameters attribute on greet")
	}
	if hasAttribute[Code](findMethod(t, cs, "nativeCall").Attributes) {
		t.Errorf("native method must not have a Code attribute")
	}
}

func TestFixtureInterface(t *testing.T) {
	cs := fixture(t, "testcases/InterfaceType")

	if want := []string{"ACC_PUBLIC", "ACC_INTERFACE", "ACC_ABSTRACT"}; !reflect.DeepEqual(cs.Access, want) {
		t.Fatalf("got class flags %v, want %v", cs.Access, want)
	}
	if f := findField(t, cs, "UNIT"); !reflect.DeepEqual(f.Access, []string{"ACC_PUBLIC", "ACC_STATIC", "ACC_FINAL"}) {
		t.Errorf("got UNIT flags %v", f.Access)
	}

	methods := map[string][]string{
		"area":     {"ACC_PUBLIC", "ACC_ABSTRACT"},
		"describe": {"ACC_PUBLIC"},
		"unit":     {"ACC_PUBLIC", "ACC_STATIC"},
		"twice":    {"ACC_PRIVATE"},
	}
	for name, want := range methods {
		if m := findMethod(t, cs, name); !reflect.DeepEqual(m.Access, want) {
			t.Errorf("method %s: got %v, want %v", name, m.Access, want)
		}
	}
	if hasAttribute[Code](findMethod(t, cs, "area").Attributes) {
		t.Errorf("abstract method must not have a Code attribute")
	}
}

func TestFixtureConstants(t *testing.T) {
	cs := fixture(t, "testcases/Constants")

	if !hasAttribute[ConstantValue](findField(t, cs, "FIELD").Attributes) {
		t.Errorf("missing ConstantValue attribute on FIELD")
	}
	checkInstructions(t, cs, map[string][]interface{}{
		"intConstant":    {Ldc{Value: int32(100000)}, Ireturn{}},
		"floatConstant":  {Ldc{Value: float32(3.25)}, Freturn{}},
		"longConstant":   {Ldc2W{Value: int64(1234567890123)}, Lreturn{}},
		"doubleConstant": {Ldc2W{Value: float64(2.5)}, Dreturn{}},
		"stringConstant": {Ldc{Value: "constant"}, Areturn{}},
		"classConstant":  {Ldc{Value: Class("java/lang/String")}, Areturn{}},
		"byteConstant":   {Bipush{Byte: -100}, Ireturn{}},
		"shortConstant":  {Sipush{Short: 1000}, Ireturn{}},
		"smallLong":      {Lconst1{}, Lreturn{}},
	})
}

func TestFixtureManyConstants(t *testing.T) {
	code := codeOf(t, findMethod(t, fixture(t, "testcases/ManyConstants"), "all")).Code

	ldcW := 0
	for _, instr := range code {
		if _, ok := instr.(LdcW); ok {
			ldcW++
		}
	}
	if ldcW != 176 {
		t.Errorf("got %d ldc_w instructions, want 176", ldcW)
	}
	if !reflect.DeepEqual(code[4], Ldc{Value: "c0"}) {
		t.Errorf("got %#v, want ldc \"c0\"", code[4])
	}
	if !reflect.DeepEqual(code[len(code)-3], LdcW{Value: "c299"}) {
		t.Errorf("got %#v, want ldc_w \"c299\"", code[len(code)-3])
	}
}

func TestFixtureInvocations(t *testing.T) {
	cs := fixture(t, "testcases/Invocations")

	toString := Methodref{Class: "java/lang/Object", Name: "toString", Type: "()Ljava/lang/String;"}
	total := Fieldref{Class: "testcases/Invocations", Name: "total", Type: "I"}
	count := Fieldref{Class: "testcases/Invocations", Name: "count", Type: "I"}

	if want := []string{"testcases/Greeter"}; !reflect.DeepEqual(cs.Interfaces, want) {
		t.Fatalf("got interfaces %v, want %v", cs.Interfaces, want)
	}
	checkInstructions(t, cs, map[string][]interface{}{
		"<init>": {
			Aload0{},
			Invokespecial{Methodref: Methodref{Class: "java/lang/Object", Name: "<init>", Type: "()V"}},
			Return{},
		},
		"fields": {
			Aload0{}, Getstatic{Fieldref: total}, Putfield{Fieldref: count},
			Aload0{}, Getfield{Fieldref: count}, Putstatic{Fieldref: total},
			Aload0{}, Getfield{Fieldref: count}, Ireturn{},
		},
		"virtualCall": {Aload0{}, Invokevirtual{Methodref: toString}, Areturn{}},
		"superCall":   {Aload0{}, Invokespecial{Methodref: toString}, Areturn{}},
		"staticCall": {
			IconstM1{},
			Invokestatic{Methodref: Methodref{Class: "java/lang/Math", Name: "abs", Type: "(I)I"}},
			Ireturn{},
		},
		"interfaceCall": {
			Aload1{},
			Invokeinterface{InterfaceMethodref: InterfaceMethodref{Class: "java/util/List", Name: "size", Type: "()I"}, Count: 1},
			Ireturn{},
		},
		"staticInterfaceCall": {
			Invokestatic{Methodref: Methodref{Class: "java/util/Comparator", Name: "naturalOrder", Type: "()Ljava/util/Comparator;"}, IsInterface: true},
			Areturn{},
		},
		"defaultSuperCall": {
			Aload0{},
			Invokespecial{Methodref: Methodref{Class: "testcases/Greeter", Name: "greet", Type: "()Ljava/lang/String;"}, IsInterface: true},
			Areturn{},
		},
		"lambda": {Invokedynamic{}, Areturn{}},
	})
	if !hasAttribute[BootstrapMethods](cs.Attributes) {
		t.Errorf("missing BootstrapMethods attribute")
	}
}

func TestFixtureSwitches(t *testing.T) {
	cs := fixture(t, "testcases/Switches")

	if got, want := codeOf(t, findMethod(t, cs, "tableSwitch")).Code[1], (Tableswitch{Default: 33, Low: 0, High: 2, Offsets: []int32{27, 29, 31}}); !reflect.DeepEqual(got, want) {
		t.Errorf("got %#v, want %#v", got, want)
	}

	pairs := []struct{ Match, Offset int32 }{{10, 35}, {50, 38}, {100, 41}}
	if got, want := codeOf(t, findMethod(t, cs, "lookupSwitch")).Code[1], (Lookupswitch{Default: 44, Pairs: pairs}); !reflect.DeepEqual(got, want) {
		t.Errorf("got %#v, want %#v", got, want)
	}

	var switches []Lookupswitch
	var calls []string
	for _, instr := range codeOf(t, findMethod(t, cs, "stringSwitch")).Code {
		switch i := instr.(type) {
		case Lookupswitch:
			switches = append(switches, i)
		case Invokevirtual:
			calls = append(calls, i.Class+"."+i.Name)
		}
	}
	want := []Lookupswitch{
		{Default: 53, Pairs: []struct{ Match, Offset int32 }{{92909918, 28}, {105858401, 42}}},
		{Default: 30, Pairs: []struct{ Match, Offset int32 }{{0, 26}, {1, 28}}},
	}
	if !reflect.DeepEqual(switches, want) {
		t.Errorf("got %#v, want %#v", switches, want)
	}
	if want := []string{"java/lang/String.hashCode", "java/lang/String.equals", "java/lang/String.equals"}; !reflect.DeepEqual(calls, want) {
		t.Errorf("got calls %v, want %v", calls, want)
	}
}

func TestFixtureArrays(t *testing.T) {
	cs := fixture(t, "testcases/ArrayOps")

	checkInstructions(t, cs, map[string][]interface{}{
		"ints":    {Iload0{}, Newarray{Type: 10}, Areturn{}},
		"strings": {Iload0{}, Anewarray("java/lang/String"), Areturn{}},
		"matrix":  {Iload0{}, Iload1{}, Multianewarray{Class: "[[I", Dimension: 2}, Areturn{}},
		"jagged":  {Iload0{}, Anewarray("[Ljava/lang/String;"), Areturn{}},
		"length":  {Aload0{}, Arraylength{}, Ireturn{}},
	})

	if m := findMethod(t, cs, "length"); !reflect.DeepEqual(m.ParamsTypes, []string{"[[I"}) || m.ReturnType != "I" {
		t.Errorf("length: got %v %s", m.ParamsTypes, m.ReturnType)
	}
	if m := findMethod(t, cs, "jagged"); m.ReturnType != "[[Ljava/lang/String;" {
		t.Errorf("jagged: got return type %s", m.ReturnType)
	}
}

func TestFixtureTypeChecks(t *testing.T) {
	cs := fixture(t, "testcases/TypeChecks")

	checkInstructions(t, cs, map[string][]interface{}{
		"isString": {Aload0{}, Instanceof("java/lang/String"), Ireturn{}},
		"cast":     {Aload0{}, Checkcast("java/lang/String"), Areturn{}},
		"locked": {
			Aload0{}, Dup{}, Astore1{}, Monitorenter{},
			Aload0{}, Invokevirtual{Methodref: Methodref{Class: "java/lang/Object", Name: "notify", Type: "()V"}},
			Aload1{}, Monitorexit{}, Goto{Offset: 8},
			Astore2{}, Aload1{}, Monitorexit{}, Aload2{}, Athrow{},
			Return{},
		},
		"fail": {
			New("java/lang/IllegalStateException"), Dup{},
			Invokespecial{Methodref: Methodref{Class: "java/lang/IllegalStateException", Name: "<init>", Type: "()V"}},
			Athrow{},
		},
	})

	want := []ExceptionTableEntry{{StartPc: 4, EndPc: 10, HandlerPc: 13}, {StartPc: 13, EndPc: 16, HandlerPc: 13}}
	if got := codeOf(t, findMethod(t, cs, "locked")).ExceptionTable; !reflect.DeepEqual(got, want) {
		t.Errorf("got exception table %+v, want %+v", got, want)
	}
}

func TestFixtureWideLocals(t *testing.T) {
	code := codeOf(t, findMethod(t, fixture(t, "testcases/WideLocals"), "wide")).Code

	found := false
	for _, instr := range code {
		if instr == (Wide{OpCode: 0x36, LocalIndex: 256}) {
			found = true
		}
	}
	if !found {
		t.Errorf("missing wide istore 256")
	}

	want := []interface{}{Wide{OpCode: 0x84, LocalIndex: 260, Const: 1}, Wide{OpCode: 0x15, LocalIndex: 260}, Ireturn{}}
	if got := code[len(code)-3:]; !reflect.DeepEqual(got, want) {
		t.Errorf("got %#v, want %#v", got, want)
	}
}

func TestFixtureExceptions(t *testing.T) {
	cf, cs := openFixture(t, "testcases/Exceptions")
	cp, err := cf.GetConstantPool()
	if err != nil {
		t.Fatalf("GetConstantPool failed: %v", err)
	}

	if !hasAttribute[Exceptions](findMethod(t, cs, "declares").Attributes) {
		t.Errorf("missing Exceptions attribute on declares")
	}

	table := codeOf(t, findMethod(t, cs, "catches")).ExceptionTable
	if len(table) != 4 {
		t.Fatalf("got %d exception table entries, want 4: %+v", len(table), table)
	}
	if cp[table[0].CatchType] != Class("java/lang/NumberFormatException") {
		t.Errorf("got catch type %#v, want java/lang/NumberFormatException", cp[table[0].CatchType])
	}
	table[0].CatchType = 0
	want := []ExceptionTableEntry{
		{StartPc: 0, EndPc: 5, HandlerPc: 10},
		{StartPc: 0, EndPc: 5, HandlerPc: 18},
		{StartPc: 10, EndPc: 13, HandlerPc: 18},
		{StartPc: 18, EndPc: 20, HandlerPc: 18},
	}
	if !reflect.DeepEqual(table, want) {
		t.Errorf("got exception table %+v, want %+v", table, want)
	}
}

func TestFixtureAnnotations(t *testing.T) {
	cs := fixture(t, "testcases/Annotations")

	targets := map[string][]Attribute{
		"class":  cs.Attributes,
		"field":  findField(t, cs, "field").Attributes,
		"method": findMethod(t, cs, "method").Attributes,
	}
	for target, attrs := range targets {
		if !hasAttribute[RuntimeVisibleAnnotations](attrs) || !hasAttribute[RuntimeInvisibleAnnotations](attrs) {
			t.Errorf("%s: missing visible or invisible annotations", target)
		}
	}

	attrs := targets["method"]
	if !hasAttribute[RuntimeVisibleParameterAnnotations](attrs) || !hasAttribute[RuntimeInvisibleParameterAnnotations](attrs) {
		t.Errorf("method: missing visible or invisible parameter annotations")
	}
}

func TestFixtureGenerics(t *testing.T) {
	cs := fixture(t, "testcases/Generics")

	cases := map[string]struct {
		got  Signature
		want Signature
	}{
		"class":  {signatureOf(cs.Attributes), "<T:Ljava/lang/Number;>Ljava/lang/Object;"},
		"field":  {signatureOf(findField(t, cs, "value").Attributes), "TT;"},
		"method": {signatureOf(findMethod(t, cs, "process").Attributes), "<U::Ljava/lang/Comparable<TU;>;>(Ljava/util/List<TU;>;TT;)Ljava/util/List<TU;>;"},
	}
	for target, c := range cases {
		if c.got != c.want {
			t.Errorf("%s: got signature %q, want %q", target, c.got, c.want)
		}
	}

	if !hasAttribute[LocalVariableTypeTable](codeOf(t, findMethod(t, cs, "process")).Attributes) {
		t.Errorf("missing LocalVariableTypeTable attribute in process")
	}
}

func TestFixtureNestedClasses(t *testing.T) {
	outer := fixture(t, "testcases/NestedClasses")
	if !hasAttribute[InnerClasses](outer.Attributes) || !hasAttribute[NestMembers](outer.Attributes) {
		t.Errorf("outer class: missing InnerClasses or NestMembers attribute")
	}

	inner := fixture(t, "testcases/NestedClasses$Inner")
	if inner.ThisClass != "testcases/NestedClasses$Inner" || !hasAttribute[NestHost](inner.Attributes) {
		t.Errorf("inner class: got %s, NestHost %v", inner.ThisClass, hasAttribute[NestHost](inner.Attributes))
	}
	if f := findField(t, inner, "this$0"); !reflect.DeepEqual(f.Access, []string{"ACC_FINAL", "ACC_SYNTHETIC"}) || f.Type != "Ltestcases/NestedClasses;" {
		t.Errorf("this$0: got %v %s", f.Access, f.Type)
	}
	if m := findMethod(t, inner, "outer"); m.Access != nil {
		t.Errorf("outer: got flags %v, want none", m.Access)
	}

	if nested := fixture(t, "testcases/NestedClasses$Nested"); nested.ThisClass != "testcases/NestedClasses$Nested" {
		t.Errorf("nested class: got %s", nested.ThisClass)
	}

	for _, name := range []string{"testcases/NestedClasses$1Local", "testcases/NestedClasses$1"} {
		if cs := fixture(t, name); !hasAttribute[EnclosingMethod](cs.Attributes) {
			t.Errorf("%s: missing EnclosingMethod attribute", name)
		}
	}
}

func TestFixtureEnum(t *testing.T) {
	cs := fixture(t, "testcases/EnumType")

	if want := []string{"ACC_PUBLIC", "ACC_FINAL", "ACC_SUPER", "ACC_ENUM"}; !reflect.DeepEqual(cs.Access, want) || cs.SuperClass != "java/lang/Enum" {
		t.Fatalf("got %v extends %s", cs.Access, cs.SuperClass)
	}
	if got, want := signatureOf(cs.Attributes), Signature("Ljava/lang/Enum<Ltestcases/EnumType;>;"); got != want {
		t.Errorf("got signature %q, want %q", got, want)
	}

	fields := map[string][]string{
		"FIRST":   {"ACC_PUBLIC", "ACC_STATIC", "ACC_FINAL", "ACC_ENUM"},
		"SECOND":  {"ACC_PUBLIC", "ACC_STATIC", "ACC_FINAL", "ACC_ENUM"},
		"$VALUES": {"ACC_PRIVATE", "ACC_STATIC", "ACC_FINAL", "ACC_SYNTHETIC"},
	}
	for name, want := range fields {
		if f := findField(t, cs, name); !reflect.DeepEqual(f.Access, want) {
			t.Errorf("field %s: got %v, want %v", name, f.Access, want)
		}
	}
	if m := findMethod(t, cs, "$values"); !reflect.DeepEqual(m.Access, []string{"ACC_PRIVATE", "ACC_STATIC", "ACC_SYNTHETIC"}) {
		t.Errorf("$values: got %v", m.Access)
	}
}

func TestFixtureSealed(t *testing.T) {
	if cs := fixture(t, "testcases/SealedType"); !hasAttribute[PermittedSubclasses](cs.Attributes) {
		t.Errorf("missing PermittedSubclasses attribute")
	}

	children := map[string][]string{
		"testcases/SealedType$First":  {"ACC_PUBLIC", "ACC_FINAL", "ACC_SUPER"},
		"testcases/SealedType$Second": {"ACC_PUBLIC", "ACC_SUPER"},
	}
	for name, want := range children {
		if cs := fixture(t, name); !reflect.DeepEqual(cs.Access, want) || cs.SuperClass != "testcases/SealedType" {
			t.Errorf("%s: got %v extends %s", name, cs.Access, cs.SuperClass)
		}
	}
}

func TestFixtureRecord(t *testing.T) {
	cs := fixture(t, "testcases/RecordType")

	if want := []string{"ACC_PUBLIC", "ACC_FINAL", "ACC_SUPER"}; !reflect.DeepEqual(cs.Access, want) || cs.SuperClass != "java/lang/Record" {
		t.Fatalf("got %v extends %s", cs.Access, cs.SuperClass)
	}
	for name, typ := range map[string]string{"x": "I", "y": "Ljava/lang/String;"} {
		if f := findField(t, cs, name); !reflect.DeepEqual(f.Access, []string{"ACC_PRIVATE", "ACC_FINAL"}) || f.Type != typ {
			t.Errorf("field %s: got %v %s", name, f.Access, f.Type)
		}
	}
	for _, name := range []string{"toString", "hashCode", "equals"} {
		if m := findMethod(t, cs, name); !reflect.DeepEqual(m.Access, []string{"ACC_PUBLIC", "ACC_FINAL"}) {
			t.Errorf("method %s: got %v", name, m.Access)
		}
	}
}

func TestFixtureDeprecated(t *testing.T) {
	cs := fixture(t, "testcases/DeprecatedMembers")

	targets := map[string][]Attribute{
		"class":  cs.Attributes,
		"field":  findField(t, cs, "OLD").Attributes,
		"method": findMethod(t, cs, "old").Attributes,
	}
	for target, attrs := range targets {
		if !hasAttribute[Deprecated](attrs) {
			t.Errorf("%s: missing Deprecated attribute", target)
		}
	}
}

func TestFixtureModule(t *testing.T) {
	cs := fixture(t, "module-info")

	if !reflect.DeepEqual(cs.Access, []string{"ACC_MODULE"}) || cs.ThisClass != "module-info" || cs.SuperClass != "" {
		t.Fatalf("got %v %s extends %q", cs.Access, cs.ThisClass, cs.SuperClass)
	}
	if !hasAttribute[ModuleInfo](cs.Attributes) {
		t.Errorf("missing Module attribute")
	}
}
