// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.


#include "textflag.h"

#define res_ptr R1
#define x_ptr R2
#define y_ptr R3
#define X0    V0
#define X1    V1
#define Y0    V2
#define Y1    V3
#define M0    V4
#define M1    V5
#define T0    V6
#define T1    V7
#define T2    V8
#define YDIG  V9

#define ADD1  V16
#define ADD1H V17
#define ADD2  V18
#define ADD2H V19
#define RED1  V20
#define RED1H V21
#define RED2  V22
#define RED2H V23
#define CAR1  V24
#define CAR1M V25

#define MK0   V30
#define K0    V31

DATA p256ordK0<>+0x00(SB)/4, $0xee00bc4f
DATA p256ord<>+0x00(SB)/8, $0xffffffff00000000
DATA p256ord<>+0x08(SB)/8, $0xffffffffffffffff
DATA p256ord<>+0x10(SB)/8, $0xbce6faada7179e84
DATA p256ord<>+0x18(SB)/8, $0xf3b9cac2fc632551
DATA p256<>+0x00(SB)/8, $0xffffffff00000001 //P256
DATA p256<>+0x08(SB)/8, $0x0000000000000000 //P256
DATA p256<>+0x10(SB)/8, $0x00000000ffffffff //P256
DATA p256<>+0x18(SB)/8, $0xffffffffffffffff //P256
DATA p256<>+0x20(SB)/8, $0x0c0d0e0f1c1d1e1f // SEL d1 d0 d1 d0
DATA p256<>+0x28(SB)/8, $0x0c0d0e0f1c1d1e1f // SEL d1 d0 d1 d0
DATA p256<>+0x30(SB)/8, $0x0000000010111213 // SEL 0  d1 d0  0
DATA p256<>+0x38(SB)/8, $0x1415161700000000 // SEL 0  d1 d0  0
DATA p256mul<>+0x00(SB)/8, $0xffffffff00000001 //P256
DATA p256mul<>+0x08(SB)/8, $0x0000000000000000 //P256
DATA p256mul<>+0x10(SB)/8, $0x00000000ffffffff //P256
DATA p256mul<>+0x18(SB)/8, $0xffffffffffffffff //P256
DATA p256mul<>+0x20(SB)/8, $0x1c1d1e1f00000000 // SEL d0  0  0 d0
DATA p256mul<>+0x28(SB)/8, $0x000000001c1d1e1f // SEL d0  0  0 d0
DATA p256mul<>+0x30(SB)/8, $0x0001020304050607 // SEL d0  0 d1 d0
DATA p256mul<>+0x38(SB)/8, $0x1c1d1e1f0c0d0e0f // SEL d0  0 d1 d0
DATA p256mul<>+0x40(SB)/8, $0x040506071c1d1e1f // SEL  0 d1 d0 d1
DATA p256mul<>+0x48(SB)/8, $0x0c0d0e0f1c1d1e1f // SEL  0 d1 d0 d1
DATA p256mul<>+0x50(SB)/8, $0x0405060704050607 // SEL  0  0 d1 d0
DATA p256mul<>+0x58(SB)/8, $0x1c1d1e1f0c0d0e0f // SEL  0  0 d1 d0
DATA p256mul<>+0x60(SB)/8, $0x0c0d0e0f1c1d1e1f // SEL d1 d0 d1 d0
DATA p256mul<>+0x68(SB)/8, $0x0c0d0e0f1c1d1e1f // SEL d1 d0 d1 d0
DATA p256mul<>+0x70(SB)/8, $0x141516170c0d0e0f // SEL 0  d1 d0  0
DATA p256mul<>+0x78(SB)/8, $0x1c1d1e1f14151617 // SEL 0  d1 d0  0
DATA p256mont<>+0x00(SB)/8,$0x1010101008090a0b
DATA p256mont<>+0x08(SB)/8,$0x0c0d0e0f10101010
GLOBL p256const0<>(SB), 8, $8
GLOBL p256const1<>(SB), 8, $8
GLOBL p256ordK0<>(SB), 8, $4
GLOBL p256ord<>(SB), 8, $32
GLOBL p256<>(SB), 8, $64
GLOBL p256mul<>(SB), 8, $128
GLOBL p256mont<>(SB),8,$16

/* ---------------------------------------*/
// func p256MovCond(res, a, b []byte, cond int)
// If cond == 0 res=b, else res=a
TEXT ·p256MovCond(SB),NOSPLIT,$0
	MOVD res+0(FP), res_ptr
	MOVD a+24(FP), x_ptr
	MOVD b+48(FP), y_ptr

	VLREPG cond+72(FP), M0
	VZERO  M1
	VCEQG  M0,M1,M0

	VL   (16*0)(x_ptr), X1
	VL   (16*0)(y_ptr), Y1
	VSEL  Y1,X1,M0,X1
	VST   X1,(16*0)(res_ptr)

	VL   (16*1)(x_ptr), X1
	VL   (16*1)(y_ptr), Y1
	VSEL  Y1,X1,M0,X1
	VST   X1,(16*1)(res_ptr)

	VL   (16*2)(x_ptr), X1
	VL   (16*2)(y_ptr), Y1
	VSEL  Y1,X1,M0,X1
	VST   X1,(16*2)(res_ptr)

	VL   (16*3)(x_ptr), X1
	VL   (16*3)(y_ptr), Y1
	VSEL  Y1,X1,M0,X1
	VST   X1,(16*3)(res_ptr)

	VL   (16*4)(x_ptr), X1
	VL   (16*4)(y_ptr), Y1
	VSEL  Y1,X1,M0,X1
	VST   X1,(16*4)(res_ptr)

	VL   (16*5)(x_ptr), X1
	VL   (16*5)(y_ptr), Y1
	VSEL  Y1,X1,M0,X1
	VST   X1,(16*5)(res_ptr)

	RET



/* ---------------------------------------*/
// func p256NegCond(val []byte, cond int)
TEXT  ·p256NegCond(SB),NOSPLIT,$0
	MOVD val+0(FP), res_ptr

	VLREPG cond+24(FP), M0
	//compare with 1 to generate and mask
	VREPIG $1, M1
	VCEQG M0, M1, M0

	VGBM $0x0fff, T0
	MOVD p256<>(SB),R6
	VL   0x10(R6),T1

	VN   T0,M0,T0
	VN   T1,M0,T1

	VL   (16*0)(res_ptr), X1
	VL   (16*1)(res_ptr), X0

	VSCBIQ  X0,T0,M1
	VSQ     X0,T0,X0
	VSBIQ   X1,T1,M1,X1

	VST     X1, (16*0)(res_ptr)
	VST     X0, (16*1)(res_ptr)

	RET

/* ---------------------------------------*/
	// Constant time point access to arbitrary point table.
	// Indexed from 1 to 15, with -1 offset
	// (index 0 is implicitly point at infinity)
	// func p256Select(point, table []byte, idx int)
TEXT ·p256Select(SB),NOSPLIT,$0
	VLREPG idx+48(FP),V16  // V16 = {idx,idx}
	MOVD table+24(FP),R2
	MOVD point+0(FP),R3

	VREPIG $1, V17   // V17 = {1,1}

	VZERO V20
	VZERO V21
	VZERO V22
	VZERO V23
	VZERO V24
	VZERO V25


	MOVD $16, R4

	VLR V17, V18  // Initialize cur_idx = 1

loop_select:

		VCEQG V16, V18, V19
		VAG   V18, V18, V17


		VL (16*0)(R2), V26
		VL (16*1)(R2), V27
		VL (16*2)(R2), V28
		VL (16*3)(R2), V29
		VL (16*4)(R2), V30
		VL (16*5)(R2), V31
		ADD $(16*6), R2

		VN   V19, V26, V26
		VN   V19, V27, V27
		VN   V19, V28, V28
		VN   V19, V29, V29
		VN   V19, V30, V30
		VN   V19, V31, V31


		VX   V26, V20, V20
		VX   V27, V21, V21
		VX   V28, V22, V22
		VX   V29, V23, V23
		VX   V30, V24, V24
		VX   V31, V25, V25

	        ADD  $-1,R4
	        BLT  loop_select
	        
		//BRCTG R4, loop_select

	VST   V20,(16*0)(R3)
	VST   V21,(16*1)(R3)
	VST   V22,(16*2)(R3)
	VST   V23,(16*3)(R3)
	VST   V24,(16*4)(R3)
	VST   V25,(16*5)(R3)

	RET

