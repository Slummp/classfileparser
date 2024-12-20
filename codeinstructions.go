package classfileparser

// nop (0x00) : Do nothing
type Nop struct{}

// aconst_null (0x01) : Push null
type AconstNull struct{}

// iconst_m1 (0x02) : Push int constant
type IconstM1 struct{}

// iconst_0 (0x03) : Push int constant
type Iconst0 struct{}

// iconst_1 (0x04) : Push int constant
type Iconst1 struct{}

// iconst_2 (0x05) : Push int constant
type Iconst2 struct{}

// iconst_3 (0x06) : Push int constant
type Iconst3 struct{}

// iconst_4 (0x07) : Push int constant
type Iconst4 struct{}

// iconst_5 (0x08) : Push int constant
type Iconst5 struct{}

// lconst_0 (0x09) : Push long constant
type Lconst0 struct{}

// lconst_1 (0x0A) : Push long constant
type Lconst1 struct{}

// fconst_0 (0x0B) : Push float
type Fconst0 struct{}

// fconst_1 (0x0C) : Push float
type Fconst1 struct{}

// fconst_2 (0x0D) : Push float
type Fconst2 struct{}

// dconst_0 (0x0E) : Push double
type Dconst0 struct{}

// dconst_1 (0x0F) : Push double
type Dconst1 struct{}

// bipush (0x10) : Push byte
type Bipush struct {
	Byte byte
}

// sipush (0x11) : Push short
type Sipush struct {
	Short int16
}

// ldc (0x12) : Push item from run-time constant pool
type Ldc struct {
	uint8 // index in constant pool (not long/double)
}

// ldc_w (0x13) : Push item from run-time constant pool (wide index)
type LdcW struct {
	uint16 // index in constant pool (not long/double)
}

// ldc2_w (0x14) : Push long or double from run-time constant pool (wide index)
type Ldc2W struct {
	uint16 // index in constant pool (not int/float/string)
}

// iload (0x15) : Load int from local variable
type Iload struct {
	LocalIndex uint8 // local variable index
}

// lload (0x16) : Load long from local variable
type Lload struct {
	LocalIndex uint8 // local variable index
}

// fload (0x17) : Load float from local variable
type Fload struct {
	LocalIndex uint8 // local variable index
}

// dload (0x18) : Load double from local variable
type Dload struct {
	LocalIndex uint8 // local variable index
}

// aload (0x19) : Load reference from local variable
type Aload struct {
	LocalIndex uint8 // local variable index
}

// iload_0 (0x1A) : Load int from local variable
type Iload0 struct{}

// iload_1 (0x1B) : Load int from local variable
type Iload1 struct{}

// iload_2 (0x1C) : Load int from local variable
type Iload2 struct{}

// iload_3 (0x1D) : Load int from local variable
type Iload3 struct{}

// lload_0 (0x1E) : Load long from local variable
type Lload0 struct{}

// lload_1 (0x1F) : Load long from local variable
type Lload1 struct{}

// lload_2 (0x20) : Load long from local variable
type Lload2 struct{}

// lload_3 (0x21) : Load long from local variable
type Lload3 struct{}

// fload_0 (0x22) : Load float from local variable
type Fload0 struct{}

// fload_1 (0x23) : Load float from local variable
type Fload1 struct{}

// fload_2 (0x24) : Load float from local variable
type Fload2 struct{}

// fload_3 (0x25) : Load float from local variable
type Fload3 struct{}

// dload_0 (0x26) : Load double from local variable
type Dload0 struct{}

// dload_1 (0x27) : Load double from local variable
type Dload1 struct{}

// dload_2 (0x28) : Load double from local variable
type Dload2 struct{}

// dload_3 (0x29) : Load double from local variable
type Dload3 struct{}

// aload_0 (0x2A) : Load reference from local variable
type Aload0 struct{}

// aload_1 (0x2B) : Load reference from local variable
type Aload1 struct{}

// aload_2 (0x2C) : Load reference from local variable
type Aload2 struct{}

// aload_3 (0x2D) : Load reference from local variable
type Aload3 struct{}

// iaload (0x2E) : Load int from array
type Iaload struct{}

// laload (0x2F) : Load long from array
type Laload struct{}

// faload (0x30) : Load float from array
type Faload struct{}

// daload (0x31) : Load double from array
type Daload struct{}

// aaload (0x32) : Load reference from array
type Aaload struct{}

// baload (0x33) : Load byte or boolean from array
type Baload struct{}

// caload (0x34) : Load char from array
type Caload struct{}

// saload (0x35) : Load short from array
type Saload struct{}

// istore (0x36) : Store int into local variable
type Istore struct {
	LocalIndex uint8 // local variable index
}

