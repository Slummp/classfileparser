package classfileparser

// Nop - nop (0x00) : Do nothing
type Nop struct{}

// AconstNull - aconst_null (0x01) : Push null
type AconstNull struct{}

// IconstM1 - iconst_m1 (0x02) : Push int constant
type IconstM1 struct{}

// Iconst0 - iconst_0 (0x03) : Push int constant
type Iconst0 struct{}

// Iconst1 - iconst_1 (0x04) : Push int constant
type Iconst1 struct{}

// Iconst2 - iconst_2 (0x05) : Push int constant
type Iconst2 struct{}

// Iconst3 - iconst_3 (0x06) : Push int constant
type Iconst3 struct{}

// Iconst4 - iconst_4 (0x07) : Push int constant
type Iconst4 struct{}

// Iconst5 - iconst_5 (0x08) : Push int constant
type Iconst5 struct{}

// Lconst0 - lconst_0 (0x09) : Push long constant
type Lconst0 struct{}

// Lconst1 - lconst_1 (0x0A) : Push long constant
type Lconst1 struct{}

// Fconst0 - fconst_0 (0x0B) : Push float
type Fconst0 struct{}

// Fconst1 - fconst_1 (0x0C) : Push float
type Fconst1 struct{}

// Fconst2 - fconst_2 (0x0D) : Push float
type Fconst2 struct{}

// Dconst0 - dconst_0 (0x0E) : Push double
type Dconst0 struct{}

// Dconst1 - dconst_1 (0x0F) : Push double
type Dconst1 struct{}

// Bipush - bipush (0x10) : Push byte
type Bipush struct {
	Byte byte
}

// Sipush - sipush (0x11) : Push short
type Sipush struct {
	Short int16
}

// Ldc - ldc (0x12) : Push item from run-time constant pool
type Ldc interface{} // index in constant pool (not long/double)

// LdcW - ldc_w (0x13) : Push item from run-time constant pool (wide index)
type LdcW interface{} // index in constant pool (not long/double)

// Ldc2W - ldc2_w (0x14) : Push long or double from run-time constant pool (wide index)
type Ldc2W interface{} // index in constant pool (not int/float/string)

// Iload - iload (0x15) : Load int from local variable
type Iload struct {
	LocalIndex uint8 // local variable index
}

// Lload - lload (0x16) : Load long from local variable
type Lload struct {
	LocalIndex uint8 // local variable index
}

// Fload - fload (0x17) : Load float from local variable
type Fload struct {
	LocalIndex uint8 // local variable index
}

// Dload - dload (0x18) : Load double from local variable
type Dload struct {
	LocalIndex uint8 // local variable index
}

// Aload - aload (0x19) : Load reference from local variable
type Aload struct {
	LocalIndex uint8 // local variable index
}

// Iload0 - iload_0 (0x1A) : Load int from local variable
type Iload0 struct{}

// Iload1 - iload_1 (0x1B) : Load int from local variable
type Iload1 struct{}

// Iload2 - iload_2 (0x1C) : Load int from local variable
type Iload2 struct{}

// Iload3 - iload_3 (0x1D) : Load int from local variable
type Iload3 struct{}

// Lload0 - lload_0 (0x1E) : Load long from local variable
type Lload0 struct{}

// Lload1 - lload_1 (0x1F) : Load long from local variable
type Lload1 struct{}

// Lload2 - lload_2 (0x20) : Load long from local variable
type Lload2 struct{}

// Lload3 - lload_3 (0x21) : Load long from local variable
type Lload3 struct{}

// Fload0 - fload_0 (0x22) : Load float from local variable
type Fload0 struct{}

// Fload1 - fload_1 (0x23) : Load float from local variable
type Fload1 struct{}

// Fload2 - fload_2 (0x24) : Load float from local variable
type Fload2 struct{}

// Fload3 - fload_3 (0x25) : Load float from local variable
type Fload3 struct{}

// Dload0 - dload_0 (0x26) : Load double from local variable
type Dload0 struct{}

// Dload1 - dload_1 (0x27) : Load double from local variable
type Dload1 struct{}

// Dload2 - dload_2 (0x28) : Load double from local variable
type Dload2 struct{}

// Dload3 - dload_3 (0x29) : Load double from local variable
type Dload3 struct{}

