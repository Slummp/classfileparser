package classfileparser

type Attributes struct {
	Code                                 []Code
	ConstantValue                        []ConstantValue
	Deprecated                           []Deprecated
	Exceptions                           []Exceptions
	InnerClasses                         []InnerClasses
	LineNumberTable                      []LineNumberTable
	LocalVariableTable                   []LocalVariableTable
	LocalVariableTypeTable               []LocalVariableTypeTable
	MethodParameters                     []MethodParameters
	RuntimeVisibleAnnotations            []RuntimeVisibleAnnotations
	RuntimeInvisibleAnnotations          []RuntimeInvisibleAnnotations
	RuntimeVisibleParameterAnnotations   []RuntimeVisibleParameterAnnotations
	RuntimeInvisibleParameterAnnotations []RuntimeInvisibleParameterAnnotations
	SourceFile                           []SourceFile
	SourceDebugExtension                 []SourceDebugExtension
	Signature                            []Signature
	StackMapTable                        []StackMapTable
	Synthetic                            []Synthetic
	EnclosingMethod                      []EnclosingMethod
	BootstrapMethods                     []BootstrapMethods
	Module                               []Module
	ModulePackages                       []ModulePackages
	NestHost                             []NestHost
	NestMembers                          []NestMembers
	PermittedSubclasses                  []PermittedSubclasses
}

// Code : Contient le bytecode et les informations relatives à la méthode.
type Code struct {
	MaxStack       uint16
	MaxLocals      uint16
	Code           []byte
	ExceptionTable []ExceptionTableEntry
	Attributes     []AttributeInfo
}

// ConstantValue : Attribut d'un champ constant
type ConstantValue struct {
	ConstantValueIndex uint16
}

// Deprecated : Attribut indiquant qu'une méthode est obsolète
type Deprecated struct{}

// Exceptions : Attribut indiquant les exceptions vérifiées lancées par une méthode
type Exceptions struct {
	NumberOfExceptions  uint16
	ExceptionIndexTable []uint16
}

// InnerClasses : Attribut contenant des informations sur les classes internes
type InnerClasses struct {
	NumberOfClasses uint16
	InnerClassInfo  []InnerClassInfo
}

// LineNumberTable : Attribut de correspondance entre les lignes source et le bytecode
type LineNumberTable struct {
	LineNumberTableLength uint16
	LineNumberTable       []LineNumberTableEntry
}

// LocalVariableTable : Attribut contenant des informations sur les variables locales
type LocalVariableTable struct {
	LocalVariableTableLength uint16
	LocalVariableTable       []LocalVariableEntry
}

// LocalVariableTypeTable : Attribut contenant des informations sur les types des variables locales
type LocalVariableTypeTable struct {
	LocalVariableTypeTableLength uint16
	LocalVariableTypeTable       []LocalVariableTypeEntry
}

// MethodParameters : Attribut contenant des informations sur les paramètres de la méthode
type MethodParameters struct {
	MethodParametersCount uint16
	MethodParameters      []MethodParameter
}

// RuntimeVisibleAnnotations : Attribut contenant des annotations visibles à l'exécution
type RuntimeVisibleAnnotations struct {
	NumAnnotations uint16
	Annotations    []Annotation
}

// RuntimeInvisibleAnnotations : Attribut contenant des annotations invisibles à l'exécution
type RuntimeInvisibleAnnotations struct {
	NumAnnotations uint16
	Annotations    []Annotation
}

// RuntimeVisibleParameterAnnotations : Attribut contenant des annotations visibles pour les paramètres
type RuntimeVisibleParameterAnnotations struct {
	NumParameters        uint16
	ParameterAnnotations []ParameterAnnotation
}

// RuntimeInvisibleParameterAnnotations : Attribut contenant des annotations invisibles pour les paramètres
type RuntimeInvisibleParameterAnnotations struct {
	NumParameters        uint16
	ParameterAnnotations []ParameterAnnotation
}

// SourceFile : Attribut contenant le nom du fichier source
type SourceFile struct {
	SourcefileIndex uint16
}

// SourceDebugExtension : Attribut contenant des données de débogage source
type SourceDebugExtension struct {
	DebugExtension []byte
}

// Signature : Attribut contenant la signature d'une classe, méthode ou champ
type Signature struct {
	SignatureIndex uint16
}

// StackMapTable : Attribut contenant des informations sur la pile et les variables locales
type StackMapTable struct {
	NumberOfEntries uint16
	Entries         []StackMapFrame
}

// Synthetic : Attribut indiquant que l'élément est synthétique
type Synthetic struct{}

// EnclosingMethod : Attribut indiquant la méthode englobante
type EnclosingMethod struct {
	ClassIndex  uint16
	MethodIndex uint16
}

// BootstrapMethods : Attribut contenant des méthodes de bootstrap
type BootstrapMethods struct {
	NumBootstrapMethods uint16
	BootstrapMethods    []BootstrapMethod
}

// Module : Attribut contenant des informations sur un module
type Module struct {
	NameIndex uint16
	Flags     uint16
	Requires  []ModuleRequire
}

// ModulePackages : Attribut contenant des informations sur les packages d'un module
type ModulePackages struct {
	NumberOfPackages uint16
	PackageIndex     []uint16
}

// NestHost : Attribut contenant la classe hôte
type NestHost struct {
	HostClassIndex uint16
}

// NestMembers : Attribut contenant des informations sur les classes membres
type NestMembers struct {
	NumberOfMembers uint16
	ClassIndex      []uint16
}

// PermittedSubclasses : Attribut contenant des informations sur les sous-classes autorisées
type PermittedSubclasses struct {
	NumberOfSubclasses uint16
	SubclassIndex      []uint16
}