// lstore (0x37) : Store long into local variable
type Lstore struct {
	LocalIndex uint8 // local variable index
}

// fstore (0x38) : Store float into local variable
type Fstore struct {
	LocalIndex uint8 // local variable index
}

// dstore (0x39) : Store double into local variable
type Dstore struct {
	LocalIndex uint8 // local variable index
}

// astore (0x3A) : Store reference into local variable
type Astore struct {
	LocalIndex uint8 // local variable index
}

// istore_0 (0x3B) : Store int into local variable
type Istore0 struct{}

// istore_1 (0x3C) : Store int into local variable
type Istore1 struct{}

// istore_2 (0x3D) : Store int into local variable
type Istore2 struct{}

// istore_3 (0x3E) : Store int into local variable
type Istore3 struct{}

// lstore_0 (0x3F) : Store long into local variable
type Lstore0 struct{}

// lstore_1 (0x40) : Store long into local variable
type Lstore1 struct{}

// lstore_2 (0x41) : Store long into local variable
type Lstore2 struct{}

// lstore_3 (0x42) : Store long into local variable
type Lstore3 struct{}

// fstore_0 (0x43) : Store float into local variable
type Fstore0 struct{}

// fstore_1 (0x44) : Store float into local variable
type Fstore1 struct{}

// fstore_2 (0x45) : Store float into local variable
type Fstore2 struct{}

// fstore_3 (0x46) : Store float into local variable
type Fstore3 struct{}

// dstore_0 (0x47) : Store double into local variable
type Dstore0 struct{}

// dstore_1 (0x48) : Store double into local variable
type Dstore1 struct{}

// dstore_2 (0x49) : Store double into local variable
type Dstore2 struct{}

// dstore_3 (0x4A) : Store double into local variable
type Dstore3 struct{}

// astore_0 (0x4B) : Store reference into local variable
type Astore0 struct{}

// astore_1 (0x4C) : Store reference into local variable
type Astore1 struct{}

// astore_2 (0x4D) : Store reference into local variable
type Astore2 struct{}

// astore_3 (0x4E) : Store reference into local variable
type Astore3 struct{}

// iastore (0x4F) : Store into int array
type Iastore struct{}

// lastore (0x50) : Store into long array
type Lastore struct{}

// fastore (0x51) : Store into float array
type Fastore struct{}

// dastore (0x52) : Store into double array
type Dastore struct{}

// aastore (0x53) : Store into reference array
type Aastore struct{}

// bastore (0x54) : Store into byte or boolean array
type Bastore struct{}

// castore (0x55) : Store into char array
type Castore struct{}

// sastore (0x56) : Store into short array
type Sastore struct{}

// pop (0x57) : Pop the top operand stack value
type Pop struct{}

// pop2 (0x58) : Pop the top one or two operand stack values
type Pop2 struct{}

// dup (0x59) : Duplicate the top operand stack value
type Dup struct{}

// dup_x1 (0x5A) : Duplicate the top operand stack value and insert two values down
type DupX1 struct{}

// dup_x2 (0x5B) : Duplicate the top operand stack value and insert two or three values down
type DupX2 struct{}

// dup2 (0x5C) : Duplicate the top one or two operand stack values
type Dup2 struct{}

// dup2_x1 (0x5D) : Duplicate the top one or two operand stack values and insert two or three values down
type Dup2X1 struct{}

// dup2_x2 (0x5E) : Duplicate the top one or two operand stack values and insert two three or four values down
type Dup2X2 struct{}

// swap (0x5F) : Swap the top two operand stack values
type Swap struct{}

// iadd (0x60) : Add int
type Iadd struct{}

// ladd (0x61) : Add long
type Ladd struct{}

// fadd (0x62) : Add float
type Fadd struct{}

// dadd (0x63) : Add double
type Dadd struct{}

// isub (0x64) : Subtract int
type Isub struct{}

// lsub (0x65) : Subtract long
type Lsub struct{}

// fsub (0x66) : Subtract float
type Fsub struct{}

// dsub (0x67) : Subtract double
type Dsub struct{}

// imul (0x68) : Multiply int
type Imul struct{}

// lmul (0x69) : Multiply long
type Lmul struct{}

// fmul (0x6A) : Multiply float
type Fmul struct{}

// dmul (0x6B) : Multiply double
type Dmul struct{}

// idiv (0x6C) : Divide int
type Idiv struct{}

// ldiv (0x6D) : Divide long
type Ldiv struct{}

// fdiv (0x6E) : Divide float
type Fdiv struct{}

// ddiv (0x6F) : Divide double
type Ddiv struct{}

// irem (0x70) : Remainder int
type Irem struct{}

// lrem (0x71) : Remainder long
type Lrem struct{}

