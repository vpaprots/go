// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.


#include "textflag.h"

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
DATA p256mul<>+0x80(SB)/8, $0x00000000fffffffe // (1*2^256)%P256
DATA p256mul<>+0x88(SB)/8, $0xffffffffffffffff // (1*2^256)%P256
DATA p256mul<>+0x90(SB)/8, $0xffffffff00000000 // (1*2^256)%P256
DATA p256mul<>+0x98(SB)/8, $0x0000000000000001 // (1*2^256)%P256
GLOBL p256ordK0<>(SB), 8, $4
GLOBL p256ord<>(SB), 8, $32
GLOBL p256<>(SB), 8, $64
GLOBL p256mul<>(SB), 8, $160


/* ---------------------------------------*/
// func p256OrdMul(res, in1, in2 []byte)
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

#define p256SubInternal(T1,T0, X1,X0, Y1,Y0) \
	VZERO   ZER\
	VSCBIQ  Y0,X0, CAR1\
	VSQ     Y0,X0, T0\
	VSBCBIQ X1,Y1, CAR1, SEL1\
	VSBIQ   X1,Y1, CAR1, T1\
	VSQ     SEL1,ZER, SEL1\
	\
	VACCQ    T0,PL,CAR1\
	VAQ	     T0,PL,TT0\
	VACQ	 T1,PH,CAR1,TT1\
	\
	VSEL     T0,TT0,SEL1,T0\
	VSEL     T1,TT1,SEL1,T1\

#define p256AddInternal(T1,T0, X1,X0, Y1,Y0)\
	VACCQ X0, Y0, CAR1\
	VAQ   X0, Y0, T0\
	VACCCQ X1, Y1, CAR1, T2\
	VACQ   X1, Y1, CAR1, T1\
	\
	VZERO   ZER\
	VSCBIQ  PL,T0, CAR1\
	VSQ     PL,T0, TT0\
	VSBCBIQ T1,PH, CAR1, CAR2\
	VSBIQ   T1,PH, CAR1, TT1\
	VSBIQ   T2,ZER,CAR2, SEL1\
	\
	VSEL    T0,TT0,SEL1, T0\
	VSEL    T1,TT1,SEL1, T1

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

// p256PointAddAffineAsm(P3, P1, P2 *p256Point, sign, sel, zero int)
#define P3ptr   R1
#define P1ptr   R2
#define P2ptr   R3
#define CPOOL   R4

// Temporaries in REGs
#define Y2L    V15
#define Y2H    V16
#define T1L    V17
#define T1H    V18
#define T2L    V19
#define T2H    V20
#define T3L    V21
#define T3H    V22
#define T4L    V23
#define T4H    V24

// Temps for Sub and Add
#define TT0  V11
#define TT1  V12
#define T2   V13

// P256Mul Parameters
#define X0    V0
#define X1    V1
#define Y0    V2
#define Y1    V3
#define T0    V4
#define T1    V5

#define PL    V30
#define PH    V31

// Names for zero/sel selects
#define X1L    V0
#define X1H    V1
#define Y1L    V2 // P256MulParmY
#define Y1H    V3 // P256MulParmY
#define Z1L    V4
#define Z1H    V5
#define X2L    V0
#define X2H    V1
#define Z2L    V4
#define Z2H    V5
#define X3L    V17 //T1L
#define X3H    V18 //T1H
#define Y3L    V21 //T3L
#define Y3H    V22 //T3H
#define Z3L    V23
#define Z3H    V24

#define ZER   V26
#define SEL1  V27
#define CAR1  V28
#define CAR2  V29

TEXT ·p256PointAddAffineAsm(SB),NOSPLIT,$0
	MOVD P3+0(FP),  P3ptr
	MOVD P1+8(FP),  P1ptr
	MOVD P2+16(FP), P2ptr

	MOVD $p256mul<>+0x00(SB), CPOOL
	VL	16(CPOOL), PL
	VL	0(CPOOL),  PH