// Aload0 - aload_0 (0x2A) : Load reference from local variable
type Aload0 struct{}

// Aload1 - aload_1 (0x2B) : Load reference from local variable
type Aload1 struct{}

// Aload2 - aload_2 (0x2C) : Load reference from local variable
type Aload2 struct{}

// Aload3 - aload_3 (0x2D) : Load reference from local variable
type Aload3 struct{}

// Iaload - iaload (0x2E) : Load int from array
type Iaload struct{}

// Laload - laload (0x2F) : Load long from array
type Laload struct{}

// Faload - faload (0x30) : Load float from array
type Faload struct{}

// Daload - daload (0x31) : Load double from array
type Daload struct{}

// Aaload - aaload (0x32) : Load reference from array
type Aaload struct{}

// Baload - baload (0x33) : Load byte or boolean from array
type Baload struct{}

// Caload - caload (0x34) : Load char from array
type Caload struct{}

// Saload - saload (0x35) : Load short from array
type Saload struct{}

// Istore - istore (0x36) : Store int into local variable
type Istore struct {
	LocalIndex uint8 // local variable index
}

// Lstore - lstore (0x37) : Store long into local variable
type Lstore struct {
	LocalIndex uint8 // local variable index
}

// Fstore - fstore (0x38) : Store float into local variable
type Fstore struct {
	LocalIndex uint8 // local variable index
}

// Dstore - dstore (0x39) : Store double into local variable
type Dstore struct {
	LocalIndex uint8 // local variable index
}

// Astore - astore (0x3A) : Store reference into local variable
type Astore struct {
	LocalIndex uint8 // local variable index
}

// Istore0 - istore_0 (0x3B) : Store int into local variable
type Istore0 struct{}

// Istore1 - istore_1 (0x3C) : Store int into local variable
type Istore1 struct{}

// Istore2 - istore_2 (0x3D) : Store int into local variable
type Istore2 struct{}

// Istore3 - istore_3 (0x3E) : Store int into local variable
type Istore3 struct{}

// Lstore0 - lstore_0 (0x3F) : Store long into local variable
type Lstore0 struct{}

// Lstore1 - lstore_1 (0x40) : Store long into local variable
type Lstore1 struct{}

// Lstore2 - lstore_2 (0x41) : Store long into local variable
type Lstore2 struct{}

// Lstore3 - lstore_3 (0x42) : Store long into local variable
type Lstore3 struct{}

// Fstore0 - fstore_0 (0x43) : Store float into local variable
type Fstore0 struct{}

// Fstore1 - fstore_1 (0x44) : Store float into local variable
type Fstore1 struct{}

// Fstore2 - fstore_2 (0x45) : Store float into local variable
type Fstore2 struct{}

// Fstore3 - fstore_3 (0x46) : Store float into local variable
type Fstore3 struct{}

// Dstore0 - dstore_0 (0x47) : Store double into local variable
type Dstore0 struct{}

// Dstore1 - dstore_1 (0x48) : Store double into local variable
type Dstore1 struct{}

// Dstore2 - dstore_2 (0x49) : Store double into local variable
type Dstore2 struct{}

// Dstore3 - dstore_3 (0x4A) : Store double into local variable
type Dstore3 struct{}

// Astore0 - astore_0 (0x4B) : Store reference into local variable
type Astore0 struct{}

// Astore1 - astore_1 (0x4C) : Store reference into local variable
type Astore1 struct{}

// Astore2 - astore_2 (0x4D) : Store reference into local variable
type Astore2 struct{}

// Astore3 - astore_3 (0x4E) : Store reference into local variable
type Astore3 struct{}

// Iastore - iastore (0x4F) : Store into int array
type Iastore struct{}

// Lastore - lastore (0x50) : Store into long array
type Lastore struct{}

// Fastore - fastore (0x51) : Store into float array
type Fastore struct{}

// Dastore - dastore (0x52) : Store into double array
type Dastore struct{}

// Aastore - aastore (0x53) : Store into reference array
type Aastore struct{}

// Bastore - bastore (0x54) : Store into byte or boolean array
type Bastore struct{}

// Castore - castore (0x55) : Store into char array
type Castore struct{}

// Sastore - sastore (0x56) : Store into short array
type Sastore struct{}

// Pop - pop (0x57) : Pop the top operand stack value
type Pop struct{}

