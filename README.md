# classfileparser

`classfileparser` is a pure Go library for reading and decoding Java `.class` files. It exposes the raw structures defined by the JVM specification, adds typed helpers on top of the constant pool and attributes, and gives you a convenient high-level snapshot of a compiled class.

## Features

- Parse `.class` files directly from any `io.Reader` with `Open`
- Navigate the raw `ClassFile` data structure or work with the friendly `ClassStruct` snapshot
- Resolve constant pool entries into typed Go values (`Utf8`, `Class`, `Methodref`, and more)
- Decode most standard JVM attributes, including `Code`, `LineNumberTable`, module metadata, and annotations
- Represent bytecode instructions with dedicated Go types so you can pattern-match opcodes safely

## Installation

```bash
go get github.com/Slummp/classfileparser
```

The module targets Go 1.25 or newer (as declared in `go.mod`).

## Quick start

```go
package main

import (
    "fmt"
    "log"
    "os"

    classfileparser "github.com/Slummp/classfileparser"
)

func main() {
    fh, err := os.Open("Hello.class")
    if err != nil {
        log.Fatal(err)
    }
    defer fh.Close()

    cf, err := classfileparser.Open(fh)
    if err != nil {
        log.Fatal(err)
    }

    cp, err := cf.GetConstantPool()
    if err != nil {
        log.Fatal(err)
    }

    snapshot, err := cf.GetClassFile()
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Class %s extends %s\n", snapshot.ThisClass, snapshot.SuperClass)
    for _, method := range snapshot.Methods {
        fmt.Printf("%v %s(%v) -> %s\n", method.Access, method.Name, method.ParamsTypes, method.ReturnType)
    }

    // Constant pool entries are typed – no need for manual tag switches.
    if utf8, ok := cp[1].(classfileparser.Utf8); ok {
        fmt.Println("#1 is a UTF-8 entry:", utf8)
    }
}
```

## API overview

### Opening a class file

`Open(io.Reader)` returns a populated `*ClassFile`. The struct mirrors the JVM specification: magic number, version, constant pool, access flags, interfaces, fields, methods, and attributes.

### Working with the constant pool

`(*ClassFile).GetConstantPool()` converts raw pool entries into idiomatic Go types for easier use. Some highlights:

- `Utf8` and primitive literals resolve to `string`, `int32`, `int64`, `float32`, or `float64`
- `Class`, `Module`, and `Package` resolve to their internal names as `string`
- Member references (`Fieldref`, `Methodref`, `InterfaceMethodref`) expand into structs with `Class`, `Name`, and `Type`
- Invoke-dynamic and method handles become `InvokeDynamic`, `Dynamic`, and `MethodHandle` structs with decoded metadata

The map is indexed by the original JVM slot number, so `cp[7]` corresponds to entry `#7` in the class file.

### High-level snapshot

If you prefer a condensed view, call `(*ClassFile).GetClassFile()`. The returned `ClassStruct` contains:

- Parsed access flags as slices of strings (`[]string`)
- Resolved class and interface names
- Field and method descriptors split into name, return type, and parameter descriptors
- Attributes already decoded via `parseAttributes`

This snapshot is perfect for rendering class summaries, generating documentation, or feeding higher-level tooling.

## Attribute decoding

`parseAttributes` recognises a broad range of standard JVM attributes. The library currently decodes:

- `Code`
- `ConstantValue`
- `Deprecated`
- `Exceptions`
- `InnerClasses`
- `LineNumberTable`
- `LocalVariableTable`
- `LocalVariableTypeTable`
- `MethodParameters`
- `RuntimeVisibleAnnotations`
- `RuntimeInvisibleAnnotations`
- `RuntimeVisibleParameterAnnotations`
- `RuntimeInvisibleParameterAnnotations`
- `SourceFile`
- `SourceDebugExtension`
- `Signature`
- `StackMapTable`
- `Synthetic`
- `EnclosingMethod`
- `BootstrapMethods`
- `Module`
- `ModulePackages`
- `NestHost`
- `NestMembers`
- `PermittedSubclasses`.

New attribute types can be added by extending the switch in `parseAttributes`.

## Bytecode representation

Inside the `Code` attribute, opcodes are converted into dedicated Go structs. For example:

```go
for _, attr := range snapshot.Methods[0].Attributes {
    code, ok := attr.(classfileparser.Code)
    if !ok {
        continue
    }
    for _, instruction := range code.Code {
        switch instr := instruction.(type) {
        case classfileparser.Iconst1:
            fmt.Println("push constant 1")
        case classfileparser.Invokevirtual:
            fmt.Println("invoke virtual", instr.Class, instr.Name)
        // ...
        }
    }
}
```

This design keeps decoding logic out of your application and lets you focus on the semantics you care about.

## Error handling and panics

- `Open` and `GetConstantPool` return descriptive errors for malformed files or unsupported tags.
- Attribute parsing uses helper readers that panic when faced with unexpected lengths or unsupported structures (e.g. stack-map frames above type 63, or unimplemented switch instructions). Wrap calls with `recover` if you need to sandbox untrusted input.

## Testing

```bash
go test ./...
```

The comprehensive test suite builds synthetic class files covering constant pool cases, attributes, and bytecode instructions.

## Roadmap and known limitations

- `tableswitch` and `lookupswitch` opcodes (`0xAA`, `0xAB`) are not yet decoded.
- Method invocation and dynamic call opcodes `invokeinterface` (`0xB9`) and `invokedynamic` (`0xBA`) are placeholders.
- Array allocation and wide index opcodes `newarray` (`0xBC`) and `wide` (`0xC4`) are not fully implemented yet.
- Some advanced StackMapTable frame types and BootstrapMethod details are placeholders.
- Module attribute parsing skips several sub-sections after counting entries; fill in the TODOs if you need full fidelity.
- Attribute parsing panics on unsupported tags instead of returning errors – consider hardening if you expect hostile input.

## License

Released under the [MIT License](LICENSE).