/* ---------------------------------------*/
	// Constant time point access to base point table.
	// func p256SelectBase(point, table []uint64, idx int)
TEXT ·p256SelectBase(SB),NOSPLIT,$0
	VLREPG idx+48(FP),V16  
	MOVD table+24(FP),R2
	MOVD point+0(FP),R3

	VREPIG   $1, V17

	VZERO V20
	VZERO V21
	VZERO V22
	VZERO V23

	MOVD $32, R4

	VLR V17, V18

loop_select_base:

		VCEQG V16, V18, V19
		VAG   V18, V17, V18

		VL (16*0)(R2), V24
		VL (16*1)(R2), V25
		VL (16*2)(R2), V26
		VL (16*3)(R2), V27

		VL (16*4)(R2), V28
		VL (16*5)(R2), V29
		VL (16*6)(R2), V30
		VL (16*7)(R2), V31

		ADD $(16*8), R2

		VN   V19, V24, V24
		VN   V19, V25, V25
		VN   V19, V26, V26
		VN   V19, V27, V27

		VCEQG V16, V18, V19
		VAG   V18, V17, V18

		VN   V19, V28, V28
		VN   V19, V29, V29
		VN   V19, V30, V30
		VN   V19, V31, V31

		VX   V24, V24, V20
		VX   V25, V25, V21
	        VX   V26, V26, V22
	        VX   V27, V27, V23

		VX   V28, V28, V20
		VX   V29, V29, V21
	        VX   V30, V30, V22
	        VX   V31, V31, V23
	        ADD  $-1,R4
	        BLT  loop_select_base
	
		//BRCTG R4, loop_select_base

	VST   V20, (16*0)(R3)
	VST   V21, (16*1)(R3)
	VST   V22, (16*2)(R3)
	VST   V23, (16*3)(R3)

	RET
	
	
/* ---------------------------------------*/
// func p256OrdMul(res, in1, in2 []byte)
TEXT ·p256OrdMul(SB),NOSPLIT,$0
	MOVD res+0(FP), res_ptr
	MOVD in1+24(FP), x_ptr
	MOVD in2+48(FP), y_ptr

	VZERO T2
	MOVD	$p256ordK0<>+0x00(SB), R4
	//VLEF    $3, 0(R4), K0
	WORD	$0xE7F40000
	BYTE    $0x38
	BYTE    $0x03 //
	MOVD	$p256ord<>+0x00(SB), R4
	VL	16(R4), M0
	VL	0(R4),  M1

	VL (1*16)(x_ptr), X0
	VL (0*16)(x_ptr), X1
	VL (1*16)(y_ptr), Y0
	VL (0*16)(y_ptr), Y1

	//---------------------------------------------------------------------------/
	VREPF $3,  Y0,   YDIG
	VMLF  X0,  YDIG, ADD1
	VMLF  ADD1,K0,   MK0
	VREPF $3,  MK0,  MK0

	VMLF  X1,YDIG, ADD2
	VMLHF X0,YDIG, ADD1H
	VMLHF X1,YDIG, ADD2H

	VMALF  M0,MK0, ADD1, RED1
	VMALHF M0,MK0, ADD1, RED1H
	VMALF  M1,MK0, ADD2, RED2
	VMALHF M1,MK0, ADD2, RED2H

	VSLDB  $12, RED2,RED1, RED1
	VSLDB  $12, T2,  RED2, RED2

	VACCQ  RED1, ADD1H,CAR1
	VAQ    RED1, ADD1H,T0
	VACCQ  RED1H,T0,   CAR1M
	VAQ    RED1H,T0,   T0

	//<< ready for next MK0

	VACQ    RED2, ADD2H,CAR1,  T1
	VACCCQ  RED2, ADD2H,CAR1,  CAR1
	VACCCQ  RED2H,T1, CAR1M, T2
	VACQ    RED2H,T1, CAR1M, T1
	VAQ     CAR1, T2,T2

	//---------------------------------------------------
    /**
   	 * ---+--------+--------+
	 *  T2|   T1   |   T0   |
	 * ---+--------+--------+
	 *           *(add)*
	 *    +--------+--------+
	 *    |   X1   |   X0   |
	 *    +--------+--------+
	 *           *(mul)*
 	 *    +--------+--------+
	 *    |  YDIG  |  YDIG  |
	 *    +--------+--------+
 	 *           *(add)*
	 *    +--------+--------+
	 *    |   M1   |   M0   |
	 *    +--------+--------+
	 *           *(mul)*
 	 *    +--------+--------+
	 *    |   MK0  |   MK0  |
	 *    +--------+--------+
	 *
	 *   ---------------------
   	 *
   	 *    +--------+--------+
	 *    |  ADD2  |  ADD1  |
	 *    +--------+--------+
   	 *  +--------+--------+
	 *  | ADD2H  | ADD1H  |
	 *  +--------+--------+
  	 *    +--------+--------+
	 *    |  RED2  |  RED1  |
	 *    +--------+--------+
   	 *  +--------+--------+
	 *  | RED2H  | RED1H  |
	 *  +--------+--------+
	 */
	VREPF $2,  Y0,   YDIG
	VMALF X0,  YDIG, T0, ADD1
	VMLF  ADD1,K0,   MK0
	VREPF $3,  MK0,  MK0

	VMALF  X1,YDIG, T1, ADD2
	VMALHF X0,YDIG, T0, ADD1H
	VMALHF X1,YDIG, T1, ADD2H

	VMALF  M0,MK0, ADD1, RED1
	VMALHF M0,MK0, ADD1, RED1H
	VMALF  M1,MK0, ADD2, RED2
	VMALHF M1,MK0, ADD2, RED2H

	VSLDB  $12, RED2,RED1, RED1
	VSLDB  $12, T2,  RED2, RED2

	VACCQ  RED1, ADD1H,CAR1
	VAQ    RED1, ADD1H,T0
	VACCQ  RED1H,T0,   CAR1M
	VAQ    RED1H,T0,   T0

	//<< ready for next MK0

	VACQ    RED2, ADD2H,CAR1,  T1
	VACCCQ  RED2, ADD2H,CAR1,  CAR1
	VACCCQ  RED2H,T1, CAR1M, T2
	VACQ    RED2H,T1, CAR1M, T1
	VAQ     CAR1, T2,T2

	//---------------------------------------------------
	VREPF $1,  Y0,   YDIG
	VMALF X0,  YDIG, T0, ADD1
	VMLF  ADD1,K0,   MK0
	VREPF $3,  MK0,  MK0

	VMALF  X1,YDIG, T1, ADD2
	VMALHF X0,YDIG, T0, ADD1H
	VMALHF X1,YDIG, T1, ADD2H

	VMALF  M0,MK0, ADD1, RED1
	VMALHF M0,MK0, ADD1, RED1H
	VMALF  M1,MK0, ADD2, RED2
	VMALHF M1,MK0, ADD2, RED2H

	VSLDB  $12, RED2,RED1, RED1
	VSLDB  $12, T2,  RED2, RED2

	VACCQ  RED1, ADD1H,CAR1
	VAQ    RED1, ADD1H,T0
	VACCQ  RED1H,T0,   CAR1M
	VAQ    RED1H,T0,   T0

	//<< ready for next MK0

	VACQ    RED2, ADD2H,CAR1,  T1
	VACCCQ  RED2, ADD2H,CAR1,  CAR1
	VACCCQ  RED2H,T1, CAR1M, T2
	VACQ    RED2H,T1, CAR1M, T1
	VAQ     CAR1, T2,T2

	//---------------------------------------------------
	VREPF $0,  Y0,   YDIG
	VMALF X0,  YDIG, T0, ADD1
	VMLF  ADD1,K0,   MK0
	VREPF $3,  MK0,  MK0

	VMALF  X1,YDIG, T1, ADD2
	VMALHF X0,YDIG, T0, ADD1H
	VMALHF X1,YDIG, T1, ADD2H

	VMALF  M0,MK0, ADD1, RED1
	VMALHF M0,MK0, ADD1, RED1H
	VMALF  M1,MK0, ADD2, RED2
	VMALHF M1,MK0, ADD2, RED2H

	VSLDB  $12, RED2,RED1, RED1
	VSLDB  $12, T2,  RED2, RED2

	VACCQ  RED1, ADD1H,CAR1
	VAQ    RED1, ADD1H,T0
	VACCQ  RED1H,T0,   CAR1M
	VAQ    RED1H,T0,   T0

	//<< ready for next MK0

	VACQ    RED2, ADD2H,CAR1,  T1
	VACCCQ  RED2, ADD2H,CAR1,  CAR1
	VACCCQ  RED2H,T1, CAR1M, T2
	VACQ    RED2H,T1, CAR1M, T1
	VAQ     CAR1, T2,T2

	//---------------------------------------------------
	VREPF $3,  Y1,   YDIG
	VMALF X0,  YDIG, T0, ADD1
	VMLF  ADD1,K0,   MK0
	VREPF $3,  MK0,  MK0

	VMALF  X1,YDIG, T1, ADD2
	VMALHF X0,YDIG, T0, ADD1H
	VMALHF X1,YDIG, T1, ADD2H

	VMALF  M0,MK0, ADD1, RED1
	VMALHF M0,MK0, ADD1, RED1H
	VMALF  M1,MK0, ADD2, RED2
	VMALHF M1,MK0, ADD2, RED2H

	VSLDB  $12, RED2,RED1, RED1
	VSLDB  $12, T2,  RED2, RED2

	VACCQ  RED1, ADD1H,CAR1
	VAQ    RED1, ADD1H,T0
	VACCQ  RED1H,T0,   CAR1M
	VAQ    RED1H,T0,   T0

	//<< ready for next MK0

	VACQ    RED2, ADD2H,CAR1,  T1
	VACCCQ  RED2, ADD2H,CAR1,  CAR1
	VACCCQ  RED2H,T1, CAR1M, T2
	VACQ    RED2H,T1, CAR1M, T1
	VAQ     CAR1, T2,T2

	//---------------------------------------------------
	VREPF $2,  Y1,   YDIG
	VMALF X0,  YDIG, T0, ADD1
	VMLF  ADD1,K0,   MK0
	VREPF $3,  MK0,  MK0

	VMALF  X1,YDIG, T1, ADD2
	VMALHF X0,YDIG, T0, ADD1H
	VMALHF X1,YDIG, T1, ADD2H

	VMALF  M0,MK0, ADD1, RED1
	VMALHF M0,MK0, ADD1, RED1H
	VMALF  M1,MK0, ADD2, RED2
	VMALHF M1,MK0, ADD2, RED2H

	VSLDB  $12, RED2,RED1, RED1
	VSLDB  $12, T2,  RED2, RED2

	VACCQ  RED1, ADD1H,CAR1
	VAQ    RED1, ADD1H,T0
	VACCQ  RED1H,T0,   CAR1M
	VAQ    RED1H,T0,   T0

	//<< ready for next MK0

	VACQ    RED2, ADD2H,CAR1,  T1
	VACCCQ  RED2, ADD2H,CAR1,  CAR1
	VACCCQ  RED2H,T1, CAR1M, T2
	VACQ    RED2H,T1, CAR1M, T1
	VAQ     CAR1, T2,T2

	//---------------------------------------------------
	VREPF $1,  Y1,   YDIG
	VMALF X0,  YDIG, T0, ADD1
	VMLF  ADD1,K0,   MK0
	VREPF $3,  MK0,  MK0

	VMALF  X1,YDIG, T1, ADD2
	VMALHF X0,YDIG, T0, ADD1H
	VMALHF X1,YDIG, T1, ADD2H

	VMALF  M0,MK0, ADD1, RED1
	VMALHF M0,MK0, ADD1, RED1H
	VMALF  M1,MK0, ADD2, RED2
	VMALHF M1,MK0, ADD2, RED2H

	VSLDB  $12, RED2,RED1, RED1
	VSLDB  $12, T2,  RED2, RED2

	VACCQ  RED1, ADD1H,CAR1
	VAQ    RED1, ADD1H,T0
	VACCQ  RED1H,T0,   CAR1M
	VAQ    RED1H,T0,   T0

	//<< ready for next MK0

	VACQ    RED2, ADD2H,CAR1,  T1
	VACCCQ  RED2, ADD2H,CAR1,  CAR1
	VACCCQ  RED2H,T1, CAR1M, T2
	VACQ    RED2H,T1, CAR1M, T1
	VAQ     CAR1, T2,T2

	//---------------------------------------------------
	VREPF $0,  Y1,   YDIG
	VMALF X0,  YDIG, T0, ADD1
	VMLF  ADD1,K0,   MK0
	VREPF $3,  MK0,  MK0

	VMALF  X1,YDIG, T1, ADD2
	VMALHF X0,YDIG, T0, ADD1H
	VMALHF X1,YDIG, T1, ADD2H

	VMALF  M0,MK0, ADD1, RED1
	VMALHF M0,MK0, ADD1, RED1H
	VMALF  M1,MK0, ADD2, RED2
	VMALHF M1,MK0, ADD2, RED2H

	VSLDB  $12, RED2,RED1, RED1
	VSLDB  $12, T2,  RED2, RED2

	VACCQ  RED1, ADD1H,CAR1
	VAQ    RED1, ADD1H,T0
	VACCQ  RED1H,T0,   CAR1M
	VAQ    RED1H,T0,   T0

	//<< ready for next MK0

	VACQ    RED2, ADD2H,CAR1,  T1
	VACCCQ  RED2, ADD2H,CAR1,  CAR1
	VACCCQ  RED2H,T1, CAR1M, T2
	VACQ    RED2H,T1, CAR1M, T1
	VAQ     CAR1, T2,T2

	//---------------------------------------------------

	VZERO   RED1
	VSCBIQ  M0,T0, CAR1
	VSQ     M0,T0, ADD1
	VSBCBIQ T1,M1,CAR1, CAR1M
	VSBIQ   T1,M1,CAR1, ADD2
	VSBIQ   T2,RED1,CAR1M, T2

	//what output to use, ADD2||ADD1 or T1||T0?
	VSEL    T0,ADD1,T2, T0
	VSEL    T1,ADD2,T2, T1

	VST T0, (1*16)(res_ptr)
	VST T1, (0*16)(res_ptr)
	RET