// Pop2 - pop2 (0x58) : Pop the top one or two operand stack values
type Pop2 struct{}

// Dup - dup (0x59) : Duplicate the top operand stack value
type Dup struct{}

// DupX1 - dup_x1 (0x5A) : Duplicate the top operand stack value and insert two values down
type DupX1 struct{}

// DupX2 - dup_x2 (0x5B) : Duplicate the top operand stack value and insert two or three values down
type DupX2 struct{}

// Dup2 - dup2 (0x5C) : Duplicate the top one or two operand stack values
type Dup2 struct{}

// Dup2X1 - dup2_x1 (0x5D) : Duplicate the top one or two operand stack values and insert two or three values down
type Dup2X1 struct{}

// Dup2X2 - dup2_x2 (0x5E) : Duplicate the top one or two operand stack values and insert two three or four values down
type Dup2X2 struct{}

// Swap - swap (0x5F) : Swap the top two operand stack values
type Swap struct{}

// Iadd - iadd (0x60) : Add int
type Iadd struct{}

// Ladd - ladd (0x61) : Add long
type Ladd struct{}

// Fadd - fadd (0x62) : Add float
type Fadd struct{}

// Dadd - dadd (0x63) : Add double
type Dadd struct{}

// Isub - isub (0x64) : Subtract int
type Isub struct{}

// Lsub - lsub (0x65) : Subtract long
type Lsub struct{}

// Fsub - fsub (0x66) : Subtract float
type Fsub struct{}

// Dsub - dsub (0x67) : Subtract double
type Dsub struct{}

// Imul - imul (0x68) : Multiply int
type Imul struct{}

// Lmul - lmul (0x69) : Multiply long
type Lmul struct{}

// Fmul - fmul (0x6A) : Multiply float
type Fmul struct{}

// Dmul - dmul (0x6B) : Multiply double
type Dmul struct{}

// Idiv - idiv (0x6C) : Divide int
type Idiv struct{}

// Ldiv - ldiv (0x6D) : Divide long
type Ldiv struct{}

// Fdiv - fdiv (0x6E) : Divide float
type Fdiv struct{}

// Ddiv - ddiv (0x6F) : Divide double
type Ddiv struct{}

// Irem - irem (0x70) : Remainder int
type Irem struct{}

// Lrem - lrem (0x71) : Remainder long
type Lrem struct{}

// Frem - frem (0x72) : Remainder float
type Frem struct{}

// Drem - drem (0x73) : Remainder double
type Drem struct{}

// Ineg - ineg (0x74) : Negate int
type Ineg struct{}

// Lneg - lneg (0x75) : Negate long
type Lneg struct{}

// Fneg - fneg (0x76) : Negate float
type Fneg struct{}

// Dneg - dneg (0x77) : Negate double
type Dneg struct{}

// Ishl - ishl (0x78) : Shift left int
type Ishl struct{}

// Lshl - lshl (0x79) : Shift left long
type Lshl struct{}

// Ishr - ishr (0x7A) : Arithmetic shift right int
type Ishr struct{}

// Lshr - lshr (0x7B) : Arithmetic shift right long
type Lshr struct{}

// Iushr - iushr (0x7C) : Logical shift right int
type Iushr struct{}

// Lushr - lushr (0x7D) : Logical shift right long
type Lushr struct{}

// Iand - iand (0x7E) : Boolean AND int
type Iand struct{}

// Land - land (0x7F) : Boolean AND long
type Land struct{}

// Ior - ior (0x80) : Boolean OR int
type Ior struct{}

// Lor - lor (0x81) : Boolean OR long
type Lor struct{}

// Ixor - ixor (0x82) : Boolean XOR int
type Ixor struct{}

// Lxor - lxor (0x83) : Boolean XOR long
type Lxor struct{}

// Iinc - iinc (0x84) : Increment local variable by constant
type Iinc struct {
	LocalIndex uint8 // local variable index
	Const      int8  // const value
}

// I2l - i2l (0x85) : Convert int to long
type I2l struct{}

// I2f - i2f (0x86) : Convert int to float
type I2f struct{}

// I2d - i2d (0x87) : Convert int to double
type I2d struct{}

// L2i - l2i (0x88) : Convert long to int
type L2i struct{}

// L2f - l2f (0x89) : Convert long to float
type L2f struct{}