// frem (0x72) : Remainder float
type Frem struct{}

// drem (0x73) : Remainder double
type Drem struct{}

// ineg (0x74) : Negate int
type Ineg struct{}

// lneg (0x75) : Negate long
type Lneg struct{}

// fneg (0x76) : Negate float
type Fneg struct{}

// dneg (0x77) : Negate double
type Dneg struct{}

// ishl (0x78) : Shift left int
type Ishl struct{}

// lshl (0x79) : Shift left long
type Lshl struct{}

// ishr (0x7A) : Arithmetic shift right int
type Ishr struct{}

// lshr (0x7B) : Arithmetic shift right long
type Lshr struct{}

// iushr (0x7C) : Logical shift right int
type Iushr struct{}

// lushr (0x7D) : Logical shift right long
type Lushr struct{}

// iand (0x7E) : Boolean AND int
type Iand struct{}

// land (0x7F) : Boolean AND long
type Land struct{}

// ior (0x80) : Boolean OR int
type Ior struct{}

// lor (0x81) : Boolean OR long
type Lor struct{}

// ixor (0x82) : Boolean XOR int
type Ixor struct{}

// lxor (0x83) : Boolean XOR long
type Lxor struct{}

// iinc (0x84) : Increment local variable by constant
type Iinc struct {
	LocalIndex uint8 // local variable index
	Const      int8  // const value
}

// i2l (0x85) : Convert int to long
type I2l struct{}

// i2f (0x86) : Convert int to float
type I2f struct{}

// i2d (0x87) : Convert int to double
type I2d struct{}

// l2i (0x88) : Convert long to int
type L2i struct{}

// l2f (0x89) : Convert long to float
type L2f struct{}

// l2d (0x8A) : Convert long to double
type L2d struct{}

// f2i (0x8B) : Convert float to int
type F2i struct{}

// f2l (0x8C) : Convert float to long
type F2l struct{}

// f2d (0x8D) : Convert float to double
type F2d struct{}

// d2i (0x8E) : Convert double to int
type D2i struct{}

// d2l (0x8F) : Convert double to long
type D2l struct{}

// d2f (0x90) : Convert double to float
type D2f struct{}

// i2b (0x91) : Convert int to byte
type I2b struct{}

// i2c (0x92) : Convert int to char
type I2c struct{}

// i2s (0x93) : Convert int to short
type I2s struct{}

// lcmp (0x94) : Compare long
type Lcmp struct{}

// fcmpl (0x95) : Compare float
type Fcmpl struct{}

// fcmpg (0x96) : Compare float
type Fcmpg struct{}

// dcmpl (0x97) : Compare double
type Dcmpl struct{}

// dcmpg (0x98) : Compare double
type Dcmpg struct{}

// ifeq (0x99) : Branch if int comparison with zero succeeds
type Ifeq struct {
	Offset int16 // jump offset (branch)
}

// ifne (0x9A) : Branch if int comparison with zero succeeds
type Ifne struct {
	Offset int16 // jump offset (branch)
}

// iflt (0x9B) : Branch if int comparison with zero succeeds
type Iflt struct {
	Offset int16 // jump offset (branch)
}

// ifge (0x9C) : Branch if int comparison with zero succeeds
type Ifge struct {
	Offset int16 // jump offset (branch)
}

// ifgt (0x9D) : Branch if int comparison with zero succeeds
type Ifgt struct {
	Offset int16 // jump offset (branch)
}

// ifle (0x9E) : Branch if int comparison with zero succeeds
type Ifle struct {
	Offset int16 // jump offset (branch)
}

// if_icmpeq (0x9F) : Branch if int comparison succeeds
type IfIcmpeq struct {
	Offset int16 // jump offset (branch)
}

// if_icmpne (0xA0) : Branch if int comparison succeeds
type IfIcmpne struct {
	Offset int16 // jump offset (branch)
}

// if_icmplt (0xA1) : Branch if int comparison succeeds
type IfIcmplt struct {
	Offset int16 // jump offset (branch)
}

// if_icmpge (0xA2) : Branch if int comparison succeeds
type IfIcmpge struct {
	Offset int16 // jump offset (branch)
}

// if_icmpgt (0xA3) : Branch if int comparison succeeds
type IfIcmpgt struct {
	Offset int16 // jump offset (branch)
}

// if_icmple (0xA4) : Branch if int comparison succeeds
type IfIcmple struct {
	Offset int16 // jump offset (branch)
}

// if_acmpeq (0xA5) : Branch if reference comparison succeeds
type IfAcmpeq struct {
	Offset int16 // jump offset (branch)
}

// if_acmpne (0xA6) : Branch if reference comparison succeeds
type IfAcmpne struct {
	Offset int16 // jump offset (branch)
}