#undef res_ptr
#undef x_ptr
#undef y_ptr
#undef X0
#undef X1
#undef Y0
#undef Y1
#undef M0
#undef M1
#undef T0
#undef T1
#undef T2
#undef YDIG

#undef ADD1
#undef ADD1H
#undef ADD2
#undef ADD2H
#undef RED1
#undef RED1H
#undef RED2
#undef RED2H
#undef CAR1
#undef CAR1M

#undef MK0
#undef K0

/* ---------------------------------------*/
// func p256Mul(res, in1, in2 []byte)
#define res_ptr R1
#define x_ptr R2
#define y_ptr R3
#define X0    V0
#define X1    V1
#define Y0    V2
#define Y1    V3
#define M0    V4
#define M1    V5
#define T0    V6
#define T1    V7
#define T2    V8
#define YDIG  V9

#define ADD1  V16
#define ADD11 V17
#define ADD1H V18
#define ADD2  V19
#define ADD2H V20
#define RED1  V21
#define RED2  V22
#define CAR1  V24
#define CAR2  V25

#define SEL1  V29
#define SEL2  V30
#define ZER   V31
  /**
    *        ---+--------+--------+
    *         T2|   T1   |   T0   |
    *        ---+--------+--------+
    *           +--------+------+-+
    *           |  ADD2  |  ADD1+ |
    *           +--------+------+-+
    *         +--------+--------+
    *         | ADD2H  | ADD1H  |
    *         +--------+--------+
    *      ========================
    *      ---+--------+--------+
    *       T2|   T1   |   T0   |
    *      ---+--------+--------+
    *         +--------+------+-+
    *         |  ADD2  |  ADD1+ |
    *         +--------+------+-+
    *       +--------+--------+
    *       | ADD2H  | ADD1H  |
    *       +--------+--------+
    *     =======================
    *    ---+--------+--------+
    *     T2|   T1   |   T0   |
    *    ---+--------+--------+
    *       +--------+--------+
    *       |  RED2  |  RED1  |
    *       +--------+--------+
    *       +--------+
    *       |  RED1  |
    *       +--------+
    */
