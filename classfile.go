package classfileparser

import (
	"fmt"
	"regexp"
)

// ClassStruct represents the entire structure of a .class file
type ClassStruct struct {
	Version struct {
		MinorVersion uint16 // Minor version
		MajorVersion uint16 // Major version
	}
	Access     []string // Access flags (e.g., public, final)
	ThisClass  string   // Current class
	SuperClass string   // Superclass (empty for java/lang/Object and module-info)
	Interfaces []string // Interfaces

	Fields     []Field     // Field structures
	Methods    []Method    // Method structures
	Attributes []Attribute // Attribute structures
}

// Field represents a field in the class
type Field struct {
	Access     []string    // Access flags for the field
	Name       string      // Name of the field
	Type       string      // Type of the field
	Attributes []Attribute // Field attributes
}

// Method represents a method in the class
type Method struct {
	Access      []string    // Access flags for the method
	Name        string      // Name of the method
	ReturnType  string      // Return type of the method
	ParamsTypes []string    // Type of the params of the method
	Attributes  []Attribute // Method attributes
}

// GetClassFile converts the parsed binary data into a structured ClassStruct snapshot
func (cf *ClassFile) GetClassFile() (*ClassStruct, error) {
	cp, err := cf.GetConstantPool()
	if err != nil {
		return nil, err
	}

	interfaces := []string{}
	for _, i := range cf.Interfaces {
		interfaces = append(interfaces, string(cp[i].(Class)))
	}

	fields := []Field{}
	for _, f := range cf.Fields {
		name := string(cp[f.NameIndex].(Utf8))
		attributes, err := parseAttributes(f.Attributes, cp)
		if err != nil {
			return nil, fmt.Errorf("failed to read field %s: %w", name, err)
		}
		fields = append(fields, Field{
			Access:     findFlags(FieldT, f.AccessFlags),
			Name:       name,
			Type:       string(cp[f.DescriptorIndex].(Utf8)),
			Attributes: attributes,
		})
	}

	methods := []Method{}
	for _, m := range cf.Methods {
		name := string(cp[m.NameIndex].(Utf8))
		attributes, err := parseAttributes(m.Attributes, cp)
		if err != nil {
			return nil, fmt.Errorf("failed to read method %s: %w", name, err)
		}
		paramsTypes, returnType := readSignature(string(cp[m.DescriptorIndex].(Utf8)))
		methods = append(methods, Method{
			Access:      findFlags(MethodT, m.AccessFlags),
			Name:        name,
			ReturnType:  returnType,
			ParamsTypes: paramsTypes,
			Attributes:  attributes,
		})
	}

	attributes, err := parseAttributes(cf.Attributes, cp)
	if err != nil {
		return nil, err
	}

	var superClass string
	if cf.SuperClass != 0 {
		superClass = string(cp[cf.SuperClass].(Class))
	}

	return &ClassStruct{
		Version: struct {
			MinorVersion uint16
			MajorVersion uint16
		}{
			MinorVersion: cf.MinorVersion,
			MajorVersion: cf.MajorVersion,
		},
		Access:     findFlags(ClassT, cf.AccessFlags),
		ThisClass:  string(cp[cf.ThisClass].(Class)),
		SuperClass: superClass,
		Interfaces: interfaces,
		Fields:     fields,
		Methods:    methods,
		Attributes: attributes,
	}, nil
}

func readSignature(signature string) ([]string, string) {
	re := regexp.MustCompile(`\((.*?)\)(.*)`)
	matches := re.FindStringSubmatch(signature)
	paramString, returnType := matches[1], matches[2]

	re = regexp.MustCompile(`(\[?(?:L.+?;|.))`)
	params := re.FindAllString(paramString, -1)

	return params, returnType
}
