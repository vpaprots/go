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
GLOBL p256const0<>(SB), 8, $8
GLOBL p256const1<>(SB), 8, $8
GLOBL p256ordK0<>(SB), 8, $4
GLOBL p256ord<>(SB), 8, $32
GLOBL p256<>(SB), 8, $64
GLOBL p256mul<>(SB), 8, $128


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
#define x_ptr   R2
#define y_ptr   R3
#define CPOOL   R4

// Parameters
#define X0    V0
#define X1    V1
#define Y0    V2
#define Y1    V3
#define T0    V4
#define T1    V5
#define P0    V30
#define P1    V31

// Constants
#define SEL1  V26
#define SEL2  V27
#define SEL3  V28
#define SEL4  V29
#define SEL5  V28  // Overloaded with SEL4
#define SEL6  V29  // Overloaded with SEL5

// Temporaries
#define YDIG  V6   //Overloaded with CAR2, ZER
#define ADD1H V7   //Overloaded with ADD3H
#define ADD2H V8   //Overloaded with ADD4H
#define ADD3  V9
#define ADD4  V10
#define RED1  V11  //Overloaded with CAR2
#define RED2  V12
#define RED3  V13
#define T2    V14
// Overloaded temporaries
#define ADD1  V4  //Overloaded with T0
#define ADD2  V5  //Overloaded with T1
#define ADD3H V7  //Overloaded with ADD1H
#define ADD4H V8  //Overloaded with ADD2H
#define ZER   V6  //Overloaded with YDIG, CAR2
#define CAR1  V6  //Overloaded with YDIG, ZER
#define CAR2  V11 //Overloaded with RED1


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

 TEXT p256MulInternal(SB),NOSPLIT,$0
	VL	32(CPOOL), SEL1
	VL	48(CPOOL), SEL2
	VL	64(CPOOL), SEL3
	VL	80(CPOOL), SEL4

	//---------------------------------------------------

	VREPF $3,Y0,   YDIG
	VMLHF X0,YDIG, ADD1H
	VMLHF X1,YDIG, ADD2H
	VMLF  X0,YDIG, ADD1
	VMLF  X1,YDIG, ADD2

	VREPF  $2,Y0,   YDIG
	VMALF  X0,YDIG, ADD1H, ADD3
	VMALF  X1,YDIG, ADD2H, ADD4
	VMALHF X0,YDIG, ADD1H, ADD3H  // ADD1H Free
	VMALHF X1,YDIG, ADD2H, ADD4H  // ADD2H Free

	VZERO ZER
	VPERM ZER, ADD1, SEL1, RED3  // [d0 0 0 d0]

	VSLDB  $12, ADD2,ADD1, T0  //ADD1 Free
	VSLDB  $12, ZER, ADD2, T1  //ADD2 Free

	VACCQ  T0, ADD3,CAR1
	VAQ    T0, ADD3,T0         //ADD3 Free
	VACCCQ T1, ADD4,CAR1, T2
	VACQ   T1, ADD4,CAR1, T1   //ADD4 Free

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

	VREPF  $1,Y0,   YDIG
	VMALHF X0,YDIG, T0, ADD1H
	VMALHF X1,YDIG, T1, ADD2H
	VMALF  X0,YDIG, T0, ADD1      // T0 Free->ADD1
	VMALF  X1,YDIG, T1, ADD2      // T1 Free->ADD2

	VREPF  $0,Y0,   YDIG
	VMALF  X0,YDIG, ADD1H, ADD3
	VMALF  X1,YDIG, ADD2H, ADD4
	VMALHF X0,YDIG, ADD1H, ADD3H  // ADD1H Free->ADD3H
	VMALHF X1,YDIG, ADD2H, ADD4H  // ADD2H Free->ADD4H , YDIG Free->ZER

	VZERO ZER
	VPERM ZER, ADD1, SEL1, RED3  // [d0 0 0 d0]

	VSLDB  $12, ADD2,ADD1, T0     // ADD1 Free->T0
	VSLDB  $12, T2,  ADD2, T1     // ADD2 Free->T1, T2 Free

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

	VREPF  $3,Y1,   YDIG
	VMALHF X0,YDIG, T0, ADD1H
	VMALHF X1,YDIG, T1, ADD2H
	VMALF  X0,YDIG, T0, ADD1
	VMALF  X1,YDIG, T1, ADD2

	VREPF  $2,Y1,   YDIG
	VMALF  X0,YDIG, ADD1H, ADD3
	VMALF  X1,YDIG, ADD2H, ADD4
	VMALHF X0,YDIG, ADD1H, ADD3H  // ADD1H Free
	VMALHF X1,YDIG, ADD2H, ADD4H  // ADD2H Free

	VZERO ZER
	VPERM ZER, ADD1, SEL1, RED3  // [d0 0 0 d0]

	VSLDB  $12, ADD2,ADD1, T0     // ADD1 Free
	VSLDB  $12, T2,  ADD2, T1     // ADD2 Free

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

	VL	96(CPOOL),  SEL5
	VL	112(CPOOL), SEL6

	VREPF  $1,Y1,   YDIG
	VMALHF X0,YDIG, T0, ADD1H
	VMALHF X1,YDIG, T1, ADD2H
	VMALF  X0,YDIG, T0, ADD1
	VMALF  X1,YDIG, T1, ADD2

	VREPF  $0,Y1, YDIG
	VMALF  X0,YDIG, ADD1H, ADD3
	VMALF  X1,YDIG, ADD2H, ADD4
	VMALHF X0,YDIG, ADD1H, ADD3H
	VMALHF X1,YDIG, ADD2H, ADD4H

	VZERO ZER
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

	VZERO   RED3
	VSCBIQ  P0,T0, CAR1
	VSQ     P0,T0, ADD1H
	VSBCBIQ T1,P1, CAR1, CAR2
	VSBIQ   T1,P1, CAR1, ADD2H
	VSBIQ   T2,RED3,CAR2, T2

	//what output to use, ADD2H||ADD1H or T1||T0?
	VSEL    T0,ADD1H,T2, T0
	VSEL    T1,ADD2H,T2, T1
 	RET