TEXT ·p256MulSane(SB),NOSPLIT,$0
	MOVD res+0(FP), res_ptr
	MOVD in1+24(FP), x_ptr
	MOVD in2+48(FP), y_ptr

	VZERO T2
	VZERO ZER
	MOVD $p256<>+0x00(SB), R4
	VL	16(R4), M0
	VL	0(R4),  M1
	VL	48(R4), SEL2
	VL	32(R4), SEL1

	VL (1*16)(x_ptr), X0
	VL (0*16)(x_ptr), X1

	//---------------------------------------------------------------------------/

	VLREPF (12+1*16)(y_ptr), YDIG

	VMLF  X0,YDIG, ADD1
	VMLF  X1,YDIG, ADD2
	VMLHF X0,YDIG, ADD1H
	VMLHF X1,YDIG, ADD2H

	VSLDB  $12, ADD2,ADD1, T0
	VSLDB  $12, T2,  ADD2, T1

	VACCQ  T0, ADD1H,CAR1
	VAQ    T0, ADD1H,T0
	VACCCQ T1, ADD2H,CAR1, T2
	VACQ   T1, ADD2H,CAR1, T1

	//---------------------------------------------------

	VLREPF (8+1*16)(y_ptr), YDIG

	VMALF  X0,YDIG, T0, ADD11
	VMALF  X1,YDIG, T1, ADD2
	VMALHF X0,YDIG, T0, ADD1H
	VMALHF X1,YDIG, T1, ADD2H

	VSLDB  $12, ADD2,ADD11, T0
	VSLDB  $12, T2,  ADD2,  T1

	VPERM ADD11,ADD1, SEL1, RED2   // d1 d0 d1 d0
	VPERM ZER,  RED2, SEL2, RED1   // 0  d1 d0  0
	VSQ   RED1, RED2, RED2         // Guaranteed not to underflow

	VACCQ  T0, ADD1H,CAR1
	VAQ    T0, ADD1H,T0
	VACCCQ T1, ADD2H,CAR1, T2
	VACQ   T1, ADD2H,CAR1, T1

	VACCQ  T0, RED1,CAR1
	VAQ    T0, RED1,T0
	VACCCQ T1, RED2,CAR1, CAR2
	VACQ   T1, RED2,CAR1, T1
	VAQ    T2, CAR2,T2

	//---------------------------------------------------

	VLREPF (4+1*16)(y_ptr), YDIG

	VMALF  X0,YDIG, T0, ADD1
	VMALF  X1,YDIG, T1, ADD2
	VMALHF X0,YDIG, T0, ADD1H
	VMALHF X1,YDIG, T1, ADD2H

	VSLDB  $12, ADD2,ADD1, T0
	VSLDB  $12, T2,  ADD2, T1

	VACCQ  T0, ADD1H,CAR1
	VAQ    T0, ADD1H,T0
	VACCCQ T1, ADD2H,CAR1, T2
	VACQ   T1, ADD2H,CAR1, T1

	//---------------------------------------------------

	VLREPF (0+1*16)(y_ptr), YDIG

	VMALF  X0,YDIG, T0, ADD11
	VMALF  X1,YDIG, T1, ADD2
	VMALHF X0,YDIG, T0, ADD1H
	VMALHF X1,YDIG, T1, ADD2H

	VSLDB  $12, ADD2,ADD11, T0
	VSLDB  $12, T2,  ADD2,  T1

	VPERM ADD11,ADD1, SEL1, RED2   // d1 d0 d1 d0
	VPERM ZER,  RED2, SEL2, RED1   // 0  d1 d0  0
	VSQ   RED1, RED2, RED2         // Guaranteed not to underflow

	VACCQ  T0, ADD1H,CAR1
	VAQ    T0, ADD1H,T0
	VACCCQ T1, ADD2H,CAR1, T2
	VACQ   T1, ADD2H,CAR1, T1

	VACCQ  T0, RED1,CAR1
	VAQ    T0, RED1,T0
	VACCCQ T1, RED2,CAR1, CAR2
	VACQ   T1, RED2,CAR1, T1
	VAQ    T2, CAR2,T2

	//---------------------------------------------------

	VLREPF (12+0*16)(y_ptr), YDIG

	VMALF  X0,YDIG, T0, ADD1
	VMALF  X1,YDIG, T1, ADD2
	VMALHF X0,YDIG, T0, ADD1H
	VMALHF X1,YDIG, T1, ADD2H

	VSLDB  $12, ADD2,ADD1, T0
	VSLDB  $12, T2,  ADD2, T1

	VACCQ  T0, ADD1H,CAR1
	VAQ    T0, ADD1H,T0
	VACCCQ T1, ADD2H,CAR1, T2
	VACQ   T1, ADD2H,CAR1, T1

	//---------------------------------------------------

	VLREPF (8+0*16)(y_ptr), YDIG

	VMALF  X0,YDIG, T0, ADD11
	VMALF  X1,YDIG, T1, ADD2
	VMALHF X0,YDIG, T0, ADD1H
	VMALHF X1,YDIG, T1, ADD2H

	VSLDB  $12, ADD2,ADD11, T0
	VSLDB  $12, T2,  ADD2,  T1

	VPERM ADD11,ADD1, SEL1, RED2   // d1 d0 d1 d0
	VPERM ZER,  RED2, SEL2, RED1   // 0  d1 d0  0
	VSQ   RED1, RED2, RED2         // Guaranteed not to underflow

	VACCQ  T0, ADD1H,CAR1
	VAQ    T0, ADD1H,T0
	VACCCQ T1, ADD2H,CAR1, T2
	VACQ   T1, ADD2H,CAR1, T1

	VACCQ  T0, RED1,CAR1
	VAQ    T0, RED1,T0
	VACCCQ T1, RED2,CAR1, CAR2
	VACQ   T1, RED2,CAR1, T1
	VAQ    T2, CAR2,T2

	//---------------------------------------------------

	VLREPF (4+0*16)(y_ptr), YDIG

	VMALF  X0,YDIG, T0, ADD1
	VMALF  X1,YDIG, T1, ADD2
	VMALHF X0,YDIG, T0, ADD1H
	VMALHF X1,YDIG, T1, ADD2H

	VSLDB  $12, ADD2,ADD1, T0
	VSLDB  $12, T2,  ADD2, T1

	VACCQ  T0, ADD1H,CAR1
	VAQ    T0, ADD1H,T0
	VACCCQ T1, ADD2H,CAR1, T2
	VACQ   T1, ADD2H,CAR1, T1

	//---------------------------------------------------

	VLREPF (0+0*16)(y_ptr), YDIG

	VMALF  X0,YDIG, T0, ADD11
	VMALF  X1,YDIG, T1, ADD2
	VMALHF X0,YDIG, T0, ADD1H
	VMALHF X1,YDIG, T1, ADD2H

	VSLDB  $12, ADD2,ADD11, T0
	VSLDB  $12, T2,  ADD2,  T1

	VPERM ADD11,ADD1, SEL1, RED2   // d1 d0 d1 d0
	VPERM ZER,  RED2, SEL2, RED1   // 0  d1 d0  0
	VSQ   RED1, RED2, RED2         // Guaranteed not to underflow

	VACCQ  T0, ADD1H,CAR1
	VAQ    T0, ADD1H,T0
	VACCCQ T1, ADD2H,CAR1, T2
	VACQ   T1, ADD2H,CAR1, T1

	VACCQ  T0, RED1,CAR1
	VAQ    T0, RED1,T0
	VACCCQ T1, RED2,CAR1, CAR2
	VACQ   T1, RED2,CAR1, T1
	VAQ    T2, CAR2,T2

	//---------------------------------------------------

	VSCBIQ  M0,T0, CAR1
	VSQ     M0,T0, ADD1
	VSBCBIQ T1,M1, CAR1, CAR2
	VSBIQ   T1,M1, CAR1, ADD2
	VSBIQ   T2,ZER,CAR2, T2

	//what output to use, ADD2||ADD1 or T1||T0?
	VSEL    T0,ADD1,T2, T0
	VSEL    T1,ADD2,T2, T1

	VST T0, (1*16)(res_ptr)
	VST T1, (0*16)(res_ptr)
	RET

