// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build s390x

#include "textflag.h"

// void runtime·memclr(void*, uintptr)
TEXT runtime·memclr(SB),NOSPLIT,$0-16
	MOVD	ptr+0(FP), R4
	MOVD	n+8(FP), R5

start:
	CMP	R5, $3
	BLE	clear0to3
	CMP	R5, $7
	BLE	clear4to7
	CMP	R5, $11
	BLE	clear8to11
	CMP	R5, $15
	BLE	clear12to15
	CMP	R5, $32
	BGE	clearmt32
	MOVD	R0, 0(R4)
	MOVD	R0, 8(R4)
	ADD	$16, R4
	SUB	$16, R5
	BR	start

clear0to3:
	CMP	R5, $0
	BEQ	done
	CMP	R5, $1
	BNE	clear2
	MOVB	R0, 0(R4)
	RET
clear2:
	CMP	R5, $2
	BNE	clear3
	MOVH	R0, 0(R4)
	RET
clear3:
	MOVH	R0, 0(R4)
	MOVB	R0, 2(R4)
	RET

clear4to7:
	CMP	R5, $4
	BNE	clear5
	MOVW	R0, 0(R4)
	RET
clear5:
	CMP	R5, $5
	BNE	clear6
	MOVW	R0, 0(R4)
	MOVB	R0, 4(R4)
	RET
clear6:
	CMP	R5, $6
	BNE	clear7
	MOVW	R0, 0(R4)
	MOVH	R0, 4(R4)
	RET
clear7:
	MOVW	R0, 0(R4)
	MOVH	R0, 4(R4)
	MOVB	R0, 6(R4)
	RET

clear8to11:
	CMP	R5, $8
	BNE	clear9
	MOVD	R0, 0(R4)
	RET
clear9:
	CMP	R5, $9
	BNE	clear10
	MOVD	R0, 0(R4)
	MOVB	R0, 8(R4)
	RET
clear10:
	CMP	R5, $10
	BNE	clear11
	MOVD	R0, 0(R4)
	MOVH	R0, 8(R4)
	RET
clear11:
	MOVD	R0, 0(R4)
	MOVH	R0, 8(R4)
	MOVB	R0, 10(R4)
	RET

clear12to15:
	CMP	R5, $12
	BNE	clear13
	MOVD	R0, 0(R4)
	MOVW	R0, 8(R4)
	RET
clear13:
	CMP	R5, $13
	BNE	clear14
	MOVD	R0, 0(R4)
	MOVW	R0, 8(R4)
	MOVB	R0, 12(R4)
	RET
clear14:
	CMP	R5, $14
	BNE	clear15
	MOVD	R0, 0(R4)
	MOVW	R0, 8(R4)
	MOVH	R0, 12(R4)
	RET
clear15:
	MOVD	R0, 0(R4)
	MOVW	R0, 8(R4)
	MOVH	R0, 12(R4)
	MOVB	R0, 14(R4)
	RET

clearmt32:
	CMP	R5, $256
	BLT	clearlt256
	XC	$256, 0(R4), 0(R4)
	ADD	$256, R4
	ADD	$-256, R5
	BR	clearmt32
clearlt256:
	CMP	R5, $0
	BEQ	done
	ADD	$-1, R5
	EXRL	$runtime·memclr_s390x_exrl_xc(SB), R5
done:
	RET

// DO NOT CALL - target for exrl (execute relative long) instruction.
TEXT runtime·memclr_s390x_exrl_xc(SB),NOSPLIT, $0-0
	XC	$1, 0(R4), 0(R4)
	MOVD	R0, 0(R0)
	RET

