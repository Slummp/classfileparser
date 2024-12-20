package classfileparser

import (
	"regexp"
	"strings"
)

// Class represents the entire structure of a .class file
type Class struct {
	Version struct {
		MinorVersion uint16 // Minor version
		MajorVersion uint16 // Major version
	}
	Access     []string // Access flags (e.g., public, final)
	ThisClass  string   // Current class
	SuperClass string   // Superclass
	Interfaces []string // Interfaces

	Fields     []Field    // Field structures
	Methods    []Method   // Method structures
	Attributes Attributes // Attribute structures
}

// Field represents a field in the class
type Field struct {
	Access     []string   // Access flags for the field
	Name       string     // Name of the field
	Type       string     // Type of the field
	Attributes Attributes // Field attributes
}

// Method represents a method in the class
type Method struct {
	Access      []string   // Access flags for the method
	Name        string     // Name of the method
	ReturnType  string     // Return type of the method
	ParamsTypes []string   // Type of the params of the method
	Attributes  Attributes // Method attributes
}

func (cf *ClassFile) GetClassFile() (*Class, error) {
	cp, err := cf.GetConstantPool()
	if err != nil {
		panic(err)
	}

	interfaces := []string{}
	for _, i := range cf.Interfaces {
		interfaces = append(interfaces, cp.Class[i])
	}

	fields := []Field{}
	for _, f := range cf.Fields {
		fields = append(fields, Field{
			Access:     findFlags(FieldT, f.AccessFlags),
			Name:       cp.Utf8[f.NameIndex],
			Type:       cp.Utf8[f.DescriptorIndex],
			Attributes: parseAttributes(f.Attributes, cp),
		})
	}

	methods := []Method{}
	for _, m := range cf.Methods {
		paramsTypes, returnType := readSignature(cp.Utf8[m.DescriptorIndex])
		methods = append(methods, Method{
			Access:      findFlags(MethodT, m.AccessFlags),
			Name:        cp.Utf8[m.NameIndex],
			ReturnType:  returnType,
			ParamsTypes: paramsTypes,
			Attributes:  parseAttributes(m.Attributes, cp),
		})
	}

	return &Class{
		Version: struct {
			MinorVersion uint16
			MajorVersion uint16
		}{
			MinorVersion: cf.MinorVersion,
			MajorVersion: cf.MajorVersion,
		},
		Access:     findFlags(ClassT, cf.AccessFlags),
		ThisClass:  cp.Class[cf.ThisClass],
		SuperClass: cp.Class[cf.SuperClass],
		Interfaces: interfaces,
		Fields:     fields,
		Methods:    methods,
		Attributes: parseAttributes(cf.Attributes, cp),
	}, nil
}

func readSignature(signature string) ([]string, string) {
	re := regexp.MustCompile(`\((.*?)\)(.*)`)
	matches := re.FindStringSubmatch(signature)
	paramString, returnType := matches[1], matches[2]
	params := extractParamTypes(paramString)

	return params, returnType
}

func extractParamTypes(paramString string) []string {
	var params []string
	for len(paramString) > 0 {
		if paramString[0] == 'L' {
			end := strings.Index(paramString, ";") + 1
			params = append(params, paramString[:end])
			paramString = paramString[end:]
		} else {
			params = append(params, string(paramString[0]))
			paramString = paramString[1:]
		}
	}
	return params
}