#undef res_ptr
#undef x_ptr
#undef y_ptr
#undef X0
#undef X1
#undef Y0
#undef Y1
#undef M0
#undef M1
#undef T0
#undef T1
#undef T2
#undef YDIG

#undef ADD1
#undef ADD11
#undef ADD1H
#undef ADD2
#undef ADD2H
#undef RED1
#undef RED2
#undef CAR1
#undef CAR2

#undef SEL1
#undef SEL2
#undef ZER

/* ---------------------------------------*/
// func p256Mul(res, in1, in2 []byte)
#define res_ptr R1
#define x_ptr R2
#define y_ptr R3
#define X0    V0
#define X1    V1
#define Y0    V2
#define Y1    V3
#define M0    V4
#define M1    V5
#define T0    V6
#define T1    V7
#define T2    V8
#define YDIG  V9

#define ADD1  V12
#define ADD1H V13
#define ADD2  V14
#define ADD2H V15
#define ADD3  V16
#define ADD3H V17
#define ADD4  V18
#define ADD4H V19
#define RED1  V20
#define RED2  V21
#define RED3  V22
#define CAR1  V23
#define CAR2  V24

#define SEL1  V25
#define SEL2  V26
#define SEL3  V27
#define SEL4  V28
#define SEL5  V29
#define SEL6  V30
#define ZER   V31

  /**
    * To follow the flow of bits, for your own sanity a stiff drink, need you shall.
    * Of a single round, a 'helpful' picture, here is. Meaning, column position has.
    * With you, SIMD be...
    *
    *                                           +--------+--------+
    *                                  +--------|  RED2  |  RED1  |
    *                                  |        +--------+--------+
    *                                  |       ---+--------+--------+
    *                                  |  +---- T2|   T1   |   T0   |--+
    *                                  |  |    ---+--------+--------+  |
    *                                  |  |                            |
    *                                  |  |    ======================= |
    *                                  |  |                            |
    *                                  |  |       +--------+--------+<-+
    *                                  |  +-------|  ADD2  |  ADD1  |--|-----+
    *                                  |  |       +--------+--------+  |     |
    *                                  |  |     +--------+--------+<---+     |
    *                                  |  |     | ADD2H  | ADD1H  |--+       |
    *                                  |  |     +--------+--------+  |       |
    *                                  |  |     +--------+--------+<-+       |
    *                                  |  |     |  ADD4  |  ADD3  |--|-+     |
    *                                  |  |     +--------+--------+  | |     |
    *                                  |  |   +--------+--------+<---+ |     |
    *                                  |  |   | ADD4H  | ADD3H  |------|-+   |(+vzero)
    *                                  |  |   +--------+--------+      | |   V
    *                                  |  | ------------------------   | | +--------+
    *                                  |  |                            | | |  RED3  |  [d0 0 0 d0]
    *                                  |  |                            | | +--------+
    *                                  |  +---->+--------+--------+    | |   |
    *   (T2[1w]||ADD2[4w]||ADD1[3w])   +--------|   T1   |   T0   |    | |   |
    *                                  |        +--------+--------+    | |   |
    *                                  +---->---+--------+--------+    | |   |
    *                                         T2|   T1   |   T0   |----+ |   |
    *                                        ---+--------+--------+    | |   |
    *                                        ---+--------+--------+<---+ |   |
    *                                    +--- T2|   T1   |   T0   |----------+
    *                                    |   ---+--------+--------+      |   |
    *                                    |  +--------+--------+<-------------+
    *                                    |  |  RED2  |  RED1  |-----+    |   | [0 d1 d0 d1] [d0 0 d1 d0]
    *                                    |  +--------+--------+     |    |   |
    *                                    |  +--------+<----------------------+
    *                                    |  |  RED3  |--------------+    |     [0 0 d1 d0]
    *                                    |  +--------+              |    |
    *                                    +--->+--------+--------+   |    |
    *                                         |   T1   |   T0   |--------+
    *                                         +--------+--------+   |    |
    *                                   --------------------------- |    |
    *                                                               |    |
    *                                       +--------+--------+<----+    |
    *                                       |  RED2  |  RED1  |          |
    *                                       +--------+--------+          |
    *                                      ---+--------+--------+<-------+
    *                                       T2|   T1   |   T0   |            (H1P-H1P-H00RRAY!)
    *                                      ---+--------+--------+
    *
    *                                                                *Mi obra de arte de siglo XXI @vpaprots
    *
    *
    * First group is special, doesnt get the two inputs:
    *                                             +--------+--------+<-+
    *                                     +-------|  ADD2  |  ADD1  |--|-----+
    *                                     |       +--------+--------+  |     |
    *                                     |     +--------+--------+<---+     |
    *                                     |     | ADD2H  | ADD1H  |--+       |
    *                                     |     +--------+--------+  |       |
    *                                     |     +--------+--------+<-+       |
    *                                     |     |  ADD4  |  ADD3  |--|-+     |
    *                                     |     +--------+--------+  | |     |
    *                                     |   +--------+--------+<---+ |     |
    *                                     |   | ADD4H  | ADD3H  |------|-+   |(+vzero)
    *                                     |   +--------+--------+      | |   V
    *                                     | ------------------------   | | +--------+
    *                                     |                            | | |  RED3  |  [d0 0 0 d0]
    *                                     |                            | | +--------+
    *                                     +---->+--------+--------+    | |   |
    *   (T2[1w]||ADD2[4w]||ADD1[3w])            |   T1   |   T0   |----+ |   |
    *                                           +--------+--------+    | |   |
    *                                        ---+--------+--------+<---+ |   |
    *                                    +--- T2|   T1   |   T0   |----------+
    *                                    |   ---+--------+--------+      |   |
    *                                    |  +--------+--------+<-------------+
    *                                    |  |  RED2  |  RED1  |-----+    |   | [0 d1 d0 d1] [d0 0 d1 d0]
    *                                    |  +--------+--------+     |    |   |
    *                                    |  +--------+<----------------------+
    *                                    |  |  RED3  |--------------+    |     [0 0 d1 d0]
    *                                    |  +--------+              |    |
    *                                    +--->+--------+--------+   |    |
    *                                         |   T1   |   T0   |--------+
    *                                         +--------+--------+   |    |
    *                                   --------------------------- |    |
    *                                                               |    |
    *                                       +--------+--------+<----+    |
    *                                       |  RED2  |  RED1  |          |
    *                                       +--------+--------+          |
    *                                      ---+--------+--------+<-------+
    *                                       T2|   T1   |   T0   |            (H1P-H1P-H00RRAY!)
    *                                      ---+--------+--------+
    *
    * Last 'group' needs to RED2||RED1 shifted less
    */

