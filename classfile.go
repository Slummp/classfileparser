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
	for i, index := range cf.Interfaces {
		name, err := getConstant[Class](cp, index)
		if err != nil {
			return nil, fmt.Errorf("failed to read interface %d: %w", i, err)
		}
		interfaces = append(interfaces, string(name))
	}

	fields := []Field{}
	for i, f := range cf.Fields {
		name, err := getConstant[Utf8](cp, f.NameIndex)
		if err != nil {
			return nil, fmt.Errorf("failed to read field %d: %w", i, err)
		}
		descriptor, err := getConstant[Utf8](cp, f.DescriptorIndex)
		if err != nil {
			return nil, fmt.Errorf("failed to read field %s: %w", name, err)
		}
		attributes, err := parseAttributes(f.Attributes, cp)
		if err != nil {
			return nil, fmt.Errorf("failed to read field %s: %w", name, err)
		}
		fields = append(fields, Field{
			Access:     findFlags(FieldT, f.AccessFlags),
			Name:       string(name),
			Type:       string(descriptor),
			Attributes: attributes,
		})
	}

	methods := []Method{}
	for i, m := range cf.Methods {
		name, err := getConstant[Utf8](cp, m.NameIndex)
		if err != nil {
			return nil, fmt.Errorf("failed to read method %d: %w", i, err)
		}
		descriptor, err := getConstant[Utf8](cp, m.DescriptorIndex)
		if err != nil {
			return nil, fmt.Errorf("failed to read method %s: %w", name, err)
		}
		attributes, err := parseAttributes(m.Attributes, cp)
		if err != nil {
			return nil, fmt.Errorf("failed to read method %s: %w", name, err)
		}
		paramsTypes, returnType, err := readSignature(string(descriptor))
		if err != nil {
			return nil, fmt.Errorf("failed to read method %s: %w", name, err)
		}
		methods = append(methods, Method{
			Access:      findFlags(MethodT, m.AccessFlags),
			Name:        string(name),
			ReturnType:  returnType,
			ParamsTypes: paramsTypes,
			Attributes:  attributes,
		})
	}

	attributes, err := parseAttributes(cf.Attributes, cp)
	if err != nil {
		return nil, err
	}

	thisClass, err := getConstant[Class](cp, cf.ThisClass)
	if err != nil {
		return nil, fmt.Errorf("failed to read this class: %w", err)
	}

	var superClass Class
	if cf.SuperClass != 0 {
		if superClass, err = getConstant[Class](cp, cf.SuperClass); err != nil {
			return nil, fmt.Errorf("failed to read super class: %w", err)
		}
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
		ThisClass:  string(thisClass),
		SuperClass: string(superClass),
		Interfaces: interfaces,
		Fields:     fields,
		Methods:    methods,
		Attributes: attributes,
	}, nil
}

func readSignature(signature string) ([]string, string, error) {
	re := regexp.MustCompile(`\((.*?)\)(.*)`)
	matches := re.FindStringSubmatch(signature)
	if matches == nil {
		return nil, "", fmt.Errorf("%w: invalid method descriptor: %s", ErrMalformed, signature)
	}
	paramString, returnType := matches[1], matches[2]

	re = regexp.MustCompile(`(\[?(?:L.+?;|.))`)
	params := re.FindAllString(paramString, -1)

	return params, returnType, nil
}