// L2d - l2d (0x8A) : Convert long to double
type L2d struct{}

// F2i - f2i (0x8B) : Convert float to int
type F2i struct{}

// F2l - f2l (0x8C) : Convert float to long
type F2l struct{}

// F2d - f2d (0x8D) : Convert float to double
type F2d struct{}

// D2i - d2i (0x8E) : Convert double to int
type D2i struct{}

// D2l - d2l (0x8F) : Convert double to long
type D2l struct{}

// D2f - d2f (0x90) : Convert double to float
type D2f struct{}

// I2b - i2b (0x91) : Convert int to byte
type I2b struct{}

// I2c - i2c (0x92) : Convert int to char
type I2c struct{}

// I2s - i2s (0x93) : Convert int to short
type I2s struct{}

// Lcmp - lcmp (0x94) : Compare long
type Lcmp struct{}

// Fcmpl - fcmpl (0x95) : Compare float
type Fcmpl struct{}

// Fcmpg - fcmpg (0x96) : Compare float
type Fcmpg struct{}

// Dcmpl - dcmpl (0x97) : Compare double
type Dcmpl struct{}

// Dcmpg - dcmpg (0x98) : Compare double
type Dcmpg struct{}

// Ifeq - ifeq (0x99) : Branch if int comparison with zero succeeds
type Ifeq struct {
	Offset int16 // jump offset (branch)
}

// Ifne - ifne (0x9A) : Branch if int comparison with zero succeeds
type Ifne struct {
	Offset int16 // jump offset (branch)
}

// Iflt - iflt (0x9B) : Branch if int comparison with zero succeeds
type Iflt struct {
	Offset int16 // jump offset (branch)
}

// Ifge - ifge (0x9C) : Branch if int comparison with zero succeeds
type Ifge struct {
	Offset int16 // jump offset (branch)
}

// Ifgt - ifgt (0x9D) : Branch if int comparison with zero succeeds
type Ifgt struct {
	Offset int16 // jump offset (branch)
}

// Ifle - ifle (0x9E) : Branch if int comparison with zero succeeds
type Ifle struct {
	Offset int16 // jump offset (branch)
}

// IfIcmpeq - if_icmpeq (0x9F) : Branch if int comparison succeeds
type IfIcmpeq struct {
	Offset int16 // jump offset (branch)
}

// IfIcmpne - if_icmpne (0xA0) : Branch if int comparison succeeds
type IfIcmpne struct {
	Offset int16 // jump offset (branch)
}

// IfIcmplt - if_icmplt (0xA1) : Branch if int comparison succeeds
type IfIcmplt struct {
	Offset int16 // jump offset (branch)
}

// IfIcmpge - if_icmpge (0xA2) : Branch if int comparison succeeds
type IfIcmpge struct {
	Offset int16 // jump offset (branch)
}

// IfIcmpgt - if_icmpgt (0xA3) : Branch if int comparison succeeds
type IfIcmpgt struct {
	Offset int16 // jump offset (branch)
}

// IfIcmple - if_icmple (0xA4) : Branch if int comparison succeeds
type IfIcmple struct {
	Offset int16 // jump offset (branch)
}

// IfAcmpeq - if_acmpeq (0xA5) : Branch if reference comparison succeeds
type IfAcmpeq struct {
	Offset int16 // jump offset (branch)
}

// IfAcmpne - if_acmpne (0xA6) : Branch if reference comparison succeeds
type IfAcmpne struct {
	Offset int16 // jump offset (branch)
}

// Goto - goto (0xA7) : Branch always
type Goto struct {
	Offset int16 // jump offset (branch)
}

// Jsr - jsr (0xA8) : Jump subroutine
type Jsr struct {
	Offset int16 // jump offset (branch)
}

// Ret - ret (0xA9) : Return from subroutine
type Ret struct {
	LocalIndex uint8 // local variable index
}

// // Tableswitch - tableswitch (0xAA) : Access jump table by index and jump // TODO
// type Tableswitch struct {
//   [0] // 0 à 3 octets vide pour le padding
//   int32 // default offset
//   int32 // low (<= high)
//   int32 // high (>= low)
//   [int32] // high - low + 1 offsets
// }