// Définition des structures supplémentaires nécessaires pour les attributs
type ExceptionTableEntry struct {
	StartPc   uint16
	EndPc     uint16
	HandlerPc uint16
	CatchType uint16
}

type InnerClassInfo struct {
	InnerClassIndex       uint16
	OuterClassIndex       uint16
	InnerNameIndex        uint16
	InnerClassAccessFlags uint16
}

type LineNumberTableEntry struct {
	StartPc    uint16
	LineNumber uint16
}

type LocalVariableEntry struct {
	StartPc        uint16
	Length         uint16
	NameIndex      uint16
	SignatureIndex uint16
	Index          uint16
}

type LocalVariableTypeEntry struct {
	StartPc        uint16
	Length         uint16
	NameIndex      uint16
	SignatureIndex uint16
	Index          uint16
}

type MethodParameter struct {
	NameIndex   uint16
	AccessFlags uint16
}

type Annotation struct {
	TypeIndex            uint16
	NumElementValuePairs uint16
	ElementValuePairs    []ElementValuePair
}

type ParameterAnnotation struct {
	NumAnnotations uint16
	Annotations    []Annotation
}

type ElementValuePair struct {
	ElementNameIndex uint16
	Value            ElementValue
}

type ElementValue struct {
	Tag   uint8
	Value uint16
}

type StackMapFrame struct {
	FrameType   uint8
	OffsetDelta uint16
	StackItems  []StackItem
}

type StackItem struct {
	Tag   uint8
	Value uint16
}

type BootstrapMethod struct {
	MethodRefIndex uint16
	ArgumentsCount uint16
	Arguments      []uint16
}

type ModuleRequire struct {
	NameIndex    uint16
	Flags        uint16
	VersionIndex uint16
}

func parseAttributes(attributes []AttributeInfo, cp *ConstantPool) Attributes {
	attr := &Attributes{}

	for _, a := range attributes {
		reader := bytes.NewReader(a.Info)
		switch cp.Utf8[a.AttributeNameIndex] {
		case "Code": // TODO
			attr.Code = append(attr.Code, Code{
				MaxStack:       0,
				MaxLocals:      0,
				Code:           a.Info,
				ExceptionTable: []ExceptionTableEntry{},
				Attributes:     []AttributeInfo{},
			})
			break
		case "ConstantValue": // TODO
			attr.ConstantValue = append(attr.ConstantValue, ConstantValue{})
			break
		case "Deprecated": // TODO
			attr.Deprecated = append(attr.Deprecated, Deprecated{})
			break
		case "Exceptions": // TODO
			attr.Exceptions = append(attr.Exceptions, Exceptions{})
			break
		case "InnerClasses": // TODO
			attr.InnerClasses = append(attr.InnerClasses, InnerClasses{})
			break
		case "LineNumberTable": // TODO
			attr.LineNumberTable = append(attr.LineNumberTable, LineNumberTable{})
			break
		case "LocalVariableTable": // TODO
			attr.LocalVariableTable = append(attr.LocalVariableTable, LocalVariableTable{})
			break
		case "LocalVariableTypeTable": // TODO
			attr.LocalVariableTypeTable = append(attr.LocalVariableTypeTable, LocalVariableTypeTable{})
			break
		case "MethodParameters": // TODO
			attr.MethodParameters = append(attr.MethodParameters, MethodParameters{})
			break
		case "RuntimeVisibleAnnotations": // TODO
			attr.RuntimeVisibleAnnotations = append(attr.RuntimeVisibleAnnotations, RuntimeVisibleAnnotations{})
			break
		case "RuntimeInvisibleAnnotations": // TODO
			attr.RuntimeInvisibleAnnotations = append(attr.RuntimeInvisibleAnnotations, RuntimeInvisibleAnnotations{})
			break
		case "RuntimeVisibleParameterAnnotations": // TODO
			attr.RuntimeVisibleParameterAnnotations = append(attr.RuntimeVisibleParameterAnnotations, RuntimeVisibleParameterAnnotations{})
			break
		case "RuntimeInvisibleParameterAnnotations": // TODO
			attr.RuntimeInvisibleParameterAnnotations = append(attr.RuntimeInvisibleParameterAnnotations, RuntimeInvisibleParameterAnnotations{})
			break
		case "SourceFile": // TODO
			attr.SourceFile = append(attr.SourceFile, SourceFile{})
			break
		case "SourceDebugExtension": // TODO
			attr.SourceDebugExtension = append(attr.SourceDebugExtension, SourceDebugExtension{})
			break
		case "Signature": // TODO
			attr.Signature = append(attr.Signature, Signature{})
			break
		case "StackMapTable": // TODO
			attr.StackMapTable = append(attr.StackMapTable, StackMapTable{})
			break
		case "Synthetic": // TODO
			attr.Synthetic = append(attr.Synthetic, Synthetic{})
			break
		case "EnclosingMethod": // TODO
			attr.EnclosingMethod = append(attr.EnclosingMethod, EnclosingMethod{})
			break
		case "BootstrapMethods": // TODO
			attr.BootstrapMethods = append(attr.BootstrapMethods, BootstrapMethods{})
			break
		case "Module": // TODO
			attr.Module = append(attr.Module, Module{})
			break
		case "ModulePackages": // TODO
			attr.ModulePackages = append(attr.ModulePackages, ModulePackages{})
			break
		case "NestHost": // TODO
			attr.NestHost = append(attr.NestHost, NestHost{})
			break
		case "NestMembers": // TODO
			attr.NestMembers = append(attr.NestMembers, NestMembers{})
			break
		case "PermittedSubclasses": // TODO
			attr.PermittedSubclasses = append(attr.PermittedSubclasses, PermittedSubclasses{})
			break
		}
	}

	return *attr
}