#undef res_ptr
#undef x_ptr
#undef y_ptr
#undef CPOOL

#undef X0
#undef X1
#undef Y0
#undef Y1
#undef T0
#undef T1
#undef P0
#undef P1

#undef SEL1
#undef SEL2
#undef SEL3
#undef SEL4
#undef SEL5
#undef SEL6

#undef YDIG
#undef ADD1H
#undef ADD2H
#undef ADD3
#undef ADD4
#undef RED1
#undef RED2
#undef RED3
#undef T2
#undef ADD1
#undef ADD2
#undef ADD3H
#undef ADD4H
#undef ZER
#undef CAR1
#undef CAR2

/* ---------------------------------------*/
// func p256Mul(res, in1, in2 []byte)
#define res_ptr R1
#define x_ptr   R2
#define y_ptr   R3
#define CPOOL   R4

// Parameters
#define X0    V0
#define X1    V1
#define Y0    V2
#define Y1    V3
#define T0    V4
#define T1    V5

// Constants
#define P0    V30
#define P1    V31
 TEXT ·p256Mul(SB),NOSPLIT,$0
	MOVD res+0(FP), res_ptr
	MOVD in1+24(FP), x_ptr
	MOVD in2+48(FP), y_ptr

	VL (1*16)(x_ptr), X0
	VL (0*16)(x_ptr), X1
	VL (1*16)(y_ptr), Y0
	VL (0*16)(y_ptr), Y1

	MOVD $p256mul<>+0x00(SB), CPOOL
	VL	16(CPOOL), P0
	VL	0(CPOOL),  P1

	CALL p256MulInternal(SB)

	VST T0, (1*16)(res_ptr)
	VST T1, (0*16)(res_ptr)
	RET