//	if (sign == 1) {
//		Y2 = fromBig(new(big.Int).Mod(new(big.Int).Sub(p256.P, new(big.Int).SetBytes(Y2)), p256.P)) // Y2  = P-Y2
//	}

	VL	32(P2ptr), Y2H
	VL	48(P2ptr), Y2L

	VLREPG sign+24(FP), SEL1
	VZERO  ZER
	VCEQG  SEL1,ZER,SEL1

	VSCBIQ  Y2L, PL, CAR1
	VSQ     Y2L, PL, T1L
	VSBIQ   PH,  Y2H, CAR1, T1H

	VSEL Y2L,T1L,SEL1,Y2L
	VSEL Y2H,T1H,SEL1,Y2H

	/**
	 * Three operand formula:
	 * Source: 2004 Hankerson–Menezes–Vanstone, page 91.
	 */
	// X=Z1; Y=Z1; MUL; T-   // T1 = Z1²      T1
		VL	 64(P1ptr), X1  //Z1H
		VL	 80(P1ptr), X0  //Z1L
		VLR  X0, Y0
		VLR  X1, Y1
		CALL p256MulInternal(SB)

	// X=T ; Y-  ; MUL; T2=T // T2 = T1*Z1    T1   T2
		VLR T0, X0
		VLR T1, X1
		CALL p256MulInternal(SB)
		VLR T0, T2L
		VLR T1, T2H

	// X-  ; Y=X2; MUL; T1=T // T1 = T1*X2    T1   T2
		VL	 0(P2ptr), Y1 //X2H
		VL	16(P2ptr), Y0 //X2L
		CALL p256MulInternal(SB)
		VLR T0, T1L
		VLR T1, T1H

	// X=T2; Y=Y2; MUL; T-   // T2 = T2*Y2    T1   T2
		VLR T2L, X0
		VLR T2H, X1
		VLR Y2L, Y0
		VLR Y2H, Y1
		CALL p256MulInternal(SB)

	// SUB(T2<T-Y1)          // T2 = T2-Y1    T1   T2
		VL	32(P1ptr), Y1H
		VL	48(P1ptr), Y1L
		p256SubInternal(T2H,T2L,T1,T0,Y1H,Y1L)

	// SUB(Y<T1-X1)          // T1 = T1-X1    T1   T2
		VL	 0(P1ptr), X1H
		VL	16(P1ptr), X1L
		p256SubInternal(Y1,Y0,T1H,T1L,X1H,X1L)

	// X=Z1; Y- ;  MUL; Z3:=T// Z3 = Z1*T1         T2
		VL	 64(P1ptr), X1  //Z1H
		VL	 80(P1ptr), X0  //Z1L
		CALL p256MulInternal(SB)
		VST T1, 64(P3ptr)
		VST T0, 80(P3ptr)

	// X=Y;  Y- ;  MUL; X=T  // T3 = T1*T1         T2
		VLR  Y0, X0
		VLR  Y1, X1
		CALL p256MulInternal(SB)
		VLR  T0, X0
		VLR  T1, X1

	// X- ;  Y- ;  MUL; T4=T // T4 = T3*T1         T2        T4
		CALL p256MulInternal(SB)
		VLR  T0, T4L
		VLR  T1, T4H

	// X- ;  Y=X1; MUL; T3=T // T3 = T3*X1         T2   T3   T4
		VL	 0(P1ptr), Y1 //X1H
		VL	16(P1ptr), Y0 //X1L
		CALL p256MulInternal(SB)
		VLR  T0, T3L
		VLR  T1, T3H

	// ADD(T1<T+T)           // T1 = T3+T3    T1   T2   T3   T4
		p256AddInternal(T1H,T1L, T1,T0,T1,T0)

	// X=T2; Y=T2; MUL; T-   // X3 = T2*T2    T1   T2   T3   T4
		VLR  T2L, X0
		VLR  T2H, X1
		VLR  T2L, Y0
		VLR  T2H, Y1
		CALL p256MulInternal(SB)

	// SUB(T<T-T1)           // X3 = X3-T1    T1   T2   T3   T4  (T1 = X3)
		p256SubInternal(T1,T0,T1,T0,T1H,T1L)

	// SUB(T<T-T4) X3:=T     // X3 = X3-T4         T2   T3   T4
		p256SubInternal(T1,T0,T1,T0,T4H,T4L)
		VLR T0, X3L
		VLR T1, X3H

	// SUB(X<T3-T)           // T3 = T3-X3         T2   T3   T4
		p256SubInternal(X1,X0,T3H,T3L,T1,T0)

	// X- ;  Y- ;  MUL; T3=T // T3 = T3*T2         T2   T3   T4
		CALL p256MulInternal(SB)
		VLR  T0, T3L
		VLR  T1, T3H

	// X=T4; Y=Y1; MUL; T-   // T4 = T4*Y1              T3   T4
		VLR  T4L, X0
		VLR  T4H, X1
		VL	32(P1ptr), Y1 //Y1H
		VL	48(P1ptr), Y0 //Y1L
		CALL p256MulInternal(SB)

	// SUB(T<T3-T) Y3:=T     // Y3 = T3-T4              T3   T4  (T3 = Y3)
		p256SubInternal(Y3H,Y3L,T3H,T3L,T1,T0)

	VL 64(P3ptr), Z3H
	VL 80(P3ptr), Z3L
	// P3 = {x:{T1H||T1L},y:{T3H||T3L},z{Z3H||Z3L}}