// // Lookupswitch - lookupswitch (0xAB) : Access jump table by key match and jump // TODO
// type Lookupswitch struct {
//   [0] // 0 à 3 octets vide pour le padding
//   int32 // default offset
//   int32 // pairs length >= 0
//   [int32 int32] // pairs (match - offset)
// }

// Ireturn - ireturn (0xAC) : Return int from method
type Ireturn struct{}

// Lreturn - lreturn (0xAD) : Return long from method
type Lreturn struct{}

// Freturn - freturn (0xAE) : Return float from method
type Freturn struct{}

// Dreturn - dreturn (0xAF) : Return double from method
type Dreturn struct{}

// Areturn - areturn (0xB0) : Return reference from method
type Areturn struct{}

// Return - return (0xB1) : Return void from method
type Return struct{}

// Getstatic - getstatic (0xB2) : Get static field from class
type Getstatic struct {
	Class string
	Name  string
	Type  string
}

// Putstatic - putstatic (0xB3) : Set static field in class
type Putstatic struct {
	Class string
	Name  string
	Type  string
}

// Getfield - getfield (0xB4) : Fetch field from object
type Getfield struct {
	Class string
	Name  string
	Type  string
}

// Putfield - putfield (0xB5) : Set field in object
type Putfield struct {
	Class string
	Name  string
	Type  string
}

// Invokevirtual - invokevirtual (0xB6) : Invoke instance method; dispatch based on class
type Invokevirtual struct {
	Class string
	Name  string
	Type  string
}

// Invokespecial - invokespecial (0xB7) : Invoke instance method;  direct invocation of instance initialization methods and methods of the current class and its supertypes
type Invokespecial struct {
	Class string
	Name  string
	Type  string
}

// Invokestatic - invokestatic (0xB8) : Invoke a class (static) method
type Invokestatic struct {
	Class string
	Name  string
	Type  string
}

// Invokeinterface - invokeinterface (0xB9) : Invoke interface method // TODO
type Invokeinterface struct {
	InterfaceMethodref InterfaceMethodref
	Count              uint8
}

// Invokedynamic - invokedynamic (0xBA) : Invoke a dynamically-computed call site // TODO
type Invokedynamic struct {
	InvokeDynamic InvokeDynamic
}

// New - new (0xBB) : Create new object
type New string

// Newarray - newarray (0xBC) : Create new array // TODO
type Newarray struct {
	Type uint8 // type -> T_BOOLEAN=4 ; T_CHAR=5 ; T_FLOAT=6 ; T_DOUBLE=7 ; T_BYTE=8 ; T_SHORT=9 ; T_INT=10 ; T_LONG=11
}

// Anewarray - anewarray (0xBD) : Create new array of reference
type Anewarray string

// Arraylength - arraylength (0xBE) : Get length of array
type Arraylength struct{}

// Athrow - athrow (0xBF) : Throw exception or error
type Athrow struct{}

// Checkcast - checkcast (0xC0) : Check whether object is of given type
type Checkcast string

// Instanceof - instanceof (0xC1) : Determine if object is of given type
type Instanceof string

// Monitorenter - monitorenter (0xC2) : Enter monitor for object
type Monitorenter struct{}

// Monitorexit - monitorexit (0xC3) : Exit monitor for object
type Monitorexit struct{}

// Wide - wide (0xC4) : Extend local variable index by additional bytes // TODO
type Wide struct {
	OpCode     uint8  // opcode in iload fload aload lload dload istore fstore astore lstore dstore ret iinc
	LocalIndex uint16 // local variable index
	Const      uint16 // const value ONLY in iinc case
}

// Multianewarray - multianewarray (0xC5) : Create new multidimensional array
type Multianewarray struct {
	Class     string // index in constant pool (Class)
	Dimension uint8  // dimension >= 1
}

// Ifnull - ifnull (0xC6) : Branch if reference is null
type Ifnull struct {
	Offset int16 // jump offset (branch)
}

// Ifnonnull - ifnonnull (0xC7) : Branch if reference not null
type Ifnonnull struct {
	Offset int16 // jump offset (branch)
}

// GotoW - goto_w (0xC8) : Branch always (wide index)
type GotoW struct {
	Offset int32 // jump offset (branch)
}

// JsrW - jsr_w (0xC9) : Jump subroutine (wide index)
type JsrW struct {
	Offset int32 // jump offset (branch)
}