TEXT ·p256Mul(SB),NOSPLIT,$0
	MOVD res+0(FP), res_ptr
	MOVD in1+24(FP), x_ptr
	MOVD in2+48(FP), y_ptr

	VZERO ZER
	MOVD $p256mul<>+0x00(SB), R4
	VL	16(R4), M0
	VL	0(R4),  M1
	VL	32(R4), SEL1
	VL	48(R4), SEL2
	VL	64(R4), SEL3
	VL	80(R4), SEL4
	VL	96(R4), SEL5
	VL	112(R4), SEL6

	VL (1*16)(x_ptr), X0
	VL (0*16)(x_ptr), X1

	//---------------------------------------------------

	VLREPF (12+1*16)(y_ptr), YDIG
	VMLF  X0,YDIG, ADD1
	VMLF  X1,YDIG, ADD2
	VMLHF X0,YDIG, ADD1H
	VMLHF X1,YDIG, ADD2H

	VLREPF (8+1*16)(y_ptr), YDIG
	VMALF  X0,YDIG, ADD1H, ADD3
	VMALF  X1,YDIG, ADD2H, ADD4
	VMALHF X0,YDIG, ADD1H, ADD3H
	VMALHF X1,YDIG, ADD2H, ADD4H

	VZERO ZER  // {YDIG,ADD1H,ADD2H}
	VPERM ZER, ADD1, SEL1, RED3  // [d0 0 0 d0]

	VSLDB  $12, ADD2,ADD1, T0
	VSLDB  $12, ZER, ADD2, T1

	VACCQ  T0, ADD3,CAR1
	VAQ    T0, ADD3,T0
	VACCCQ T1, ADD4,CAR1, T2
	VACQ   T1, ADD4,CAR1, T1

    VPERM RED3, T0, SEL2, RED1  // [d0  0 d1 d0]
    VPERM RED3, T0, SEL3, RED2  // [ 0 d1 d0 d1]
    VPERM RED3, T0, SEL4, RED3  // [ 0  0 d1 d0]
    VSQ   RED3,RED2, RED2       // Guaranteed not to underflow

	VSLDB  $12, T1,T0, T0
	VSLDB  $12, T2,T1, T1

	VACCQ  T0, ADD3H,CAR1
	VAQ    T0, ADD3H,T0
	VACCCQ T1, ADD4H,CAR1, T2
	VACQ   T1, ADD4H,CAR1, T1

	//---------------------------------------------------

	VLREPF (4+1*16)(y_ptr), YDIG
	VMALF  X0,YDIG, T0, ADD1
	VMALF  X1,YDIG, T1, ADD2
	VMALHF X0,YDIG, T0, ADD1H
	VMALHF X1,YDIG, T1, ADD2H

	VLREPF (0+1*16)(y_ptr), YDIG
	VMALF  X0,YDIG, ADD1H, ADD3
	VMALF  X1,YDIG, ADD2H, ADD4
	VMALHF X0,YDIG, ADD1H, ADD3H
	VMALHF X1,YDIG, ADD2H, ADD4H

	VZERO ZER  // {YDIG,ADD1H,ADD2H}
	VPERM ZER, ADD1, SEL1, RED3  // [d0 0 0 d0]

	VSLDB  $12, ADD2,ADD1, T0
	VSLDB  $12, T2,  ADD2, T1

	VACCQ  T0, RED1,CAR1
	VAQ    T0, RED1,T0
	VACCCQ T1, RED2,CAR1, T2
	VACQ   T1, RED2,CAR1, T1

	VACCQ  T0, ADD3,CAR1
	VAQ    T0, ADD3,T0
	VACCCQ T1, ADD4,CAR1, CAR2
	VACQ   T1, ADD4,CAR1, T1
	VAQ    T2, CAR2,T2

    VPERM RED3, T0, SEL2, RED1  // [d0  0 d1 d0]
    VPERM RED3, T0, SEL3, RED2  // [ 0 d1 d0 d1]
    VPERM RED3, T0, SEL4, RED3  // [ 0  0 d1 d0]
    VSQ   RED3,RED2, RED2       // Guaranteed not to underflow

	VSLDB  $12, T1,T0, T0
	VSLDB  $12, T2,T1, T1

	VACCQ  T0, ADD3H,CAR1
	VAQ    T0, ADD3H,T0
	VACCCQ T1, ADD4H,CAR1, T2
	VACQ   T1, ADD4H,CAR1, T1

	//---------------------------------------------------

	VLREPF (12+0*16)(y_ptr), YDIG
	VMALF  X0,YDIG, T0, ADD1
	VMALF  X1,YDIG, T1, ADD2
	VMALHF X0,YDIG, T0, ADD1H
	VMALHF X1,YDIG, T1, ADD2H

	VLREPF (8+0*16)(y_ptr), YDIG
	VMALF  X0,YDIG, ADD1H, ADD3
	VMALF  X1,YDIG, ADD2H, ADD4
	VMALHF X0,YDIG, ADD1H, ADD3H
	VMALHF X1,YDIG, ADD2H, ADD4H

	VZERO ZER  // {YDIG,ADD1H,ADD2H}
	VPERM ZER, ADD1, SEL1, RED3  // [d0 0 0 d0]

	VSLDB  $12, ADD2,ADD1, T0
	VSLDB  $12, T2,  ADD2, T1

	VACCQ  T0, RED1,CAR1
	VAQ    T0, RED1,T0
	VACCCQ T1, RED2,CAR1, T2
	VACQ   T1, RED2,CAR1, T1

	VACCQ  T0, ADD3,CAR1
	VAQ    T0, ADD3,T0
	VACCCQ T1, ADD4,CAR1, CAR2
	VACQ   T1, ADD4,CAR1, T1
	VAQ    T2, CAR2,T2

    VPERM RED3, T0, SEL2, RED1  // [d0  0 d1 d0]
    VPERM RED3, T0, SEL3, RED2  // [ 0 d1 d0 d1]
    VPERM RED3, T0, SEL4, RED3  // [ 0  0 d1 d0]
    VSQ   RED3,RED2, RED2       // Guaranteed not to underflow

	VSLDB  $12, T1,T0, T0
	VSLDB  $12, T2,T1, T1

	VACCQ  T0, ADD3H,CAR1
	VAQ    T0, ADD3H,T0
	VACCCQ T1, ADD4H,CAR1, T2
	VACQ   T1, ADD4H,CAR1, T1

	//---------------------------------------------------

	VLREPF (4+0*16)(y_ptr), YDIG
	VMALF  X0,YDIG, T0, ADD1
	VMALF  X1,YDIG, T1, ADD2
	VMALHF X0,YDIG, T0, ADD1H
	VMALHF X1,YDIG, T1, ADD2H

	VLREPF (0+0*16)(y_ptr), YDIG
	VMALF  X0,YDIG, ADD1H, ADD3
	VMALF  X1,YDIG, ADD2H, ADD4
	VMALHF X0,YDIG, ADD1H, ADD3H
	VMALHF X1,YDIG, ADD2H, ADD4H

	VZERO ZER  // {YDIG,ADD1H,ADD2H}
	VPERM ZER, ADD1, SEL1, RED3  // [d0 0 0 d0]

	VSLDB  $12, ADD2,ADD1, T0
	VSLDB  $12, T2,  ADD2, T1

	VACCQ  T0, RED1,CAR1
	VAQ    T0, RED1,T0
	VACCCQ T1, RED2,CAR1, T2
	VACQ   T1, RED2,CAR1, T1

	VACCQ  T0, ADD3,CAR1
	VAQ    T0, ADD3,T0
	VACCCQ T1, ADD4,CAR1, CAR2
	VACQ   T1, ADD4,CAR1, T1
	VAQ    T2, CAR2,T2

    VPERM T0,  RED3, SEL5, RED2  // [d1 d0 d1 d0]
    VPERM T0,  RED3, SEL6, RED1  // [ 0 d1 d0  0]
    VSQ   RED1,RED2, RED2        // Guaranteed not to underflow

	VSLDB  $12, T1,T0, T0
	VSLDB  $12, T2,T1, T1

	VACCQ  T0, ADD3H,CAR1
	VAQ    T0, ADD3H,T0
	VACCCQ T1, ADD4H,CAR1, T2
	VACQ   T1, ADD4H,CAR1, T1

	VACCQ  T0, RED1,CAR1
	VAQ    T0, RED1,T0
	VACCCQ T1, RED2,CAR1, CAR2
	VACQ   T1, RED2,CAR1, T1
	VAQ    T2, CAR2,T2

	//---------------------------------------------------

	VSCBIQ  M0,T0, CAR1
	VSQ     M0,T0, ADD1
	VSBCBIQ T1,M1, CAR1, CAR2
	VSBIQ   T1,M1, CAR1, ADD2
	VSBIQ   T2,ZER,CAR2, T2

	//what output to use, ADD2||ADD1 or T1||T0?
	VSEL    T0,ADD1,T2, T0
	VSEL    T1,ADD2,T2, T1

	VST T0, (1*16)(res_ptr)
	VST T1, (0*16)(res_ptr)
	RET

	
#define IN1H    V0
#define IN1L	V1
#define IN2H	V2
#define IN2L	V3
#define OUTH	V6
#define OUTL	V7
#define POLYH	V4
#define POLYL	V5
	
#define in1x_ptr R1
#define in1y_ptr R2
#define in1z_ptr R3
#define in2x_ptr R4
#define in2y_ptr R5
#define in2z_ptr R6
#define sign     R7
#define sel      R8
#define zero     R9
#define GTMP	 R10

#define p256MulBy2Inline\
	VACCQ IN1L, IN2L, CAR1
	VAQ   IN1L, IN2L, OUTL
	VACCCQ IN1H, IN2H, CAR1, CAR2
	VACQ   IN1H, IN2H, CAR1, OUTH
	

	
