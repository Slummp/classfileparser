# classfileparser

[![CI](https://github.com/Slummp/classfileparser/actions/workflows/ci.yml/badge.svg)](https://github.com/Slummp/classfileparser/actions/workflows/ci.yml)
[![Coverage](https://img.shields.io/endpoint?url=https://raw.githubusercontent.com/Slummp/classfileparser/badges/coverage.json)](https://github.com/Slummp/classfileparser/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/Slummp/classfileparser.svg)](https://pkg.go.dev/github.com/Slummp/classfileparser)
[![Go version](https://img.shields.io/github/go-mod/go-version/Slummp/classfileparser)](go.mod)
[![Release](https://img.shields.io/github/v/tag/Slummp/classfileparser?label=release&sort=semver)](https://github.com/Slummp/classfileparser/tags)
[![License](https://img.shields.io/github/license/Slummp/classfileparser)](LICENSE)

`classfileparser` is a pure Go library for reading and decoding Java `.class` files. It exposes the raw structures defined by the JVM specification, adds typed helpers on top of the constant pool and attributes, and gives you a convenient high-level snapshot of a compiled class.

## Features

- Parse `.class` files directly from any `io.Reader` with `Open`
- Navigate the raw `ClassFile` data structure or work with the friendly `ClassStruct` snapshot
- Resolve constant pool entries into typed Go values (`Utf8`, `Class`, `Methodref`, and more)
- Recognize every standard JVM attribute, and fully decode `Code` (instructions, exception table) and `Signature`
- Represent bytecode instructions with dedicated Go types so you can pattern-match opcodes safely

## Installation

```bash
go get github.com/Slummp/classfileparser
```

The module targets Go 1.22 or newer (as declared in `go.mod`) and is tested against Go 1.22 and the latest stable release.

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
- `MethodHandle` entries are fully resolved (kind, target member and whether it is an interface method); `Dynamic` and `InvokeDynamic` entries are recognized but not decoded yet

The map is indexed by the original JVM slot number, so `cp[7]` corresponds to entry `#7` in the class file.

### High-level snapshot

If you prefer a condensed view, call `(*ClassFile).GetClassFile()`. The returned `ClassStruct` contains:

- Parsed access flags as slices of strings (`[]string`)
- Resolved class and interface names
- Field and method descriptors split into name, return type, and parameter descriptors
- Attributes already decoded via `parseAttributes`

This snapshot is perfect for rendering class summaries, generating documentation, or feeding higher-level tooling.

## Attribute decoding

`GetClassFile` returns the attributes of the class, of its fields, of its methods and of their `Code` attributes. Every standard attribute is recognized and returned as its own Go type, but only some of them are decoded yet:

| Status | Attributes |
|---|---|
| Fully decoded | `Code` (instructions, exception table, nested attributes), `Signature`, `Deprecated`, `Synthetic` |
| Recognized, content not decoded yet | `ConstantValue`, `Exceptions`, `InnerClasses`, `EnclosingMethod`, `LineNumberTable`, `LocalVariableTable`, `LocalVariableTypeTable`, `MethodParameters`, `RuntimeVisibleAnnotations`, `RuntimeInvisibleAnnotations`, `RuntimeVisibleParameterAnnotations`, `RuntimeInvisibleParameterAnnotations`, `SourceFile`, `SourceDebugExtension`, `StackMapTable`, `BootstrapMethods`, `Module`, `ModulePackages`, `NestHost`, `NestMembers`, `PermittedSubclasses`, `Record` |

Recognized attributes are returned as empty structs for now. Unknown attributes are ignored.

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

## Error handling

`Open`, `GetConstantPool` and `GetClassFile` never panic: every problem is returned as an error, wrapped with its context, for example:

```text
failed to read method run: failed to read attribute Code: failed to read instruction at offset 2: malformed class file: unknown opcode: 0xCB
```

Errors fall into two categories that you can check with `errors.Is`:

- **Truncated input**: the reader ends before the class file does. The error wraps `io.EOF` or `io.ErrUnexpectedEOF`.
- **Malformed content**: an invalid magic number, an unknown constant pool tag or opcode, an invalid constant pool reference, data overflowing its attribute… The error wraps `classfileparser.ErrMalformed`.

```go
cf, err := classfileparser.Open(r)
if errors.Is(err, classfileparser.ErrMalformed) {
    // not a valid class file
}
```

Lengths read from the file are checked against the available data before allocating, so a small malicious file cannot make the parser allocate gigabytes of memory.

## Testing

The test suite runs on every push against Go 1.22 and the latest stable release, and the CI enforces 100% statement coverage.

- **Unit tests** build class files byte by byte to cover every opcode, every constant pool entry and every error path.
- **Fixtures** are real classes compiled from `testdata/java/src` with `javac`, checked against `javap`. Compile them before running the tests (requires a JDK 17 or newer):

  ```bash
  testdata/java/build.sh
  go test ./...
  ```

- **Fuzzing** checks that arbitrary input never panics and always returns a classified error. It runs for 30 seconds on every push; run it longer locally with:

  ```bash
  go test -run '^$' -fuzz=FuzzOpen -fuzztime=10m .
  ```

- **JDK corpus**: every release is tested against all the classes of the JDK 21 (more than 28,000 classes). To run it locally:

  ```bash
  jimage extract --dir /tmp/jdk "$JAVA_HOME/lib/modules"
  CLASSFILEPARSER_CORPUS=/tmp/jdk go test -run TestCorpus -v .
  ```

## Roadmap and known limitations

- Decode the content of the attributes that are only recognized for now (see [Attribute decoding](#attribute-decoding)).
- Decode `CONSTANT_Dynamic` and `CONSTANT_InvokeDynamic` entries, which are returned empty for now.

## License

Released under the [MIT License](LICENSE).