//	if (sel == 0) {
//		copy(P3.x[:], X1)
//		copy(P3.y[:], Y1)
//		copy(P3.z[:], Z1)
//	}

	VL	 0(P1ptr), X1H
	VL	16(P1ptr), X1L
	//Y1 already loaded, left over from addition
	VL	64(P1ptr), Z1H
	VL	80(P1ptr), Z1L

	VLREPG sel+32(FP), SEL1
	VZERO  ZER
	VCEQG  SEL1,ZER,SEL1

	VSEL X1L,X3L,SEL1,X3L
	VSEL X1H,X3H,SEL1,X3H
	VSEL Y1L,Y3L,SEL1,Y3L
	VSEL Y1H,Y3H,SEL1,Y3H
	VSEL Z1L,Z3L,SEL1,Z3L
	VSEL Z1H,Z3H,SEL1,Z3H

//	if (zero == 0) {
//		copy(P3.x[:], X2)
//		copy(P3.y[:], Y2)
//		copy(P3.z[:], []byte{0x00, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xfe, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
//			0xff, 0xff, 0xff, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01})  //(p256.z*2^256)%p
//	}
	VL	 0(P2ptr), X2H
	VL	16(P2ptr), X2L
	//Y2 already loaded
	VL	128(CPOOL), Z2H
	VL	144(CPOOL), Z2L

	VLREPG zero+40(FP), SEL1
	VZERO  ZER
	VCEQG  SEL1,ZER,SEL1

	VSEL X2L,X3L,SEL1,X3L
	VSEL X2H,X3H,SEL1,X3H
	VSEL Y2L,Y3L,SEL1,Y3L
	VSEL Y2H,Y3H,SEL1,Y3H
	VSEL Z2L,Z3L,SEL1,Z3L
	VSEL Z2H,Z3H,SEL1,Z3H

	// All done, store out the result!!!
	VST X3H,  0(P3ptr)
	VST X3L, 16(P3ptr)
	VST Y3H, 32(P3ptr)
	VST Y3L, 48(P3ptr)
	VST Z3H, 64(P3ptr)
	VST Z3L, 80(P3ptr)

	RET

	#undef TT1
#undef TT0