TEXT ·p256MulInternal(SB),NOSPLIT,$0
	MOVD $p256mul<>+0x00(SB), R4
	VL	32(R4), SEL1
	VL	48(R4), SEL2
	VL	64(R4), SEL3
	VL	80(R4), SEL4
	VL	96(R4), SEL5
	VL	112(R4), SEL6

	

	//---------------------------------------------------

	VREPF $3,Y0, YDIG
	VMLF  X0,YDIG, ADD1
	VMLF  X1,YDIG, ADD2
	VMLHF X0,YDIG, ADD1H
	VMLHF X1,YDIG, ADD2H

	VREPF  $2, Y0, YDIG
	VMALF  X0,YDIG, ADD1H, ADD3
	VMALF  X1,YDIG, ADD2H, ADD4
	VMALHF X0,YDIG, ADD1H, ADD3H
	VMALHF X1,YDIG, ADD2H, ADD4H


	VPERM ZER, ADD1, SEL1, RED3  // [d0 0 0 d0]

	VSLDB  $12, ADD2,ADD1, T0
	VSLDB  $12, ZER, ADD2, T1

	VACCQ  T0, ADD3,CAR1
	VAQ    T0, ADD3,T0
	VACCCQ T1, ADD4,CAR1, T2
	VACQ   T1, ADD4,CAR1, T1

    VPERM RED3, T0, SEL2, RED1  // [d0  0 d1 d0]
    VPERM RED3, T0, SEL3, RED2  // [ 0 d1 d0 d1]
    VPERM RED3, T0, SEL4, RED3  // [ 0  0 d1 d0]
    VSQ   RED3,RED2, RED2       // Guaranteed not to underflow

	VSLDB  $12, T1,T0, T0
	VSLDB  $12, T2,T1, T1

	VACCQ  T0, ADD3H,CAR1
	VAQ    T0, ADD3H,T0
	VACCCQ T1, ADD4H,CAR1, T2
	VACQ   T1, ADD4H,CAR1, T1

	//---------------------------------------------------

	VREPF  $1,Y0, YDIG
	VMALF  X0,YDIG, T0, ADD1
	VMALF  X1,YDIG, T1, ADD2
	VMALHF X0,YDIG, T0, ADD1H
	VMALHF X1,YDIG, T1, ADD2H

	VREPF  $0,Y0, YDIG
	VMALF  X0,YDIG, ADD1H, ADD3
	VMALF  X1,YDIG, ADD2H, ADD4
	VMALHF X0,YDIG, ADD1H, ADD3H
	VMALHF X1,YDIG, ADD2H, ADD4H


	VPERM ZER, ADD1, SEL1, RED3  // [d0 0 0 d0]

	VSLDB  $12, ADD2,ADD1, T0
	VSLDB  $12, T2,  ADD2, T1

	VACCQ  T0, RED1,CAR1
	VAQ    T0, RED1,T0
	VACCCQ T1, RED2,CAR1, T2
	VACQ   T1, RED2,CAR1, T1

	VACCQ  T0, ADD3,CAR1
	VAQ    T0, ADD3,T0
	VACCCQ T1, ADD4,CAR1, CAR2
	VACQ   T1, ADD4,CAR1, T1
	VAQ    T2, CAR2,T2

    VPERM RED3, T0, SEL2, RED1  // [d0  0 d1 d0]
    VPERM RED3, T0, SEL3, RED2  // [ 0 d1 d0 d1]
    VPERM RED3, T0, SEL4, RED3  // [ 0  0 d1 d0]
    VSQ   RED3,RED2, RED2       // Guaranteed not to underflow

	VSLDB  $12, T1,T0, T0
	VSLDB  $12, T2,T1, T1

	VACCQ  T0, ADD3H,CAR1
	VAQ    T0, ADD3H,T0
	VACCCQ T1, ADD4H,CAR1, T2
	VACQ   T1, ADD4H,CAR1, T1

	//---------------------------------------------------

	VREPF  $3,Y1, YDIG
	VMALF  X0,YDIG, T0, ADD1
	VMALF  X1,YDIG, T1, ADD2
	VMALHF X0,YDIG, T0, ADD1H
	VMALHF X1,YDIG, T1, ADD2H

	VREPF $2,Y1, YDIG
	VMALF  X0,YDIG, ADD1H, ADD3
	VMALF  X1,YDIG, ADD2H, ADD4
	VMALHF X0,YDIG, ADD1H, ADD3H
	VMALHF X1,YDIG, ADD2H, ADD4H

	VPERM ZER, ADD1, SEL1, RED3  // [d0 0 0 d0]

	VSLDB  $12, ADD2,ADD1, T0
	VSLDB  $12, T2,  ADD2, T1

	VACCQ  T0, RED1,CAR1
	VAQ    T0, RED1,T0
	VACCCQ T1, RED2,CAR1, T2
	VACQ   T1, RED2,CAR1, T1

	VACCQ  T0, ADD3,CAR1
	VAQ    T0, ADD3,T0
	VACCCQ T1, ADD4,CAR1, CAR2
	VACQ   T1, ADD4,CAR1, T1
	VAQ    T2, CAR2,T2

    VPERM RED3, T0, SEL2, RED1  // [d0  0 d1 d0]
    VPERM RED3, T0, SEL3, RED2  // [ 0 d1 d0 d1]
    VPERM RED3, T0, SEL4, RED3  // [ 0  0 d1 d0]
    VSQ   RED3,RED2, RED2       // Guaranteed not to underflow

	VSLDB  $12, T1,T0, T0
	VSLDB  $12, T2,T1, T1

	VACCQ  T0, ADD3H,CAR1
	VAQ    T0, ADD3H,T0
	VACCCQ T1, ADD4H,CAR1, T2
	VACQ   T1, ADD4H,CAR1, T1

	//---------------------------------------------------

	VREPF  $1,Y1, YDIG
	VMALF  X0,YDIG, T0, ADD1
	VMALF  X1,YDIG, T1, ADD2
	VMALHF X0,YDIG, T0, ADD1H
	VMALHF X1,YDIG, T1, ADD2H

	VREPF  $0,Y1, YDIG
	VMALF  X0,YDIG, ADD1H, ADD3
	VMALF  X1,YDIG, ADD2H, ADD4
	VMALHF X0,YDIG, ADD1H, ADD3H
	VMALHF X1,YDIG, ADD2H, ADD4H

	VZERO ZER  // {YDIG,ADD1H,ADD2H}
	VPERM ZER, ADD1, SEL1, RED3  // [d0 0 0 d0]

	VSLDB  $12, ADD2,ADD1, T0
	VSLDB  $12, T2,  ADD2, T1

	VACCQ  T0, RED1,CAR1
	VAQ    T0, RED1,T0
	VACCCQ T1, RED2,CAR1, T2
	VACQ   T1, RED2,CAR1, T1

	VACCQ  T0, ADD3,CAR1
	VAQ    T0, ADD3,T0
	VACCCQ T1, ADD4,CAR1, CAR2
	VACQ   T1, ADD4,CAR1, T1
	VAQ    T2, CAR2,T2

    VPERM T0,  RED3, SEL5, RED2  // [d1 d0 d1 d0]
    VPERM T0,  RED3, SEL6, RED1  // [ 0 d1 d0  0]
    VSQ   RED1,RED2, RED2        // Guaranteed not to underflow

	VSLDB  $12, T1,T0, T0
	VSLDB  $12, T2,T1, T1

	VACCQ  T0, ADD3H,CAR1
	VAQ    T0, ADD3H,T0
	VACCCQ T1, ADD4H,CAR1, T2
	VACQ   T1, ADD4H,CAR1, T1

	VACCQ  T0, RED1,CAR1
	VAQ    T0, RED1,T0
	VACCCQ T1, RED2,CAR1, CAR2
	VACQ   T1, RED2,CAR1, T1
	VAQ    T2, CAR2,T2

	//---------------------------------------------------

	VSCBIQ  POLYL,T0, CAR1
	VSQ     POLYL,T0, ADD1
	VSBCBIQ T1,POLYH, CAR1, CAR2
	VSBIQ   T1,POLYH, CAR1, ADD2
	VSBIQ   T2,ZER,CAR2, T2

	//what output to use, ADD2||ADD1 or T1||T0?
	VSEL    T0,ADD1,T2, T0
	VSEL    T1,ADD2,T2, T1

	VST T0, (1*16)(res_ptr)
	VST T1, (0*16)(res_ptr)

	RET