// goto (0xA7) : Branch always
type Goto struct {
	Offset int16 // jump offset (branch)
}

// jsr (0xA8) : Jump subroutine
type Jsr struct {
	Offset int16 // jump offset (branch)
}

// ret (0xA9) : Return from subroutine
type Ret struct {
	LocalIndex uint8 // local variable index
}

// // tableswitch (0xAA) : Access jump table by index and jump // TODO
// type Tableswitch struct {
//   [0] // 0 à 3 octets vide pour le padding
//   int32 // default offset
//   int32 // low (<= high)
//   int32 // high (>= low)
//   [int32] // high - low + 1 offsets
// }

// // lookupswitch (0xAB) : Access jump table by key match and jump // TODO
// type Lookupswitch struct {
//   [0] // 0 à 3 octets vide pour le padding
//   int32 // default offset
//   int32 // pairs length >= 0
//   [int32 int32] // pairs (match - offset)
// }

// ireturn (0xAC) : Return int from method
type Ireturn struct{}

// lreturn (0xAD) : Return long from method
type Lreturn struct{}

// freturn (0xAE) : Return float from method
type Freturn struct{}

// dreturn (0xAF) : Return double from method
type Dreturn struct{}

// areturn (0xB0) : Return reference from method
type Areturn struct{}

// return (0xB1) : Return void from method
type Return struct{}

// getstatic (0xB2) : Get static field from class
type Getstatic struct {
	Class string
	Name  string
	Type  string
}

// putstatic (0xB3) : Set static field in class
type Putstatic struct {
	Class string
	Name  string
	Type  string
}

// getfield (0xB4) : Fetch field from object
type Getfield struct {
	Class string
	Name  string
	Type  string
}

// putfield (0xB5) : Set field in object
type Putfield struct {
	Class string
	Name  string
	Type  string
}

// invokevirtual (0xB6) : Invoke instance method; dispatch based on class
type Invokevirtual struct {
	Class string
	Name  string
	Type  string
}

// invokespecial (0xB7) : Invoke instance method;  direct invocation of instance initialization methods and methods of the current class and its supertypes
type Invokespecial struct {
	Class string
	Name  string
	Type  string
}

// invokestatic (0xB8) : Invoke a class (static) method
type Invokestatic struct {
	Class string
	Name  string
	Type  string
}

// // invokeinterface (0xB9) : Invoke interface method // TODO
// type Invokeinterface struct {
// 	uint16 // index in constant pool (InterfaceMethodref)
// 	uint8  // count
// 	uint8  // 0
// }

// // invokedynamic (0xBA) : Invoke a dynamically-computed call site // TODO
// type Invokedynamic struct {
// 	uint16 // index in constant pool (InvokeDynamic)
// 	uint8  // 0
// 	uint8  // 0
// }

// new (0xBB) : Create new object
type New string

// // newarray (0xBC) : Create new array // TODO
// type Newarray struct {
// 	uint8 // type -> T_BOOLEAN=4 ; T_CHAR=5 ; T_FLOAT=6 ; T_DOUBLE=7 ; T_BYTE=8 ; T_SHORT=9 ; T_INT=10 ; T_LONG=11
// }

// anewarray (0xBD) : Create new array of reference
type Anewarray string

// arraylength (0xBE) : Get length of array
type Arraylength struct{}

// athrow (0xBF) : Throw exception or error
type Athrow struct{}

// checkcast (0xC0) : Check whether object is of given type
type Checkcast string

// instanceof (0xC1) : Determine if object is of given type
type Instanceof string

// monitorenter (0xC2) : Enter monitor for object
type Monitorenter struct{}

// monitorexit (0xC3) : Exit monitor for object
type Monitorexit struct{}

// // wide (0xC4) : Extend local variable index by additional bytes // TODO
// type Wide struct {
//   uint8 // opcode in iload fload aload lload dload istore fstore astore lstore dstore ret iinc
//   LocalIndex uint16 // local variable index
//   [uint16] // const value ONLY in iinc case
// }

// multianewarray (0xC5) : Create new multidimensional array
type Multianewarray struct {
	uint16          // index in constant pool (Class)
	Dimension uint8 // dimension >= 1
}

// ifnull (0xC6) : Branch if reference is null
type Ifnull struct {
	Offset int16 // jump offset (branch)
}

// ifnonnull (0xC7) : Branch if reference not null
type Ifnonnull struct {
	Offset int16 // jump offset (branch)
}

// goto_w (0xC8) : Branch always (wide index)
type GotoW struct {
	Offset int32 // jump offset (branch)
}

// jsr_w (0xC9) : Jump subroutine (wide index)
type JsrW struct {
	Offset int32 // jump offset (branch)
}
