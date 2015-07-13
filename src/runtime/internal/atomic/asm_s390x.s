// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "textflag.h"

// bool ·Cas(uint32 *ptr, uint32 old, uint32 new)
// Atomically:
//	if(*val == old){
//		*val = new;
//		return 1;
//	} else
//		return 0;
TEXT ·Cas(SB), NOSPLIT, $0-17
	MOVD	ptr+0(FP), R3
	MOVWZ	old+8(FP), R4
	MOVWZ	new+12(FP), R5
	CS	R4, R5, 0(R3)    //  if (R4 == 0(R3)) then 0(R3)= R5
	BNE	cas_fail
	MOVD	$1, R3
	MOVB	R3, ret+16(FP)
	RET
cas_fail:
	MOVD	$0, R3
	MOVB	R3, ret+16(FP)
	RET

// bool	·Cas64(uint64 *ptr, uint64 old, uint64 new)
// Atomically:
//	if(*val == *old){
//		*val = new;
//		return 1;
//	} else {
//		return 0;
//	}
TEXT ·Cas64(SB), NOSPLIT, $0-25
	MOVD	ptr+0(FP), R3
	MOVD	old+8(FP), R4
	MOVD	new+16(FP), R5
	CSG	R4, R5, 0(R3)    //  if (R4 == 0(R3)) then 0(R3)= R5
	BNE	cas64_fail
	MOVD	$1, R3
	MOVB	R3, ret+24(FP)
	RET
cas64_fail:
	MOVD	$0, R3
	MOVB	R3, ret+24(FP)
	RET

TEXT ·Casuintptr(SB), NOSPLIT, $0-25
	BR	·Cas64(SB)

TEXT ·Loaduintptr(SB), NOSPLIT, $0-16
	BR	·Load64(SB)

TEXT ·Loaduint(SB), NOSPLIT, $0-16
	BR	·Load64(SB)

TEXT ·Storeuintptr(SB), NOSPLIT, $0-16
	BR	·Store64(SB)

TEXT ·Loadint64(SB), NOSPLIT, $0-16
	BR	·Load64(SB)

TEXT ·Xadduintptr(SB), NOSPLIT, $0-24
	BR	·Xadd64(SB)

TEXT ·Xaddint64(SB), NOSPLIT, $0-16
	BR	·Xadd64(SB)

// bool ·Casp1(void **val, void *old, void *new)
// Atomically:
//	if(*val == old){
//		*val = new;
//		return 1;
//	} else
//		return 0;
TEXT ·Casp1(SB), NOSPLIT, $0-25
	BR ·Cas64(SB)

// uint32 ·Xadd(uint32 volatile *ptr, int32 delta)
// Atomically:
//	*val += delta;
//	return *val;
TEXT ·Xadd(SB), NOSPLIT, $0-20
	MOVD	ptr+0(FP), R4
	MOVW	delta+8(FP), R5
repeat:
	MOVW	(R4), R3
	MOVD	R3, R6
	ADD	R5, R3
	CS	R6, R3, (R4)    // if (R6==(R4)) then (R4)=R3
	BNE	repeat
	MOVW	R3, ret+16(FP)
	RET

TEXT ·Xadd64(SB), NOSPLIT, $0-24
	MOVD	ptr+0(FP), R4
	MOVD	delta+8(FP), R5
repeat:
	MOVD	(R4), R3
	MOVD	R3, R6
	ADD	R5, R3
	CSG	R6, R3, (R4)    // if (R6==(R4)) then (R4)=R3
	BNE	repeat
	MOVD	R3, ret+16(FP)
	RET

TEXT ·Xchg(SB), NOSPLIT, $0-20
	MOVD	ptr+0(FP), R4
	MOVW	new+8(FP), R3
repeat:
	MOVW	(R4), R6
	CS	R6, R3, (R4)    // if (R6==(R4)) then (R4)=R3
	BNE	repeat
	MOVW	R6, ret+16(FP)
	RET

TEXT ·Xchg64(SB), NOSPLIT, $0-24
	MOVD	ptr+0(FP), R4
	MOVD	new+8(FP), R3
repeat:
	MOVD	(R4), R6
	CSG	R6, R3, (R4)    // if (R6==(R4)) then (R4)=R3
	BNE	repeat
	MOVD	R6, ret+16(FP)
	RET

TEXT ·Xchguintptr(SB), NOSPLIT, $0-24
	BR	·Xchg64(SB)

TEXT ·Storep1(SB), NOSPLIT, $0-16
	BR	·Store64(SB)

// on Z, load & store both are atomic operations
TEXT ·Store(SB), NOSPLIT, $0-12
	MOVD	ptr+0(FP), R3
	MOVW	val+8(FP), R4
	SYNC
	MOVW	R4, 0(R3)
	RET

TEXT ·Store64(SB), NOSPLIT, $0-16
	MOVD	ptr+0(FP), R3
	MOVD	val+8(FP), R4
	SYNC
	MOVD	R4, 0(R3)
	RET

// func Or8(addr *uint8, v uint8)
TEXT ·Or8(SB), NOSPLIT, $0-9
	MOVD    ptr+0(FP), R3
	MOVBZ   val+8(FP), R4
	// Calculate shift.
	AND	$3, R3, R5
	XOR	$3, R5 // big endian - flip direction
	SLD	$3, R5 // MUL $8, R5
	SLD	R5, R4
	// Align ptr down to 4 bytes so we can use 32-bit load/store.
	SRD	$2, R3, R5
	SLD	$2, R5
	MOVWZ	0(R5), R6
again:
	OR	R4, R6, R7
	CS	R6, R7, 0(R5) //  if (R6 == 0(R5)) then 0(R5)= R7 else R6= 0(R5)
	BNE	again
	RET

// func And8(addr *uint8, v uint8)
TEXT ·And8(SB), NOSPLIT, $0-9
	MOVD    ptr+0(FP), R3
	MOVBZ   val+8(FP), R4
	// Calculate shift.
	AND	$3, R3, R5
	XOR	$3, R5 // big endian - flip direction
	SLD	$3, R5 // MUL $8, R5
	OR	$-256, R4 // create 0xffffffffffffffxx
	BYTE	$0xEB // RLLG R5, R4
	BYTE	$0x44
	BYTE	$0x50
	BYTE	$0x00
	BYTE	$0x00
	BYTE	$0x1C
	// Align ptr down to 4 bytes so we can use 32-bit load/store.
	SRD	$2, R3, R5
	SLD	$2, R5
	MOVWZ	0(R5), R6
again:
	AND	R4, R6, R7
	CS	R6, R7, 0(R5) //  if (R6 == 0(R5)) then 0(R5)= R7 else R6= 0(R5)
	BNE	again
	RET