/* ---------------------------------------*/
TEXT p256SubInternal(SB),NOSPLIT,$0
	VSCBIQ   IN2L, IN1L, CAR1
	VSQ	 IN2L, IN1L, OUTL
	VSBIQ    IN2H, IN1H, CAR1, OUTH

	VACCQ    OUTL,POLYL,CAR1
	VAQ	 OUTL,POLYL,T0
	VACCCQ   OUTH,POLYH,CAR1,CAR2
	VACQ	 OUTH,POLYH,CAR1,T1

	VSQ      CAR2,ZER,CAR2
	VSEL     OUTL,T0,CAR2,OUTL
	VSEL     OUTH,T1,CAR2,OUTH

	RET


	
/* ---------------------------------------*/
#define RS R7
#define x1in(off) (32*0 + off)(RS)
#define y1in(off) (32*1 + off)(RS)
#define z1in(off) (32*2 + off)(RS)
#define x2in(off) (32*3 + off)(RS)
#define y2in(off) (32*4 + off)(RS)
#define z2in(off) (32*5 + off)(RS)

#define xout(off) (32*6 + off)(RS)
#define yout(off) (32*7 + off)(RS)
#define zout(off) (32*8 + off)(RS)

#define u1(off)    (32*9 + off)(RS)
#define u2(off)    (32*10 + off)(RS)
#define s1(off)    (32*11 + off)(RS)
#define s2(off)    (32*12 + off)(RS)
#define z1sqr(off)  z1sqrv-(16-off)(SP)
#define z2sqr(off)  z2sqrv-(32*14 + off)(SP)
#define h(off)     (32*15 + off)(RS)
#define r(off)     (32*16 + off)(RS)
#define hsqr(off)  (32*17 + off)(RS)
#define rsqr(off)  (32*18 + off)(RS)
#define hcub(off)  (32*19 + off)(RS)
#define rptr       (32*20)(RS)

/*#define in1x_ptr  R1
#define	in1y_ptr  R2
#define in1z_ptr  R3
#define in2x_ptr  R4
#define in2y_ptr  R5
#define in2z_ptr  R6
*/


// func p256PointAddAffineAsm(resx,resy,resz, in1x, in1y, in1z, in2x, in2y, in2z []byte)
// FIXME: Does the $672-72 need to be updated for longer parm list??
TEXT ·p256PointAddAsm2(SB),0,$672-72
	MOVD  in1x+(3*24)(FP), in1x_ptr
	MOVD  in1y+(4*24)(FP), in1y_ptr
	MOVD  in1z+(5*24)(FP), in1z_ptr
	MOVD  in2x+(6*24)(FP), in2x_ptr
	MOVD  in2y+(7*24)(FP), in2y_ptr
	MOVD  in2z+(8*24)(FP), in2z_ptr

	VGBM $0x0fff,POLYL
	MOVD p256<>(SB),GTMP
	VL   0x10(GTMP),POLYH
	VZERO ZER

	ADD $8, R15, RS
	
	VL  0(in2z_ptr), IN1H
	VL  16(in2z_ptr),IN1L
	VLR IN1H, IN2H
	VLR IN1L, IN2L
	CALL p256MulInternal(SB)    // z2^2
	VST OUTH, z2sqr(0)
	VST OUTL, z2sqr(16)
	VLR OUTH, IN1H
	VLR OUTL, IN1L
	CALL p256MulInternal(SB)    // z2^3
	VL  0(in1y_ptr), IN1H
	VL  16(in1y_ptr), IN1L
	VLR OUTH, IN2H
	VLR OUTL, IN2L
	CALL p256MulInternal(SB)   // z2^3 * y1
	VST OUTH, s1(0)
	VST OUTL, s1(16)
	VL  0(in1z_ptr), IN1H
	VL  16(in1z_ptr), IN1L
	VLR IN1H, IN2H
	VLR IN1L, IN2L
	CALL p256MulInternal(SB) 	//z1^2
	VST OUTH, z1sqr(0)
	VST OUTL, z1sqr(16)
	VLR OUTH, IN2H
	VLR OUTL, IN2L
	CALL p256MulInternal(SB)	//z1^3
	VL  0(in1y_ptr), IN1H
	VL  16(in2y_ptr), IN1L
	VLR OUTH, IN2H
	VLR OUTL, IN2L
	CALL p256MulInternal(SB)	//z1^3 * y2
	VST OUTH, s2(0)
	VST OUTL, s2(16)
	
	VLR OUTH, IN1H
	VLR OUTL, IN1L
	VL  s1(0), IN2H
	VL  s1(16), IN2L
	CALL p256SubInternal(SB)	// s2-s1
	VST OUTH, r(0)
	VST OUTL, r(16)

	VL z2sqr(0), IN1H
	VL z2sqr(16*1), IN1L
	VL 0(in1x_ptr), IN2H
	VL 16(in1x_ptr), IN2L
	CALL p256MulInternal(SB)	// x1 * z2^2
	VST OUTH, u1(0)
	VST OUTL, u1(16)

	VL z1sqr(0), IN1H
	VL z1sqr(16), IN1L
	VL 0(in2x_ptr), IN2H
	VL 16(in2x_ptr), IN2L
	CALL p256MulInternal(SB)	// u1 = x1 * z2^s
	VST OUTH, u2(0)
	VST OUTL, u2(16)

	VLR OUTH, IN1H
	VLR OUTL, IN1L
	VL u1(0), IN2H
	VL u1(16), IN2L
	CALL p256SubInternal(SB)	// h = u2 - u1
	VST OUTH, h(0)
	VST OUTL, h(16)

	VL r(0), IN1H
	VL r(16), IN1L
	//may be faster to just load so there is no dependency
	//Then again these go away once there is an optimized squaring
	VLR IN1H, IN2H
	VLR IN1L, IN2L
	CALL p256MulInternal(SB) 	// rsqr = r^2
	VST OUTH, rsqr(0)
	VST OUTL, rsqr(16)

	VL h(0), IN1H
	VL h(16), IN1L
	VLR IN1H, IN2H
	VLR IN1L, IN2L
	CALL p256MulInternal(SB)	// hsqr = h^2
	VST OUTH, hsqr(0)
	VST OUTL, hsqr(16)

	VLR OUTH, IN2H
	VLR OUTL, IN2L
	CALL p256MulInternal(SB)	// hcub = h^3
	VST OUTH, hcub(0)
	VST OUTL, hcub(16)

	VL s1(0), IN1H
	VL s1(16), IN1L
	VLR OUTH, IN2H
	VLR OUTL, IN2L
	CALL p256MulInternal(SB)	// s2 = s1*h^3
	VST OUTH, s2(0)
	VST OUTL, s2(16)

	VL 0(in1z_ptr), IN1H
	VL 16(in1z_ptr), IN1L
	VL 0(in2z_ptr), IN2H
	VL 16(in2z_ptr), IN2L
	CALL p256MulInternal(SB)

	VL h(0), IN1H
	VL h(16), IN1L
	VLR OUTH, IN2H
	VLR OUTL, IN2L
	CALL p256MulInternal(SB)
	VST OUTH, zout(0)
	VST OUTL, zout(16)

	VL hsqr(0), IN1H
	VL hsqr(16), IN1L
	VL u1(0), IN2H
	VL u1(16), IN2L
	CALL p256MulInternal(SB)
	VST OUTH, u2(0)
	VST OUTL, u2(16)

	



	
	RET	

// func p256PointAddAffineAsm(res, in1, in2 []uint64, sign, sel, zero int)
TEXT ·p256PointAddAffineAsm(SB),0,$512-96
	// Move input to stack in order to free registers
	MOVD res+0(FP), R1
	MOVD in1+24(FP), R2
	MOVD in2+48(FP), R3
//	MOVD sign+72(FP), R4
//	MOVD sel+80(FP), R5
//	MOVD zero+88(FP), R6

	VL   16(R3), IN2L
	VL   0(R3),  IN2H

	VGBM $0x0fff,POLYL
	MOVD p256<>(SB),GTMP
	VL   0x10(GTMP),POLYH
	VZERO ZER

	RET
